import i18next from "i18next";

import {
  getCanonicalHomeConfig,
  getCanonicalProfile,
  getCanonicalProjects,
  getCanonicalResearchEntries,
} from "../profile-content";

import type {
  ContactConfigResponse,
  ContactTopic,
  HomeChipSource,
  HomePageConfig,
  HomeQuickLink,
  LocalizedText,
  ProfileAffiliation,
  ProfileLab,
  ProfileResponse,
  ProfileTechSection,
  ProfileTheme,
  ProfileWorkHistoryItem,
  Project,
  ProjectLink,
  ResearchAsset,
  ResearchEntry,
  ResearchLink,
  SocialLink,
  TechCatalogEntry,
  TechContext,
  TechMembership,
  TechLevel,
} from "./types";
import type { SupportedLanguage } from "./types";

type RawTechCatalogEntry = {
  id?: number | string;
  slug?: string | null;
  displayName?: string | null;
  category?: string | null;
  level?: TechLevel | null;
  icon?: string | null;
  sortOrder?: number | null;
  active?: boolean | null;
};

type RawTechMembership = {
  membershipId?: number | string;
  context?: string | null;
  note?: string | null;
  sortOrder?: number | null;
  tech?: RawTechCatalogEntry | null;
};

type RawProfileAffiliation = {
  id?: number | string;
  name?: string | null;
  url?: string | null;
  description?: LocalizedText | null;
  startedAt?: string | null;
  sortOrder?: number | null;
};

type RawProfileWorkHistory = {
  id?: number | string;
  organization?: LocalizedText | null;
  role?: LocalizedText | null;
  summary?: LocalizedText | null;
  startedAt?: string | null;
  endedAt?: string | null;
  externalUrl?: string | null;
  sortOrder?: number | null;
};

type RawProfileTechSection = {
  id?: number | string;
  title?: LocalizedText | null;
  layout?: string | null;
  breakpoint?: string | null;
  sortOrder?: number | null;
  members?: RawTechMembership[] | null;
};

type RawProfileSocialLink = {
  id?: number | string;
  provider?: string | null;
  label?: LocalizedText | null;
  url?: string | null;
  isFooter?: boolean | null;
  sortOrder?: number | null;
};

type RawProfileLab = {
  name?: LocalizedText | null;
  advisor?: LocalizedText | null;
  room?: LocalizedText | null;
  url?: string | null;
};

export type RawProfileDocument = {
  id?: number | string;
  displayName?: string | null;
  headline?: LocalizedText | null;
  summary?: LocalizedText | null;
  avatarUrl?: string | null;
  location?: LocalizedText | null;
  theme?: {
    mode?: string | null;
    accentColor?: string | null;
  } | null;
  lab?: RawProfileLab | null;
  affiliations?: RawProfileAffiliation[] | null;
  communities?: RawProfileAffiliation[] | null;
  workHistory?: RawProfileWorkHistory[] | null;
  techSections?: RawProfileTechSection[] | null;
  socialLinks?: RawProfileSocialLink[] | null;
  updatedAt?: string | null;
};

type RawHomeQuickLink = {
  id?: number | string;
  section?: string | null;
  label?: LocalizedText | null;
  description?: LocalizedText | null;
  cta?: LocalizedText | null;
  targetUrl?: string | null;
  sortOrder?: number | null;
};

type RawHomeChipSource = {
  id?: number | string;
  source?: string | null;
  label?: LocalizedText | null;
  limit?: number | null;
  sortOrder?: number | null;
};

export type RawHomePageConfig = {
  heroSubtitle?: LocalizedText | null;
  quickLinks?: RawHomeQuickLink[] | null;
  chipSources?: RawHomeChipSource[] | null;
  updatedAt?: string | null;
};

type RawProjectLink = {
  id?: number | string;
  type?: string | null;
  label?: LocalizedText | null;
  url?: string | null;
  sortOrder?: number | null;
};

type RawProjectPeriod = {
  start?: string | null;
  end?: string | null;
};

export type RawProjectDocument = {
  id?: number | string;
  slug?: string | null;
  title?: LocalizedText | null;
  summary?: LocalizedText | null;
  description?: LocalizedText | null;
  coverImageUrl?: string | null;
  primaryLink?: string | null;
  links?: RawProjectLink[] | null;
  period?: RawProjectPeriod | null;
  tech?: RawTechMembership[] | null;
  highlight?: boolean | null;
  published?: boolean | null;
  sortOrder?: number | null;
  createdAt?: string | null;
  updatedAt?: string | null;
};

type RawResearchLink = {
  id?: number | string;
  type?: string | null;
  label?: LocalizedText | null;
  url?: string | null;
  sortOrder?: number | null;
};

type RawResearchAsset = {
  id?: number | string;
  url?: string | null;
  caption?: LocalizedText | null;
  sortOrder?: number | null;
};

type RawResearchTag =
  | {
      id?: number | string;
      value?: string | null;
    }
  | string;

export type RawResearchDocument = {
  id?: number | string;
  slug?: string | null;
  kind?: string | null;
  title?: LocalizedText | null;
  overview?: LocalizedText | null;
  outcome?: LocalizedText | null;
  outlook?: LocalizedText | null;
  externalUrl?: string | null;
  publishedAt?: string | null;
  updatedAt?: string | null;
  highlightImageUrl?: string | null;
  imageAlt?: LocalizedText | null;
  isDraft?: boolean | null;
  tags?: RawResearchTag[] | null;
  links?: RawResearchLink[] | null;
  assets?: RawResearchAsset[] | null;
  tech?: RawTechMembership[] | null;
};

type RawContactTopic = {
  id?: string | null;
  label?: LocalizedText | null;
  description?: LocalizedText | null;
};

export type RawContactConfig = {
  heroTitle?: LocalizedText | null;
  heroDescription?: LocalizedText | null;
  topics?: RawContactTopic[] | null;
  consentText?: LocalizedText | null;
  minimumLeadHours?: number | null;
  recaptchaSiteKey?: string | null;
  supportEmail?: string | null;
  calendarTimezone?: string | null;
  googleCalendarId?: string | null;
  bookingWindowDays?: number | null;
};

const FALLBACK_LANGUAGE: SupportedLanguage = "en";

function clone<T>(value: T): T {
  if (typeof structuredClone === "function") {
    return structuredClone(value);
  }
  return JSON.parse(JSON.stringify(value)) as T;
}

function resolveLanguage(): SupportedLanguage {
  const language = i18next.language ?? FALLBACK_LANGUAGE;
  return language.toLowerCase().startsWith("ja") ? "ja" : "en";
}

function normaliseString(value?: string | null): string | undefined {
  const trimmed = value?.trim();
  return trimmed ? trimmed : undefined;
}

function toStringId(value: number | string | null | undefined, fallback: string): string {
  if (value === null || value === undefined || value === "") {
    return fallback;
  }
  return String(value);
}

function toNumber(value: number | null | undefined, fallback = 0): number {
  if (typeof value === "number" && Number.isFinite(value)) {
    return value;
  }
  return fallback;
}

function selectLocalizedText(
  text: LocalizedText | string | null | undefined,
  language: SupportedLanguage,
  fallback?: string,
): string | undefined {
  if (!text) {
    return fallback;
  }

  if (typeof text === "string") {
    const normalised = normaliseString(text);
    return normalised ?? fallback;
  }

  const primary =
    language === "ja"
      ? normaliseString(text.ja)
      : normaliseString(text.en);
  if (primary) {
    return primary;
  }

  const secondary =
    language === "ja"
      ? normaliseString(text.en)
      : normaliseString(text.ja);
  if (secondary) {
    return secondary;
  }

  return fallback;
}

function mapTechCatalogEntry(raw?: RawTechCatalogEntry | null): TechCatalogEntry | null {
  if (!raw) {
    return null;
  }
  const displayName = normaliseString(raw.displayName);
  if (!displayName) {
    return null;
  }
  const slug = normaliseString(raw.slug) ?? displayName.toLowerCase().replace(/\s+/g, "-");

  return {
    id: toStringId(raw.id ?? slug, slug),
    slug,
    displayName,
    category: normaliseString(raw.category),
    level: (raw.level ?? "intermediate") as TechLevel,
    icon: normaliseString(raw.icon),
    sortOrder: toNumber(raw.sortOrder),
    active: raw.active ?? true,
  };
}

function mapTechMembership(raw: RawTechMembership | null | undefined): TechMembership | null {
  if (!raw) {
    return null;
  }
  const tech = mapTechCatalogEntry(raw.tech);
  if (!tech) {
    return null;
  }

  const context = raw.context === "supporting" ? "supporting" : "primary";

  return {
    id: toStringId(raw.membershipId ?? `${tech.id}-${context}`, `${tech.id}-${context}`),
    context: context as TechContext,
    note: normaliseString(raw.note),
    sortOrder: toNumber(raw.sortOrder),
    tech,
  };
}

function mapAffiliations(
  rows: RawProfileAffiliation[] | null | undefined,
  language: SupportedLanguage,
  fallback: ProfileAffiliation[],
): ProfileAffiliation[] {
  if (!rows?.length) {
    return clone(fallback);
  }

  const mapped: ProfileAffiliation[] = [];
  rows.forEach((row, index) => {
    const name = normaliseString(row.name);
    if (!name) {
      return;
    }

    const affiliation: ProfileAffiliation = {
      id: toStringId(row.id, `affiliation-${index}`),
      name,
      startedAt: row.startedAt ?? fallback[index]?.startedAt ?? "",
      sortOrder: toNumber(row.sortOrder, index),
    };

    const url = normaliseString(row.url);
    if (url) {
      affiliation.url = url;
    }

    const description = selectLocalizedText(row.description ?? undefined, language);
    if (description) {
      affiliation.description = description;
    }

    mapped.push(affiliation);
  });

  return mapped.sort((a, b) => a.sortOrder - b.sortOrder);
}

function mapWorkHistory(
  rows: RawProfileWorkHistory[] | null | undefined,
  language: SupportedLanguage,
  fallback: ProfileWorkHistoryItem[],
): ProfileWorkHistoryItem[] {
  if (!rows?.length) {
    return clone(fallback);
  }

  const mapped: ProfileWorkHistoryItem[] = [];
  rows.forEach((row, index) => {
    const organization = selectLocalizedText(row.organization ?? undefined, language);
    const role = selectLocalizedText(row.role ?? undefined, language);
    if (!organization || !role) {
      return;
    }

    const item: ProfileWorkHistoryItem = {
      id: toStringId(row.id, `work-${index}`),
      organization,
      role,
      startedAt: row.startedAt ?? fallback[index]?.startedAt ?? "",
      endedAt: row.endedAt ?? undefined,
      sortOrder: toNumber(row.sortOrder, index),
    };

    const summary = selectLocalizedText(row.summary ?? undefined, language);
    if (summary) {
      item.summary = summary;
    }

    const externalUrl = normaliseString(row.externalUrl);
    if (externalUrl) {
      item.externalUrl = externalUrl;
    }

    mapped.push(item);
  });

  return mapped.sort((a, b) => a.sortOrder - b.sortOrder);
}

function mapTechSections(
  sections: RawProfileTechSection[] | null | undefined,
  language: SupportedLanguage,
  fallback: ProfileTechSection[],
): ProfileTechSection[] {
  if (!sections?.length) {
    return clone(fallback);
  }

  const mapped: ProfileTechSection[] = [];
  sections.forEach((section, index) => {
    const title = selectLocalizedText(section.title ?? undefined, language);
    if (!title) {
      return;
    }

    const members =
      section.members
        ?.map((member) => mapTechMembership(member))
        .filter((member): member is TechMembership => Boolean(member))
        .sort((a, b) => a.sortOrder - b.sortOrder) ?? [];

    mapped.push({
      id: toStringId(section.id, `tech-section-${index}`),
      title,
      layout: section.layout ?? "grid",
      breakpoint: section.breakpoint ?? "md",
      sortOrder: toNumber(section.sortOrder, index),
      members,
    });
  });

  return mapped.sort((a, b) => a.sortOrder - b.sortOrder);
}

function mapSocialLinks(
  links: RawProfileSocialLink[] | null | undefined,
  language: SupportedLanguage,
  fallback: SocialLink[],
): SocialLink[] {
  if (!links?.length) {
    return clone(fallback);
  }

  return links
    .map((link, index) => {
      const url = normaliseString(link.url);
      const label = selectLocalizedText(link.label ?? undefined, language);
      const provider = (link.provider ?? "other").toLowerCase() as SocialLink["provider"];
      if (!url || !label) {
        return null;
      }

      return {
        id: toStringId(link.id, `social-${index}`),
        provider,
        label,
        url,
        isFooter: Boolean(link.isFooter),
        sortOrder: toNumber(link.sortOrder, index),
      };
    })
    .filter((value): value is SocialLink => value !== null)
    .sort((a, b) => a.sortOrder - b.sortOrder);
}

function mapLab(
  lab: RawProfileLab | null | undefined,
  language: SupportedLanguage,
  fallback: ProfileLab | undefined,
): ProfileLab | undefined {
  if (!lab) {
    return fallback ? clone(fallback) : undefined;
  }
  const name = selectLocalizedText(lab.name ?? undefined, language, fallback?.name);
  const advisor = selectLocalizedText(lab.advisor ?? undefined, language, fallback?.advisor);
  const room = selectLocalizedText(lab.room ?? undefined, language, fallback?.room);
  const url = normaliseString(lab.url) ?? fallback?.url;

  if (!name && !advisor && !room && !url) {
    return fallback ? clone(fallback) : undefined;
  }

  return {
    name,
    advisor,
    room,
    url,
  };
}

function mapHomeQuickLinks(
  links: RawHomeQuickLink[] | null | undefined,
  language: SupportedLanguage,
  fallback: HomeQuickLink[],
): HomeQuickLink[] {
  if (!links?.length) {
    return clone(fallback);
  }

  const mapped: HomeQuickLink[] = [];
  links.forEach((link, index) => {
    const label = selectLocalizedText(link.label ?? undefined, language);
    const cta = selectLocalizedText(link.cta ?? undefined, language);
    const targetUrl = normaliseString(link.targetUrl);
    if (!label || !cta || !targetUrl) {
      return;
    }

    const item: HomeQuickLink = {
      id: toStringId(link.id, `home-link-${index}`),
      section: (link.section ?? "profile") as HomeQuickLink["section"],
      label,
      cta,
      targetUrl,
      sortOrder: toNumber(link.sortOrder, index),
    };

    const description = selectLocalizedText(link.description ?? undefined, language);
    if (description) {
      item.description = description;
    }

    mapped.push(item);
  });

  return mapped.sort((a, b) => a.sortOrder - b.sortOrder);
}

function mapHomeChipSources(
  sources: RawHomeChipSource[] | null | undefined,
  language: SupportedLanguage,
  fallback: HomeChipSource[],
): HomeChipSource[] {
  if (!sources?.length) {
    return clone(fallback);
  }

  const mapped: HomeChipSource[] = [];
  sources.forEach((source, index) => {
    const label = selectLocalizedText(source.label ?? undefined, language);
    if (!label) {
      return;
    }

    mapped.push({
      id: toStringId(source.id, `chip-source-${index}`),
      source: (source.source ?? "tech") as HomeChipSource["source"],
      label,
      limit: Math.max(1, toNumber(source.limit, 6)),
      sortOrder: toNumber(source.sortOrder, index),
    });
  });

  return mapped.sort((a, b) => a.sortOrder - b.sortOrder);
}

function transformHomeConfig(
  raw: RawHomePageConfig | null | undefined,
  language: SupportedLanguage,
  fallback: HomePageConfig,
): HomePageConfig {
  if (!raw) {
    return clone(fallback);
  }

  const heroSubtitle = selectLocalizedText(
    raw.heroSubtitle ?? undefined,
    language,
    fallback.heroSubtitle,
  );

  const quickLinks = mapHomeQuickLinks(raw.quickLinks, language, fallback.quickLinks);
  const chipSources = mapHomeChipSources(raw.chipSources, language, fallback.chipSources);

  return {
    heroSubtitle,
    quickLinks,
    chipSources,
    updatedAt: raw.updatedAt ?? fallback.updatedAt,
  };
}

export function transformProfile(
  raw: RawProfileDocument | undefined,
  homeConfig?: RawHomePageConfig | undefined,
): ProfileResponse {
  const language = resolveLanguage();

  const canonicalProfile = getCanonicalProfile(language);
  const canonicalHome = getCanonicalHomeConfig(language);
  const profile = clone(canonicalProfile);

  profile.id = toStringId(raw?.id ?? profile.id, "profile");
  profile.displayName = normaliseString(raw?.displayName) ?? profile.displayName;
  profile.headline = selectLocalizedText(raw?.headline ?? undefined, language, profile.headline);
  profile.summary = selectLocalizedText(raw?.summary ?? undefined, language, profile.summary);
  profile.avatarUrl = normaliseString(raw?.avatarUrl) ?? profile.avatarUrl;
  profile.location = selectLocalizedText(raw?.location ?? undefined, language, profile.location);

  const themeMode = raw?.theme?.mode;
  if (themeMode === "light" || themeMode === "dark" || themeMode === "system") {
    profile.theme = {
      mode: themeMode as ProfileTheme["mode"],
      accentColor: normaliseString(raw?.theme?.accentColor) ?? profile.theme.accentColor,
    };
  } else if (raw?.theme?.accentColor) {
    profile.theme = {
      ...profile.theme,
      accentColor:
        normaliseString(raw?.theme?.accentColor) ?? profile.theme.accentColor,
    };
  }

  profile.lab = mapLab(raw?.lab ?? undefined, language, profile.lab);

  profile.affiliations = mapAffiliations(raw?.affiliations, language, profile.affiliations);
  profile.communities = mapAffiliations(raw?.communities, language, profile.communities);
  profile.workHistory = mapWorkHistory(raw?.workHistory, language, profile.workHistory);
  profile.techSections = mapTechSections(raw?.techSections, language, profile.techSections);
  profile.socialLinks = mapSocialLinks(raw?.socialLinks, language, profile.socialLinks);
  profile.footerLinks = profile.socialLinks.filter((link) => link.isFooter);
  profile.updatedAt = raw?.updatedAt ?? profile.updatedAt;
  profile.home = transformHomeConfig(homeConfig ?? null, language, canonicalHome);

  return profile;
}

function mapProjectLinks(
  links: RawProjectLink[] | null | undefined,
  language: SupportedLanguage,
  fallback: ProjectLink[],
): ProjectLink[] {
  if (!links?.length) {
    return clone(fallback);
  }

  const mapped: ProjectLink[] = [];
  links.forEach((link, index) => {
    const label = selectLocalizedText(link.label ?? undefined, language);
    const url = normaliseString(link.url);
    if (!label || !url) {
      return;
    }

    mapped.push({
      id: toStringId(link.id, `project-link-${index}`),
      type: (link.type ?? "other") as ProjectLink["type"],
      label,
      url,
      sortOrder: toNumber(link.sortOrder, index),
    });
  });

  return mapped.sort((a, b) => a.sortOrder - b.sortOrder);
}

function transformProject(
  raw: RawProjectDocument,
  language: SupportedLanguage,
  fallback: Project,
): Project {
  const project = clone(fallback);

  project.id = toStringId(raw.id ?? project.id, `project-${project.id}`);
  project.slug = normaliseString(raw.slug) ?? project.slug;
  const projectTitle =
    selectLocalizedText(raw.title ?? undefined, language, project.title) ?? project.title;
  project.title = projectTitle;
  project.summary = selectLocalizedText(raw.summary ?? undefined, language, project.summary);
  project.description = selectLocalizedText(
    raw.description ?? undefined,
    language,
    project.description,
  );
  project.coverImageUrl = normaliseString(raw.coverImageUrl) ?? project.coverImageUrl;
  project.primaryLink = normaliseString(raw.primaryLink) ?? project.primaryLink;
  project.links = mapProjectLinks(raw.links, language, project.links);
  project.period = {
    start: raw.period?.start ?? project.period?.start ?? null,
    end: raw.period?.end ?? project.period?.end ?? null,
  };
  project.tech =
    raw.tech
      ?.map((item) => mapTechMembership(item))
      .filter((item): item is TechMembership => Boolean(item))
      .sort((a, b) => a.sortOrder - b.sortOrder) ?? project.tech;
  project.highlight = Boolean(
    raw.highlight ?? project.highlight ?? project.links.some((link) => link.type === "demo"),
  );
  project.published = raw.published ?? project.published ?? true;
  project.sortOrder = toNumber(raw.sortOrder, project.sortOrder ?? 0);
  if (typeof raw.createdAt === "string") {
    project.createdAt = raw.createdAt;
  }
  if (typeof raw.updatedAt === "string") {
    project.updatedAt = raw.updatedAt;
  }

  return project;
}

export function transformProjects(
  projects: RawProjectDocument[] | undefined,
): Project[] {
  const language = resolveLanguage();
  const canonical = getCanonicalProjects(language).map(clone);

  if (!projects?.length) {
    return canonical;
  }

  return projects.map((raw, index) => {
    const fallback = canonical[index] ?? canonical[0];
    return transformProject(raw, language, fallback);
  });
}

function mapResearchLinks(
  links: RawResearchLink[] | null | undefined,
  language: SupportedLanguage,
  fallback: ResearchLink[],
): ResearchLink[] {
  if (!links?.length) {
    return clone(fallback);
  }

  const mapped: ResearchLink[] = [];
  links.forEach((link, index) => {
    const label = selectLocalizedText(link.label ?? undefined, language);
    const url = normaliseString(link.url);
    if (!label || !url) {
      return;
    }

    mapped.push({
      id: toStringId(link.id, `research-link-${index}`),
      type: (link.type ?? "external") as ResearchLink["type"],
      label,
      url,
      sortOrder: toNumber(link.sortOrder, index),
    });
  });

  return mapped.sort((a, b) => a.sortOrder - b.sortOrder);
}

function mapResearchAssets(
  assets: RawResearchAsset[] | null | undefined,
  language: SupportedLanguage,
  fallback: ResearchAsset[],
): ResearchAsset[] {
  if (!assets?.length) {
    return clone(fallback);
  }

  const mapped: ResearchAsset[] = [];
  assets.forEach((asset, index) => {
    const url = normaliseString(asset.url);
    if (!url) {
      return;
    }

    const caption = selectLocalizedText(asset.caption ?? undefined, language);
    const item: ResearchAsset = {
      id: toStringId(asset.id, `research-asset-${index}`),
      url,
      sortOrder: toNumber(asset.sortOrder, index),
    };
    if (caption) {
      item.caption = caption;
    }

    mapped.push(item);
  });

  return mapped.sort((a, b) => a.sortOrder - b.sortOrder);
}

function transformResearchEntry(
  raw: RawResearchDocument,
  language: SupportedLanguage,
  fallback: ResearchEntry,
): ResearchEntry {
  const entry = clone(fallback);

  entry.id = toStringId(raw.id ?? entry.id, `research-${entry.id}`);
  entry.slug = normaliseString(raw.slug) ?? entry.slug;
  entry.kind = (raw.kind ?? entry.kind ?? "research") as ResearchEntry["kind"];
  const resolvedTitle =
    selectLocalizedText(raw.title ?? undefined, language, entry.title) ?? entry.title;
  entry.title = resolvedTitle;
  entry.overview =
    selectLocalizedText(raw.overview ?? undefined, language, entry.overview) ?? entry.overview;
  entry.outcome =
    selectLocalizedText(raw.outcome ?? undefined, language, entry.outcome) ?? entry.outcome;
  entry.outlook =
    selectLocalizedText(raw.outlook ?? undefined, language, entry.outlook) ?? entry.outlook;
  const externalUrl = normaliseString(raw.externalUrl);
  if (externalUrl) {
    entry.externalUrl = externalUrl;
  }
  if (typeof raw.publishedAt === "string") {
    entry.publishedAt = raw.publishedAt;
  }
  if (typeof raw.updatedAt === "string") {
    entry.updatedAt = raw.updatedAt;
  }
  entry.highlightImageUrl =
    normaliseString(raw.highlightImageUrl) ?? entry.highlightImageUrl;
  entry.imageAlt = selectLocalizedText(raw.imageAlt ?? undefined, language, entry.imageAlt);
  entry.isDraft = raw.isDraft ?? entry.isDraft ?? false;
  entry.tags =
    raw.tags
      ?.map((tag) => {
        if (typeof tag === "string") {
          return normaliseString(tag);
        }
        return normaliseString(tag.value);
      })
      .filter((tag): tag is string => Boolean(tag)) ??
    entry.tags;
  entry.links = mapResearchLinks(raw.links, language, entry.links);
  entry.assets = mapResearchAssets(raw.assets, language, entry.assets);
  entry.tech =
    raw.tech
      ?.map((tech) => mapTechMembership(tech))
      .filter((tech): tech is TechMembership => Boolean(tech))
      .sort((a, b) => a.sortOrder - b.sortOrder) ?? entry.tech;

  return entry;
}

export function transformResearchEntries(
  entries: RawResearchDocument[] | undefined,
): ResearchEntry[] {
  const language = resolveLanguage();
  const canonical = getCanonicalResearchEntries(language).map(clone);

  if (!entries?.length) {
    return canonical;
  }

  return entries.map((raw, index) => {
    const fallback = canonical[index] ?? canonical[0];
    return transformResearchEntry(raw, language, fallback);
  });
}

export function transformContactConfig(
  raw: RawContactConfig | undefined,
): ContactConfigResponse {
  const language = resolveLanguage();
  const topics =
    raw?.topics?.reduce<ContactConfigResponse["topics"]>((acc, topic, index) => {
      const label = selectLocalizedText(topic.label ?? undefined, language);
      if (!label) {
        return acc;
      }
      const description = selectLocalizedText(topic.description ?? undefined, language);
      const item: ContactTopic = {
        id: topic.id ?? `topic-${index}`,
        label,
      };
      if (description) {
        item.description = description;
      }
      acc.push(item);
      return acc;
    }, []) ?? [];

  return {
    heroTitle: selectLocalizedText(raw?.heroTitle ?? undefined, language),
    heroDescription: selectLocalizedText(raw?.heroDescription ?? undefined, language),
    topics,
    consentText: selectLocalizedText(raw?.consentText ?? undefined, language),
    minimumLeadHours: Math.max(1, raw?.minimumLeadHours ?? 24),
    recaptchaSiteKey: normaliseString(raw?.recaptchaSiteKey),
    supportEmail: normaliseString(raw?.supportEmail),
    calendarTimezone: normaliseString(raw?.calendarTimezone),
    googleCalendarId: normaliseString(raw?.googleCalendarId),
    bookingWindowDays: raw?.bookingWindowDays ?? undefined,
  };
}
