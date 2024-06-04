package stoke

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"google.golang.org/grpc"
)

// Creates a new set of stoke credentials
// Users must call either Token, or Login to populate credentials with data
//
// Examples:
//    // Already have a token to use
//    client := grpc.NewClient(Credentials().Token(myToken).DialOption())
//
//    // Already have a token and refresh token to use (Will only refresh if it is used after expiry)
//    client := grpc.NewClient(
//      Credentials().
//      Token(myToken).
//      ExpiresAt(expiry). // Must come before EnableRefresh to take effect
//      EnableRefresh(myRefresh, "http://mystoke/api/refresh").
//      DialOption()
//    )
//
//    // Login as a user and let the token expire if requests are not made withing expiring time
//    creds := Credentials()
//    if err := creds.Login("myuser", "mypass", "http://mystoke", true); err != nil {
//      log.Fatalf("Could not login: %v", err)
//    }
//    client := grpc.NewClient(creds.DialOption())
//
//    // Login as a user and keep the token up to date with a goroutine
//    ctx, cancel := context.WithCancel(context.Background())
//    creds := Credentials()
//    if err := creds.Login("myuser", "mypass", "http://mystoke", true); err != nil {
//      log.Fatalf("Could not login: %v", err)
//    }
//    creds.StartRefresh(ctx)
//    defer cancel() // this will cancel the refresh goroutine
//
//    client := grpc.NewClient(creds.DialOPtion())
func Credentials() *stokeCredentials {
	return &stokeCredentials{}
}

// Specify the initial token to use
// This token will be refreshed if EnableRefresh is called.
func (c *stokeCredentials) Token(token string) *stokeCredentials {
	c.token = token
	return c
}

// Enables automatic refresh of token
// If ExpireTime is not set before calling, the token is set to expire 30 seconds after calling
//
// Call StartRefresh() to keep token up to date with a goroutine.
// Otherwise the token will only be refreshed when a request is made after the set expiration time.
func (c *stokeCredentials) EnableRefresh(refresh, refreshURL string) *stokeCredentials {
	c.refresh = refresh
	c.refreshURL = refreshURL
	if c.expireTime.Before(time.Now()) {
		c.expireTime = time.Now().Add(30 * time.Second)
	}
	return c
}

// Starts a goroutine to keep the token valid by automatically refreshing.
// You only need to use this if requests are fewer and far between and wont be called frequently enough to refresh per request.
// i.e. requests are more than expireTime apart
// Cancel the context to stop the goroutine.
func (c *stokeCredentials) StartRefresh(ctx context.Context) {
	if c.refresh == "" || c.refreshURL == "" {
		log.Fatalf("Could not start refresh. Must call EnableRefresh to set refresh token and refreshURL.")
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Printf("Context cancelled. Shutting down token refresh.")
				return

			case <-time.After(c.expireTime.Sub(time.Now())):
				err := c.refreshToken()
				if err != nil {
					log.Printf("An error occurred while refreshing token: %s", err)
					return
				}
			}
		}
	}()
}

// Sets the expire time of the token
func (c *stokeCredentials) ExpiresAt(exp time.Time) *stokeCredentials {
	c.expireTime = exp
	return c
}

// Logs user into a given stoke server at stokeURL using username and password.
// stokeURL should point to the root of a stoke server ({stokeURL}/api/login and {stokeURL}/api/refresh should be reachable.
// Useful for long lived or automated access by a service account when enableRefresh is set to true
// This method is not chainable.
func (c *stokeCredentials) Login(username, password, stokeURL string, enableRefresh bool) error {
	res, err := http.Post(
		fmt.Sprintf("%s/api/login", stokeURL),
		"application/json",
		bytes.NewBuffer(
			[]byte(fmt.Sprintf(`{"username" : "%s", "password": "%s"}`, username, password)),
		),
	)
	if err != nil {
		return err
	}

	parsedResp, err := getTokenResponse(res)
	c.updateToken(parsedResp)

	if enableRefresh {
		c.EnableRefresh(parsedResp.Refresh, fmt.Sprintf("%s/api/refresh", stokeURL))
	}

	return nil
}

// Disable requiring transport security.
// USE IN DEVELOPMENT ONLY
func (c *stokeCredentials) DisableSecurity() *stokeCredentials {
	c.disableSecurity = true
	return c
}

// Convert credentials into a grpc PerRPCCredentials Call Option
func (c *stokeCredentials) CallOption() grpc.CallOption {
	return grpc.PerRPCCredentials(c)
}

// Convert credentials into a grpc PerRPCCredentials Call Option
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

	mutex sync.RWMutex
}

// Struct to marshall json received from stoke server
type tokenResponse struct {
	Token string   `json:"token"`
	Refresh string `json:"refresh"`
}

// GetRequestMetadata implements credentials.PerRPCCredentials.
func (c *stokeCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Bearer " + c.token,
	}, nil
}

// RequireTransportSecurity implements credentials.PerRPCCredentials.
func (c *stokeCredentials) RequireTransportSecurity() bool {
	return !c.disableSecurity
}

// Concurrency safe update of token values
func (c *stokeCredentials) updateToken(res *tokenResponse) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.token = res.Token
	c.refresh = res.Refresh

	// TODO: incorperate expireTime info into/from the response
	c.expireTime = time.Now().Add(5 * time.Minute)
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

