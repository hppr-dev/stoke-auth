# Stoke Authentication Clients

Clients libraries implement the following:

1. Managing, i.e. pulling and refreshing, public keys in JWK format
    * Public keys should be pulled from the /api/pkeys endpoint
    * Clients must refresh keys at the keyset expire time returned by /api/pkeys
    * The keyset returned by /api/pkeys may be cached until the time at which they expire
    * Stoke's /api/pk
    
2. Extract the token from a given requests Authentication header using the standard "Bearer" prefix, i.e. "Authentication: Bearer <TOKEN>".

3. Verify the token's signature against the current public keys

4. Make claims available in plain text for the rest of the application

Some considerations:

  * Clients are light weight wrappers that connect JWT libraries and the public keys provided by a stoke deployment
  * Client libraries should be as framework agnostic as possible by:
    * Using standard web request/handling libraries
    * Exposing verification functionality
  * Clients should provide testing structures to allow mocking/faking during tests
