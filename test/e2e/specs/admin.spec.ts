import { test, expect } from '@playwright/test';

// US-006: Anonymous user can list login providers
test('GET /api/available_providers returns providers @US-006', async ({ request }) => {
  const res = await request.get('/api/available_providers');
  expect(res.ok()).toBeTruthy();
  const body = await res.json();
  expect(Array.isArray(body.providers)).toBe(true);
});

// Admin UI: login page loads
test('admin UI login page loads', async ({ page }) => {
  await page.goto('/admin');
  await expect(page).toHaveURL(/\/admin/);
  await expect(page.getByRole('heading', { name: /login|sign in|stoke/i })).toBeVisible({ timeout: 10000 });
});
