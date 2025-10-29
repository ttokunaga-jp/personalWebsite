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
      await route.fulfill({
        json: {
          data: {
            ...defaultBookingResponse,
            bookingId: `bk-${body?.slotId ?? "unknown"}`
          }
        }
      });
    }
  );
});

test("visitor walks through primary navigation and submits a booking", async ({ page }) => {
  await page.goto("/");

  await expect(
    page.getByRole("heading", { name: profileFixture.headline })
  ).toBeVisible();

  const primaryNavigation = page.getByLabel("Primary Navigation");

  await primaryNavigation.getByRole("link", { name: "Projects" }).click();
  await expect(page.getByRole("heading", { name: "Project archive" })).toBeVisible();

  await page.getByRole("button", { name: "TypeScript" }).click();
  await expect(page.getByRole("article", { name: projectsFixture[0]?.title ?? "" })).toBeVisible();
  await expect(
    page.getByRole("article", { name: projectsFixture[1]?.title ?? "" })
  ).not.toBeVisible();

  await primaryNavigation.getByRole("link", { name: "Contact" }).click();
  await expect(page.getByRole("heading", { name: "Get in touch" })).toBeVisible();

  await page.fill('input[name="name"]', "E2E Tester");
  await page.fill('input[name="email"]', "tester@example.com");
  await page.selectOption('select[name="topic"]', contactConfigFixture.topics[0] ?? "");
  await page.fill(
    'textarea[name="message"]',
    "Exploring possibilities for a robotics research collaboration."
  );

  const slotButton = page.getByRole("button", { name: /Ends/ }).first();
  await slotButton.click();

  await page.getByRole("button", { name: "Request booking" }).click();

  await expect(
    page.getByText(/Your request \(ID: bk-/)
  ).toBeVisible();
});
