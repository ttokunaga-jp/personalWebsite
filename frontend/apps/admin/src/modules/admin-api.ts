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

type AdminApi = {
  health: () => ApiClientPromise<{ status: string }>;
  fetchSummary: () => ApiClientPromise<AdminSummary>;
  listProjects: () => ApiClientPromise<AdminProject[]>;
  createProject: (payload: ProjectPayload) => ApiClientPromise<AdminProject>;
  updateProject: (
    id: number,
    payload: ProjectPayload,
  ) => ApiClientPromise<AdminProject>;
  deleteProject: (id: number) => ApiClientPromise<void>;

  listResearch: () => ApiClientPromise<AdminResearch[]>;
  createResearch: (payload: ResearchPayload) => ApiClientPromise<AdminResearch>;
  updateResearch: (
    id: number,
    payload: ResearchPayload,
  ) => ApiClientPromise<AdminResearch>;
  deleteResearch: (id: number) => ApiClientPromise<void>;

  listBlogs: () => ApiClientPromise<BlogPost[]>;
  createBlog: (payload: BlogPayload) => ApiClientPromise<BlogPost>;
  updateBlog: (id: number, payload: BlogPayload) => ApiClientPromise<BlogPost>;
  deleteBlog: (id: number) => ApiClientPromise<void>;

  listMeetings: () => ApiClientPromise<Meeting[]>;
  createMeeting: (payload: MeetingPayload) => ApiClientPromise<Meeting>;
  updateMeeting: (
    id: number,
    payload: MeetingPayload,
  ) => ApiClientPromise<Meeting>;
  deleteMeeting: (id: number) => ApiClientPromise<void>;

  listBlacklist: () => ApiClientPromise<BlacklistEntry[]>;
  createBlacklist: (
    payload: BlacklistPayload,
  ) => ApiClientPromise<BlacklistEntry>;
  deleteBlacklist: (id: number) => ApiClientPromise<void>;
};

export const adminApi: AdminApi = {
  health: () => apiClient.get<{ status: string }>("/admin/health"),
  fetchSummary: () => apiClient.get<AdminSummary>("/admin/summary"),
  listProjects: () => apiClient.get<AdminProject[]>("/admin/projects"),
  createProject: (payload: ProjectPayload) =>
    apiClient.post<AdminProject>("/admin/projects", payload),
  updateProject: (id: number, payload: ProjectPayload) =>
    apiClient.put<AdminProject>(`/admin/projects/${id}`, payload),
  deleteProject: (id: number) =>
    apiClient.delete<void>(`/admin/projects/${id}`),

  listResearch: () => apiClient.get<AdminResearch[]>("/admin/research"),
  createResearch: (payload: ResearchPayload) =>
    apiClient.post<AdminResearch>("/admin/research", payload),
  updateResearch: (id: number, payload: ResearchPayload) =>
    apiClient.put<AdminResearch>(`/admin/research/${id}`, payload),
  deleteResearch: (id: number) =>
    apiClient.delete<void>(`/admin/research/${id}`),

  listBlogs: () => apiClient.get<BlogPost[]>("/admin/blogs"),
  createBlog: (payload: BlogPayload) =>
    apiClient.post<BlogPost>("/admin/blogs", payload),
  updateBlog: (id: number, payload: BlogPayload) =>
    apiClient.put<BlogPost>(`/admin/blogs/${id}`, payload),
  deleteBlog: (id: number) => apiClient.delete<void>(`/admin/blogs/${id}`),

  listMeetings: () => apiClient.get<Meeting[]>("/admin/meetings"),
  createMeeting: (payload: MeetingPayload) =>
    apiClient.post<Meeting>("/admin/meetings", payload),
  updateMeeting: (id: number, payload: MeetingPayload) =>
    apiClient.put<Meeting>(`/admin/meetings/${id}`, payload),
  deleteMeeting: (id: number) =>
    apiClient.delete<void>(`/admin/meetings/${id}`),

  listBlacklist: () => apiClient.get<BlacklistEntry[]>("/admin/blacklist"),
  createBlacklist: (payload: BlacklistPayload) =>
    apiClient.post<BlacklistEntry>("/admin/blacklist", payload),
  deleteBlacklist: (id: number) =>
    apiClient.delete<void>(`/admin/blacklist/${id}`),
};
