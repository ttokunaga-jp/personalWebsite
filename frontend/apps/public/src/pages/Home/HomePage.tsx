import { apiClient } from "@shared/lib/api-client";
import { startTransition, useEffect, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";

import { getCanonicalProfile } from "../../modules/profile-content";
import { useProfileResource } from "../../modules/public-api";
import type { HomeQuickLink, ProfileResponse, SocialLink } from "../../modules/public-api";
import { formatDateRange } from "../../utils/date";
import { getSocialIcon } from "../../utils/icons";

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

export function HomePage() {
  const { t, i18n } = useTranslation();
  const {
    data: profile,
    isLoading: isProfileLoading,
    error: profileError,
  } = useProfileResource();
  const [status, setStatus] = useState<string>("loading");

  const canonicalProfile = useMemo(
    () => getCanonicalProfile(i18n.language),
    [i18n.language],
  );
  const effectiveProfile = profile ?? canonicalProfile;
  const homeConfig = effectiveProfile.home;

  useEffect(() => {
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

  const quickLinks = useMemo<HomeQuickLink[]>(() => {
    if (!homeConfig?.quickLinks?.length) {
      return [];
    }
    return [...homeConfig.quickLinks].sort(
      (a, b) => a.sortOrder - b.sortOrder,
    );
  }, [homeConfig]);

  const chipGroups = useMemo<ChipGroup[]>(() => {
    return buildChipGroups(effectiveProfile, t("common.presentLabel"));
  }, [effectiveProfile, t]);

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
    <section className="mx-auto flex w-full max-w-6xl flex-1 flex-col gap-12 px-4 py-16 sm:px-8 lg:px-12">
      <header className="grid gap-8 lg:grid-cols-[2fr,1fr]">
        <div className="space-y-6">
          <p className="text-xs font-semibold uppercase tracking-[0.3em] text-slate-500 dark:text-slate-400">
            {homeConfig?.heroSubtitle ?? t("home.hero.tagline")}
          </p>
          <div className="space-y-3">
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
            <div className="grid gap-4 sm:grid-cols-2">
              {quickLinks.map((link) => (
                <a
                  key={link.id}
                  href={link.targetUrl}
                  className="group flex flex-col gap-2 rounded-2xl border border-slate-200 bg-white/80 p-5 text-left shadow-sm transition hover:border-sky-300 hover:bg-white dark:border-slate-800 dark:bg-slate-900/60 dark:hover:border-sky-500"
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
          ) : null}
        </div>

        <aside className="flex flex-col gap-6">
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

      <div className="grid gap-8 lg:grid-cols-[2fr,1fr]">
        <section className="space-y-6">
          <div className="space-y-2">
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

          {recentWork.length ? (
            <div className="space-y-4">
              <h2 className="text-lg font-semibold text-slate-900 dark:text-slate-100">
                {t("home.work.title")}
              </h2>
              <ul className="space-y-4">
                {recentWork.map((work) => (
                  <li
                    key={work.id}
                    className="rounded-2xl border border-slate-200 bg-white/90 p-5 shadow-sm dark:border-slate-800 dark:bg-slate-900/60"
                  >
                    <div className="flex flex-col gap-1 sm:flex-row sm:items-center sm:justify-between">
                      <div>
                        <p className="text-sm font-semibold text-slate-900 dark:text-slate-100">
                          {work.organization}
                        </p>
                        <p className="text-sm text-slate-600 dark:text-slate-300">
                          {work.role}
                        </p>
                      </div>
                      <p className="text-xs uppercase tracking-wide text-slate-500 dark:text-slate-400">
                        {formatDateRange(
                          work.startedAt,
                          work.endedAt ?? undefined,
                          t("common.presentLabel"),
                        )}
                      </p>
                    </div>
                    {work.summary ? (
                      <p className="mt-2 text-sm text-slate-600 dark:text-slate-300">
                        {work.summary}
                      </p>
                    ) : null}
                  </li>
                ))}
              </ul>
            </div>
          ) : null}
        </section>

        <aside className="flex flex-col gap-6">
          {chipGroups.map((group) => (
            <section
              key={group.id}
              className="rounded-2xl border border-slate-200 bg-white/90 p-6 shadow-sm backdrop-blur dark:border-slate-800 dark:bg-slate-900/60"
            >
              <h2 className="text-sm font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-400">
                {group.label}
              </h2>
              <div className="mt-3 flex flex-wrap gap-2">
                {group.chips.length ? (
                  group.chips.map((chip) => (
                    <span
                      key={chip.id}
                      className="inline-flex flex-col rounded-xl border border-slate-200 px-3 py-2 text-xs text-slate-700 dark:border-slate-700 dark:text-slate-200"
                    >
                      <span className="font-semibold">{chip.label}</span>
                      {chip.description ? (
                        <span className="text-[11px] text-slate-500 dark:text-slate-400">
                          {chip.description}
                        </span>
                      ) : null}
                    </span>
                  ))
                ) : (
                  <span className="text-sm text-slate-500 dark:text-slate-400">
                    {t("home.chips.empty")}
                  </span>
                )}
              </div>
            </section>
          ))}

          {effectiveProfile.affiliations.length ? (
            <section className="rounded-2xl border border-slate-200 bg-white/90 p-6 shadow-sm backdrop-blur dark:border-slate-800 dark:bg-slate-900/60">
              <h2 className="text-sm font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-400">
                {t("home.affiliations.title")}
              </h2>
              <ul className="mt-3 space-y-3">
                {effectiveProfile.affiliations
                  .slice(0, 4)
                  .map((affiliation) => (
                    <li key={affiliation.id}>
                      <p className="text-sm font-semibold text-slate-900 dark:text-slate-100">
                        {affiliation.name}
                      </p>
                      {affiliation.description ? (
                        <p className="text-xs text-slate-500 dark:text-slate-400">
                          {affiliation.description}
                        </p>
                      ) : null}
                      <p className="text-xs uppercase tracking-wide text-slate-500 dark:text-slate-400">
                        {formatDateRange(
                          affiliation.startedAt,
                          null,
                          t("common.presentLabel"),
                        )}
                      </p>
                    </li>
                  ))}
              </ul>
            </section>
          ) : null}

          {profileError ? (
            <p
              role="alert"
              className="rounded-xl border border-rose-200 bg-rose-50 p-4 text-sm text-rose-700 dark:border-rose-900 dark:bg-rose-950/40 dark:text-rose-300"
            >
              {t("home.about.error")}
            </p>
          ) : null}
        </aside>
      </div>
    </section>
  );
}
