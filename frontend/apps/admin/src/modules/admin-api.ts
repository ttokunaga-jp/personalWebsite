import { apiClient } from "@shared/lib/api-client";
import type { ApiClientPromise } from "@shared/lib/api-client";

import type {
  AdminProfile,
  AdminProject,
  AdminResearch,
  AdminSummary,
  BlacklistEntry,
  ContactFormSettings,
  ContactMessage,
  HomePageConfig,
  LocalizedText,
  ResearchKind,
  ResearchLinkType,
  SocialProvider,
  TechCatalogEntry,
  TechContext,
} from "../types";

type ProfilePayload = {
  displayName: string;
  headline: LocalizedText;
  summary: LocalizedText;
  avatarUrl: string;
  location: LocalizedText;
  theme: {
    mode: "light" | "dark" | "system";
    accentColor?: string;
  };
  lab: {
    name: LocalizedText;
    advisor: LocalizedText;
    room: LocalizedText;
    url?: string;
  };
  affiliations: {
    id?: number;
    name: string;
    url?: string;
    description: LocalizedText;
    startedAt: string;
    sortOrder: number;
  }[];
  communities: {
    id?: number;
    name: string;
    url?: string;
    description: LocalizedText;
    startedAt: string;
    sortOrder: number;
  }[];
  workHistory: {
    id?: number;
    organization: LocalizedText;
    role: LocalizedText;
    summary: LocalizedText;
    startedAt: string;
    endedAt?: string | null;
    externalUrl?: string;
    sortOrder: number;
  }[];
  socialLinks: {
    id?: number;
    provider: SocialProvider;
    label: LocalizedText;
    url: string;
    isFooter: boolean;
    sortOrder: number;
  }[];
};

type ProjectPayload = {
  title: { ja?: string; en?: string };
  description: { ja?: string; en?: string };
  tech: {
    membershipId?: number;
    techId: number;
    context: TechContext;
    note: string;
    sortOrder: number;
  }[];
  linkUrl: string;
  year: number;
  published: boolean;
  sortOrder?: number | null;
};

type ResearchPayload = {
  slug: string;
  kind: ResearchKind;
  title: LocalizedText;
  overview: LocalizedText;
  outcome: LocalizedText;
  outlook: LocalizedText;
  externalUrl: string;
  highlightImageUrl: string;
  imageAlt: LocalizedText;
  publishedAt: string;
  isDraft: boolean;
  tags: { id?: number; value: string; sortOrder: number }[];
  links: {
    id?: number;
    type: ResearchLinkType;
    label: LocalizedText;
    url: string;
    sortOrder: number;
  }[];
  assets: {
    id?: number;
    url: string;
    caption: LocalizedText;
    sortOrder: number;
  }[];
  tech: {
    membershipId?: number;
    techId: number;
    context: TechContext;
    note: string;
    sortOrder: number;
  }[];
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

type ContactSettingsPayload = {
  id: number;
  heroTitle: LocalizedText;
  heroDescription: LocalizedText;
  topics: {
    id: string;
    label: LocalizedText;
    description: LocalizedText;
  }[];
  consentText: LocalizedText;
  minimumLeadHours: number;
  recaptchaSiteKey: string;
  supportEmail: string;
  calendarTimezone: string;
  googleCalendarId: string;
  bookingWindowDays: number;
  updatedAt: string;
};

type HomeSettingsPayload = {
  id: number;
  profileId: number;
  heroSubtitle: LocalizedText;
  quickLinks: {
    id?: number;
    section: "profile" | "research_blog" | "projects" | "contact";
    label: LocalizedText;
    description: LocalizedText;
    cta: LocalizedText;
    targetUrl: string;
    sortOrder: number;
  }[];
  chipSources: {
    id?: number;
    source: "affiliation" | "community" | "skill";
    label: LocalizedText;
    limit: number;
    sortOrder: number;
  }[];
  updatedAt: string;
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
  listTechCatalog: (params?: {
    includeInactive?: boolean;
  }) => ApiResult<TechCatalogEntry[]>;
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
  getContactSettings: () => ApiResult<ContactFormSettings>;
  updateContactSettings: (
    payload: ContactSettingsPayload,
  ) => ApiResult<ContactFormSettings>;
  getHomeSettings: () => ApiResult<HomePageConfig>;
  updateHomeSettings: (
    payload: HomeSettingsPayload,
  ) => ApiResult<HomePageConfig>;

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
  listTechCatalog: (params) =>
    unwrap(
      apiClient.get<TechCatalogEntry[]>("/admin/tech-catalog", {
        params: {
          includeInactive: params?.includeInactive ?? false,
        },
      }),
    ),
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
  getContactSettings: () =>
    unwrap(apiClient.get<ContactFormSettings>("/admin/contact-settings")),
  updateContactSettings: (payload: ContactSettingsPayload) =>
    unwrap(
      apiClient.put<ContactFormSettings>("/admin/contact-settings", payload),
    ),
  getHomeSettings: () => unwrap(apiClient.get<HomePageConfig>("/admin/home")),
  updateHomeSettings: (payload: HomeSettingsPayload) =>
    unwrap(apiClient.put<HomePageConfig>("/admin/home", payload)),

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
        expiresAt?: number;
        email?: string;
        roles?: string[];
        source?: string;
        refreshed?: boolean;
      }>("/admin/auth/session"),
    ),
};
