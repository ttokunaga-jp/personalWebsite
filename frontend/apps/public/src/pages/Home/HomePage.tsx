import { apiClient } from "@shared/lib/api-client";
import type { CSSProperties } from "react";
import { startTransition, useEffect, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";

import { useAdminMode } from "../../hooks/useAdminMode";
import type { LocalizedField } from "../../modules/admin-console/types";
import {
  HomeEditorPanel,
  HomeEditorProvider,
  useOptionalHomeEditorContext,
  type EditableHomeConfig,
  type EditableQuickLink,
  type EditableChipSource,
} from "../../modules/admin-editors/HomeEditor";
import { getCanonicalProfile } from "../../modules/profile-content";
import { useProfileResource } from "../../modules/public-api";
import type {
  HomeChipSource,
  HomePageConfig,
  HomeQuickLink,
  ProfileResponse,
  SocialLink,
} from "../../modules/public-api";
import { formatDateRange } from "../../utils/date";
import { getSocialIcon } from "../../utils/icons";

const DISABLE_HEALTH_CHECKS =
  (import.meta.env?.["VITE_DISABLE_HEALTH_CHECKS"] ?? "false") === "true";

const QUICK_LINK_GRID_BASE_STYLES: CSSProperties = {
  display: "grid",
  gap: "1rem",
  gridTemplateColumns: "repeat(1, minmax(0, 1fr))",
  minHeight: "220px",
};
const QUICK_LINK_HERO_LIMIT = 2;

const QUICK_LINK_CARD_BASE_STYLES: CSSProperties = {
  display: "flex",
  flexDirection: "column",
  gap: "0.5rem",
  width: "100%",
  borderRadius: "1rem",
  borderWidth: 1,
  borderStyle: "solid",
  padding: "1.25rem",
  textDecoration: "none",
  boxSizing: "border-box",
};

type HealthResponse = {
  status: string;
};

type Chip = {
  id: string;
  label: string;
  description?: string;
};

type ChipGroup = {
  id: string;
  label: string;
  chips: Chip[];
};

function buildChipGroups(profile: ProfileResponse, presentLabel: string): ChipGroup[] {
  if (!profile.home) {
    return [];
  }

  return profile.home.chipSources.map((source) => {
    const chips: Chip[] = [];

    if (source.source === "tech") {
      const memberships = profile.techSections
        .flatMap((section) => section.members)
        .sort((a, b) => a.sortOrder - b.sortOrder);
      const seen = new Set<string>();
      for (const membership of memberships) {
        const tech = membership.tech;
        if (!tech.displayName || seen.has(tech.id)) {
          continue;
        }
        seen.add(tech.id);
        chips.push({
          id: membership.id,
          label: tech.displayName,
          description: tech.level ? tech.level : tech.category ?? undefined,
        });
        if (chips.length >= source.limit) {
          break;
        }
      }
    } else if (source.source === "affiliation") {
      const affiliations = [...profile.affiliations].sort(
        (a, b) => a.sortOrder - b.sortOrder,
      );
      for (const affiliation of affiliations.slice(0, source.limit)) {
        chips.push({
          id: affiliation.id,
          label: affiliation.name,
          description: formatDateRange(affiliation.startedAt, null, presentLabel),
        });
      }
    } else if (source.source === "community") {
      const communities = [...profile.communities].sort(
        (a, b) => a.sortOrder - b.sortOrder,
      );
      for (const community of communities.slice(0, source.limit)) {
        chips.push({
          id: community.id,
          label: community.name,
          description: community.description,
        });
      }
    }

    return {
      id: source.id,
      label: source.label,
      chips,
    };
  });
}

function getFeaturedLinks(links: SocialLink[]): SocialLink[] {
  const allowed = new Set([
    "github",
    "zenn",
    "linkedin",
    "x",
    "twitter",
    "email",
    "website",
  ]);
  return links.filter((link) => allowed.has(link.provider));
}

function resolveLocaleKey(language: string): "ja" | "en" {
  return language.startsWith("ja") ? "ja" : "en";
}

function localizedText(
  field: LocalizedField,
  locale: "ja" | "en",
  fallback?: string,
): string {
  const primary = locale === "ja" ? field.ja : field.en;
  const secondary = locale === "ja" ? field.en : field.ja;
  return primary ?? secondary ?? fallback ?? "";
}

function draftToPreviewConfig(
  draft: EditableHomeConfig,
  locale: "ja" | "en",
): HomePageConfig {
  const toQuickLink = (link: EditableQuickLink): HomeQuickLink => ({
    id: link.id ? String(link.id) : link.clientId,
    section: link.section,
    label: localizedText(link.label, locale).trim(),
    description: localizedText(link.description, locale).trim() || undefined,
    cta: localizedText(link.cta, locale).trim(),
    targetUrl: link.targetUrl,
    sortOrder: link.sortOrder,
  });

  const toChipSource = (source: EditableChipSource): HomeChipSource => ({
    id: source.id ? String(source.id) : source.clientId,
    source: source.source,
    label: localizedText(source.label, locale).trim(),
    limit: source.limit,
    sortOrder: source.sortOrder,
  });

  return {
    heroSubtitle: localizedText(draft.heroSubtitle, locale).trim() || undefined,
    quickLinks: draft.quickLinks.map(toQuickLink).sort((a, b) => a.sortOrder - b.sortOrder),
    chipSources: draft.chipSources
      .map(toChipSource)
      .sort((a, b) => a.sortOrder - b.sortOrder),
    updatedAt: draft.updatedAt ?? new Date().toISOString(),
  };
}

export function HomePage() {
  const { isAdminMode, sessionActive } = useAdminMode();
  const {
    data: profile,
    isLoading: isProfileLoading,
    error: profileError,
  } = useProfileResource();
  const [status, setStatus] = useState<string>(
    DISABLE_HEALTH_CHECKS ? "healthy" : "loading",
  );

  useEffect(() => {
    if (DISABLE_HEALTH_CHECKS) {
      return () => {};
    }

    let subscribed = true;

    const fetchHealthStatus = async () => {
      try {
        const { data } = await apiClient.get<HealthResponse>("/health");
        if (subscribed) {
          startTransition(() => {
            setStatus(data.status ?? "ok");
          });
        }
      } catch {
        if (subscribed) {
          startTransition(() => {
            setStatus("unreachable");
          });
        }
      }
    };

    void fetchHealthStatus();

    return () => {
      subscribed = false;
    };
  }, []);

  const editorEnabled = isAdminMode && sessionActive;

  return (
    <HomeEditorProvider enabled={editorEnabled}>
      <HomePageContent
        profile={profile ?? null}
        status={status}
        isProfileLoading={isProfileLoading}
        profileError={profileError}
        editorEnabled={editorEnabled}
      />
    </HomeEditorProvider>
  );
}

type HomePageContentProps = {
  profile: ProfileResponse | null;
  status: string;
  isProfileLoading: boolean;
  profileError: unknown;
  editorEnabled: boolean;
};

function HomePageContent({
  profile,
  status,
  isProfileLoading,
  profileError,
  editorEnabled,
}: HomePageContentProps) {
  const { t, i18n } = useTranslation();
  const editor = useOptionalHomeEditorContext();

  const canonicalProfile = useMemo(
    () => getCanonicalProfile(i18n.language),
    [i18n.language],
  );
  const effectiveProfile = profile ?? canonicalProfile;

  const localeKey = useMemo(() => resolveLocaleKey(i18n.language), [i18n.language]);

  const previewHomeConfig = useMemo<HomePageConfig | null>(() => {
    if (!editorEnabled || !editor?.draft) {
      return null;
    }
    return draftToPreviewConfig(editor.draft, localeKey);
  }, [editorEnabled, editor?.draft, localeKey]);

  const homeConfig = previewHomeConfig ?? effectiveProfile.home ?? null;

  const quickLinks = useMemo<HomeQuickLink[]>(() => {
    const fallbackQuickLinks = canonicalProfile.home?.quickLinks ?? [];
    const activeQuickLinks =
      homeConfig?.quickLinks?.length && homeConfig.quickLinks[0] != null
        ? homeConfig.quickLinks
        : fallbackQuickLinks;

    if (!activeQuickLinks.length) {
      return [];
    }

    return [...activeQuickLinks]
      .sort((a, b) => a.sortOrder - b.sortOrder)
      .slice(0, QUICK_LINK_HERO_LIMIT);
  }, [canonicalProfile, homeConfig]);

  const chipGroups = useMemo<ChipGroup[]>(() => {
    const previewProfile =
      previewHomeConfig != null
        ? ({
            ...effectiveProfile,
            home: previewHomeConfig,
          } satisfies ProfileResponse)
        : effectiveProfile;
    return buildChipGroups(previewProfile, t("common.presentLabel"));
  }, [effectiveProfile, previewHomeConfig, t]);

  const featuredLinks = useMemo(
    () => getFeaturedLinks(effectiveProfile.socialLinks ?? []),
    [effectiveProfile.socialLinks],
  );

  const techSections = useMemo(
    () =>
      [...effectiveProfile.techSections].sort(
        (a, b) => a.sortOrder - b.sortOrder,
      ),
    [effectiveProfile.techSections],
  );

  const recentWork = useMemo(() => {
    return [...effectiveProfile.workHistory]
      .sort((a, b) => a.sortOrder - b.sortOrder)
      .slice(0, 3);
  }, [effectiveProfile.workHistory]);

  return (
    <>
      {editorEnabled && (
        <div className="mx-auto mb-10 w-full max-w-6xl px-4 sm:px-8 lg:px-12">
          <HomeEditorPanel />
        </div>
      )}
      <section className="mx-auto flex w-full max-w-6xl flex-1 flex-col gap-12 px-4 py-16 sm:px-8 lg:px-12">
        <header
          className="grid gap-8 lg:grid-cols-[2fr,1fr]"
          style={{ display: "grid", gap: "2rem" }}
        >
          <div
            className="flex flex-col gap-6"
            style={{ display: "flex", flexDirection: "column", rowGap: "1.5rem" }}
          >
            <p className="text-xs font-semibold uppercase tracking-[0.3em] text-slate-500 dark:text-slate-400">
              {homeConfig?.heroSubtitle ?? t("home.hero.tagline")}
            </p>
            <div
              className="flex flex-col gap-3"
              style={{ display: "flex", flexDirection: "column", rowGap: "0.75rem" }}
            >
              <h1 className="text-4xl font-bold text-slate-900 dark:text-slate-50 sm:text-5xl lg:text-6xl">
                {effectiveProfile.displayName}
              </h1>
              {effectiveProfile.headline ? (
                <p className="text-xl font-semibold text-slate-700 dark:text-slate-200">
                  {effectiveProfile.headline}
                </p>
              ) : null}
            </div>
            {effectiveProfile.summary ? (
              <p className="text-lg leading-relaxed text-slate-600 dark:text-slate-300 md:max-w-3xl">
                {effectiveProfile.summary}
              </p>
            ) : null}

            {quickLinks.length ? (
              <div className="relative" style={{ minHeight: "220px" }}>
                <div
                  className="grid gap-4 sm:grid-cols-2"
                  style={{
                    ...QUICK_LINK_GRID_BASE_STYLES,
                    position: "absolute",
                    inset: 0,
                  }}
                >
                  {quickLinks.map((link) => (
                    <a
                      key={link.id}
                      href={link.targetUrl}
                      className="group flex flex-col gap-2 rounded-2xl border border-slate-200 bg-white/80 p-5 text-left shadow-sm transition hover:border-sky-300 hover:bg-white dark:border-slate-800 dark:bg-slate-900/60 dark:hover:border-sky-500"
                      style={QUICK_LINK_CARD_BASE_STYLES}
                    >
                      <span className="text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-400">
                        {link.section.replace("_", " ")}
                      </span>
                      <span className="text-lg font-semibold text-slate-900 transition group-hover:text-sky-600 dark:text-slate-100 dark:group-hover:text-sky-300">
                        {link.label}
                      </span>
                      {link.description ? (
                        <span className="text-sm text-slate-600 dark:text-slate-300">
                          {link.description}
                        </span>
                      ) : null}
                      <span className="text-sm font-medium text-sky-600 transition group-hover:underline dark:text-sky-400">
                        {link.cta}
                      </span>
                    </a>
                  ))}
                </div>
              </div>
            ) : null}
          </div>

          <aside
            className="flex flex-col gap-6"
            style={{ display: "flex", flexDirection: "column", rowGap: "1.5rem" }}
          >
            <div className="rounded-2xl border border-slate-200 bg-white/90 p-6 shadow-sm backdrop-blur dark:border-slate-800 dark:bg-slate-900/60">
              <h2 className="text-sm font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-400">
                {t("home.health.title")}
              </h2>
              <p className="mt-4 text-2xl font-semibold text-emerald-500 dark:text-emerald-400">
                {status}
              </p>
              <p className="mt-4 text-sm text-slate-500 dark:text-slate-400">
                {t("home.health.caption")}
              </p>
            </div>

            <div className="flex flex-col gap-4 rounded-2xl border border-slate-200 bg-white/90 p-6 shadow-sm backdrop-blur dark:border-slate-800 dark:bg-slate-900/60">
              <h2 className="text-sm font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-400">
                {t("home.social.title")}
              </h2>
              <div className="flex flex-wrap gap-3">
                {isProfileLoading && featuredLinks.length === 0 ? (
                  <>
                    <span className="inline-flex h-9 w-28 animate-pulse rounded-full bg-slate-200 dark:bg-slate-700" />
                    <span className="inline-flex h-9 w-24 animate-pulse rounded-full bg-slate-200 dark:bg-slate-700" />
                  </>
                ) : null}
                {featuredLinks.map((link) => {
                  const Icon = getSocialIcon(link.provider);
                  return (
                    <a
                      key={link.id}
                      href={link.url}
                      target={link.provider === "email" ? "_self" : "_blank"}
                      rel={link.provider === "email" ? undefined : "noreferrer"}
                      className="inline-flex items-center gap-2 rounded-full border border-slate-300 px-4 py-2 text-sm font-medium text-slate-700 transition hover:border-sky-400 hover:text-sky-600 dark:border-slate-700 dark:text-slate-200 dark:hover:border-sky-400 dark:hover:text-sky-300"
                    >
                      <Icon aria-hidden className="h-4 w-4" />
                      <span>{link.label}</span>
                    </a>
                  );
                })}
              </div>
            </div>

            <div className="rounded-2xl border border-slate-200 bg-white/90 p-6 shadow-sm backdrop-blur dark:border-slate-800 dark:bg-slate-900/60">
              <h2 className="text-sm font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-400">
                {t("home.about.lastUpdated")}
              </h2>
              <p className="mt-2 text-sm text-slate-600 dark:text-slate-300">
                {new Date(effectiveProfile.updatedAt).toLocaleString(i18n.language, {
                  year: "numeric",
                  month: "short",
                  day: "numeric",
                })}
              </p>
            </div>
          </aside>
        </header>

        <div
          className="grid gap-8 lg:grid-cols-[2fr,1fr]"
          style={{ display: "grid", gap: "2rem" }}
        >
          <section
            className="flex flex-col gap-6"
            style={{ display: "flex", flexDirection: "column", rowGap: "1.5rem" }}
          >
            <div
              className="flex flex-col gap-2"
              style={{ display: "flex", flexDirection: "column", rowGap: "0.5rem" }}
            >
              <h2 className="text-lg font-semibold text-slate-900 dark:text-slate-100">
                {t("home.tech.title")}
              </h2>
              <p className="text-sm text-slate-600 dark:text-slate-300">
                {t("home.tech.description")}
              </p>
            </div>
            <div className="grid gap-4 md:grid-cols-2">
              {techSections.map((section) => (
                <article
                  key={section.id}
                  className="rounded-2xl border border-slate-200 bg-white/90 p-5 shadow-sm backdrop-blur transition hover:border-sky-300 dark:border-slate-800 dark:bg-slate-900/60 dark:hover:border-sky-500"
                >
                  <h3 className="text-sm font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-400">
                    {section.title}
                  </h3>
                  <div className="mt-3 flex flex-wrap gap-2">
                    {section.members.map((member) => (
                      <span
                        key={member.id}
                        className="inline-flex items-center gap-2 rounded-full border border-slate-300 px-3 py-1 text-xs font-medium text-slate-700 dark:border-slate-700 dark:text-slate-200"
                      >
                        <span>{member.tech.displayName}</span>
                        <span className="text-[10px] uppercase text-slate-500 dark:text-slate-400">
                          {member.tech.level}
                        </span>
                      </span>
                    ))}
                  </div>
                </article>
              ))}
            </div>
          </section>

          <aside
            className="flex flex-col gap-6"
            style={{ display: "flex", flexDirection: "column", rowGap: "1.5rem" }}
          >
            <article className="rounded-2xl border border-slate-200 bg-white/90 p-6 shadow-sm backdrop-blur dark:border-slate-800 dark:bg-slate-900/60">
              <h2 className="text-sm font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-400">
                {t("home.work.title")}
              </h2>
              <p className="mt-2 text-sm text-slate-600 dark:text-slate-300">
                {t("home.work.description")}
              </p>
              <ul className="mt-4 space-y-4">
                {isProfileLoading ? (
                  <li className="space-y-2">
                    <span className="block h-4 w-52 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
                    <span className="block h-3 w-24 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
                  </li>
                ) : null}
                {!isProfileLoading && !recentWork.length ? (
                  <li className="text-sm text-slate-500 dark:text-slate-400">
                    {t("home.work.empty")}
                  </li>
                ) : null}
                {recentWork.map((item) => (
                  <li
                    key={item.id}
                    className="rounded-xl border border-slate-200 p-4 dark:border-slate-700"
                  >
                    <div className="flex flex-col gap-1 sm:flex-row sm:items-center sm:justify-between">
                      <div>
                        <p className="text-sm font-semibold text-slate-900 dark:text-slate-100">
                          {item.organization}
                        </p>
                        <p className="text-sm text-slate-600 dark:text-slate-300">
                          {item.role}
                        </p>
                      </div>
                      <p className="text-xs uppercase tracking-wide text-slate-500 dark:text-slate-400">
                        {formatDateRange(
                          item.startedAt,
                          item.endedAt ?? undefined,
                          t("common.presentLabel"),
                        )}
                      </p>
                    </div>
                    {item.summary ? (
                      <p className="mt-2 text-sm text-slate-600 dark:text-slate-300">
                        {item.summary}
                      </p>
                    ) : null}
                  </li>
                ))}
              </ul>
            </article>

            <article className="rounded-2xl border border-slate-200 bg-white/90 p-6 shadow-sm backdrop-blur dark:border-slate-800 dark:bg-slate-900/60">
              <h2 className="text-sm font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-400">
                {t("home.affiliations.title")}
              </h2>
              {chipGroups.length ? (
                <ul className="mt-4 space-y-3">
                  {chipGroups.map((group) => (
                    <li key={group.id} className="space-y-2 rounded-xl border border-slate-200 p-4 dark:border-slate-700">
                      <p className="text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-400">
                        {group.label}
                      </p>
                      <div className="flex flex-wrap gap-2">
                        {group.chips.map((chip) => (
                          <span
                            key={chip.id}
                            className="inline-flex items-center gap-2 rounded-full border border-slate-300 px-3 py-1 text-xs font-medium text-slate-700 dark:border-slate-700 dark:text-slate-200"
                          >
                            <span>{chip.label}</span>
                            {chip.description ? (
                              <span className="text-[10px] uppercase text-slate-500 dark:text-slate-400">
                                {chip.description}
                              </span>
                            ) : null}
                          </span>
                        ))}
                      </div>
                    </li>
                  ))}
                </ul>
              ) : (
                <p className="mt-4 text-sm text-slate-500 dark:text-slate-400">
                  {t("home.affiliations.empty")}
                </p>
              )}
            </article>
          </aside>
        </div>

        <section className="grid gap-6 lg:grid-cols-2">
          {profileError ? (
            <article className="rounded-2xl border border-rose-300 bg-rose-50 p-6 text-sm text-rose-700 dark:border-rose-900 dark:bg-rose-950/40 dark:text-rose-300">
              {t("home.errors.profileLoadFailed")}
            </article>
          ) : null}
        </section>
      </section>
    </>
  );
}
