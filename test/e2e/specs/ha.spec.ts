import { test, expect } from '@playwright/test';

// US-013: With HA enabled, each replica serves a merged JWKS at /api/pkeys
test('GET /api/pkeys returns 200 and valid JWKS @US-013', async ({ request }) => {
  const res = await request.get('/api/pkeys');
  expect(res.status()).toBe(200);
  const body = await res.json();
  expect(body).toHaveProperty('keys');
  expect(Array.isArray(body.keys)).toBe(true);
  expect(body.keys.length).toBeGreaterThanOrEqual(1);
  for (const k of body.keys) {
    expect(k).toHaveProperty('kid');
    expect(k).toHaveProperty('kty');
  }
});

// HA profile only: replica 2's /api/pkeys returns merged JWKS (at least 2 keys from both replicas)
test('GET replica 2 /api/pkeys returns merged JWKS when HA @US-013', async ({ request }) => {
  const replica2Base = process.env.STOKE_BASE_URL_REPLICA_2;
  test.skip(!replica2Base, 'STOKE_BASE_URL_REPLICA_2 not set (run with E2E_SERVER_PROFILE=ha)');
  const res = await request.get(replica2Base + '/api/pkeys');
  expect(res.status()).toBe(200);
  const body = await res.json();
  expect(body).toHaveProperty('keys');
  expect(Array.isArray(body.keys)).toBe(true);
  expect(body.keys.length).toBeGreaterThanOrEqual(2);
});
