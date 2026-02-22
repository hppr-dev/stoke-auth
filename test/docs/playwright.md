# Running Playwright E2E tests

This document describes how to run the Playwright end-to-end tests for the Stoke auth server. The tests live in **test/e2e/** and exercise the admin UI and API against a running Stoke instance.

## Prerequisites

- **Node.js** 22 or later (for `test/e2e`)
- **Docker** (when using automatic server start)
- **Stoke server** running and reachable, **or** use automatic start/teardown (see below)
- **npm** (comes with Node)

Playwright will install browser binaries (Chromium by default) on first run.

## Does the task start the server?

**Yes, unless you pass `-n`.** From the repo root:

- **`task test e2e`** — Builds the Stoke image (if missing), starts a single container with `test/e2e/configs/`, runs Playwright, then stops the container.
- **`task test e2e --ha`** — Starts the HA stack (Postgres + two Stoke replicas) via `docker-compose.e2e-ha.yaml`, runs Playwright with `STOKE_BASE_URL` and `STOKE_BASE_URL_REPLICA_2` set, then tears down the stack.
- **`task test e2e -n`** — Runs Playwright only (no server start/stop). Use when you already have a server running. Set `STOKE_BASE_URL` if it’s not `http://localhost:8080`.

Run the test suite on different configs by starting the server yourself (single or HA, or custom), then run **`task test e2e -n`** or from **test/e2e**: **`STOKE_BASE_URL=... npm run test`**.

## Quick start

1. **From repo root (single server):** `task test e2e`
2. **From repo root (HA):** `task test e2e --ha`
3. **No server (use existing):** `task test e2e -n` or `cd test/e2e && STOKE_BASE_URL=http://localhost:8080 npm run test`

## Environment variables

| Variable | Default | Description |
|----------|---------|-------------|
| `STOKE_BASE_URL` | `http://localhost:8080` | Base URL of the Stoke server. |
| `STOKE_BASE_URL_REPLICA_2` | — | Set by `task test e2e --ha`; used by HA specs. |

## Headed vs headless

Tests run **headless** (no visible browser) by default.

- **Headless (default):**  
  `npm run test` or `task test e2e`

- **Headed (browser window visible):**  
  `npm run test:headed` from `test/e2e/`, or:
  ```bash
  task test e2e -h
  ```
  (or `task test e2e --headed`). You can combine with no-server: `task test e2e -n -h`.

- **UI mode (interactive):**  
  `npm run test:ui` from `test/e2e/`

All commands are run from **test/e2e/** unless you use `task test e2e` (which changes into that directory for you).

## Running a subset of tests

- **Single file:**  
  `npx playwright test specs/admin.spec.ts`

- **Single test by name:**  
  `npx playwright test -g "available_providers"`

- **By story tag (e.g. US-001):**  
  `npx playwright test -g "@US-001"`

- **UI mode (interactive):**  
  `npm run test:ui`

All commands above are run from **test/e2e/** unless you use `task test e2e` (which changes into that directory for you).

## Viewing the report

After a run, an HTML report is generated. To open it:

```bash
cd test/e2e
npm run report
```

Or open **test/e2e/playwright-report/index.html** in a browser (path after default `reporter: 'html'`).

## What the tests expect

- **Admin UI login page:** Tests expect the admin app at `$STOKE_BASE_URL/admin` to load and show a heading containing "login", "sign in", or "stoke" (case-insensitive).
- **API tests:**  
  - `GET /api/available_providers` must return 200 and JSON with a `providers` array.  
  - `POST /api/login` with `{"username":"tester","password":"tester"}` must return 200 with `token` and `refresh` when the server has a local user `tester` with password `tester` (e.g. from dbinit such as `test/configs/dbinit/smoke_test.yaml`).

If the server is not running or the base URL is wrong, tests will fail (e.g. connection refused or 404).

## User stories and traceability

User stories are documented in **test/stories/** with stable IDs (e.g. US-001, US-006). E2E tests reference them in the test name with `@US-xxx`. See:

- [test/stories/README.md](../stories/README.md) – index of all stories  
- [test/stories/retrieve-user-token.md](../stories/retrieve-user-token.md) – token/login stories  
- [test/stories/admin-api.md](../stories/admin-api.md) – capabilities and available_providers  

## CI

In CI, run E2E after starting the Stoke server (and any required backing services). Set `STOKE_BASE_URL` to the URL the job uses to reach the server. The task script uses `npm ci` and installs Chromium so that the run is reproducible.

## Troubleshooting

- **Tests fail with connection refused:** Ensure Stoke is running and `STOKE_BASE_URL` matches the URL (host and port) where it listens.
- **Login test fails (401 or missing token):** Ensure the server is started with dbinit that creates user `tester` with password `tester` (e.g. smoke_test.yaml).
- **Admin UI test fails (heading not found):** The app may use different text; adjust the test in `test/e2e/specs/admin.spec.ts` to match the actual heading, or ensure the admin UI is built and served at `/admin`.
