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
  tech: AdminTechMembership[];
  linkUrl: string;
  year: number;
  published: boolean;
  sortOrder?: number | null;
  createdAt: string;
  updatedAt: string;
};

export type ResearchKind = "research" | "blog";

export type ResearchLinkType = "paper" | "slides" | "video" | "code" | "external";

export type TechContext = "primary" | "supporting";

export type ResearchTag = {
  id: number;
  entryId: number;
  value: string;
  sortOrder: number;
};

export type ResearchLink = {
  id: number;
  entryId: number;
  type: ResearchLinkType;
  label: LocalizedText;
  url: string;
  sortOrder: number;
};

export type ResearchAsset = {
  id: number;
  entryId: number;
  url: string;
  caption: LocalizedText;
  sortOrder: number;
};

export type TechCatalogEntry = {
  id: number;
  slug: string;
  displayName: string;
  category?: string;
  level: "beginner" | "intermediate" | "advanced";
  icon?: string;
  sortOrder: number;
  active: boolean;
  createdAt?: string;
  updatedAt?: string;
};

export type AdminTechMembership = {
  membershipId: number;
  entityType: string;
  entityId: number;
  tech: TechCatalogEntry;
  context: TechContext;
  note: string;
  sortOrder: number;
};

export type AdminResearch = {
  id: number;
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
  createdAt: string;
  updatedAt: string;
  tags: ResearchTag[];
  links: ResearchLink[];
  assets: ResearchAsset[];
  tech: AdminTechMembership[];
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
