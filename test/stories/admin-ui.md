# Admin UI (smoke, list flows, create flow, dbinit)

## US-009: Admin can view Users, Groups, and Claims pages (smoke)

**Acceptance criteria:** After logging in, the admin can navigate to the Users, Groups, and Claims pages and each page loads and shows expected content (e.g. list or table headers).

**Covered by:** E2E [test/e2e/specs/admin-pages.spec.ts](../e2e/specs/admin-pages.spec.ts) (smoke tests for /user, /group, /claim).

---

## US-010: Admin can list users, list groups, list claims (key flows)

**Acceptance criteria:** After logging in, the Users, Groups, and Claims pages show list/table data (at least one row or a clear empty state).

**Covered by:** E2E [test/e2e/specs/admin-pages.spec.ts](../e2e/specs/admin-pages.spec.ts) (list users, list groups, list claims tests).

---

## US-011: Admin can create claim → group (with claim) → user (with group) and that user can retrieve a JWT from the login endpoint

**Acceptance criteria:** Admin creates a claim, creates a group and assigns the claim to it, creates a user and assigns the group to them; then the new user can call the login API and receive a JWT (token and refresh).

**Covered by:** E2E [test/e2e/specs/admin-create-flow.spec.ts](../e2e/specs/admin-create-flow.spec.ts) (full create flow then POST /api/login assertion).

---

## US-012: Generated dbinit file creates a loginable user

**Acceptance criteria:** A minimal dbinit YAML file can be generated (with one user and required groups/claims); when the server is started with that dbinit, the user can log in via the login endpoint.

**Covered by:** E2E [test/e2e/specs/dbinit.spec.ts](../e2e/specs/dbinit.spec.ts) (when server is run with generated dbinit).
