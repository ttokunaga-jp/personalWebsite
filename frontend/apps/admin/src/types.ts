export type LocalizedText = {
  ja?: string;
  en?: string;
};

export type AdminProfile = {
  name: LocalizedText;
  title: LocalizedText;
  affiliation: LocalizedText;
  lab: LocalizedText;
  summary: LocalizedText;
  skills: LocalizedText[];
  focusAreas: LocalizedText[];
  updatedAt?: string | null;
};

export type AdminProject = {
  id: number;
  title: LocalizedText;
  description: LocalizedText;
  techStack: string[];
  linkUrl: string;
  year: number;
  published: boolean;
  sortOrder?: number | null;
  createdAt: string;
  updatedAt: string;
};

export type AdminResearch = {
  id: number;
  title: LocalizedText;
  summary: LocalizedText;
  contentMd: LocalizedText;
  year: number;
  published: boolean;
  createdAt: string;
  updatedAt: string;
};

export type ContactStatus = "pending" | "in_review" | "resolved" | "archived";

export type ContactMessage = {
  id: string;
  name: string;
  email: string;
  topic: string;
  message: string;
  status: ContactStatus;
  adminNote: string;
  createdAt: string;
  updatedAt: string;
};

export type BlacklistEntry = {
  id: number;
  email: string;
  reason: string;
  createdAt: string;
};

export type AdminSummary = {
  profileUpdatedAt?: string | null;
  skillCount: number;
  focusAreaCount: number;
  publishedProjects: number;
  draftProjects: number;
  publishedResearch: number;
  draftResearch: number;
  pendingContacts: number;
  blacklistEntries: number;
};
