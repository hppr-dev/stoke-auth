import { test, expect } from '@playwright/test';
import { loginAsTester } from '../helpers/auth';

const baseURL = () => process.env.STOKE_BASE_URL || 'http://localhost:8080';

test('create claim → group (with claim) → user (with group) then user can retrieve JWT @US-011', async ({
  page,
  request,
}) => {
  const suffix = Date.now();
  const claimName = `e2e-claim-${suffix}`;
  const groupName = `e2e-group-${suffix}`;
  const userName = `e2e-user-${suffix}`;
  const userPassword = 'E2eCreateFlow1!';

  await loginAsTester(page, baseURL());

  // 1. Create claim
  await page.goto('/admin/claim');
  await page.getByRole('button', { name: /add claim/i }).click();
  await page.getByLabel('Claim Name').fill(claimName);
  await page.getByLabel('Claim Short Name').fill('e2e');
  await page.getByLabel('Claim Value').fill('acc');
  await page.getByLabel('Claim Description').fill('E2E test claim');
  await page.getByRole('button', { name: 'Save' }).click();
  await expect(page.getByRole('button', { name: /add claim/i })).toBeVisible({ timeout: 5000 });

  // 2. Create group and assign claim to it
  await page.goto('/admin/group');
  await page.getByRole('button', { name: /add group/i }).click();
  await page.getByLabel('Name').fill(groupName);
  await page.getByLabel('Description').fill('E2E test group');
  await page.getByText(claimName).first().click();
  await page.getByRole('button', { name: 'Save' }).click();
  await expect(page.getByRole('button', { name: /add group/i })).toBeVisible({ timeout: 5000 });

  // 3. Create user
  await page.goto('/admin/user');
  await page.getByRole('button', { name: /add user/i }).click();
  await page.getByLabel('Username').fill(userName);
  await page.getByLabel('First Name').fill('E2E');
  await page.getByLabel('Last Name').fill('User');
  await page.getByLabel('Email').fill('e2e@test.example');
  await page.getByLabel('Password').first().fill(userPassword);
  await page.getByLabel('Repeat Password').fill(userPassword);
  await page.getByRole('button', { name: 'Save' }).click();
  await expect(page.getByRole('button', { name: /add user/i })).toBeVisible({ timeout: 5000 });

  // 4. Select the new user and assign group (Edit User → click group → Save)
  await page.getByText(userName).first().click();
  await page.getByRole('button', { name: /edit user/i }).click();
  await page.getByText(groupName).first().click();
  await page.getByRole('button', { name: 'Save' }).click();

  // 5. Assert user can retrieve JWT from login endpoint
  const loginRes = await request.post(baseURL() + '/api/login', {
    data: { username: userName, password: userPassword },
    headers: { 'Content-Type': 'application/json' },
  });
  expect(loginRes.status()).toBe(200);
  const body = await loginRes.json();
  expect(body.token).toBeDefined();
  expect(body.refresh).toBeDefined();
});
