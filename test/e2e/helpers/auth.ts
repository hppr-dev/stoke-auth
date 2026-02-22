import { expect, type Page } from '@playwright/test';

/**
 * Logs in as a user on the admin UI login page.
 * E2E config (test/e2e/configs/dbinit.yaml) provides user stoke / password admin.
 * @param page - Playwright page
 * @param baseURL - Base URL (e.g. process.env.STOKE_BASE_URL || 'http://localhost:8080')
 * @param username - Username (default 'stoke')
 * @param password - Password (default 'admin')
 */
export async function loginAsTester(
  page: Page,
  baseURL: string,
  username = 'stoke',
  password = 'admin'
): Promise<void> {
  const url = baseURL.replace(/\/$/, '') + '/admin';
  await page.goto(url);

  await page.getByLabel('Username').fill(username);
  await page.locator('input[type="password"]').fill(password);
  await page.getByRole('button', { name: /login/i }).click();

  await expect(page).toHaveURL(/\/user/, { timeout: 10000 });
}
