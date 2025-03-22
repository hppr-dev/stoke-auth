package usr

import (
	"net/url"
	"strings"
)

type ClaimSourceType uint8
const (
	IDENTITY_TOKEN ClaimSourceType = 1 << iota
	USER_INFO
)

func (c ClaimSourceType) String() string {
	switch c {
	case IDENTITY_TOKEN:
		return "identity token"
	case USER_INFO:
		return "user info endpoint"
	}
	return "undefined"
}

type AuthFlowType uint8
const (
	_ID_TOKEN AuthFlowType = 1 << iota
	_ACCESS_TOKEN
	_AUTH_CODE

	// An auth code is received upon redirect from the provider
	CODE_FLOW = _AUTH_CODE

	// An id token is returned from the provider
	IMPLICIT_FLOW      = _ID_TOKEN
	// An id token and an access token are returned from the provider
	IMPLICIT_USER_INFO = _ID_TOKEN | _ACCESS_TOKEN

	// An auth code and id token are returned from the provider
	HYBRID_FLOW      = _AUTH_CODE | _ID_TOKEN
	// An auth code and an access token are returned from the provider
	HYBRID_USER_INFO = _AUTH_CODE | _ACCESS_TOKEN
	// Unsupported. An auth code, id token and  access token are returned from the provider
	HYBRID_FLOW_ALL  = _AUTH_CODE | _ID_TOKEN | _ACCESS_TOKEN // UNSUPPORTED
)

func (a AuthFlowType) CompatibleWithSourceType(c ClaimSourceType) bool {
	// AUTH_CODE is compatible with either
	// An ACCESS_TOKEN is required to access the USER_INFO endpoint
	// ID_TOKEN requires that we receive an identity token
	return a & _AUTH_CODE == _AUTH_CODE || uint8(a) & uint8(c) > 0
}

func (a AuthFlowType) HasAuthCode() bool { return a & _AUTH_CODE == _AUTH_CODE }
func (a AuthFlowType) HasIDToken() bool { return a & _ID_TOKEN == _ID_TOKEN }
func (a AuthFlowType) HasAccessToken() bool { return a & _ACCESS_TOKEN == _ACCESS_TOKEN }

func (a AuthFlowType) GetResponseMap(q url.Values) map[string]string {
	m := make(map[string]string)
	if a.HasAuthCode() {
		m["code"] = q.Get("code")
	}
	if a.HasIDToken() {
		m["id_token"] = q.Get("id_token")
	}
	if a.HasAccessToken() {
		m["token"] = q.Get("token")
	}
	m["state"] = q.Get("state")
	m["nonce"] = q.Get("nonce")
	return m
}

func (a AuthFlowType) String() string {
	b := strings.Builder{}
	if a.HasAuthCode() {
		b.WriteString("code ")
	}
	if a.HasIDToken() {
		b.WriteString("id_token ")
	}
	if a.HasAccessToken() {
		b.WriteString("token ")
	}
	if b.Len() == 0 {
		return ""
	}
	return b.String()[0:b.Len() - 1]
}
