package stoke

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
)

// Creates a new set of stoke credentials
// Optional arguments are interpreted as token, refreshToken, refreshURL (in that order)
// Users must call either include a Token as an argument, set the token using Token or Login to populate credentials before using
//
// Examples:
//    // Already have a token to use for multiple grpc server calls
//    client := grpc.NewClient(Credentials(myToken).DialOption())
//    client := grpc.NewClient(Credentials().Token(myToken).DialOption()) // Equivalent
//
//    // Already have a token to use for a single grpc server call
//    client := grpc.NewClient(...)
//    customConn := proto.NewCustomGrpcClient(client)
//    res := customConn.DoTheThing(Credentials(myToken).CallOption())
//
//    // Already have a token and refresh token to use (Will only refresh if it is used 5 seconds before expire time)
//    client := grpc.NewClient(
//      Credentials().                                           // Credentials(myToken, myRefresh, "http://mystoke/api/refresh") may also be used
//      	Token(myToken).
//      	EnableRefresh(myRefresh, "http://mystoke/api/refresh").
//				RefreshWindow(5 * time.Second).
//      	DialOption(),
//    )
//
//    // Login as a user and let the token expire if requests are not made within 10 seconds of expiring time
//    client := grpc.NewClient(
//      Credentials().
//      	Login("myuser", "mypass", "http://mystoke", false).
//				RefreshWindow(10 * time.Second).
//      	DialOption(),
//    )
//
//    // Login as a user and keep the token up to date with a goroutine (no success check)
//    ctx, cancel := context.WithCancel(context.Background())
//    defer cancel() // this will cancel the refresh goroutine
//    client := grpc.NewClient(
//			Credentials().
//				Login("myuser", "mypass", "http://mystoke", true)
//				StartRefresh(ctx).
//				DialOption(),
//		)
//
//    // Login as a user and keep the token up to date with a goroutine (check for login success)
//    ctx, cancel := context.WithCancel(context.Background())
//    defer cancel() // this will cancel the refresh goroutine
//    creds := Credentials().Login("myuser", "mypass", "http://mystoke", true)
//    if creds.Initialized() { // Check if login succeeded
//    	creds.StartRefresh(ctx)
//    } else {
// 			panic("Could not obtain token!")
//    }
//
//    client := grpc.NewClient(creds.DialOption())
//
func Credentials(optionalTokens ...string) *stokeCredentials {
	creds := &stokeCredentials{
		refreshWindow : 2 * time.Second,
	}
	if len(optionalTokens) > 0 {
		creds.token = optionalTokens[0]
		creds.expireTime = getExpiryFromToken(creds.token)
	}
	if len(optionalTokens) > 1 {
		creds.refresh = optionalTokens[1]
	}
	if len(optionalTokens) > 2 {
		creds.refreshURL = optionalTokens[2]
	}
	return creds
}

// Specify the initial token to use
// This token will be refreshed if EnableRefresh is called.
func (c *stokeCredentials) Token(token string) *stokeCredentials {
	c.token = token
	c.expireTime = getExpiryFromToken(token)
	return c
}

// Enables automatic refresh of token
// If ExpireTime is not set before calling, the token is set to expire 30 seconds after calling
//
// Call StartRefresh() to keep token up to date with a goroutine.
// Otherwise the token will only be refreshed when a request is made within the refresh window of the expiration time.
func (c *stokeCredentials) EnableRefresh(refresh, refreshURL string) *stokeCredentials {
	c.refresh = refresh
	c.refreshURL = refreshURL
	if c.expireTime.Before(time.Now()) {
		c.expireTime = time.Now().Add(30 * time.Second)
	}
	return c
}

// Returns whether these credentials have been initialized
// Useful for checking if login succeeded
func (c *stokeCredentials) Initialized() bool {
	return c.token != ""
}

// Starts a goroutine to keep credential tokens valid by automatically refreshing.
// Cancel the context to stop the goroutine.
// Panics if EnableRefresh has not been called
func (c *stokeCredentials) StartRefresh(ctx context.Context) *stokeCredentials {
	if c.refresh == "" || c.refreshURL == "" {
		panic("Must call EnableRefresh to set refresh token and refreshURL.")
	}
	c.managed = true

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case <-time.After(c.expireTime.Sub(time.Now())):
				err := c.refreshToken()
				if err != nil {
					log.Printf("An error occurred while refreshing token: %s", err)
				}
			}
		}
	}()
	return c
}

// Logs user into a given stoke server at stokeURL using username and password.
// stokeURL should point to the root of a stoke server ({stokeURL}/api/login and {stokeURL}/api/refresh should be reachable.
// Useful for long lived or automated access by a service account when enableRefresh is set to true
// Sets all values to empty if unable to login
// Call Initialized to check if a the login succeeded
func (c *stokeCredentials) Login(username, password, stokeURL string, enableRefresh bool) *stokeCredentials {
	res, err := http.Post(
		fmt.Sprintf("%s/api/login", stokeURL),
		"application/json",
		bytes.NewBuffer(
			[]byte(fmt.Sprintf(`{"username" : "%s", "password": "%s"}`, username, password)),
		),
	)
	if err != nil {
		c.token = ""
		c.refresh = ""
		c.expireTime = time.Now()
		return c
	}

	parsedResp, err := getTokenResponse(res)
	c.updateToken(parsedResp)

	if enableRefresh {
		c.EnableRefresh(parsedResp.Refresh, fmt.Sprintf("%s/api/refresh", stokeURL))
	}

	return c
}

// Disable requiring transport security.
// USE IN DEVELOPMENT ONLY
func (c *stokeCredentials) DisableSecurity() *stokeCredentials {
	c.disableSecurity = true
	return c
}

// Set the refresh window for credentials that have not been 'Start'ed
// This must be less than the token duration and defaults to 2 seconds
// For example, if a token is set to expire at 12:34:56 with a refreshWindow of 5 seconds any request after 12:34:51 will cause a refresh.
// Does not apply to tokens that have not had refresh enabled
func (c *stokeCredentials) RefreshWindow(dur time.Duration) *stokeCredentials {
	c.refreshWindow = dur
	return c
}

// Convert credentials into a grpc PerRPCCredentials Call Option
func (c *stokeCredentials) CallOption() grpc.CallOption {
	return grpc.PerRPCCredentials(c)
}

// Convert credentials into a grpc PerRPCCredentials Dial Option
func (c *stokeCredentials) DialOption() grpc.DialOption {
	return grpc.WithPerRPCCredentials(c)
}

// Private credentials struct. Use Credentials() to create.
type stokeCredentials struct {
	token string
	refresh string
	refreshURL string
	expireTime time.Time
	disableSecurity bool
	managed bool
	refreshWindow time.Duration

	mutex sync.RWMutex
}

// Struct to marshall json received from stoke server
type tokenResponse struct {
	Token string   `json:"token"`
	Refresh string `json:"refresh"`
}

// GetRequestMetadata implements credentials.PerRPCCredentials.
// Refreshes the token if the current time is within the RefreshWindow before expiration.
func (c *stokeCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	if !c.managed && c.refresh != "" && c.refreshURL != "" && c.expireTime.Before(time.Now().Add(-1 * c.refreshWindow)){
		if err := c.refreshToken(); err != nil {
			return nil, err
		}
	}
	return map[string]string{
		"authorization": "Bearer " + c.token,
	}, nil
}

// RequireTransportSecurity implements credentials.PerRPCCredentials.
// Call DisableSecurity to set this to false
func (c *stokeCredentials) RequireTransportSecurity() bool {
	return !c.disableSecurity
}

// Concurrency safe update of token values
func (c *stokeCredentials) updateToken(res *tokenResponse) {
	expTime := getExpiryFromToken(res.Token)

	c.mutex.Lock()
	c.token = res.Token
	c.refresh = res.Refresh
	c.expireTime = expTime
	c.mutex.Unlock()
}

// Refreshes token using the refresh token
func (c *stokeCredentials) refreshToken() error {
	res, err := http.Post(
		c.refreshURL,
		"application/json",
		bytes.NewBuffer(
			[]byte(fmt.Sprintf(`{"refresh" : "%s"}`, c.refresh)),
		),
	)
	if err != nil {
		return err
	}

	parsedResp, err := getTokenResponse(res)

	c.updateToken(parsedResp)
	return nil
}

// Extract the token response from a http response
func getTokenResponse(res *http.Response) (*tokenResponse, error) {
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Unexpected status code: %d", res.StatusCode)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	parsedResp := &tokenResponse{}
	if err := json.Unmarshal(bodyBytes, parsedResp); err != nil {
		return nil, err
	}

	return parsedResp, nil
}


func getExpiryFromToken(token string) time.Time {
	expTime := time.Now().Add(5 * time.Minute)
	usedDefaultTime := true

	jwtBody := strings.Split(token, ".")[1]
	padding := strings.Repeat("=", 4 - len(jwtBody) % 4)

	if claimsBytes, err := base64.URLEncoding.DecodeString(jwtBody + padding); err == nil {
		var claimsMap map[string]interface{}
		if err := json.Unmarshal(claimsBytes, &claimsMap); err == nil {
			if expClaim, ok := claimsMap["exp"]; ok {
				if expFloat, ok := expClaim.(float64); ok {
					expTime = time.Unix(int64(expFloat), 0).Add(-750 * time.Millisecond)
					usedDefaultTime = false
				}
			}
		}
	}

	if usedDefaultTime {
		log.Println("Could not parse expire time out of token. Using 5 minutes as default")
	}

	return expTime
}
