package usr_test

import (
	"context"
	"crypto/tls"
	"errors"
	"regexp"
	"time"

	"github.com/go-ldap/ldap/v3"
)

var (
	LookupError = errors.New("Lookup Error")
	NotBoundError = errors.New("LDAP not bound before operation requiring bind")
)

type MockLDAPConnector struct{
	mockClient *MockLDAPClient
	connectError error
}

type MockLDAPClient struct{
	userPassMap      map[string]string
	searchDB         MockSearchDB
	bindError        error
	userSearchError  error
	groupSearchError error
	boundAs          string
}

type MockAttributes map[string][]string
type MockEntry map[string]MockAttributes
type MockSearchResult []MockEntry
type MockSearchDB map[string]MockSearchResult

type MockLDAPOption func(*MockLDAPConnector)

// Creates a new mock connector that "Connects" to a mock client
// The connector and client can be customized using MockLDAPOptions functions that modify connector internals
func NewMockLDAPServer(opts ...MockLDAPOption) *MockLDAPConnector {
	conn := &MockLDAPConnector{
		mockClient : &MockLDAPClient{
			userPassMap: make(map[string]string),
			searchDB:    make(MockSearchDB),
		},
	}

	for _, opt := range opts {
		opt(conn)
	}
	return conn
}

// "Connects" to a mock client
func (c MockLDAPConnector) Connect(string, ...ldap.DialOpt) (ldap.Client, error) {
	if c.connectError != nil {
		return nil, c.connectError
	}
	return c.mockClient, nil
}

// Return err on call to Connect
func ConnectError(err error) MockLDAPOption {
	return func (c *MockLDAPConnector) {
		c.connectError = err
	}
}

// Return err on call to Bind
func BindError(err error) MockLDAPOption {
	return func (c *MockLDAPConnector) {
		c.mockClient.bindError = err
	}
}

// Return err on call to Search for users
func UserSearchError(err error) MockLDAPOption {
	return func (c *MockLDAPConnector) {
		c.mockClient.userSearchError = err
	}
}

// Return err on call to Search for groups
func GroupSearchError(err error) MockLDAPOption {
	return func (c *MockLDAPConnector) {
		c.mockClient.groupSearchError = err
	}
}

// Create an User in the mock LDAP
func LDAPUser(username, first, last, email, password string) MockLDAPOption {
	return func(c *MockLDAPConnector) {
		c.mockClient.userPassMap[username] = password
		searchResult := MockSearchResult{
			MockEntry{
				username: MockAttributes{
					"username"  : []string{ username },
					"first_name": []string{ first },
					"last_name" : []string{ last },
					"email"     : []string{ email },
				},
			},
		}
		c.mockClient.searchDB["user:" + username] = searchResult
	}
}

// Create a group in the mock LDAP
func LDAPGroup(username, groupName string) MockLDAPOption {
	return func(c *MockLDAPConnector) {
		searchResult := MockSearchResult{
			MockEntry{
				username : MockAttributes{
					"group_name" : []string{ groupName },
				},
			},
		}
		c.mockClient.searchDB["group:" + username] = searchResult
	}
}

// Convert a mock entry to an ldap entry
func (e MockEntry) toEntry() *ldap.Entry {
	// Entries should only have one key
	for k, v := range e {
		return ldap.NewEntry(k, v)
	}
	return nil
}

// Convert a mock search result to an ldap search result
func (r MockSearchResult) toSearchResult() *ldap.SearchResult {
	res := &ldap.SearchResult{}
	for _, v := range r {
		res.Entries = append(res.Entries, v.toEntry())
	}
	return res
}

/*
================ ldap.Client methods ========================
*/

// Bind implements ldap.Client.
// Checks whether the password matches
func (m *MockLDAPClient) Bind(username string, password string) error {
	if m.bindError != nil{
		return m.bindError
	}

	if val, ok := m.userPassMap[username]; ok && val == password {
		m.boundAs = username
		return nil
	}

	return LookupError
}

var filterReg = regexp.MustCompile(`groupFilter:(?P<group>.*)|userFilter:(?P<user>.*)`)

// Search implements ldap.Client.
// Does not use ldap-like search filter/base dn
// The filter should be either groupFilter:Name or userFilter:Name
func (m *MockLDAPClient) Search(req *ldap.SearchRequest) (*ldap.SearchResult, error) {
	if m.boundAs == "" {
		return nil, NotBoundError
	}

	matches := filterReg.FindStringSubmatch(req.Filter)
	group := "group:" + matches[1]
	user := "user:" + matches[2]

	var result MockSearchResult
	if group != "group:"  {
		result, _ = m.searchDB[group]
		if m.groupSearchError != nil {
			return nil, m.groupSearchError
		}
	}
	if user != "user:" {
		result, _ = m.searchDB[user]
		if m.userSearchError != nil {
			return nil, m.userSearchError
		}
	}
	return result.toSearchResult(), nil
}

// Close implements ldap.Client.
func (m *MockLDAPClient) Close() error {
	m.boundAs = ""
	return nil
}

// Start implements ldap.Client.
func (m *MockLDAPClient) Start() {}
// StartTLS implements ldap.Client.
func (m *MockLDAPClient) StartTLS(*tls.Config) error { return nil }
// SetTimeout implements ldap.Client.
func (m *MockLDAPClient) SetTimeout(time.Duration) {}
// Unbind implements ldap.Client.
func (m *MockLDAPClient) Unbind() error { return nil}

/*
		Below methods are unimplemented because they are not needed for the current implementation
*/

// SearchAsync implements ldap.Client.
func (m *MockLDAPClient) SearchAsync(ctx context.Context, searchRequest *ldap.SearchRequest, bufferSize int) ldap.Response {
	panic("unimplemented")
}

// SearchWithPaging implements ldap.Client.
func (m *MockLDAPClient) SearchWithPaging(searchRequest *ldap.SearchRequest, pagingSize uint32) (*ldap.SearchResult, error) {
	panic("unimplemented")
}

// Compare implements ldap.Client.
func (m *MockLDAPClient) Compare(dn string, attribute string, value string) (bool, error) {
	panic("unimplemented")
}

// DirSync implements ldap.Client.
func (m *MockLDAPClient) DirSync(searchRequest *ldap.SearchRequest, flags int64, maxAttrCount int64, cookie []byte) (*ldap.SearchResult, error) {
	panic("unimplemented")
}

// DirSyncAsync implements ldap.Client.
func (m *MockLDAPClient) DirSyncAsync(ctx context.Context, searchRequest *ldap.SearchRequest, bufferSize int, flags int64, maxAttrCount int64, cookie []byte) ldap.Response {
	panic("unimplemented")
}

// ExternalBind implements ldap.Client.
func (m *MockLDAPClient) ExternalBind() error {
	panic("unimplemented")
}

// GetLastError implements ldap.Client.
func (m *MockLDAPClient) GetLastError() error {
	panic("unimplemented")
}

// IsClosing implements ldap.Client.
func (m *MockLDAPClient) IsClosing() bool {
	panic("unimplemented")
}

// NTLMUnauthenticatedBind implements ldap.Client.
func (m *MockLDAPClient) NTLMUnauthenticatedBind(domain string, username string) error {
	panic("unimplemented")
}

// SimpleBind implements ldap.Client.
func (m *MockLDAPClient) SimpleBind(*ldap.SimpleBindRequest) (*ldap.SimpleBindResult, error) {
	panic("unimplemented")
}

// Syncrepl implements ldap.Client.
func (m *MockLDAPClient) Syncrepl(ctx context.Context, searchRequest *ldap.SearchRequest, bufferSize int, mode ldap.ControlSyncRequestMode, cookie []byte, reloadHint bool) ldap.Response {
	panic("unimplemented")
}

// TLSConnectionState implements ldap.Client.
func (m *MockLDAPClient) TLSConnectionState() (tls.ConnectionState, bool) {
	panic("unimplemented")
}

// UnauthenticatedBind implements ldap.Client.
func (m *MockLDAPClient) UnauthenticatedBind(username string) error {
	panic("unimplemented")
}

/*
	Below methods are not implemented because we do not intend to modify the LDAP server
*/

// Del implements ldap.Client.
func (m *MockLDAPClient) Del(*ldap.DelRequest) error {
	panic("Deleting from LDAP is not allowed (Del)")
}

// PasswordModify implements ldap.Client.
func (m *MockLDAPClient) PasswordModify(*ldap.PasswordModifyRequest) (*ldap.PasswordModifyResult, error) {
	panic("LDAP Password Modification is not allowed (PasswordModify)")
}

// ModifyWithResult implements ldap.Client.
func (m *MockLDAPClient) ModifyWithResult(*ldap.ModifyRequest) (*ldap.ModifyResult, error) {
	panic("LDAP Modification is not allowed (ModifyWithResult)")
}

// Modify implements ldap.Client.
func (m *MockLDAPClient) Modify(*ldap.ModifyRequest) error {
	panic("LDAP Modification is not allowed (Modify)")
}

// ModifyDN implements ldap.Client.
func (m *MockLDAPClient) ModifyDN(*ldap.ModifyDNRequest) error {
	panic("LDAP Modification is not allowed (ModifyDN)")
}

// Add implements ldap.Client.
func (m *MockLDAPClient) Add(*ldap.AddRequest) error {
	panic("Adding to LDAP is not allowed (Add)")
}

