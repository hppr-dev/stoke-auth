# Go Client Library

# Create the Stoke Client

Public key stores are used to keep the public signing keys up to date and verify received tokens.
These are passed to the middleware.

```go
// Application context (used for tracing)
ctx := context.Background()

// Key store that checks signing key expiration on every request.
keyStore, err := stoke.NewPerRequestPublicKeyStore("http://172.17.0.1:8080", ctx)
if err != nil {
    panic("Unable to pull public keys!")
}

// Key store that actively refreshes public keys within a go routine
// Cancel the context to stop the go routine
keyStore, err := stoke.NewWebCachePublicKeyStore("http://172.17.0.1:8080", ctx)
if err != nil {
    panic("Unable to pull public keys!")
}

// Use TLS
keyStore, err := stoke.NewPerRequestPublicKeyStore("https://172.17.0.1:8080", ctx, stoke.ConfigureTLS("ca.crt"))
if err != nil {
    panic("Unable to pull public keys!")
}
```

# Claim Requirement DSL

A claim DSL is used to specify what token claims are required to access a certain resource.
In an effort to make it more human readable, methods may be chained together to build claim requirements.
Claim requirements start with `stoke.RequireToken()` and follow with any number of:
    * `WithClaim(claim_name, claim_value)` -- require claim with the exact value
    * `WithClaimMatch(claim_name, claim_regex)` -- require claim to match a regular expression
    * `WithClaimContains(claim_name, claim_sub)` -- require claim to have a given substring
    * `Or(stoke.RequireToken()...)` -- provide an alternative claim requirement

All chained "With" functions are required, i.e. all conditions chained off of a single RequireToken are required together (AND).

```go
tokenRequired := stoke.RequireToken()                       // Requires any valid signed token
fooRequired := stoke.RequireToken().WithClaim("foo", "bar") // Requires a valid signed token with the "foo" : "bar

complexRequirement := stoke.RequireToken(). // Complex rules can be created using chaining
    WithClaim("hello", "world").            // Token must have "hello" : "world"
    WithClaimMatch("foo", "^bar.*").        // And must have a claim "foo" that matches the regex "^bar.*"
    Or(stoke.RequireToken().                // Or
        WithClaim("hello", "you").          // Token must have "hello" : "you"
        WithClaimContains("bar", "sup"),    // and must have a claim "bar" that contains the letters "sup"
    )
// The above allows tokens with the following claims:
//    { "hello" : "world" , "foo" : "barific" }
//    { "hello" : "you", "bar" : "heysupdog" }
//    { "hello" : "world,you" , "bar" : "super,sticky" }
// But none of the following:
//    { "hello" : "world", "bar" : "super" }
//    { "hello" : "dog" }
//    { "hello" : "you", "foo" : "barf" }

// NOTE: chained functions modify the existing claim requirement.
// For example requirement1 has the same effect as requirement2:
requirement1 := stoke.RequireToken().WithClaim("foo", "bar")

requirement2 := stoke.RequireToken()
// Modifies the original
requirement2.WithClaim("foo", "bar")

```

# HTTP Handler Wrappers

```go
// Anonymous function
handler := stoke.AuthFunc(
    func (res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte("Hello world!"))
    }
    keyStore,
    stoke.RequireToken().WithClaim("foo", "bar"), // Requires token with a claim "foo":"bar"
)

// Handler
handler := stoke.Auth(
    innerHandler,
    keyStore,
    stoke.RequireToken() // Require any valid signed token
)
```


# GRPC

```go
// Server
server := grpc.NewServer(
	grpc.Creds(grpcCreds),
	stoke.NewUnaryTokenInterceptor(keyStore, stoke.RequireToken().WithClaim("un", "y")).Opt(),  // Intercept all unary endpoints and require "un" : "y"
	stoke.NewStreamTokenInterceptor(keyStore, stoke.RequireToken().WithClaim("st", "y")).Opt(), // Intercept all streaming endpoint and require "st" : "y"
)

// Client
tokenStr := "MYJWTTOKEN"

// Use token for single call with CallOption()
result, err := grpc.NewClient(grpcConn).FooBarTest(ctx, stoke.Credentials().Token(tokenStr).CallOption())
if err != nil {
    panic("could not create grpc client")
}

// Use same token for multiple grpc calls with DialOption
client := grpc.NewClient(Credentials(myToken).DialOption())
```
