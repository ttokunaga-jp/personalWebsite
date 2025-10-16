export type LocalizedText = {
  ja?: string;
  en?: string;
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

export type BlogPost = {
  id: number;
  title: LocalizedText;
  summary: LocalizedText;
  contentMd: LocalizedText;
  tags: string[];
  published: boolean;
  publishedAt?: string | null;
  createdAt: string;
  updatedAt: string;
};

export type MeetingStatus = "pending" | "confirmed" | "cancelled";

export type Meeting = {
  id: number;
  name: string;
  email: string;
  datetime: string;
  durationMinutes: number;
  meetUrl: string;
  status: MeetingStatus;
  notes: string;
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
  publishedProjects: number;
  draftProjects: number;
  publishedResearch: number;
  draftResearch: number;
  publishedBlogs: number;
  draftBlogs: number;
  pendingMeetings: number;
  blacklistEntries: number;
};

