const path = require("path");
const { test, expect } = require("@playwright/test");

test("avatar button appears only after choosing an image", async ({ page }) => {
  const suffix = `${Date.now()}-${Math.random().toString(16).slice(2)}`;
  const email = `avatar-e2e-${suffix}@example.com`;
  const password = "e2e-password";

  await page.goto("/user/signup");
  await page.locator('input[name="name"]').fill("Avatar E2E");
  await page.locator('input[name="email"]').fill(email);
  await page.locator('input[name="password"]').fill(password);
  await page.locator('input[type="submit"]').click();

  await expect(page).toHaveURL(/\/user\/login$/);
  await page.locator('input[name="email"]').fill(email);
  await page.locator('input[name="password"]').fill(password);
  await page.locator('input[type="submit"]').click();

  await expect(page).toHaveURL(/\/account\/view$/);

  const submitButton = page.locator("[data-avatar-submit]");
  await expect(submitButton).toBeHidden();
  await expect(submitButton).toBeDisabled();

  const avatarPath = path.resolve("ui/static/img/avatars/owl.png");
  await page.locator('input[name="avatar"]').setInputFiles(avatarPath);

  await expect(page.locator("[data-avatar-hint]")).toHaveText("Ready to upload.");
  await expect(submitButton).toBeVisible();
  await expect(submitButton).toBeEnabled();

  await submitButton.click();

  await expect(page).toHaveURL(/\/account\/view$/);
  await expect(page.getByText("Avatar updated successfully")).toBeVisible();
  await expect(page.locator(".account-avatar-image")).toHaveAttribute(
    "src",
    /^\/upload\/avatars\/user_\d+_\d+\.png$/,
  );
});
