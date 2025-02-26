package usr

import "context"

// https://openid.net/specs/openid-connect-core-1_0.html
// Need to verify the user is who they say they are by sending username/password to oidc
// The received JWT has claims that are used to build a stoke token
// Should be able to:
// 		* Pass through claims from provider to issued token
//    * Map claims from provider to issued token

type OIDCAuthRequest struct {
	Scope        string
	ResponseType string
	ClientID     string
	RedirectURI  string
	State        string
	Nonce        string
	ExtraArgs    map[string]string
}

type oidcUserProvider struct {
	// URL to use to get the OIDC token from the provider
	TokenURL          string
	// URL to use to authenticate users with the provider
	AuthenticationURL string
	// The authentication request to use when authenticating to the AuthenticationURL
	Request           OIDCAuthRequest
}

func NewOIDCUserProvider(tokenURL, authURL string, req OIDCAuthRequest) *oidcUserProvider {
	return &oidcUserProvider{
		TokenURL : tokenURL,
		AuthenticationURL: authURL,
		Request: req,
	}
}

func (o *oidcUserProvider) UpdateUserClaims(username, password string, ctx context.Context) error {
	return nil
}
