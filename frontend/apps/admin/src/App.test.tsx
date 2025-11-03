import { render, screen, waitFor } from "@testing-library/react";
import { vi } from "vitest";

import App from "./App";
import { DomainError } from "./modules/admin-api";
import {
  AuthSessionProvider,
  clearToken,
  setToken as persistToken,
} from "./modules/auth-session";

const apiMocks = vi.hoisted(() => ({
  health: vi.fn(),
  fetchSummary: vi.fn(),
  getProfile: vi.fn(),
  listProjects: vi.fn(),
  listResearch: vi.fn(),
  listContacts: vi.fn(),
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
  listProjects: projectsMock,
  listResearch: researchMock,
  listContacts: contactsMock,
  listBlacklist: blacklistMock,
  session: sessionMock,
} = apiMocks;

describe("Admin App", () => {
  beforeEach(() => {
    clearToken();
    window.sessionStorage.clear();
    window.location.hash = "";
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
    projectsMock.mockResolvedValue({
      data: [
        {
          id: 1,
          title: { ja: "タイトル", en: "Title" },
          description: { ja: "説明", en: "Description" },
          techStack: ["Go"],
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
          title: { ja: "研究", en: "Research" },
          summary: { ja: "", en: "" },
          contentMd: { ja: "", en: "" },
          year: 2023,
          published: true,
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
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
    persistToken("test-token");

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
