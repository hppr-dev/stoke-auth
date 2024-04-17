# Stoke Authentication Server

A drop in solution for authentication.

Stoke Authentication Server is a server to authenticate users with Json Web Tokens.
It is meant to be a lightweight way of bringing authentication to microservices.

Stoke aims to be:
  * Lightweight
  * Simple
  * Secure

Stoke Features:
  * Simple deployment
    * All you need is a docker image and a config file
  * Keys included
    * Go HTTP/GRPC client middleware available
  * Use customizable login page or api
  * Automatic public key distribution
  * Token renewal
  * Admin console
    * Manage claims
      * Create/Update/Delete claims
      * Issue claims by group or user
    * Manage keys
      * Configure automatic key rotation
      * Configure algorithm
      * Manual key rotation
    * Manage token expiration
    * Configure credential sources
      * LDAP and Local DB credential sources

Non-goals:
  * Support different authentication schemes
  * Replace identity providers
  * All encompassing SSO server

# Configuration

TODO

# Client Usage

TODO

## Key Verification

Clients must use the public keys from the /api/pkeys endpoint to authenticate requests.
All public keys have an expiry date and a renew date.
The expire date indicates when the key will expire and any JWTs signed with it will no longer be valid.
The renew date indicates when a new key will be generated and any new JWTs will be signed with the new key.
The expire date must always be at least 1 JWT expiration after the renew date.
For example, a public key may have an expiration date in T+3hours and JWTs are issued with an expiration date of T+30minutes.
The renewal date will fall before T+2.5hours.

Clients should pull keys on start up and maintain the list of public keys.
JWTs remain valid if they were signed with any of the given pkeys.
Clients should remove keys on their expiration date and update keys on the renewal date.
Public keys are guarenteed to be available at the renewal date.
If a manual rotation is performed, clients can wait until the next renewal/expiration date to pick up the change or perform a manual sync of the keys.
A periodic check of the keys may be used to automatically update keys.
In other words, it is not necessary to check verification keys on every request.

JWTs remain valid for their full lifetime.
A manual rotation and sync of the keys may be used to invalidate any tokens signed with a key.

# HTTP Endpoints

TODO document endpoint usage

- /login -- configured login page (HTML)
- /admin -- the admin console (HTML)
- /api -- JSON api
  - /api/pkeys -- current valid public verification keys
  - /api/login -- JSON login
  - /api/renew -- renew a given JWT

# Performance/Load Benchmarks

Tests were run with the following failure conditions:
    * If more than 1% of requests fail, the tests stop
    * If more than 1% of requests take more than 300 ms, the tests stop

The server was configured as follows
    * Log level Info
    * Pretty logging disabled
    * Log to stdout only
    * Tracing enabled
    * Local sqlite database
    * ECDSA keys
    * 1s request timeout
    * Connected to a test LDAP container instance
    * Token issued for a user with 3 groups and 30 seperate claims
    * Aggressive key/token rotation:
        * Key Duration: 5m
        * Token duration 1m

Default run environment:
    * ~5 milliseconds / token in isolation
    * Breakpoint: ~230 tokens / second
    * High Load Test:        200 tokens / second for 10m with no loss
    * Medium Load Test:      150 tokens / second for 20m with no loss
    * Max Nominal Load Test: 100 tokens / second for 2h with no loss

GOMAXPROCS=1:
    * ~20 milliseconds / token in isolation
    * Breakpoint: ~58 tokens / second
    * High Load Test:        50 tokens / second for 10m with no loss
    * Medium Load Test:      25 tokens / second for 20m with no loss
    * Max Nominal Load Test: 15 tokens / second for 2h with no loss

In resource constrained environments with low traffic expectations,GOMAXPROCS=1 can be used to limit the memory/cpu footprint of the server.
Otherwise it is recommended to run the server without modifying GOMAXPROCS.

# High Availablility

TODO
Load balancing and high availability are features that are planned for the future.
