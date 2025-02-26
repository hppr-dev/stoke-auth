package cfg

import (
	"context"
	"stoke/internal/usr"
	"strings"
)

type OIDCProviderConfig struct {
	// URL to use to get the OIDC token from the provider
	TokenURL          string
	// URL to use to authenticate users with the provider
	AuthenticationURL string
	// The authentication request to use when authenticating to the AuthenticationURL
	Request AuthenticationRequestConfig
}

type AuthenticationRequestConfig struct {
	// Scopes to include in Authentication requests
	Scopes        []string
	// Client ID that matches what is registered with the OpenID Provider
	ClientID     string
	// Redirect URI that matches what is registered at the OpenID Provider
	RedirectURI  string
	// Extra authentication requirements that are added to the request
	ExtraArguments map[string]string
}

func (o OIDCProviderConfig) CreateProvider(ctx context.Context) foreignProvider {
	return usr.NewOIDCUserProvider(
		o.TokenURL,
		o.AuthenticationURL,
		usr.OIDCAuthRequest{
			Scope: strings.Join(append(o.Request.Scopes, "openid"), " "),
			ClientID: o.Request.ClientID,
			ResponseType: "code",
			RedirectURI: o.Request.RedirectURI,
			ExtraArgs: o.Request.ExtraArguments,
		},
	)
}
