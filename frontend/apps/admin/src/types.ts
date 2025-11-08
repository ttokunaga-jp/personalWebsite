export type LocalizedText = {
  ja?: string;
  en?: string;
};

export type ProfileThemeMode = "light" | "dark" | "system";

export type ProfileTheme = {
  mode: ProfileThemeMode;
  accentColor?: string | null;
};

export type ProfileLab = {
  name: LocalizedText;
  advisor: LocalizedText;
  room: LocalizedText;
  url?: string | null;
};

export type ProfileAffiliation = {
  id: number;
  kind: "affiliation" | "community";
  name: string;
  url?: string | null;
  description: LocalizedText;
  startedAt: string;
  sortOrder: number;
};

export type ProfileWorkExperience = {
  id: number;
  organization: LocalizedText;
  role: LocalizedText;
  summary: LocalizedText;
  startedAt: string;
  endedAt?: string | null;
  externalUrl?: string | null;
  sortOrder: number;
};

export type SocialProvider =
  | "github"
  | "zenn"
  | "linkedin"
  | "x"
  | "email"
  | "other";

export type ProfileSocialLink = {
  id: number;
  provider: SocialProvider;
  label: LocalizedText;
  url: string;
  isFooter: boolean;
  sortOrder: number;
};

export type ProfileTechSection = {
  id: number;
  title: LocalizedText;
  layout: "chips" | "list";
  breakpoint: string;
  sortOrder: number;
  members: AdminTechMembership[];
};

export type HomeQuickLinkSection =
  | "profile"
  | "research_blog"
  | "projects"
  | "contact";

export type HomeQuickLink = {
  id: number;
  section: HomeQuickLinkSection;
  label: LocalizedText;
  description: LocalizedText;
  cta: LocalizedText;
  targetUrl: string;
  sortOrder: number;
};

export type HomeChipSourceType = "affiliation" | "community" | "skill";

export type HomeChipSource = {
  id: number;
  source: HomeChipSourceType;
  label: LocalizedText;
  limit: number;
  sortOrder: number;
};

export type HomePageConfig = {
  id: number;
  profileId: number;
  heroSubtitle: LocalizedText;
  quickLinks: HomeQuickLink[];
  chipSources: HomeChipSource[];
  updatedAt: string;
};

export type AdminProfile = {
  id: number;
  displayName: string;
  headline: LocalizedText;
  summary: LocalizedText;
  avatarUrl?: string | null;
  location: LocalizedText;
  theme: ProfileTheme;
  lab: ProfileLab;
  affiliations: ProfileAffiliation[];
  communities: ProfileAffiliation[];
  workHistory: ProfileWorkExperience[];
  socialLinks: ProfileSocialLink[];
  techSections: ProfileTechSection[];
  home?: HomePageConfig | null;
  updatedAt: string;
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

export type ContactTopic = {
  id: string;
  label: LocalizedText;
  description: LocalizedText;
};

export type ContactFormSettings = {
  id: number;
  heroTitle: LocalizedText;
  heroDescription: LocalizedText;
  topics: ContactTopic[];
  consentText: LocalizedText;
  minimumLeadHours: number;
  recaptchaSiteKey: string;
  supportEmail: string;
  calendarTimezone: string;
  googleCalendarId: string;
  bookingWindowDays: number;
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
