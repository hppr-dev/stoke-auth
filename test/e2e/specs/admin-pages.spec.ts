import { test, expect } from '@playwright/test';
import { loginAsTester } from '../helpers/auth';

const baseURL = () => process.env.STOKE_BASE_URL || 'http://localhost:8080';

test.describe('admin pages when logged in', () => {
  test.beforeEach(async ({ page }) => {
    await loginAsTester(page, baseURL());
  });

  // US-009: Smoke – Users, Groups, Claims pages load
  test('Users page loads @US-009', async ({ page }) => {
    await page.goto('/admin/user');
    await expect(page.getByText('Username').first()).toBeVisible({ timeout: 10000 });
  });

  test('Groups page loads @US-009', async ({ page }) => {
    await page.goto('/admin/group');
    await expect(page.getByText('Group Name').first()).toBeVisible({ timeout: 10000 });
  });

  test('Claims page loads @US-009', async ({ page }) => {
    await page.goto('/admin/claim');
    await expect(page.getByText('Claim Name').first()).toBeVisible({ timeout: 10000 });
  });

  // US-010: List flows – list users, list groups, list claims
  test('list users shows at least one row or empty state @US-010', async ({ page }) => {
    await page.goto('/admin/user');
    await expect(page.getByText('Username').first()).toBeVisible({ timeout: 10000 });
    const rows = page.locator('table tbody tr');
    const count = await rows.count();
    expect(count >= 0).toBe(true);
    if (count === 0) {
      await expect(page.getByRole('button', { name: /add user/i })).toBeVisible({ timeout: 5000 });
    } else {
      await expect(rows.first()).toBeVisible();
    }
  });

  test('list groups shows at least one row or empty state @US-010', async ({ page }) => {
    await page.goto('/admin/group');
    await expect(page.getByText('Group Name').first()).toBeVisible({ timeout: 10000 });
    const rows = page.locator('table tbody tr');
    const count = await rows.count();
    expect(count >= 0).toBe(true);
    if (count === 0) {
      await expect(page.getByRole('button', { name: /add group/i })).toBeVisible({ timeout: 5000 });
    } else {
      await expect(rows.first()).toBeVisible();
    }
  });

  test('list claims shows at least one row or empty state @US-010', async ({ page }) => {
    await page.goto('/admin/claim');
    await expect(page.getByText('Claim Name').first()).toBeVisible({ timeout: 10000 });
    const rows = page.locator('table tbody tr');
    const count = await rows.count();
    expect(count >= 0).toBe(true);
    if (count === 0) {
      await expect(page.getByRole('button', { name: /add claim/i })).toBeVisible({ timeout: 5000 });
    } else {
      await expect(rows.first()).toBeVisible();
    }
  });
});
