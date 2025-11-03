import { apiClient } from "@shared/lib/api-client";
import type { ApiClientPromise } from "@shared/lib/api-client";

import type {
  AdminProfile,
  AdminProject,
  AdminResearch,
  AdminSummary,
  BlacklistEntry,
  ContactMessage,
} from "../types";

type ProfilePayload = {
  name: { ja?: string; en?: string };
  title: { ja?: string; en?: string };
  affiliation: { ja?: string; en?: string };
  lab: { ja?: string; en?: string };
  summary: { ja?: string; en?: string };
  skills: { ja?: string; en?: string }[];
  focusAreas: { ja?: string; en?: string }[];
};

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

type ContactUpdatePayload = {
  topic: string;
  message: string;
  status: string;
  adminNote: string;
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
  getProfile: () => ApiResult<AdminProfile>;
  updateProfile: (payload: ProfilePayload) => ApiResult<AdminProfile>;
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

  listContacts: () => ApiResult<ContactMessage[]>;
  getContact: (id: string) => ApiResult<ContactMessage>;
  updateContact: (
    id: string,
    payload: ContactUpdatePayload,
  ) => ApiResult<ContactMessage>;
  deleteContact: (id: string) => ApiResult<void>;

  listBlacklist: () => ApiResult<BlacklistEntry[]>;
  createBlacklist: (
    payload: BlacklistPayload,
  ) => ApiResult<BlacklistEntry>;
  updateBlacklist: (
    id: number,
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
  getProfile: () => unwrap(apiClient.get<AdminProfile>("/admin/profile")),
  updateProfile: (payload: ProfilePayload) =>
    unwrap(apiClient.put<AdminProfile>("/admin/profile", payload)),
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

  listContacts: () => unwrap(apiClient.get<ContactMessage[]>("/admin/contacts")),
  getContact: (id: string) =>
    unwrap(apiClient.get<ContactMessage>(`/admin/contacts/${id}`)),
  updateContact: (id: string, payload: ContactUpdatePayload) =>
    unwrap(apiClient.put<ContactMessage>(`/admin/contacts/${id}`, payload)),
  deleteContact: (id: string) =>
    unwrap(apiClient.delete<void>(`/admin/contacts/${id}`)),

  listBlacklist: () => unwrap(apiClient.get<BlacklistEntry[]>("/admin/blacklist")),
  createBlacklist: (payload: BlacklistPayload) =>
    unwrap(apiClient.post<BlacklistEntry>("/admin/blacklist", payload)),
  updateBlacklist: (id: number, payload: BlacklistPayload) =>
    unwrap(apiClient.put<BlacklistEntry>(`/admin/blacklist/${id}`, payload)),
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
