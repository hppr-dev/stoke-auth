import { test, expect } from '@playwright/test';
import { E2E_DBINIT_USERNAME, E2E_DBINIT_PASSWORD } from '../fixtures/dbinit';

const baseURL = () => process.env.STOKE_BASE_URL || 'http://localhost:8080';

// US-012: Run only when server was started with generated dbinit (E2E_USE_GENERATED_DBINIT=1).
test('generated dbinit user can retrieve JWT from login endpoint @US-012', async ({ request }) => {
  test.skip(!process.env.E2E_USE_GENERATED_DBINIT, 'Set E2E_USE_GENERATED_DBINIT=1 and start server with generated dbinit');

  const loginRes = await request.post(baseURL() + '/api/login', {
    data: { username: E2E_DBINIT_USERNAME, password: E2E_DBINIT_PASSWORD },
    headers: { 'Content-Type': 'application/json' },
  });
  expect(loginRes.status()).toBe(200);
  const body = await loginRes.json();
  expect(body.token).toBeDefined();
  expect(body.refresh).toBeDefined();
});
