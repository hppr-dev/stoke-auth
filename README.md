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

The following results were found using the scripts in the benchmark and benchmark/k6 directory.
Tests were run with a sqlite file database.
Logins (AKA token requests) were done with a user with 5 groups with 50 seperate claims.

Tests were run with the following failure conditions:
    * If more than 1% of requests fail, the tests stop
    * If more than 1% of requests take more than 300 ms , the tests stop

Default run environment:
    * X nanoseconds / token in isolation
    * Breakpoint:       X tokens / second
    * High Load Test:   X tokens / second
    * Medium Load Test: X tokens / second
    * Suggested Max Nominal Load: X tokens / second

GOMAXPROCS=1:
    * X nanoseconds / token in isolation
    * Breakpoint:       X tokens / second
    * High Load Test:   X tokens / second
    * Medium Load Test: X tokens / second
    * Suggested Max Nominal Load: X tokens / second


