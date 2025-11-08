import { render, screen, waitFor } from "@testing-library/react";
import { vi } from "vitest";

import App from "./App";
import { DomainError } from "./modules/admin-api";
import { AuthSessionProvider } from "./modules/auth-session";

const apiMocks = vi.hoisted(() => ({
  health: vi.fn(),
  fetchSummary: vi.fn(),
  getProfile: vi.fn(),
  updateProfile: vi.fn(),
  listTechCatalog: vi.fn(),
  listProjects: vi.fn(),
  listResearch: vi.fn(),
  listContacts: vi.fn(),
  getHomeSettings: vi.fn(),
  updateHomeSettings: vi.fn(),
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
  updateProfile: updateProfileMock,
  listTechCatalog: techCatalogMock,
  listProjects: projectsMock,
  listResearch: researchMock,
  listContacts: contactsMock,
  getHomeSettings: homeSettingsMock,
  updateHomeSettings: updateHomeSettingsMock,
  getContactSettings: contactSettingsMock,
  updateContactSettings: updateContactSettingsMock,
  listBlacklist: blacklistMock,
  session: sessionMock,
} = apiMocks;

describe("Admin App", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    const now = new Date().toISOString();
    const profileResponse = {
      id: 1,
      displayName: "Takumi Takami",
      headline: { ja: "エンジニア", en: "Engineer" },
      summary: { ja: "概要", en: "Summary" },
      avatarUrl: null,
      location: { ja: "東京都", en: "Tokyo" },
      theme: { mode: "light", accentColor: "#00a0e9" },
      lab: {
        name: { ja: "未来創造研究室", en: "Future Lab" },
        advisor: { ja: "山田 教授", en: "Prof. Yamada" },
        room: { ja: "305", en: "305" },
        url: "https://example.com/lab",
      },
      affiliations: [
        {
          id: 1,
          kind: "affiliation" as const,
          name: "東京大学",
          url: "https://example.com",
          description: { ja: "所属説明", en: "Affiliation description" },
          startedAt: now,
          sortOrder: 1,
        },
      ],
      communities: [
        {
          id: 2,
          kind: "community" as const,
          name: "OSS Community",
          url: "https://community.example.com",
          description: { ja: "コミュニティ", en: "Community" },
          startedAt: now,
          sortOrder: 1,
        },
      ],
      workHistory: [
        {
          id: 1,
          organization: { ja: "株式会社テスト", en: "Test Inc." },
          role: { ja: "ソフトウェアエンジニア", en: "Software Engineer" },
          summary: { ja: "開発業務を担当。", en: "Handled development tasks." },
          startedAt: now,
          endedAt: null,
          externalUrl: "https://example.com/work",
          sortOrder: 1,
        },
      ],
      socialLinks: [
        {
          id: 1,
          provider: "github" as const,
          label: { ja: "GitHub", en: "GitHub" },
          url: "https://github.com/example",
          isFooter: true,
          sortOrder: 1,
        },
        {
          id: 2,
          provider: "zenn" as const,
          label: { ja: "Zenn", en: "Zenn" },
          url: "https://zenn.dev/example",
          isFooter: true,
          sortOrder: 2,
        },
        {
          id: 3,
          provider: "linkedin" as const,
          label: { ja: "LinkedIn", en: "LinkedIn" },
          url: "https://linkedin.com/in/example",
          isFooter: true,
          sortOrder: 3,
        },
      ],
      techSections: [],
      home: null,
      updatedAt: now,
    };

    const homeSettingsResponse = {
      id: 1,
      profileId: profileResponse.id,
      heroSubtitle: { ja: "研究と実践の融合", en: "Crafting research and practice" },
      quickLinks: [
        {
          id: 1,
          section: "profile" as const,
          label: { ja: "プロフィール", en: "Profile" },
          description: { ja: "略歴を見る", en: "Read biography" },
          cta: { ja: "詳細", en: "Details" },
          targetUrl: "/profile",
          sortOrder: 1,
        },
      ],
      chipSources: [
        {
          id: 1,
          source: "affiliation" as const,
          label: { ja: "所属", en: "Affiliations" },
          limit: 3,
          sortOrder: 1,
        },
      ],
      updatedAt: now,
    };

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
        profileUpdatedAt: now,
      },
    });
    profileMock.mockResolvedValue({
      data: profileResponse,
    });
    updateProfileMock.mockResolvedValue({
      data: profileResponse,
    });
    homeSettingsMock.mockResolvedValue({ data: homeSettingsResponse });
    updateHomeSettingsMock.mockResolvedValue({ data: homeSettingsResponse });
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
      createdAt: now,
      updatedAt: now,
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
          createdAt: now,
          updatedAt: now,
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
          publishedAt: now,
          isDraft: false,
          createdAt: now,
          updatedAt: now,
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
          createdAt: now,
          updatedAt: now,
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
