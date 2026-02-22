# Retrieve a user token

Stories for obtaining a JWT (and refresh token) via the login API. These are **already covered** by existing integration tests; this file documents the coverage so E2E can align or extend later.

## API

- **Endpoint:** `POST /api/login`
- **Body:** `{"username":"...","password":"...","provider":"..."}` (provider optional for local)

## Stories and current test coverage

### US-001: User can obtain a token with valid local credentials

**Acceptance criteria:** Given a user created via dbinit (local), when they POST to `/api/login` with correct username and password, the response is 200 and includes `token` and `refresh`.

**Tested by:**


| Test                     | Location                                               | How                                                                               |
| ------------------------ | ------------------------------------------------------ | --------------------------------------------------------------------------------- |
| k6 smoke "issued"        | [test/k6/smoke.js](test/k6/smoke.js)                   | `okLogin("tester", "tester", 0.01)`                                               |
| k6 common                | [test/k6/common.js](test/k6/common.js)                 | `okLogin()` checks status 200 and `token.token && token.refresh`                  |
| k6 provider (local)      | [test/k6/provider_test.js](test/k6/provider_test.js)   | `localIssued` → `okLogin("tester", "tester", 0.01)`                               |
| k6 database/cert configs | [test/README.md](test/README.md)                       | Same smoke.js with different server configs                                       |
| Go benchmark             | [test/login_test.go](test/login_test.go)               | `BenchmarkParallelSingleLogin` (tester/tester)                                    |
| Go capabilities          | [test/capabilities_test.go](test/capabilities_test.go) | `loginForToken()` uses login API then calls `/api/capabilities` with Bearer token |
| E2E                      | [test/e2e/specs/login-api.spec.ts](test/e2e/specs/login-api.spec.ts) | POST /api/login with valid local credentials returns token and refresh @US-001 |


**Dbinit (local user):** [test/configs/dbinit/smoke_test.yaml](test/configs/dbinit/smoke_test.yaml) defines user `tester` with password hash/salt.

---

### US-002: User is rejected when credentials are invalid

**Acceptance criteria:** When POST to `/api/login` with wrong password (or unknown user), the response is 401 and does not include token or refresh.

**Tested by:**


| Test                | Location                                             | How                                                         |
| ------------------- | ---------------------------------------------------- | ----------------------------------------------------------- |
| k6 smoke "rejected" | [test/k6/smoke.js](test/k6/smoke.js)                 | `badLogin("tester", "badpass", 0.01)`                       |
| k6 common           | [test/k6/common.js](test/k6/common.js)               | `badLogin()` checks status 401 and absence of token/refresh |
| k6 provider         | [test/k6/provider_test.js](test/k6/provider_test.js) | `localRejected`, `providerRejected`                         |


---

### US-003: User can obtain a token via LDAP provider

**Acceptance criteria:** When an LDAP provider is configured and the user exists in LDAP, POST to `/api/login` with that user’s credentials (and provider name if required) returns 200 with token and refresh.

**Tested by:**


| Test        | Location                                             | How                                                                                                                                                                                                                                                                                                                                 |
| ----------- | ---------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| k6 provider | [test/k6/provider_test.js](test/k6/provider_test.js) | `providerIssuedUser` → `okLogin("hermes", "hermes", 0.01)`, `providerIssuedAdmin` → `okLogin("fry", "fry", 0.01)`; setup creates users via login against server backed by [test/configs/provider_type/ldap.yaml](test/configs/provider_type/ldap.yaml) and LDAP container from [test/docker-compose.yaml](test/docker-compose.yaml) |


**Config:** Provider-type integration uses [test/configs/provider_type/](test/configs/provider_type/) and [test/configs/dbinit/provider_test.yaml](test/configs/dbinit/provider_test.yaml). Run with `task test int -p` (provider).

---

### US-004: Admin and user tokens differ by provider (local vs LDAP)

**Acceptance criteria:** Tokens issued for local user vs LDAP user both work; admin vs non-admin users get tokens with different claims (e.g. stk=S for admin).

**Tested by:**


| Test        | Location                                             | How                                                                                                                                                              |
| ----------- | ---------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| k6 provider | [test/k6/provider_test.js](test/k6/provider_test.js) | Separate scenarios for `localIssued`, `providerIssuedUser`, `providerIssuedAdmin`; capabilities test uses login with `required_claims`/`filter_claims` for admin |


**Note:** [test/capabilities_test.go](test/capabilities_test.go) obtains a token with admin claims (stk) and calls `/api/capabilities` and `/api/available_providers`, demonstrating token use for admin API.

---

## Running these tests

- **All integration tests (includes token flows):** `task test int -a`
- **Smoke only (cert/database + token):** `task test int -c`, `task test int -d`
- **Provider (LDAP + token):** `task test int -p`
- **Capabilities (Go, needs server):** run Stoke then `go test ./test/... -run Capabilities` from repo root (or via task if wired).

