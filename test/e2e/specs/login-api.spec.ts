import { test, expect } from '@playwright/test';

// US-001: User can obtain a token with valid local credentials
test('POST /api/login with valid local credentials returns token and refresh @US-001', async ({ request }) => {
  const res = await request.post('/api/login', {
    data: { username: 'tester', password: 'tester' },
    headers: { 'Content-Type': 'application/json' },
  });
  expect(res.status()).toBe(200);
  const body = await res.json();
  expect(body.token).toBeDefined();
  expect(body.refresh).toBeDefined();
});
