package usr

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"slices"
	"stoke/internal/ent"
	"stoke/internal/ent/grouplink"
	"stoke/internal/ent/predicate"
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


type postRedirectData struct {
	IDToken string
	AccessToken string
	NextURL string
	LocalStorage bool
	ChildWindow bool
	LoginURL string
}

var (
	// If xfer is set to window, next must be set to the original owners origin
	// If xfer is set to local, next must match the same origin and share sessionStorage
	// Not sure what happens when next is not set
	POST_TEMPLATE = `
<html>
	<head>
		<script lang="javascript">
			window.onload = function() {
				{{ if .LocalStorage }}
					window.sessionStorage.setItem("id_token", "{{ .IDToken }}")
					window.sessionStorage.setItem("access_code", "{{ .AccessCode }}")
					{{ if ne .NextURL "" }}
						window.location = "{{ .NextURL }}"
					{{ end }}
				{{ else if .ChildWindow }}
					var message = JSON.stringify({
						id_token : "{{ .IDToken }}",
						access_code : "{{ .AccessToken }}",
					})
					window.opener.postMessage(message, "{{ .NextURL }}")
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

	// First Name Claim
	FNameClaim string
	// Last Name Claim
	LNameClaim string
	// Email Claim
	EmailClaim string

	postRedirectTempl *template.Template
	dbSourceName string
}

func NewOIDCUserProvider(
	name, scopes, redirectURI,
	fNameClaim, lNameClaim, emailClaim,
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
		FNameClaim: fNameClaim,
		LNameClaim: lNameClaim,
		EmailClaim: emailClaim,
		postRedirectTempl: prt,
		dbSourceName: "OIDC:" + name,
	}

	mux.Handle("/oidc/" + name, provider)

	return provider
}

// Handles redirect to and from provider
// Users should navigate to this endpoint to authenticate with the provider
// Include the following query parameters to control redirect behavior:
// 		* next -- the url the user wants to goto after authenticating
//    * xfer -- the transfer method of the request. May be:
//			* local
//			* window
//
// Clients should register to be returned back to this endpoint after the auth with the provider,
// i.e. the redirect uri should be registered at /oidc/<PROVIDER_NAME>.
// Users MUST NOT include a state query parameter when requesting because that indicates a return request from the provider
//
// This function serves as steps 2-5 in the full authentication flow above
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

	ctx, span := tel.GetTracer().Start(ctx, "oidcUserProvider.ServeHTTP")
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
	if nonceCookie, err := req.Cookie(o.cookieName("nonce")); err == nil {
		nonceStr = nonceCookie.Value
	}


	respValues := postRedirectData{
		LoginURL: "/api/login",
		IDToken: idToken,
		AccessToken: accessToken,
		NextURL: next,
		LocalStorage: xferMethod == "local",
		ChildWindow: xferMethod == "window",
	}
	if o.ClaimSource == IDENTITY_TOKEN {
		respValues.AccessToken = fmt.Sprintf("%s$%s$%s", urlState, req.RemoteAddr, nonceStr)
	} 

	if err := o.postRedirectTempl.Execute(res, respValues); err != nil {
		logger.Error().Err(err).Msg("Could not fill template")
		res.WriteHeader(http.StatusInternalServerError)
	}
}

// Update user claims in the database with the claims from the provider.
// This function should be called with the idToken and accessToken returned from the provider after the user authorizes access.
// When the ClaimsSource is ID_TOKEN, the accessToken is used to verify the idToken using the state, network address and nonce.
//
// This function finishes the process (steps 6 and 7 above) and will result in the user getting a token with the most up-to-date claims available from the provider
func (o *oidcUserProvider) UpdateUserClaims(idToken, accessToken string, ctx context.Context) (*ent.User, error) {
	logger := zerolog.Ctx(ctx).With().
		Str("component", "OIDCProvider.UpdateUserClaims").
		Stringer("flow_type", o.FlowType).
		Stringer("claim_source", o.ClaimSource).
		Logger()

	ctx, span := tel.GetTracer().Start(ctx, "oidcUserProvider.UpdateUserClaims")
	defer span.End()

	var claimMap jwt.MapClaims
	var err error
	var ok bool

	jParser := jwt.NewParser()
	t, _, err := jParser.ParseUnverified(idToken, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}
	claimMap, ok = t.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrTokenMalformed
	}

	logger.Debug().
		Interface("claim_map", claimMap).
		Msg("Parsed user claim map.")

	if o.ClaimSource == USER_INFO {
		// trust the provider to authenticate the user.
		infoClaimMap, err := o.getUserInfo(accessToken, ctx)
		if err != nil {
			logger.Error().Err(err).Msg("Could not get user info from endpoint")
			return nil, err
		}
		for k,v := range infoClaimMap {
			claimMap[k] = v
		}
	} else {
		//o.ClaimSource == IDENTITY_TOKEN, verify accessToken as state$addr$nonce
		accessParts := strings.Split(accessToken, "$")
		if len(accessParts) != 3 {
			logger.Debug().Str("token", accessToken).Msg("Received bad access token")
			return nil, AuthenticationError
		}
		state := accessParts[0]
		addr := accessParts[1]
		nonceStr := accessParts[2]

		if nonceClaim, ok := claimMap["nonce"]; ok && nonceClaim != nonceStr {
			logger.Debug().Interface("nonce", nonceClaim).Msg("Received bad nonce")
			return nil, AuthenticationError
		}

		// The iat should be within the same interval as the generated state (10min)
		tint := timePeriod(10)
		if iatClaim, err := claimMap.GetIssuedAt(); err == nil{
			if time.Now().After(iatClaim.Add(10 * time.Minute)){
				logger.Debug().Time("issued_at", iatClaim.Time).Msg("Got stale id token")
				return nil, AuthenticationError
			}
			tint = timeToBytes(iatClaim.Truncate(10 * time.Minute))
		}

		nonce, err := base64.URLEncoding.DecodeString(nonceStr)
		if err != nil {
			logger.Error().Err(err).Str("nonce", nonceStr).Msg("could not decode nonce")
			return nil, err
		}
		expectedState := stateHash(o.StateSalt, nonce, []byte(addr), tint) 
		if state != expectedState {
			logger.Debug().
				Str("state", state).
				Str("expected_state", expectedState).
				Msg("State did not match expected")
			return nil, AuthenticationError
		}
	}

	logger.Debug().
		Interface("claim_map", claimMap).
		Msg("Received/validated all claims")

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
func (o *oidcUserProvider) persistClaims(claimMap jwt.MapClaims, ctx context.Context) (*ent.User, error) {
	logger := zerolog.Ctx(ctx).With().
		Str("component", "OIDCProvider.persistClaims").
		Interface("provider_claims", claimMap).
		Logger()

	u, err := o.getOrCreateUser(claimMap, ctx)
	if err != nil {
		return nil, err
	}
	logger.Debug().Interface("user", u).Msg("Saving user claims")

	db := ent.FromContext(ctx)

	claimLinks := []predicate.GroupLink{}
	
	for cKey, cValue := range claimMap {
		claimLinks = append(claimLinks, grouplink.ResourceSpecEQ(fmt.Sprintf("%s=%s", cKey, cValue)))
	}

	foundLinks, err := db.GroupLink.Query().
		Where(
			grouplink.And(
				grouplink.TypeEQ(o.dbSourceName),
				grouplink.Or(claimLinks...),
			),
		).
		WithClaimGroup(func (q *ent.ClaimGroupQuery) {
			q.WithClaims()
		}).
		All(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Could not get group links.")
		return nil, err
	} else if len(foundLinks) == 0 {
		logger.Error().Msg("No group links found")
		return nil, NoLinkedGroupsError
	}

	logger.Debug().Interface("found_links", foundLinks).Msg("Found group links")

	//TODO allow user claim passthrough with or without persistance?
	add, del := findGroupChanges(u, foundLinks)
	if u, err = applyGroupChanges(add, del, u, ctx) ; err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to update oidc groups to local user")
		return nil, err
	}

	logger.Debug().Interface("user", u).Msg("Updated user groups")

	return retreiveLocalUser(u.Username, ctx)
}

func (o *oidcUserProvider) getOrCreateUser(claimMap jwt.MapClaims, ctx context.Context) (*ent.User, error) {
	fname := safeGetClaim(o.FNameClaim, claimMap)
	lname := safeGetClaim(o.LNameClaim, claimMap)
	email := safeGetClaim(o.EmailClaim, claimMap)

	logger := zerolog.Ctx(ctx).With().
		Str("component", "OIDCProvider.getOrCreateUser").
		Str("fname", fname).
		Str("lname", lname).
		Str("email", email).
		Logger()

	if fname == "" || lname == "" || email == "" {
		logger.Error().Msg("Could not determine first name, last name or email")
		return nil, jwt.ErrTokenMalformed
	}

	u, err := retreiveLocalUser(email, ctx)
	if ent.IsNotFound(err) {
		logger.Info().Msg("User not found, creating in database.")
		return ent.FromContext(ctx).User.Create().
			SetFname(fname).
			SetLname(lname).
			SetEmail(email).
			SetUsername(email).
			SetSource(o.dbSourceName).
			Save(ctx)
	}
	logger.Debug().Interface("user", u).Msg("User found.")
	return u, nil

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
		return "", "", OIDCTokenRetrievalError
	}

	if ( idToken == "" && o.ClaimSource == IDENTITY_TOKEN ) || ( accessToken == "" && o.ClaimSource == USER_INFO ) {
		return "", "", OIDCTokenRetrievalError
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

func safeGetClaim(key string, claims jwt.MapClaims) string {
	if valInt, ok := claims[key]; ok {
		if val, ok := valInt.(string); ok {
			return val
		}
	}
	return ""
}
