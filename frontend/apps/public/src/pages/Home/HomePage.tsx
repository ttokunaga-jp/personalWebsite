import { apiClient } from "@shared/lib/api-client";
import { startTransition, useEffect, useMemo, useState } from "react";
import { Trans, useTranslation } from "react-i18next";

import { useProfileResource } from "../../modules/public-api";
import type { SocialLink } from "../../modules/public-api";
import { formatDateRange } from "../../utils/date";
import { getSocialIcon } from "../../utils/icons";

type HealthResponse = {
  status: string;
};

export function HomePage() {
  const { t } = useTranslation();
  const { data: profile, isLoading: isProfileLoading, error: profileError } = useProfileResource();
  const [status, setStatus] = useState<string>("loading");

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

  const primaryAffiliation = useMemo(() => {
    const affiliations = profile?.affiliations ?? [];
    if (!affiliations.length) {
      return null;
    }
    return (
      affiliations.find((affiliation) => affiliation.isCurrent) ?? affiliations.at(0) ?? null
    );
  }, [profile]);

  const featuredLinks = useMemo<SocialLink[]>(() => {
    if (!profile?.socialLinks) {
      return [];
    }

    return profile.socialLinks.filter((link) =>
      ["github", "x", "twitter", "linkedin", "email", "website"].includes(link.platform)
    );
  }, [profile]);

  const showSocialSkeleton = isProfileLoading && featuredLinks.length === 0;

  return (
    <section className="mx-auto flex w-full max-w-6xl flex-1 flex-col gap-12 px-4 py-16 sm:px-8 lg:px-12">
      <header className="flex flex-col gap-6 text-center md:text-left">
        <p className="text-sm font-semibold uppercase tracking-[0.3em] text-slate-900 dark:text-slate-100">
          {t("home.hero.tagline")}
        </p>
        <h1 className="text-4xl font-bold text-slate-900 dark:text-slate-50 sm:text-5xl lg:text-6xl">
          {profile?.headline ?? t("home.hero.title")}
        </h1>
        <p className="text-lg leading-relaxed text-slate-600 dark:text-slate-300 md:max-w-3xl md:text-left md:leading-relaxed">
          {profile?.summary ? profile.summary : <Trans i18nKey="home.hero.description" />}
        </p>
      </header>

      <div className="grid gap-8 lg:grid-cols-[2fr,1fr]">
        <div className="flex flex-col gap-6 rounded-2xl border border-slate-200 bg-white/90 p-6 shadow-sm backdrop-blur dark:border-slate-800 dark:bg-slate-900/60">
          <h2 className="text-sm font-semibold uppercase tracking-wider text-slate-500 dark:text-slate-400">
            {t("home.about.title")}
          </h2>
          <div className="grid gap-4 sm:grid-cols-2">
            <div>
              <p className="text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-400">
                {t("home.about.name")}
              </p>
              <p className="mt-1 text-lg font-semibold text-slate-900 dark:text-slate-100">
                {profile?.name ?? t("home.about.fallbackName")}
              </p>
            </div>
            <div>
              <p className="text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-400">
                {t("home.about.location")}
              </p>
              <p className="mt-1 text-lg font-semibold text-slate-900 dark:text-slate-100">
                {profile?.location ?? t("home.about.fallbackLocation")}
              </p>
            </div>
          </div>
          <div>
            <p className="text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-400">
              {t("home.about.affiliation")}
            </p>
            <div className="mt-1 text-sm text-slate-600 dark:text-slate-300">
              {isProfileLoading && !primaryAffiliation ? (
                <span className="inline-flex h-4 w-24 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
              ) : primaryAffiliation ? (
                <>
                  <span className="font-semibold text-slate-900 dark:text-slate-100">
                    {primaryAffiliation.organization}
                  </span>
                  {primaryAffiliation.department ? ` · ${primaryAffiliation.department}` : ""}
                  <br />
                  <span>{primaryAffiliation.role}</span>
                  <br />
                  <span className="text-xs uppercase tracking-wide text-slate-500 dark:text-slate-400">
                    {formatDateRange(
                      primaryAffiliation.startDate,
                      primaryAffiliation.endDate,
                      t("common.presentLabel")
                    )}
                  </span>
                </>
              ) : (
                t("home.about.affiliationFallback")
              )}
            </div>
          </div>
          <div>
            <p className="text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-400">
              {t("home.about.communities")}
            </p>
            <div className="mt-2 flex flex-wrap gap-2">
              {isProfileLoading && !profile && (
                <>
                  <span className="inline-flex h-6 w-16 animate-pulse rounded-full bg-slate-200 dark:bg-slate-700" />
                  <span className="inline-flex h-6 w-20 animate-pulse rounded-full bg-slate-200 dark:bg-slate-700" />
                </>
              )}
              {profile?.communities?.map((community) => (
                <span
                  key={community}
                  className="inline-flex items-center rounded-full border border-slate-300 px-3 py-1 text-xs font-medium text-slate-700 dark:border-slate-700 dark:text-slate-200"
                >
                  {community}
                </span>
              ))}
              {!isProfileLoading && !profile?.communities?.length ? (
                <span className="text-sm text-slate-500 dark:text-slate-400">
                  {t("home.about.communitiesFallback")}
                </span>
              ) : null}
            </div>
          </div>
          {profileError ? (
            <p role="alert" className="text-sm text-rose-500 dark:text-rose-400">
              {t("home.about.error")}
            </p>
          ) : null}
        </div>

        <aside className="flex flex-col gap-6">
          <div className="rounded-2xl border border-slate-200 bg-white/90 p-6 shadow-sm backdrop-blur dark:border-slate-800 dark:bg-slate-900/60">
            <h2 className="text-sm font-semibold uppercase tracking-wider text-slate-500 dark:text-slate-400">
              {t("home.health.title")}
            </h2>
            <p className="mt-2 text-2xl font-semibold text-emerald-500 dark:text-emerald-400">
              {status}
            </p>
            <p className="mt-4 text-sm text-slate-500 dark:text-slate-400">
              {t("home.health.caption")}
            </p>
          </div>

          <div className="flex flex-col gap-4 rounded-2xl border border-slate-200 bg-white/90 p-6 shadow-sm backdrop-blur dark:border-slate-800 dark:bg-slate-900/60">
            <h2 className="text-sm font-semibold uppercase tracking-wider text-slate-500 dark:text-slate-400">
              {t("home.social.title")}
            </h2>
            <div className="flex min-h-[3.25rem] flex-wrap gap-3">
              {showSocialSkeleton
                ? Array.from({ length: 3 }).map((_, index) => (
                    <span
                      key={`social-skeleton-${index}`}
                      className="inline-flex h-9 w-28 animate-pulse rounded-full bg-slate-200 dark:bg-slate-700"
                    />
                  ))
                : null}
              {featuredLinks.map((link) => {
                const Icon = getSocialIcon(link.platform);
                const label = `${t("home.social.connectWith", { label: link.label })}`;
                return (
                  <a
                    key={link.id}
                    href={link.url}
                    target={link.platform === "email" ? "_self" : "_blank"}
                    rel={link.platform === "email" ? undefined : "noreferrer"}
                    aria-label={label}
                    title={label}
                    className="inline-flex items-center gap-2 rounded-full border border-slate-300 px-4 py-2 text-sm font-medium text-slate-700 transition hover:border-sky-400 hover:text-sky-600 dark:border-slate-700 dark:text-slate-200 dark:hover:border-sky-400 dark:hover:text-sky-300"
                  >
                    <Icon aria-hidden className="h-4 w-4" />
                    <span>{link.label}</span>
                  </a>
                );
              })}
              {!featuredLinks.length && !isProfileLoading ? (
                <p className="text-sm text-slate-500 dark:text-slate-400">
                  {t("home.social.placeholder")}
                </p>
              ) : null}
            </div>
          </div>

          <div className="flex flex-col gap-4 rounded-2xl border border-slate-200 bg-white/90 p-6 shadow-sm backdrop-blur dark:border-slate-800 dark:bg-slate-900/60">
            <h2 className="text-sm font-semibold uppercase tracking-wider text-slate-500 dark:text-slate-400">
              {t("home.quickLinks.title")}
            </h2>
            <div className="flex flex-wrap gap-3">
              <a
                href="#projects"
                className="rounded-full border border-slate-300 px-4 py-2 text-sm font-medium text-slate-700 transition hover:border-sky-400 hover:text-sky-600 dark:border-slate-700 dark:text-slate-200 dark:hover:border-sky-400 dark:hover:text-sky-300"
              >
                {t("home.quickLinks.projects")}
              </a>
              <a
                href="#research"
                className="rounded-full border border-slate-300 px-4 py-2 text-sm font-medium text-slate-700 transition hover:border-sky-400 hover:text-sky-600 dark:border-slate-700 dark:text-slate-200 dark:hover:border-sky-400 dark:hover:text-sky-300"
              >
                {t("home.quickLinks.research")}
              </a>
              <a
                href="#contact"
                className="rounded-full border border-slate-300 px-4 py-2 text-sm font-medium text-slate-700 transition hover:border-sky-400 hover:text-sky-600 dark:border-slate-700 dark:text-slate-200 dark:hover:border-sky-400 dark:hover:text-sky-300"
              >
                {t("home.quickLinks.contact")}
              </a>
            </div>
            <p className="text-sm text-slate-500 dark:text-slate-400">
              {t("home.quickLinks.supporting")}
            </p>
          </div>
        </aside>
      </div>
    </section>
  );
}
