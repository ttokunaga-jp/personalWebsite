export type ReservationStatus = "pending" | "confirmed" | "cancelled";

export type ReservationNotificationStatus = "pending" | "sent" | "failed";

export type ReservationNotification = {
  id: number;
  reservationId: number;
  type: string;
  status: ReservationNotificationStatus;
  errorMessage?: string;
  createdAt: string;
};

export type Reservation = {
  id: number;
  name: string;
  email: string;
  topic?: string;
  message?: string;
  startAt: string;
  endAt: string;
  status: ReservationStatus;
  googleCalendarStatus?: string;
  lookupHash: string;
  cancellationReason?: string;
  createdAt: string;
  updatedAt: string;
  notifications?: ReservationNotification[];
};

export type ReservationUpdatePayload = {
  status: ReservationStatus;
  cancellationReason?: string;
};

export type BlacklistEntry = {
  id: number;
  email: string;
  reason: string;
  createdAt: string;
};

export type BlacklistInput = {
  email: string;
  reason: string;
};

export type TechLevel = "beginner" | "intermediate" | "advanced";

export type TechCatalogEntry = {
  id: number;
  slug: string;
  displayName: string;
  category?: string;
  level: TechLevel;
  icon?: string;
  sortOrder: number;
  active: boolean;
  createdAt?: string;
  updatedAt?: string;
};

export type TechCatalogInput = {
  slug: string;
  displayName: string;
  category?: string;
  level: TechLevel;
  icon?: string;
  sortOrder: number;
  active: boolean;
};

export type LocalizedField = {
  ja: string;
  en: string;
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
  id: number;
  provider: SocialProvider;
  label: LocalizedField;
  url: string;
  isFooter: boolean;
  sortOrder: number;
};

export type SocialLinkInput = {
  provider: SocialProvider;
  label: LocalizedField;
  url: string;
  isFooter: boolean;
  sortOrder: number;
};

export type MeetingUrlTemplate = {
  template: string;
  updatedAt?: string;
};

export type MeetingUrlPayload = {
  template: string;
};

export type HomeQuickLinkSection = "profile" | "research_blog" | "projects" | "contact";

export type HomeQuickLinkItem = {
  id: number;
  section: HomeQuickLinkSection;
  label: LocalizedField;
  description: LocalizedField;
  cta: LocalizedField;
  targetUrl: string;
  sortOrder: number;
};

export type HomeChipSourceKind = "tech" | "affiliation" | "community";

export type HomeChipSourceItem = {
  id: number;
  source: HomeChipSourceKind;
  label: LocalizedField;
  limit: number;
  sortOrder: number;
};

export type HomeConfigDocument = {
  id: number;
  heroSubtitle: LocalizedField;
  quickLinks: HomeQuickLinkItem[];
  chipSources: HomeChipSourceItem[];
  updatedAt?: string;
};

export type ProfileAffiliationKind = "affiliation" | "community";

export type ProfileAffiliationItem = {
  id: number;
  kind: ProfileAffiliationKind;
  name: string;
  url?: string;
  description: LocalizedField;
  startedAt: string;
  sortOrder: number;
};

export type ProfileWorkHistoryItem = {
  id: number;
  organization: LocalizedField;
  role: LocalizedField;
  summary: LocalizedField;
  startedAt: string;
  endedAt?: string | null;
  externalUrl?: string;
  sortOrder: number;
};

export type ProfileSocialLinkItem = {
  id: number;
  provider: SocialProvider;
  label: LocalizedField;
  url: string;
  isFooter: boolean;
  sortOrder: number;
};

export type ProfileDocument = {
  id: number;
  displayName: string;
  headline: LocalizedField;
  summary: LocalizedField;
  avatarUrl?: string;
  location: LocalizedField;
  affiliations: ProfileAffiliationItem[];
  communities: ProfileAffiliationItem[];
  workHistory: ProfileWorkHistoryItem[];
  socialLinks: ProfileSocialLinkItem[];
  updatedAt?: string;
};
