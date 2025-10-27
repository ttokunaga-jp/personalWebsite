import { render, screen, waitFor } from "@testing-library/react";
import { vi } from "vitest";

import App from "./App";

const apiMocks = vi.hoisted(() => ({
  health: vi.fn(),
  fetchSummary: vi.fn(),
  listProjects: vi.fn(),
  listResearch: vi.fn(),
  listBlogs: vi.fn(),
  listMeetings: vi.fn(),
  listBlacklist: vi.fn(),
  createProject: vi.fn(),
  createResearch: vi.fn(),
  createBlog: vi.fn(),
  createMeeting: vi.fn(),
  createBlacklist: vi.fn(),
  updateProject: vi.fn(),
  updateResearch: vi.fn(),
  updateBlog: vi.fn(),
  updateMeeting: vi.fn(),
  deleteProject: vi.fn(),
  deleteResearch: vi.fn(),
  deleteBlog: vi.fn(),
  deleteMeeting: vi.fn(),
  deleteBlacklist: vi.fn()
}));

vi.mock("./modules/admin-api", () => ({
  adminApi: apiMocks
}));

const {
  health: healthMock,
  fetchSummary: summaryMock,
  listProjects: projectsMock,
  listResearch: researchMock,
  listBlogs: blogsMock,
  listMeetings: meetingsMock,
  listBlacklist: blacklistMock
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
        publishedBlogs: 3,
        draftBlogs: 0,
        pendingMeetings: 1,
        blacklistEntries: 1
      }
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
          updatedAt: new Date().toISOString()
        }
      ]
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
          updatedAt: new Date().toISOString()
        }
      ]
    });
    blogsMock.mockResolvedValue({
      data: [
        {
          id: 1,
          title: { ja: "ブログ", en: "Blog" },
          summary: { ja: "", en: "" },
          contentMd: { ja: "", en: "" },
          tags: [],
          published: true,
          publishedAt: new Date().toISOString(),
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString()
        }
      ]
    });
    meetingsMock.mockResolvedValue({
      data: [
        {
          id: 1,
          name: "Tester",
          email: "tester@example.com",
          datetime: new Date().toISOString(),
          durationMinutes: 30,
          meetUrl: "",
          status: "pending",
          notes: "",
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString()
        }
      ]
    });
    blacklistMock.mockResolvedValue({
      data: [
        {
          id: 1,
          email: "blocked@example.com",
          reason: "test",
          createdAt: new Date().toISOString()
        }
      ]
    });
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  it("renders dashboard title", async () => {
    render(<App />);

    await waitFor(() => {
      expect(healthMock).toHaveBeenCalled();
      expect(summaryMock).toHaveBeenCalled();
      expect(screen.getByRole("heading", { name: /Admin console/i })).toBeInTheDocument();
    });
  });

  it("blocks access when health check returns 401", async () => {
    healthMock.mockRejectedValueOnce({ response: { status: 401 } });

    render(<App />);

    await waitFor(() => {
      expect(screen.getByText(/Sign in required/i)).toBeInTheDocument();
    });

    expect(summaryMock).not.toHaveBeenCalled();
  });
});
