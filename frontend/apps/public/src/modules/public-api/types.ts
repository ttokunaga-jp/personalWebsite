export type SocialPlatform =
  | "github"
  | "x"
  | "twitter"
  | "linkedin"
  | "youtube"
  | "email"
  | "website"
  | "other";

export type SocialLink = {
  id: string;
  platform: SocialPlatform;
  label: string;
  url: string;
  handle?: string;
};

export type Affiliation = {
  id: string;
  organization: string;
  department?: string;
  role: string;
  startDate: string;
  endDate?: string | null;
  location?: string;
  isCurrent: boolean;
};

export type WorkHistoryItem = {
  id: string;
  organization: string;
  role: string;
  startDate: string;
  endDate?: string | null;
  achievements?: string[];
  description?: string;
  location?: string;
};

export type SkillGroup = {
  id: string;
  category: string;
  items: Array<{
    id: string;
    name: string;
    level: "beginner" | "intermediate" | "advanced" | "expert";
    description?: string;
  }>;
};

export type LabProfile = {
  name: string;
  advisor?: string;
  researchFocus?: string;
  websiteUrl?: string;
};

export type ProfileResponse = {
  name: string;
  headline: string;
  summary: string;
  avatarUrl?: string;
  location?: string;
  affiliations: Affiliation[];
  lab?: LabProfile;
  workHistory: WorkHistoryItem[];
  skillGroups: SkillGroup[];
  communities: string[];
  socialLinks: SocialLink[];
};

export type ResearchAsset = {
  alt: string;
  url: string;
  caption?: string;
};

export type ResearchEntry = {
  id: string;
  title: string;
  slug: string;
  summary: string;
  publishedOn: string;
  updatedOn?: string;
  tags: string[];
  contentMarkdown: string;
  contentHtml?: string;
  assets?: ResearchAsset[];
  links?: Array<{
    label: string;
    url: string;
    type: "paper" | "slide" | "code" | "video" | "other";
  }>;
};

export type ProjectLink = {
  label: string;
  url: string;
  type: "repo" | "demo" | "article" | "paper" | "other";
};

export type Project = {
  id: string;
  title: string;
  subtitle?: string;
  description: string;
  techStack: string[];
  category?: string;
  tags?: string[];
  period?: {
    start: string;
    end?: string | null;
  };
  links: ProjectLink[];
  coverImageUrl?: string;
  highlight?: boolean;
};

export type AvailabilitySlot = {
  id: string;
  start: string;
  end: string;
  isBookable: boolean;
};

export type ContactAvailabilityResponse = {
  timezone: string;
  generatedAt: string;
  slots: AvailabilitySlot[];
};

export type ContactConfigResponse = {
  topics: string[];
  recaptchaSiteKey: string;
  minimumLeadHours: number;
  consentText: string;
};

export type CreateBookingPayload = {
  name: string;
  email: string;
  topic: string;
  message: string;
  slotId: string;
  recaptchaToken: string;
};

export type CreateBookingResponse = {
  bookingId: string;
  status: "pending" | "confirmed";
  calendarUrl?: string;
};
