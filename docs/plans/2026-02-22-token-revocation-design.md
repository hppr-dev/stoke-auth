# Token revocation – design

**Date:** 2026-02-22  
**Status:** Design approved

## Goals

- Support **user-level** revocation: (a) user-initiated “log me out everywhere” / invalidate my sessions, and (b) admin revoking a specific user’s access.
- **Enable** revocation rather than mandate it: server enforces at refresh (401 when revoked) and, when configured, at Stoke’s own API; resource-server clients may optionally respect revocation.
- Keep state minimal and HA-friendly: per-user revocation timestamp only, no per-token storage.

## Out of scope

- Per-token or per-session revocation (revoke only one device/session).
- Mandating that every resource server enforce revocation; clients opt in.

## Architecture

### Revocation model

- **State:** One revocation record per user: “this user was revoked at time T” (e.g. `username` + `revoked_at`). Stored in the shared database so all replicas see the same state in HA.
- **Rule:** A token is considered revoked if it belongs to that user and was issued before the revocation time: `token.iat < revoked_at` (or equivalent when `iat` is derived). User identity is taken from the token’s username claim (Stoke already puts this in token claims via `user_info.username`).
- **Revoke user:** Set or update `revoked_at = now()` for that user. “Log me out everywhere” and “admin revokes user” both map to this.
- **No jti storage:** Issued tokens are not stored; revocation is purely by user + issued-before time.

### Where revocation is enforced

1. **Refresh (mandatory when revocation is enabled):** Before issuing a new token, the server checks the user and `iat` from the current token against the revocation store; if revoked → 401.
2. **Stoke’s own API (optional, config-driven):** When revocation is enabled, after validating the Bearer token (signature + claims), the server checks revocation before accepting the request; if revoked → 401.
3. **Resource servers (optional for clients):** A cacheable revocation API (e.g. list of revoked users + `revoked_at`) allows resource-server clients that wish to respect revocation to fetch and cache this data and reject tokens whose user+iat is revoked. Use is optional.

### Configuration

- **Enable revocation:** A config flag (e.g. under `tokens` or `users`) turns on revocation: persistence of revocation state, checks at refresh, and optionally at Stoke API. When disabled, no revocation table is used and no revocation checks are performed.
- **Token requirements when revocation is on:** Tokens must carry a stable user identity (username claim) and an issued-at time. If `include_issued_at` is not already set, the server should require it when revocation is enabled (or derive “issued before” from `exp` and configured token duration). Design assumes the username claim key is known from existing token config (`user_info.username`).

## Data model

- **Revocation store:** Table (or equivalent) keyed by user identity. Suggested shape:
  - `username` (string, primary key or unique) — matches the value in the token’s username claim (e.g. local user or LDAP username).
  - `revoked_at` (timestamp) — time at which all tokens issued before this time for this user are considered revoked.
- **Index:** Lookup by username; no need to list by time for the core flow. Optional: index on `revoked_at` for cleanup or analytics.

## APIs and flows

### Revocation actions

- **User “log me out everywhere”:** New endpoint or existing auth flow (e.g. `POST /api/revoke` or `POST /api/logout` with Bearer token) that sets `revoked_at = now()` for the authenticated user. Requires a valid token (so the server knows the username). Returns 204 or 200.
- **Admin revokes user:** Admin API (e.g. “revoke tokens for user X”) that sets `revoked_at = now()` for the given user. Protected by existing admin claims (e.g. `stk=U` or `stk=S`). Returns 204 or 200.

### Refresh flow

- Existing `POST /api/refresh` with `token` + `refresh`: after verifying the refresh token and parsing the JWT, the server resolves username and `iat` from the token. If revocation is enabled, look up the user in the revocation store; if present and `revoked_at > iat`, return 401. Otherwise proceed to issue a new token pair.

### Stoke API (admin, capabilities, etc.)

- When revocation is enabled, the security handler (or the layer that validates the Bearer token and injects claims) performs a revocation check after successful signature/claims validation: resolve username and `iat` from the token, check revocation store; if revoked, return 401. When revocation is disabled, no check.

### Client-optional revocation (resource servers)

- **Revocation data endpoint:** e.g. `GET /api/revocation` (or similar) returning a small payload: list of `{ username, revoked_at }` (or equivalent) for all users with an active revocation. Cacheable (e.g. short TTL, or ETag). Authenticated if desired (e.g. same Bearer as for other APIs) or unauthenticated with rate limiting, depending on policy.
- **Client behavior:** Resource-server clients that want to respect revocation fetch this payload (periodically or on demand), cache it, and before accepting a token check that the token’s username and `iat` are not revoked (e.g. `revoked_at > iat` for that username). Not required for all clients; enabling this is the goal.

## Error handling

- **Refresh:** Revoked user → 401 Unauthorized (same as invalid or expired refresh). No distinct error code required for “revoked” to avoid leaking information; logging can distinguish for diagnostics.
- **Stoke API:** Revoked token → 401 Unauthorized.
- **Revocation endpoints:** Invalid or missing auth → 401. Admin revoke for non-existent user → 404 or 204 (idempotent “revoked” is acceptable). User revoke (logout) with invalid token → 401.

## HA and replication

- Revocation state lives in the same shared database used for users, groups, and claims. All replicas read and write the same revocation table; no separate sync mechanism. Consistent with the existing HA design (shared DB, federated keys).

## Testing

- **Unit:** Revocation store: set revoked_at for user; lookup returns correct value; “is revoked” for (username, iat) is true when iat < revoked_at, false otherwise.
- **Integration:** With revocation enabled: (1) User revokes self (logout) → subsequent refresh returns 401; (2) Admin revokes user → that user’s refresh returns 401; (3) New login after revoke returns new tokens; (4) Stoke API (e.g. capabilities or admin) returns 401 for revoked token; (5) Revocation endpoint returns expected list. With revocation disabled: no revocation table usage and no 401s due to revocation.
- **Optional:** Resource-server client test that, when configured to use revocation data, rejects a token whose user+iat is revoked and accepts when not revoked.

## Documentation

- README or docs section on “Token revocation”: how to enable it, that it is user-level, that refresh and Stoke API enforce it when enabled, and that resource servers can optionally use the revocation API to enforce it. Note requirement for username and issued-at in tokens when revocation is on.

## Next step

After this design is committed, invoke the writing-plans skill to produce a detailed implementation plan (schema, config, refresh/API checks, revocation endpoints, client-optional revocation API, and tests).
