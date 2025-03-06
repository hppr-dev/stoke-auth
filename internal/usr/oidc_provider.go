package usr

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"slices"
	"stoke/internal/tel"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
)

// https://openid.net/specs/openid-connect-core-1_0.html
// Need to verify the user is who they say they are by sending username/password to oidc
// The received JWT has claims that are used to build a stoke token
// Should be able to:
// 		* Passthrough claims from provider to issued token
//    * Map claims from provider to issued token
//
// The process should go as follows
// 1. User goes to /oidc/{provider-name}
// 2. Stoke sends an authentication request to the provider
// 3. Provider responds with a combination of  the following:
//     i.   If the provider is using "code" in it's AuthFlowType then an access code for the TokenURL is returned
//     ii.  If the provider is using "id_token" in it's AuthFlowType then an identity token is returned
//     iii. If the provider is using "token" in it's AuthFlowType then an access token for the UserInfoURL is returned
// 4. Stoke sends a token request to the provider using the access code
// 5. Provider returns an id_token and/or an access token
// 6. Stoke sends a UserInfo request to the provider using the access token
// 7. Stoke uses claims to update database
//
// Depending on the AuthFlowType and the ClaimSource, some steps may be skipped.
// The ClaimSource defines the goal of the process, so once it has been obtained no additional steps are needed


var (
	TokenRetrievalError = errors.New("Could not retrieve token from token url")
	ProviderAuthError = errors.New("Could not authenticate user")
)

type postRedirectData struct {
	IDToken string
	AccessToken string
	NextURL string
	LocalStorage bool
	ChildWindow bool
}

var (
	POST_TEMPLATE = `
<html>
	<head>
		<script lang="javascript">
			window.onload = function() {
			{{ if .LocalStorage }}
				window.sessionStorage.setItem("id_token", "{{ .IDToken }}")
				window.sessionStorage.setItem("access_token", "{{ .AccessToken }}")
				{{ if ne .NextURL "" }}
					window.location = "{{ .NextURL }}"
				{{ end }}
			{{ else if .ChildWindow }}
				window.opener.sessionStorage.setItem("id_token", "{{ .IDToken }}")
				window.opener.sessionStorage.setItem("access_token", "{{ .AccessToken }}")
				window.close()
			{{ end }}
			}
		</script>
	</head>
</html>
`
)

type oidcAuthRequest struct {
	Scope        string
	ClientID     string
	RedirectURI  string
	ExtraArgs    map[string]string

	responseType string
}

type oidcUserProvider struct {
	// Unique name of this user provider
	Name string
	// URL to use to authenticate users with the provider
	AuthenticationURL *url.URL
	// URL to use to get the tokens from the provider. Requires an auth code
	TokenURL          *url.URL
	// URL to use to get UserInfo. Requires an access token
	UserInfoURL       *url.URL

	// Whether to retreive user claims from the id token or the UserInfo endpoint
	ClaimSource       ClaimSourceType
	// The authentication request to use when authenticating to the AuthenticationURL
	Request           oidcAuthRequest
	// Salt to be used to generate state hashes. Must be the same accross shared stoke servers
	StateSalt         []byte
	// Authenctication flow type
	FlowType          AuthFlowType
	// Client Secret, given by provider
	ClientSecret      string

	postRedirectTempl *template.Template
}

func NewOIDCUserProvider(
	name, scopes, redirectURI,
	stateSecret, clientID, clientSecret string,
	extraArgs map[string]string,
	tokenURL, authURL, userInfoURL *url.URL,
	mux *http.ServeMux,
	flowType AuthFlowType,
	claimSource ClaimSourceType,
) *oidcUserProvider {
	stateSalt := []byte(stateSecret)
	if stateSecret == "" {
		stateSalt = make([]byte, 16)
		rand.Read(stateSalt)
	}
	prt, _ := template.New("postRedirect").Parse(POST_TEMPLATE)
	provider := &oidcUserProvider{
		Name: name,
		TokenURL : tokenURL,
		AuthenticationURL: authURL,
		UserInfoURL: userInfoURL,
		Request: oidcAuthRequest{
			Scope: scopes,
			RedirectURI: redirectURI,
			ClientID: clientID,
			ExtraArgs: extraArgs,
			responseType: flowType.String(),
		},
		StateSalt: stateSalt,
		FlowType: flowType,
		ClaimSource: claimSource,
		ClientSecret: clientSecret,
		postRedirectTempl: prt,
	}

	mux.Handle("/oidc/" + name, provider)

	return provider
}

// Handles redirect to and from provider
// 2. Stoke sends an authentication request to the provider (by redirecting to the AuthenticationURL with params)
// 3. Provider responds with a combination of  the following:
//   i.   If the provider is using "code" in it's AuthFlowType then an access code for the TokenURL is returned
//   ii.  If the provider is using "id_token" in it's AuthFlowType then an identity token is returned
//   iii. If the provider is using "token" in it's AuthFlowType then an access token for the UserInfoURL is returned
// 4. Stoke sends a token request to the provider using the access code
// 5. Provider returns an id_token and/or an access token
func (o *oidcUserProvider) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	urlParams := req.URL.Query()
	urlState := urlParams.Get("state")

	ctx := req.Context()
	logger := zerolog.Ctx(ctx).With().
		Str("component", "OIDCProvider").
		Str("state", urlState).
		Stringer("claim_source", o.ClaimSource).
		Stringer("flow_type", o.FlowType).
		Logger()

	ctx, span := tel.GetTracer().Start(ctx, "oidcUserProvider.UpdateUserClaims")
	defer span.End()

	if urlError := urlParams.Get("error"); urlError != "" {
		res.WriteHeader(http.StatusNotAcceptable)
		return
	}

	// State is used to determine which side of the process we are on.
	if urlState == "" {
		nonce := newNonce()
		logger.Info().Bytes("nonce", nonce).Msg("Redirecting to AuthURL")
		// Save the generated nonce
		http.SetCookie(res, &http.Cookie{
			Name:  o.cookieName("nonce"),
			Value: base64.URLEncoding.EncodeToString(nonce),
		})
		// Save the supplied transfer method
		xferMethod := "local"
		if urlXfer := urlParams.Get("xfer"); urlXfer != "" && urlXfer == "window"{
			xferMethod = urlXfer
		}
		http.SetCookie(res, &http.Cookie{
			Name:  o.cookieName("xfer"),
			Value: xferMethod,
		})
		// Save the next url if it was specified
		if urlNext := urlParams.Get("next"); urlNext != "" {
			http.SetCookie(res, &http.Cookie{
				Name:  o.cookieName("next-url"),
				Value: urlNext,
			})
		}
		http.Redirect(res, req, o.addParamsToAuthURL(req, nonce).String(), http.StatusTemporaryRedirect)
		return
	}

	// Validate that the state was generated by us for the given address
	if !o.hasValidState(req) {
		logger.Error().Msg("Could not verify state")
		res.WriteHeader(http.StatusConflict)
		return
	}

	accessToken := urlParams.Get("token")
	idToken := urlParams.Get("id_token")
	authCode := urlParams.Get("code")

	idToken, accessToken, err := o.getTokens(idToken, accessToken, authCode, ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Could not get tokens")
		res.WriteHeader(http.StatusInternalServerError)
		return 
	}

	xferMethod := "local"
	next := ""
	nonceStr := ""
	if nextCookie, err := req.Cookie(o.cookieName("next-url")) ; err == nil {
		next = nextCookie.Value
	}
	if xferCookie, err := req.Cookie(o.cookieName("xfer")); err == nil {
		xferMethod = xferCookie.Value
	}
	if nonceCookie, err := req.Cookie(o.cookieName("xfer")); err == nil {
		nonceStr = nonceCookie.Value
	}

	respValues := postRedirectData{
		IDToken: idToken,
		AccessToken: accessToken,
		NextURL: next,
		LocalStorage: xferMethod == "local",
		ChildWindow: xferMethod == "window",
	}
	if o.ClaimSource == IDENTITY_TOKEN {
		respValues.AccessToken = fmt.Sprintf("%s.%s.%s", urlState, req.RemoteAddr, nonceStr)
	} 

	if err := o.postRedirectTempl.Execute(res, respValues); err != nil {
		logger.Error().Err(err).Msg("Could not fill template")
		res.WriteHeader(http.StatusInternalServerError)
	}
}

// Update user claims
// 6a. Stoke sends a UserInfo request to the provider using the access token
// 6b. Stoke gets state from idToken and verifies the access token is the state that we gave for that idToken
// 7. Stoke uses claims to update database
func (o *oidcUserProvider) UpdateUserClaims(idToken, accessToken string, ctx context.Context) error {
	ctx, span := tel.GetTracer().Start(ctx, "oidcUserProvider.UpdateUserClaims")
	defer span.End()

	var claimMap jwt.MapClaims
	var err error
	var ok bool

	if o.ClaimSource == IDENTITY_TOKEN {
		jParser := jwt.NewParser()
		t, _, err := jParser.ParseUnverified(idToken, jwt.MapClaims{})
		if err != nil {
			return err
		}
		claimMap, ok = t.Claims.(jwt.MapClaims)
		if !ok {
			return jwt.ErrTokenMalformed
		}

		// state.addr.nonce
		accessParts := strings.Split(accessToken, ".")
		if len(accessParts) != 3 {
			return ProviderAuthError
		}
		state := accessParts[0]
		addr := accessParts[1]
		nonceStr := accessParts[2]

		if nonceClaim, ok := claimMap["nonce"]; ok && nonceClaim != nonceStr {
			return ProviderAuthError
		}

		// The iat should be within the same interval as the generated state (10min)
		tint := timePeriod(10)
		if iatClaim, err := claimMap.GetIssuedAt(); err == nil{
			if time.Now().After(iatClaim.Add(time.Second)){
				return ProviderAuthError
			}
			tint = timeToBytes(iatClaim.Truncate(10 * time.Minute))
		}

		nonce, err := base64.URLEncoding.DecodeString(nonceStr)
		if err != nil {
			return err
		}

		if state != stateHash(o.StateSalt, nonce, []byte(addr), tint) {
			return ProviderAuthError
		}
	} else {
		// o.ClaimSource == USER_INFO, trust the provider to authenticate the user.
		claimMap, err = o.getUserInfo(accessToken, ctx)
		if err != nil {
			return err
		}
	}

	return o.persistClaims(claimMap, ctx)
}

func (o *oidcUserProvider) addParamsToAuthURL(req *http.Request, nonce []byte) *url.URL {
	u, _ := url.Parse(o.AuthenticationURL.String())
	q := u.Query()

	q.Add("scope", o.Request.Scope)
	q.Add("response_type", o.FlowType.String())
	q.Add("client_id",o.Request.ClientID)
	q.Add("redirect_uri",o.Request.RedirectURI)
	for key, val := range o.Request.ExtraArgs {
		q.Add(key, val)
	}

	q.Add("state", o.computeState(req, nonce))

	u.RawQuery = q.Encode()
	return u
}

// Translates/persists claims to the database from the given claimsString
func (o *oidcUserProvider) persistClaims(claimMap jwt.MapClaims, ctx context.Context) error {
	logger := zerolog.Ctx(ctx).With().
		Str("component", "OIDCProvider").
		Interface("provider_claims", claimMap).
		Logger()
	logger.Info().Msg("TODO: Pushing claims to database")
	return nil
}

// Gets user claims from the UserInfoURL
func (o *oidcUserProvider) getUserInfo(accessToken string, ctx context.Context) (jwt.MapClaims, error) {
	logger := zerolog.Ctx(ctx).With().
		Str("component", "OIDCProvider").
		Stringer("claim_source", o.ClaimSource).
		Stringer("flow_type", o.FlowType).
		Stringer("user_info_url", o.UserInfoURL).
		Logger()

	logger.Info().Msg("Requesting user info")
	req, err := http.NewRequest(
		http.MethodGet, 
		o.UserInfoURL.String(),
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer " + accessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Could not get user info")
		return nil, err
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Could not read user info")
		return nil, err
	}

	claims := make(jwt.MapClaims)
	err = json.Unmarshal(bodyBytes, &claims)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Could not unmarshal user info")
		return nil, err
	}

	logger.Debug().
		Interface("claims", claims).
		Msg("Retrieved user claims from user info endpoint")

	return claims, nil
}

// Gets id token and access token from tokenURL.
// If we already have the tokens we need, it will do nothing
func (o *oidcUserProvider) getTokens(idToken, accessToken, authCode string, ctx context.Context) (string, string, error) {
	logger := zerolog.Ctx(ctx).With().
		Str("component", "OIDCProvider").
		Stringer("claim_source", o.ClaimSource).
		Stringer("flow_type", o.FlowType).
		Stringer("token_url", o.TokenURL).
		Str("id_token", idToken).
		Logger()

	if o.ClaimSource == IDENTITY_TOKEN && idToken != "" {
		return idToken, accessToken, nil
	}
	if o.ClaimSource == USER_INFO && accessToken != "" {
		return idToken, accessToken, nil
	}

	logger.Info().Msg("Requesting tokens")

	values := make(url.Values)
	values.Add("grant_type", "authorization_code")
	values.Add("code", authCode)
	values.Add("redirect_uri", o.Request.RedirectURI)

	req, _ := http.NewRequest(http.MethodPost, o.TokenURL.String(), strings.NewReader(values.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(o.Request.ClientID, o.ClientSecret)

	logger = logger.With().Str("post_values", values.Encode()).Logger()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Debug().Msg("Could not post form to token url")
		return idToken, accessToken, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return idToken, accessToken, err
	}

	tokens := make(map[string]interface{})
	err = json.Unmarshal(bodyBytes, &tokens)
	if err != nil {
		logger.Error().
			Bytes("response_bytes", bodyBytes).
			Msg("Could not unmarshal json from token endpoint")
		return idToken, accessToken, err
	}

	if acc, ok := tokens["access_token"]; accessToken == "" && ok {
		logger.Debug().Msg("Received access token")
		accessToken, _ = acc.(string)
	}
	if id, ok := tokens["id_token"]; idToken == "" && ok {
		logger.Debug().Msg("Received identity token")
		idToken, _ = id.(string)
	}
	if respErr, ok := tokens["error"]; ok {
		errStr, _ := respErr.(string)
		logger.Error().
			Bytes("response_bytes", bodyBytes).
			Str("response_error", errStr).
			Msg("Received an error from the token endpoint")
		return "", "", TokenRetrievalError
	}

	if ( idToken == "" && o.ClaimSource == IDENTITY_TOKEN ) || ( accessToken == "" && o.ClaimSource == USER_INFO ) {
		return "", "", TokenRetrievalError
	}

	logger.Debug().
		Str("access_token", accessToken).
		Str("id_token", idToken).
		Msg("Retrieved tokens from token url")

	return idToken, accessToken, nil
}

// Computes a state hash that should be specific to a host and a request
func (o *oidcUserProvider) computeState(req *http.Request, nonce []byte) string {
	return requestHash(o.StateSalt, nonce, req)
}

// Validates whether the state variable was generated by us
func (o *oidcUserProvider) hasValidState(req *http.Request) bool {
	nonceB64, err := req.Cookie(o.cookieName("nonce"))
	if err != nil {
		return false
	}
	nonce, err := base64.URLEncoding.DecodeString(nonceB64.Value)
	return req.URL.Query().Get("state") == requestHash(o.StateSalt, nonce, req)
}

func (o *oidcUserProvider) cookieName(c string) string {
	return fmt.Sprintf("stoke-oidc-%s-%s", o.Name, c)
}

// Creates a new random 32 byte string
func newNonce() []byte {
	nonce := make([]byte, 32)
	rand.Read(nonce)
	return nonce
}

// Returns the time stamp for the given division of time in minutes
func timePeriod(m int) []byte {
	return timeToBytes(time.Now().Truncate(time.Duration(m) * time.Minute))
}

func timeToBytes(t time.Time) []byte {
	tsBytes := make([]byte, 8)
	ts := t.Unix()
	binary.NativeEndian.PutUint64(tsBytes, uint64(ts))
	return tsBytes
}

// Creates a request specific hash. Only valid for 10 minutes.
func requestHash(salt, nonce []byte, req *http.Request) string {
	return stateHash(salt, nonce, []byte(req.RemoteAddr), timePeriod(10))
}

func stateHash(salt, nonce, addr, tint []byte) string {
	hash := sha256.Sum256(slices.Concat( addr, salt, tint, nonce ))
	return base64.URLEncoding.EncodeToString(hash[:])
}
