import { expect, type Page } from '@playwright/test';

/**
 * Logs in as a user on the admin UI login page.
 * @param page - Playwright page
 * @param baseURL - Base URL (e.g. process.env.STOKE_BASE_URL || 'http://localhost:8080')
 * @param username - Username (default 'tester')
 * @param password - Password (default 'tester')
 */
export async function loginAsTester(
  page: Page,
  baseURL: string,
  username = 'tester',
  password = 'tester'
): Promise<void> {
  const url = baseURL.replace(/\/$/, '') + '/admin';
  await page.goto(url);

  await page.getByLabel('Username').fill(username);
  await page.getByLabel('Password').fill(password);
  await page.getByRole('button', { name: /login/i }).click();

  await expect(page).toHaveURL(/\/user/, { timeout: 10000 });
}
