export type LocalizedText = {
  ja?: string | null;
  en?: string | null;
};

export type SupportedLanguage = "ja" | "en";

export type TechLevel = "beginner" | "intermediate" | "advanced";

export type TechContext = "primary" | "supporting";

export type TechCatalogEntry = {
  id: string;
  slug: string;
  displayName: string;
  category?: string;
  level: TechLevel;
  icon?: string;
  sortOrder: number;
  active: boolean;
};

export type TechMembership = {
  id: string;
  context: TechContext;
  note?: string;
  sortOrder: number;
  tech: TechCatalogEntry;
};

export type ProfileThemeMode = "light" | "dark" | "system";

export type ProfileTheme = {
  mode: ProfileThemeMode;
  accentColor?: string;
};

export type ProfileLab = {
  name?: string;
  advisor?: string;
  room?: string;
  url?: string;
};

export type ProfileAffiliation = {
  id: string;
  name: string;
  url?: string;
  description?: string;
  startedAt: string;
  sortOrder: number;
};

export type ProfileWorkHistoryItem = {
  id: string;
  organization: string;
  role: string;
  summary?: string;
  startedAt: string;
  endedAt?: string | null;
  externalUrl?: string;
  sortOrder: number;
};

export type ProfileTechSection = {
  id: string;
  title: string;
  layout: string;
  breakpoint: string;
  sortOrder: number;
  members: TechMembership[];
};

export type SocialProvider =
  | "github"
  | "zenn"
  | "linkedin"
  | "x"
  | "email"
  | "website"
  | "other";

export type SocialLink = {
  id: string;
  provider: SocialProvider;
  label: string;
  url: string;
  isFooter: boolean;
  sortOrder: number;
};

export type HomeQuickLinkSection =
  | "profile"
  | "research_blog"
  | "projects"
  | "contact";

export type HomeQuickLink = {
  id: string;
  section: HomeQuickLinkSection;
  label: string;
  description?: string;
  cta: string;
  targetUrl: string;
  sortOrder: number;
};

export type HomeChipSourceKind = "affiliation" | "community" | "tech";

export type HomeChipSource = {
  id: string;
  source: HomeChipSourceKind;
  label: string;
  limit: number;
  sortOrder: number;
};

export type HomePageConfig = {
  heroSubtitle?: string;
  quickLinks: HomeQuickLink[];
  chipSources: HomeChipSource[];
  updatedAt: string;
};

export type ProfileResponse = {
  id: string;
  displayName: string;
  headline?: string;
  summary?: string;
  avatarUrl?: string;
  location?: string;
  theme: ProfileTheme;
  lab?: ProfileLab;
  affiliations: ProfileAffiliation[];
  communities: ProfileAffiliation[];
  workHistory: ProfileWorkHistoryItem[];
  techSections: ProfileTechSection[];
  socialLinks: SocialLink[];
  footerLinks: SocialLink[];
  updatedAt: string;
  home?: HomePageConfig;
};

export type ProjectLinkType =
  | "repo"
  | "demo"
  | "article"
  | "slides"
  | "other";

export type ProjectLink = {
  id: string;
  type: ProjectLinkType;
  label: string;
  url: string;
  sortOrder: number;
};

export type Project = {
  id: string;
  slug: string;
  title: string;
  summary?: string;
  description?: string;
  coverImageUrl?: string;
  primaryLink?: string;
  links: ProjectLink[];
  period: {
    start?: string | null;
    end?: string | null;
  };
  tech: TechMembership[];
  highlight: boolean;
  published: boolean;
  sortOrder: number;
  createdAt: string;
  updatedAt: string;
};

export type ResearchLinkType =
  | "paper"
  | "slides"
  | "video"
  | "code"
  | "external";

export type ResearchLink = {
  id: string;
  type: ResearchLinkType;
  label: string;
  url: string;
  sortOrder: number;
};

export type ResearchAsset = {
  id: string;
  url: string;
  caption?: string;
  sortOrder: number;
};

export type ResearchEntry = {
  id: string;
  slug: string;
  kind: "research" | "blog";
  title: string;
  overview?: string;
  outcome?: string;
  outlook?: string;
  externalUrl: string;
  publishedAt: string;
  updatedAt: string;
  highlightImageUrl?: string;
  imageAlt?: string;
  isDraft: boolean;
  tags: string[];
  links: ResearchLink[];
  assets: ResearchAsset[];
  tech: TechMembership[];
};

export type ContactAvailabilitySlot = {
  id: string;
  start: string;
  end: string;
  isBookable: boolean;
};

export type ContactAvailabilityDay = {
  date: string;
  slots: ContactAvailabilitySlot[];
};

export type ContactAvailabilityResponse = {
  timezone: string;
  generatedAt: string;
  days: ContactAvailabilityDay[];
};

export type ContactTopic = {
  id: string;
  label: string;
  description?: string;
};

export type ContactConfigResponse = {
  heroTitle?: string;
  heroDescription?: string;
  topics: ContactTopic[];
  consentText?: string;
  minimumLeadHours: number;
  recaptchaSiteKey?: string;
  supportEmail?: string;
  calendarTimezone?: string;
  googleCalendarId?: string;
  bookingWindowDays?: number;
};

export type CreateBookingPayload = {
  name: string;
  email: string;
  topic: string;
  agenda: string;
  startTime: string;
  durationMinutes: number;
  recaptchaToken: string;
};

export type MeetingStatus = "pending" | "confirmed" | "cancelled";

export type Meeting = {
  id: string;
  name: string;
  email: string;
  datetime: string;
  durationMinutes: number;
  meetUrl?: string;
  calendarEventId?: string;
  status: MeetingStatus | string;
  notes?: string;
  confirmationSentAt?: string | null;
  lastNotificationSentAt?: string | null;
  lookupHash?: string;
  googleCalendarStatus?: string;
  cancellationReason?: string;
};

export type BookingResult = {
  meeting: Meeting;
  calendarEventId?: string;
  supportEmail?: string;
  calendarTimezone?: string;
};
