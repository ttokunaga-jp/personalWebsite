import { apiClient } from "@shared/lib/api-client";
import type { ApiClientPromise } from "@shared/lib/api-client";

import type {
  AdminProject,
  AdminResearch,
  AdminSummary,
  BlogPost,
  BlacklistEntry,
  Meeting,
} from "../types";

type ProjectPayload = {
  title: { ja?: string; en?: string };
  description: { ja?: string; en?: string };
  techStack: string[];
  linkUrl: string;
  year: number;
  published: boolean;
  sortOrder?: number | null;
};

type ResearchPayload = {
  title: { ja?: string; en?: string };
  summary: { ja?: string; en?: string };
  contentMd: { ja?: string; en?: string };
  year: number;
  published: boolean;
};

type BlogPayload = {
  title: { ja?: string; en?: string };
  summary: { ja?: string; en?: string };
  contentMd: { ja?: string; en?: string };
  tags: string[];
  published: boolean;
  publishedAt?: string | null;
};

type MeetingPayload = {
  name: string;
  email: string;
  datetime: string;
  durationMinutes: number;
  meetUrl: string;
  status: string;
  notes: string;
};

type BlacklistPayload = {
  email: string;
  reason: string;
};

export class DomainError extends Error {
  readonly status: number;

  constructor(status: number, message = "domain error") {
    super(message);
    this.name = "DomainError";
    this.status = status;
  }
}

type ApiResult<T> = Promise<Awaited<ApiClientPromise<T>>>;

async function unwrap<T>(promise: ApiClientPromise<T>): ApiResult<T> {
  try {
    return await promise;
  } catch (error) {
    const status = (error as { response?: { status?: number } })?.response
      ?.status;
    if (status && status === 401) {
      throw new DomainError(401, "unauthorized");
    }
    throw error;
  }
}

type AdminApi = {
  health: () => ApiResult<{ status: string }>;
  fetchSummary: () => ApiResult<AdminSummary>;
  listProjects: () => ApiResult<AdminProject[]>;
  createProject: (payload: ProjectPayload) => ApiResult<AdminProject>;
  updateProject: (
    id: number,
    payload: ProjectPayload,
  ) => ApiResult<AdminProject>;
  deleteProject: (id: number) => ApiResult<void>;

  listResearch: () => ApiResult<AdminResearch[]>;
  createResearch: (payload: ResearchPayload) => ApiResult<AdminResearch>;
  updateResearch: (
    id: number,
    payload: ResearchPayload,
  ) => ApiResult<AdminResearch>;
  deleteResearch: (id: number) => ApiResult<void>;

  listBlogs: () => ApiResult<BlogPost[]>;
  createBlog: (payload: BlogPayload) => ApiResult<BlogPost>;
  updateBlog: (id: number, payload: BlogPayload) => ApiResult<BlogPost>;
  deleteBlog: (id: number) => ApiResult<void>;

  listMeetings: () => ApiResult<Meeting[]>;
  createMeeting: (payload: MeetingPayload) => ApiResult<Meeting>;
  updateMeeting: (
    id: number,
    payload: MeetingPayload,
  ) => ApiResult<Meeting>;
  deleteMeeting: (id: number) => ApiResult<void>;

  listBlacklist: () => ApiResult<BlacklistEntry[]>;
  createBlacklist: (
    payload: BlacklistPayload,
  ) => ApiResult<BlacklistEntry>;
  deleteBlacklist: (id: number) => ApiResult<void>;
  session: () => ApiResult<{
    active: boolean;
    token?: string;
    expiresAt?: number;
  }>;
};

export const adminApi: AdminApi = {
  health: () => unwrap(apiClient.get<{ status: string }>("/admin/health")),
  fetchSummary: () => unwrap(apiClient.get<AdminSummary>("/admin/summary")),
  listProjects: () => unwrap(apiClient.get<AdminProject[]>("/admin/projects")),
  createProject: (payload: ProjectPayload) =>
    unwrap(apiClient.post<AdminProject>("/admin/projects", payload)),
  updateProject: (id: number, payload: ProjectPayload) =>
    unwrap(apiClient.put<AdminProject>(`/admin/projects/${id}`, payload)),
  deleteProject: (id: number) =>
    unwrap(apiClient.delete<void>(`/admin/projects/${id}`)),

  listResearch: () => unwrap(apiClient.get<AdminResearch[]>("/admin/research")),
  createResearch: (payload: ResearchPayload) =>
    unwrap(apiClient.post<AdminResearch>("/admin/research", payload)),
  updateResearch: (id: number, payload: ResearchPayload) =>
    unwrap(apiClient.put<AdminResearch>(`/admin/research/${id}`, payload)),
  deleteResearch: (id: number) =>
    unwrap(apiClient.delete<void>(`/admin/research/${id}`)),

  listBlogs: () => unwrap(apiClient.get<BlogPost[]>("/admin/blogs")),
  createBlog: (payload: BlogPayload) =>
    unwrap(apiClient.post<BlogPost>("/admin/blogs", payload)),
  updateBlog: (id: number, payload: BlogPayload) =>
    unwrap(apiClient.put<BlogPost>(`/admin/blogs/${id}`, payload)),
  deleteBlog: (id: number) =>
    unwrap(apiClient.delete<void>(`/admin/blogs/${id}`)),

  listMeetings: () => unwrap(apiClient.get<Meeting[]>("/admin/meetings")),
  createMeeting: (payload: MeetingPayload) =>
    unwrap(apiClient.post<Meeting>("/admin/meetings", payload)),
  updateMeeting: (id: number, payload: MeetingPayload) =>
    unwrap(apiClient.put<Meeting>(`/admin/meetings/${id}`, payload)),
  deleteMeeting: (id: number) =>
    unwrap(apiClient.delete<void>(`/admin/meetings/${id}`)),

  listBlacklist: () => unwrap(apiClient.get<BlacklistEntry[]>("/admin/blacklist")),
  createBlacklist: (payload: BlacklistPayload) =>
    unwrap(apiClient.post<BlacklistEntry>("/admin/blacklist", payload)),
  deleteBlacklist: (id: number) =>
    unwrap(apiClient.delete<void>(`/admin/blacklist/${id}`)),
  session: () =>
    unwrap(
      apiClient.get<{
        active: boolean;
        token?: string;
        expiresAt?: number;
      }>("/admin/auth/session"),
    ),
};
