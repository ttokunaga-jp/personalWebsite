import { expect, test } from "@playwright/test";

import {
  contactAvailabilityFixture,
  contactConfigFixture,
  defaultBookingResponse,
  profileFixture,
  projectsFixture,
  researchEntriesFixture
} from "../src/test-utils/fixtures";

test.beforeEach(async ({ page }) => {
  await page.addInitScript(() => {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    (window as any).grecaptcha = {
      ready: (callback: () => void) => callback(),
      execute: async () => "playwright-recaptcha-token"
    };
  });

  await page.route("**/api/health", async (route) => {
    await route.fulfill({ json: { status: "ok" } });
  });

  await page.route("**/api/v1/public/profile", async (route) => {
    await route.fulfill({ json: { data: profileFixture } });
  });

  await page.route("**/api/v1/public/research", async (route) => {
    await route.fulfill({ json: { data: researchEntriesFixture } });
  });

  await page.route("**/api/v1/public/projects", async (route) => {
    await route.fulfill({ json: { data: projectsFixture } });
  });

  await page.route("**/api/v1/public/contact/availability", async (route) => {
    await route.fulfill({ json: { data: contactAvailabilityFixture } });
  });

  await page.route("**/api/v1/public/contact/config", async (route) => {
    await route.fulfill({ json: { data: contactConfigFixture } });
  });

  await page.route(
    "**/api/v1/public/contact/bookings",
    async (route, request) => {
      const body = await request.postDataJSON();
      const meetingId = `bk-${body?.slotId ?? "unknown"}`;
      await route.fulfill({
        json: {
          data: {
            ...defaultBookingResponse,
            reservation: {
              ...defaultBookingResponse.reservation,
              id: meetingId
            }
          }
        }
      });
    }
  );
});

test("visitor walks through primary navigation and submits a booking", async ({ page }) => {
  await page.goto("/");

  await expect(
    page.getByRole("heading", {
      name: profileFixture.displayName,
      exact: false
    })
  ).toBeVisible();

  const primaryNavigation = page.getByLabel("Primary Navigation");

  await primaryNavigation.getByRole("link", { name: "Projects" }).click();
  await expect(page.getByRole("heading", { name: "Project archive" })).toBeVisible();

  await page.getByRole("button", { name: "TypeScript" }).click();
  const firstProjectCard = page
    .getByRole("article")
    .filter({ hasText: projectsFixture[0]?.title ?? "" })
    .first();
  const secondProjectCard = page
    .getByRole("article")
    .filter({ hasText: projectsFixture[1]?.title ?? "" })
    .first();
  await expect(firstProjectCard).toBeVisible();
  await expect(secondProjectCard).not.toBeVisible();

  await primaryNavigation.getByRole("link", { name: "Contact" }).click();
  await expect(
    page.getByRole("heading", { name: contactConfigFixture.heroTitle }),
  ).toBeVisible();

  const nameField = page.locator('input[name="name"]');
  const emailField = page.locator('input[name="email"]');
  const topicField = page.locator('select[name="topic"]');
  const messageField = page.locator('textarea[name="agenda"]');
  const firstAvailableSlot = page.getByTestId("availability-slot-available").first();
  const submitButton = page.getByRole("button", { name: "Request booking" });

  await Promise.all([
    expect(nameField).toBeVisible(),
    expect(emailField).toBeVisible(),
    expect(topicField).toBeVisible(),
    expect(messageField).toBeVisible(),
    expect(firstAvailableSlot).toBeVisible(),
  ]);

  await nameField.fill("E2E Tester");
  await emailField.fill("tester@example.com");
  await topicField.selectOption(contactConfigFixture.topics[0]?.id ?? "");
  await messageField.fill(
    "Exploring possibilities for a robotics research collaboration.",
  );

  await firstAvailableSlot.click();
  await expect(firstAvailableSlot).toHaveClass(/bg-sky-600/);

  await expect(submitButton).toBeEnabled();
  await submitButton.click();

  await expect(page.getByText(/Your request \(ID: bk-/)).toBeVisible();
});
