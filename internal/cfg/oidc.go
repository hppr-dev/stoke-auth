package cfg

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"stoke/internal/usr"
	"strings"

	"github.com/rs/zerolog"
)

type OIDCProviderConfig struct {
	// Name of this OIDC provider. Used in the login URL
	Name              string `json:"name"`
	// Discovery url that conforms to https://openid.net/specs/openid-connect-discovery-1_0.html
	// Automatically sets the TokenURL, AuthenticationURL, UserInfoURL
	DiscoveryURL      string `json:"discovery_url"`

	// URL to use to get the OIDC token from the provider
	TokenURL          string `json:"token_url"`
	// URL to use to authenticate users with the provider
	AuthenticationURL string `json:"authentication_url"`
	// URL of the UserInfo Endpoint
	UserInfoURL       string `json:"user_info_url"`
	// State secret, must be the same accross deployments. Random 16 bytes if not specified
	StateSecret       string `json:"state_secret"`
	// Authentication flow type to use. May be code, implicit or hybrid
	AuthFlowType      string `json:"auth_flow_type"`
	// Where to pull claims from. May be id_token or user_info
	ClaimsSource      string `json:"claims_source"`
	// The authentication request to use when authenticating to the AuthenticationURL
	Request AuthenticationRequestConfig `json:"request_config"`
}

type AuthenticationRequestConfig struct {
	// Scopes to include in Authentication requests
	Scopes       []string `json:"scopes"`
	// Client ID that matches what is registered with the OpenID Provider
	ClientID     string   `json:"client_id"`
	// Redirect URI that matches what is registered at the OpenID Provider
	RedirectURI  string   `json:"redirect_uri"`
	// Extra authentication requirements that are added to the request
	ExtraArguments map[string]string `json:"extra_arguments"`
}

func (o OIDCProviderConfig) CreateProvider(ctx context.Context) foreignProvider {
	logger := zerolog.Ctx(ctx).With().
		Str("provider_name", o.Name).
		Str("token_url", o.TokenURL).
		Str("authentication_url", o.AuthenticationURL).
		Str("user_info_url", o.UserInfoURL).
		Str("auth_flow_type", o.AuthFlowType).
		Str("claims_source", o.ClaimsSource).
		Logger()

	if o.DiscoveryURL != "" {
		o.retrieveConfigFromDiscovery(ctx)
	}
	authURL, err := url.Parse(o.AuthenticationURL)
	if err != nil {
		logger.Fatal().Err(err).Msg("Authentication URL is required")
	}

	tokenURL, _ := url.Parse(o.TokenURL)
	userInfoURL, _ := url.Parse(o.UserInfoURL)

	var sourceType usr.ClaimSourceType
	switch strings.ToLower(o.ClaimsSource) {
	case "id_token", "token", "id token":
		sourceType = usr.IDENTITY_TOKEN
	case "user_info", "user info", "endpoint":
		sourceType = usr.USER_INFO
	default:
		logger.Fatal().Msg("Unsupported claims_source. Must be id token or endpoint.")
	}

	var flowType usr.AuthFlowType
	switch strings.ToLower(o.AuthFlowType) {
		case "code", "auth_code", "auth code": 
			flowType = usr.CODE_FLOW
		case "implicit":
			if sourceType == usr.IDENTITY_TOKEN {
				flowType = usr.IMPLICIT_FLOW
			} else {
				flowType = usr.IMPLICIT_USER_INFO
			}
		case "hybrid":
			if sourceType == usr.IDENTITY_TOKEN {
				flowType = usr.HYBRID_FLOW
			} else {
				flowType = usr.HYBRID_USER_INFO
			}
		default:
			logger.Fatal().Msg("Unsupported auth_flow_type. Must be code, implicit or hybrid")
	}

	return usr.NewOIDCUserProvider(
		o.Name,
		o.StateSecret,
		tokenURL,
		authURL,
		userInfoURL,
		usr.OIDCAuthRequest{
			Scope: strings.Join(append(o.Request.Scopes, "openid"), " "),
			ClientID: o.Request.ClientID,
			RedirectURI: o.Request.RedirectURI,
			ExtraArgs: o.Request.ExtraArguments,
		},
		MuxFromContext(ctx),
		flowType,
		sourceType,
	)
}

func (o *OIDCProviderConfig) retrieveConfigFromDiscovery(ctx context.Context) {
	logger := zerolog.Ctx(ctx).With().
		Str("discovery_url", o.DiscoveryURL).
		Logger()

	logger.Info().Msg("Retrieving config from discovery url")
	resp, err := http.Get(o.DiscoveryURL)
	if err != nil {
		logger.Fatal().Err(err).Msg("Could not retrieve info from discovery url")
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Fatal().Err(err).Msg("Could not retrieve info from discovery url")
	}
	discInfo := make(map[string]interface{})
	err = json.Unmarshal(bodyBytes, &discInfo)
	if err != nil {
		logger.Fatal().Err(err).Msg("Could not retrieve info from discovery url")
	}

	if u, ok := discInfo["authorization_endpoint"]; ok {
		o.AuthenticationURL = u.(string)
	}
	if u, ok := discInfo["token_endpoint"]; ok {
		o.TokenURL = u.(string)
	}
	if u, ok := discInfo["userinfo_endpoint"]; ok {
		o.UserInfoURL = u.(string)
	}

	logger.Info().
		Str("auth_url", o.AuthenticationURL).
		Str("token_url", o.TokenURL).
		Str("user_info_url", o.UserInfoURL).
		Msg("Done dicovering config.")
}
