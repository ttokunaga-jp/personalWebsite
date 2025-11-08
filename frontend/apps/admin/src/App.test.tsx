import { render, screen, waitFor } from "@testing-library/react";
import { vi } from "vitest";

import App from "./App";
import { DomainError } from "./modules/admin-api";
import { AuthSessionProvider } from "./modules/auth-session";

const apiMocks = vi.hoisted(() => ({
  health: vi.fn(),
  fetchSummary: vi.fn(),
  getProfile: vi.fn(),
  listTechCatalog: vi.fn(),
  listProjects: vi.fn(),
  listResearch: vi.fn(),
  listContacts: vi.fn(),
  getContactSettings: vi.fn(),
  updateContactSettings: vi.fn(),
  listBlacklist: vi.fn(),
  session: vi.fn(),
}));

vi.mock("./modules/admin-api", async (importOriginal) => {
  const actual = await importOriginal<typeof import("./modules/admin-api")>();
  return {
    ...actual,
    adminApi: apiMocks,
  };
});

const {
  health: healthMock,
  fetchSummary: summaryMock,
  getProfile: profileMock,
  listTechCatalog: techCatalogMock,
  listProjects: projectsMock,
  listResearch: researchMock,
  listContacts: contactsMock,
  getContactSettings: contactSettingsMock,
  updateContactSettings: updateContactSettingsMock,
  listBlacklist: blacklistMock,
  session: sessionMock,
} = apiMocks;

describe("Admin App", () => {
  beforeEach(() => {
    healthMock.mockResolvedValue({ data: { status: "ok" } });
    summaryMock.mockResolvedValue({
      data: {
        publishedProjects: 1,
        draftProjects: 0,
        publishedResearch: 2,
        draftResearch: 1,
        pendingContacts: 1,
        blacklistEntries: 1,
        skillCount: 4,
        focusAreaCount: 2,
        profileUpdatedAt: new Date().toISOString(),
      },
    });
    profileMock.mockResolvedValue({
      data: {
        name: { ja: "高見 拓実", en: "Takumi Takami" },
        title: { ja: "エンジニア", en: "Engineer" },
        affiliation: { ja: "", en: "" },
        lab: { ja: "", en: "" },
        summary: { ja: "概要", en: "Summary" },
        skills: [],
        focusAreas: [],
        updatedAt: new Date().toISOString(),
      },
    });
    const contactSettingsResponse = {
      id: 1,
      heroTitle: { ja: "お問い合わせ", en: "Contact" },
      heroDescription: {
        ja: "研究や講演のご相談を受け付けています。",
        en: "Reach out for collaborations or speaking engagements.",
      },
      topics: [
        {
          id: "general",
          label: { ja: "一般", en: "General" },
          description: {
            ja: "一般的なお問い合わせはこちら。",
            en: "General inquiries.",
          },
        },
      ],
      consentText: {
        ja: "送信によりプライバシーポリシーに同意したものとみなします。",
        en: "By submitting you agree to the privacy policy.",
      },
      minimumLeadHours: 24,
      recaptchaSiteKey: "site-key",
      supportEmail: "support@example.com",
      calendarTimezone: "Asia/Tokyo",
      googleCalendarId: "primary",
      bookingWindowDays: 30,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    };
    contactSettingsMock.mockResolvedValue({ data: contactSettingsResponse });
    updateContactSettingsMock.mockResolvedValue({
      data: contactSettingsResponse,
    });
    techCatalogMock.mockResolvedValue({
      data: [
        {
          id: 1,
          slug: "go",
          displayName: "Go",
          category: "language",
          level: "advanced",
          icon: "",
          sortOrder: 1,
          active: true,
        },
        {
          id: 2,
          slug: "react",
          displayName: "React",
          category: "frontend",
          level: "advanced",
          icon: "",
          sortOrder: 2,
          active: true,
        },
      ],
    });
    projectsMock.mockResolvedValue({
      data: [
        {
          id: 1,
          title: { ja: "タイトル", en: "Title" },
          description: { ja: "説明", en: "Description" },
          tech: [
            {
              membershipId: 1,
              entityType: "project",
              entityId: 1,
              tech: {
                id: 1,
                slug: "go",
                displayName: "Go",
                category: "language",
                level: "advanced",
                icon: "",
                sortOrder: 1,
                active: true,
              },
              context: "primary",
              note: "",
              sortOrder: 1,
            },
          ],
          linkUrl: "",
          year: 2024,
          published: true,
          sortOrder: null,
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
        },
      ],
    });
    researchMock.mockResolvedValue({
      data: [
        {
          id: 1,
          slug: "demo-research",
          kind: "research",
          title: { ja: "研究", en: "Research" },
          overview: { ja: "概要", en: "Overview" },
          outcome: { ja: "", en: "" },
          outlook: { ja: "", en: "" },
          externalUrl: "https://example.com",
          highlightImageUrl: "",
          imageAlt: { ja: "", en: "" },
          publishedAt: new Date().toISOString(),
          isDraft: false,
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
          tags: [],
          links: [],
          assets: [],
          tech: [],
        },
      ],
    });
    contactsMock.mockResolvedValue({
      data: [
        {
          id: "contact-1",
          name: "Tester",
          email: "tester@example.com",
          topic: "相談",
          message: "テストメッセージ",
          status: "pending",
          adminNote: "",
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
        },
      ],
    });
    blacklistMock.mockResolvedValue({
      data: [
        {
          id: 1,
          email: "blocked@example.com",
          reason: "test",
          createdAt: new Date().toISOString(),
        },
      ],
    });
    sessionMock.mockResolvedValue({ data: { active: true } });
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  it("renders dashboard title", async () => {
    render(
      <AuthSessionProvider>
        <App />
      </AuthSessionProvider>,
    );

    await waitFor(() => {
      expect(healthMock).toHaveBeenCalled();
      expect(summaryMock).toHaveBeenCalled();
      const heading = screen.getByRole("heading", {
        name: /Admin console|管理コンソール/i,
      });
      expect(heading).toBeInTheDocument();
    });
  });

  it("blocks access when health check returns 401", async () => {
    healthMock.mockRejectedValueOnce(new DomainError(401, "unauthorized"));
    sessionMock.mockResolvedValueOnce({ data: { active: false } });

    render(
      <AuthSessionProvider>
        <App />
      </AuthSessionProvider>,
    );

    await waitFor(() => {
      expect(screen.getByText(/Sign in required/i)).toBeInTheDocument();
    });

    expect(summaryMock).not.toHaveBeenCalled();
  });

});
