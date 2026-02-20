# Client integration (token used for downstream services)

Stories for using a Stoke-issued token to call downstream services. Covered by test/k6/client_integration.js with client examples and Docker Compose.

## US-007: User with correct claims can access protected endpoints

**Acceptance criteria:** After obtaining a token via `/api/login`, when the user calls downstream HTTP and WebSocket endpoints with `Authorization: Bearer <token>` (or token query for WS), endpoints return 200 (HTTP) or 101 (WS) when the token has the required claims.

**Tested by:**

| Test | Location | How |
|------|----------|-----|
| k6 | test/k6/client_integration.js | access_granted: login (leela/fry), then GET/WS to client examples with Bearer; check 200/101 |

---

## US-008: User without required claims receives 403

**Acceptance criteria:** When a token lacks the claims required by a downstream service, requests with that token return 403 (or equivalent denial).

**Tested by:**

| Test | Location | How |
|------|----------|-----|
| k6 | test/k6/client_integration.js | access_denied: login (hermes/professor), call same services; check 403/denial |
