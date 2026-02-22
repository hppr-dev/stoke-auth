import { test, expect } from '@playwright/test';
import { loginAsTester } from '../helpers/auth';

// US-006: Anonymous user can list login providers
test('GET /api/available_providers returns providers @US-006', async ({ request }) => {
  const res = await request.get('/api/available_providers');
  expect(res.ok()).toBeTruthy();
  const body = await res.json();
  expect(Array.isArray(body.providers)).toBe(true);
});

// Admin UI: login page loads (card shows "Stoke" and a Login button; no semantic heading)
test('admin UI login page loads', async ({ page }) => {
  await page.goto('/admin');
  await expect(page).toHaveURL(/\/admin/);
  await expect(page.getByText('Stoke').first()).toBeVisible({ timeout: 10000 });
  await expect(page.getByRole('button', { name: /login/i })).toBeVisible({ timeout: 5000 });
});

const baseURL = () => process.env.STOKE_BASE_URL || 'http://localhost:8080';

test('admin UI login as stoke shows user area', async ({ page }) => {
  await loginAsTester(page, baseURL());
  await expect(page).toHaveURL(/\/user/);
  await expect(
    page.locator('[data-testid="nav-users"], [data-testid="nav-groups"], [data-testid="nav-claims"]').first()
  ).toBeVisible({ timeout: 5000 });
});
