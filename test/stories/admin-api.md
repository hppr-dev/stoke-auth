# Admin API (capabilities and available providers)

## US-005: Authenticated admin can get capabilities

**Acceptance criteria:** When a user with admin claims (e.g. stk) has a valid token and GETs `/api/capabilities` with `Authorization: Bearer <token>`, the response is 200 and JSON includes a `capabilities` array.

**Tested by:**

| Test | Location | How |
|------|----------|-----|
| Go | test/capabilities_test.go | TestCapabilities_ReturnsCapabilities; loginForToken() then GET capabilities with Bearer |

---

## US-006: Anonymous user can list login providers

**Acceptance criteria:** When GETting `/api/available_providers` without auth, the response is 200 and JSON includes a `providers` array; optionally `base_admin_path` is present.

**Tested by:**

| Test | Location | How |
|------|----------|-----|
| Go | test/capabilities_test.go | TestAvailableProviders_ReturnsProvidersAndBaseAdminPath |
| E2E | test/e2e/specs/admin.spec.ts | GET /api/available_providers returns providers @US-006 |
