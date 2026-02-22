# User stories

User stories for stoke-auth live here. Each story has a stable **ID** used to trace implementation and automated tests.

## Format

- One file per theme or feature area (e.g. `retrieve-user-token.md`).
- Each story: **ID**, **title**, **acceptance criteria**, and **test coverage** (where it is tested: E2E in `test/e2e/`, k6 in `test/k6/`, Go in `test/*_test.go`, etc.).

## Index

| ID | Title | File | Covered by |
|----|--------|------|------------|
| US-001 | User can obtain a token with valid local credentials | [retrieve-user-token.md](retrieve-user-token.md) | k6, Go |
| US-002 | User is rejected when credentials are invalid | [retrieve-user-token.md](retrieve-user-token.md) | k6 |
| US-003 | User can obtain a token via LDAP provider | [retrieve-user-token.md](retrieve-user-token.md) | k6 |
| US-004 | Admin and user tokens differ by provider | [retrieve-user-token.md](retrieve-user-token.md) | k6 |
| US-005 | Authenticated admin can get capabilities | [admin-api.md](admin-api.md) | Go |
| US-006 | Anonymous user can list login providers | [admin-api.md](admin-api.md) | Go |
| US-007 | User with correct claims can access protected endpoints | [client-integration.md](client-integration.md) | k6 |
| US-008 | User without required claims receives 403 | [client-integration.md](client-integration.md) | k6 |
| US-009 | Admin can view Users, Groups, and Claims pages (smoke) | [admin-ui.md](admin-ui.md) | E2E |
| US-010 | Admin can list users, list groups, list claims | [admin-ui.md](admin-ui.md) | E2E |
| US-011 | Admin can create claim → group → user and that user can retrieve a JWT | [admin-ui.md](admin-ui.md) | E2E |
| US-012 | Generated dbinit file creates a loginable user | [admin-ui.md](admin-ui.md) | E2E |
| US-013 | With HA enabled, each replica serves merged JWKS at /api/pkeys | [high-availability.md](high-availability.md) | Go, E2E |
| US-014 | Token issued by one replica is verified with federated JWKS | [high-availability.md](high-availability.md) | Go |

(Add rows as more stories are added.)

## E2E traceability

When Playwright (or other E2E) is added under `test/e2e/`, tests should reference story IDs in `describe` blocks or tags, e.g. `@US-001`, so that "story implemented" means the corresponding E2E scenario(s) pass.
