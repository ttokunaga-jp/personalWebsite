import { expect, test } from "@playwright/test";

import {
  contactAvailabilityFixture,
  contactConfigFixture,
  profileFixture,
  projectsFixture,
  researchEntriesFixture,
} from "../src/test-utils/fixtures";

type MeetingUpdateCapture = {
  headers: Record<string, string>;
  body: unknown;
};

const ADMIN_EMAIL = "admin@example.com";

const meetingTemplateInitial = {
  template: "Initial meeting template with {{meeting_url}}",
  updatedAt: new Date().toISOString(),
};

let lastMeetingUpdate: MeetingUpdateCapture | null;

test.beforeEach(async ({ page }) => {
  lastMeetingUpdate = null;

  // Simulate reCAPTCHA availability (shared with public tests)
  await page.addInitScript(() => {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    (window as any).grecaptcha = {
      ready: (callback: () => void) => callback(),
      execute: async () => "playwright-recaptcha-token",
    };
  });

  // Public API routes
  await page.route("**/api/health", async (route) => {
    await route.fulfill({ json: { status: "ok" } });
  });
  await page.route("**/api/v1/public/profile", async (route) => {
    await route.fulfill({ json: { data: profileFixture } });
  });
  await page.route("**/api/v1/public/projects", async (route) => {
    await route.fulfill({ json: { data: projectsFixture } });
  });
  await page.route("**/api/v1/public/research", async (route) => {
    await route.fulfill({ json: { data: researchEntriesFixture } });
  });
  await page.route("**/api/v1/public/contact/availability", async (route) => {
    await route.fulfill({ json: { data: contactAvailabilityFixture } });
  });
  await page.route("**/api/v1/public/contact/config", async (route) => {
    await route.fulfill({ json: { data: contactConfigFixture } });
  });

  // CSRF token endpoint
  await page.route("**/api/security/csrf", async (route) => {
    await route.fulfill({
      headers: {
        "set-cookie": "ps_csrf=csrf-token; Path=/; HttpOnly",
      },
      json: {
        data: {
          token: "csrf-token",
          expires_at: new Date(Date.now() + 10 * 60_000).toISOString(),
        },
      },
    });
  });

  // Admin session and data stubs
  await page.route("**/api/admin/auth/session", async (route) => {
    await route.fulfill({
      json: {
        active: true,
        email: ADMIN_EMAIL,
        roles: ["admin"],
      },
    });
  });

  const homeConfig = {
    id: 1,
    heroSubtitle: { ja: "管理用サブタイトル", en: "Admin subtitle" },
    quickLinks: [
      {
        id: 10,
        section: "projects",
        label: { ja: "プロジェクト", en: "Projects" },
        description: { ja: "最新プロジェクト", en: "Latest projects" },
        cta: { ja: "見る", en: "View" },
        targetUrl: "/projects",
        sortOrder: 1,
      },
    ],
    chipSources: [
      {
        id: 20,
        source: "tech",
        label: { ja: "技術スタック", en: "Tech stack" },
        limit: 4,
        sortOrder: 1,
      },
    ],
    updatedAt: new Date().toISOString(),
  };

  await page.route("**/api/admin/home?**", async (route) => {
    if (route.request().method() === "GET") {
      await route.fulfill({ json: { data: homeConfig } });
      return;
    }
    // Allow other operations (e.g., PUT) to succeed with a simple echo
    const body = await route.request().postDataJSON();
    await route.fulfill({ json: { data: body } });
  });

  await page.route("**/api/admin/reservations?**", async (route) => {
    await route.fulfill({
      json: {
        data: [
          {
            id: 1,
            name: "Alice Example",
            email: "alice@example.com",
            topic: "Research call",
            message: "Discuss collaboration opportunities.",
            startAt: new Date(Date.now() + 3600_000).toISOString(),
            endAt: new Date(Date.now() + 5400_000).toISOString(),
            status: "pending",
            googleCalendarStatus: "tentative",
            lookupHash: "lookup-1",
            cancellationReason: "",
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
            notifications: [],
          },
        ],
      },
    });
  });

  await page.route("**/api/admin/blacklist?**", async (route) => {
    await route.fulfill({
      json: {
        data: [
          {
            id: 1,
            email: "blocked@example.com",
            reason: "Spam submissions",
            createdAt: new Date().toISOString(),
          },
        ],
      },
    });
  });

  await page.route("**/api/admin/tech-catalog?**", async (route) => {
    await route.fulfill({
      json: {
        data: [
          {
            id: 1,
            slug: "nextjs",
            displayName: "Next.js",
            category: "Framework",
            level: "advanced",
            icon: "⚡",
            sortOrder: 1,
            active: true,
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
          },
        ],
      },
    });
  });

  await page.route("**/api/admin/social-links?**", async (route) => {
    await route.fulfill({
      json: {
        data: [
          {
            id: 1,
            provider: "github",
            label: { ja: "GitHub", en: "GitHub" },
            url: "https://github.com/example",
            isFooter: true,
            sortOrder: 1,
          },
        ],
      },
    });
  });

  await page.route("**/api/admin/meeting-url?**", async (route) => {
    const request = route.request();
    if (request.method() === "GET") {
      await route.fulfill({ json: { data: meetingTemplateInitial } });
      return;
    }

    if (request.method() === "PUT") {
      const headers = request.headers();
      const body = await request.postDataJSON();
      lastMeetingUpdate = {
        headers,
        body,
      };
      await route.fulfill({
        json: {
          data: {
            template: body.template,
            updatedAt: new Date().toISOString(),
          },
        },
      });
      return;
    }

    await route.continue();
  });
});

test("toggles admin mode and preserves query parameters across navigation", async ({ page }) => {
  await page.goto("/");

  const toggle = page.getByRole("button", { name: /View mode/i });
  await expect(toggle).toBeVisible();

  await toggle.click();
  await expect(page).toHaveURL(/mode=admin/);

  // Ensure home editor is rendered
  await expect(page.locator("#hero-subtitle-ja")).toBeVisible();

  // Navigate via primary navigation; mode should persist
  await page
    .getByLabel("Primary Navigation")
    .getByRole("link", { name: "Projects" })
    .click();

  await expect(page).toHaveURL(/\/projects\?mode=admin/);
});

test("shows unsaved changes warning before navigation", async ({ page }) => {
  await page.goto("/?mode=admin");

  const subtitleInput = page.locator("#hero-subtitle-ja");
  await subtitleInput.fill("管理モードのテスト編集");
  await expect(subtitleInput).toHaveValue("管理モードのテスト編集");
  await page.waitForTimeout(100);

  const dialogPromise = page.waitForEvent("dialog");
  const clickPromise = page
    .getByLabel("Primary Navigation")
    .getByRole("link", { name: "Projects" })
    .click({ noWaitAfter: true });

  const dialog = await dialogPromise;
  await expect(dialog.message()).toContain("unsaved changes");
  await dialog.accept();
  await clickPromise;

  await expect(page).toHaveURL(/\/projects\?mode=admin/);
});

test("updates meeting URL template and sends CSRF headers", async ({ page }) => {
  await page.goto("/admin?mode=admin");

  await page.getByRole("button", { name: "Meeting URL" }).click();

  const templateField = page.getByLabel("Message body");
  await expect(templateField).toHaveValue(meetingTemplateInitial.template);

  const updatedTemplate = "Updated template for {{guest_name}} at {{meeting_url}}";
  await templateField.fill(updatedTemplate);

  await page.getByRole("button", { name: "Save template" }).click();

  await expect.poll(() => lastMeetingUpdate).not.toBeNull();
  await expect(lastMeetingUpdate!.headers["x-requested-with"]).toBe("XMLHttpRequest");
  await expect(lastMeetingUpdate!.headers["x-csrf-token"]).toBe("csrf-token");
  await expect(lastMeetingUpdate!.body).toEqual({ template: updatedTemplate });

  await expect(page.getByRole("button", { name: "Save template" })).toBeDisabled();
});
