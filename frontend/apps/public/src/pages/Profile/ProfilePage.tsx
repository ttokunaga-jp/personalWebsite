import { useMemo } from "react";
import { useTranslation } from "react-i18next";

import { getCanonicalProfile } from "../../modules/profile-content";
import { useProfileResource } from "../../modules/public-api";
import { formatDateRange } from "../../utils/date";
import { getSocialIcon } from "../../utils/icons";

export function ProfilePage() {
  const { t, i18n } = useTranslation();
  const { data: profile, isLoading, error } = useProfileResource();

  const effectiveProfile = profile ?? getCanonicalProfile(i18n.language);

  const affiliations = useMemo(
    () =>
      [...effectiveProfile.affiliations].sort(
        (a, b) => a.sortOrder - b.sortOrder,
      ),
    [effectiveProfile.affiliations],
  );

  const communities = useMemo(
    () =>
      [...effectiveProfile.communities].sort(
        (a, b) => a.sortOrder - b.sortOrder,
      ),
    [effectiveProfile.communities],
  );

  const workHistory = useMemo(
    () =>
      [...effectiveProfile.workHistory].sort(
        (a, b) => a.sortOrder - b.sortOrder,
      ),
    [effectiveProfile.workHistory],
  );

  const techSections = useMemo(
    () =>
      [...effectiveProfile.techSections].sort(
        (a, b) => a.sortOrder - b.sortOrder,
      ),
    [effectiveProfile.techSections],
  );

  const socialLinks = useMemo(
    () =>
      [...effectiveProfile.socialLinks].sort(
        (a, b) => a.sortOrder - b.sortOrder,
      ),
    [effectiveProfile.socialLinks],
  );

  return (
    <section className="mx-auto flex w-full max-w-5xl flex-col gap-8 px-4 py-12 sm:px-8">
      <header className="space-y-3 text-center md:text-left">
        <p className="text-sm font-semibold uppercase tracking-[0.3em] text-slate-500 dark:text-slate-400">
          {t("profile.tagline")}
        </p>
        <h1 className="text-4xl font-bold text-slate-900 dark:text-slate-50 sm:text-5xl">
          {effectiveProfile.displayName}
        </h1>
        {effectiveProfile.headline ? (
          <p className="text-xl font-semibold text-slate-600 dark:text-slate-300">
            {effectiveProfile.headline}
          </p>
        ) : null}
        {effectiveProfile.summary ? (
          <p className="text-base text-slate-600 dark:text-slate-300">
            {effectiveProfile.summary}
          </p>
        ) : (
          <p className="text-base text-slate-600 dark:text-slate-300">
            {t("profile.description")}
          </p>
        )}
      </header>

      <div className="grid gap-6 lg:grid-cols-[2fr,1fr]">
        <div className="space-y-6">
          <article className="rounded-2xl border border-slate-200 bg-white/80 p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900/60">
            <h2 className="text-lg font-semibold text-slate-900 dark:text-slate-100">
              {t("profile.sections.affiliations.title")}
            </h2>
            <p className="mt-2 text-sm text-slate-600 dark:text-slate-300">
              {t("profile.sections.affiliations.description")}
            </p>
            <ul className="mt-4 space-y-4">
              {isLoading ? (
                <li className="space-y-2">
                  <span className="block h-4 w-40 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
                  <span className="block h-3 w-32 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
                </li>
              ) : null}
              {!isLoading && !affiliations.length ? (
                <li className="text-sm text-slate-500 dark:text-slate-400">
                  {t("profile.sections.affiliations.empty")}
                </li>
              ) : null}
              {affiliations.map((affiliation) => (
                <li
                  key={affiliation.id}
                  className="rounded-xl border border-slate-200 p-4 dark:border-slate-700"
                >
                  <p className="text-sm font-semibold text-slate-900 dark:text-slate-100">
                    {affiliation.name}
                  </p>
                  {affiliation.description ? (
                    <p className="text-sm text-slate-600 dark:text-slate-300">
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
          </article>

          <article className="rounded-2xl border border-slate-200 bg-white/80 p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900/60">
            <h2 className="text-lg font-semibold text-slate-900 dark:text-slate-100">
              {t("profile.sections.work.title")}
            </h2>
            <p className="mt-2 text-sm text-slate-600 dark:text-slate-300">
              {t("profile.sections.work.description")}
            </p>
            <ul className="mt-4 space-y-4">
              {isLoading ? (
                <li className="space-y-2">
                  <span className="block h-4 w-52 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
                  <span className="block h-3 w-24 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
                </li>
              ) : null}
              {!isLoading && !workHistory.length ? (
                <li className="text-sm text-slate-500 dark:text-slate-400">
                  {t("profile.sections.work.empty")}
                </li>
              ) : null}
              {workHistory.map((item) => (
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
        </div>

        <aside className="space-y-6">
          <article className="rounded-2xl border border-slate-200 bg-white/80 p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900/60">
            <h2 className="text-lg font-semibold text-slate-900 dark:text-slate-100">
              {t("profile.sections.communities.title")}
            </h2>
            <p className="mt-2 text-sm text-slate-600 dark:text-slate-300">
              {t("profile.sections.communities.description")}
            </p>
            <ul className="mt-4 space-y-3">
              {communities.map((community) => (
                <li key={community.id} className="rounded-xl border border-slate-200 p-3 dark:border-slate-700">
                  <p className="text-sm font-semibold text-slate-900 dark:text-slate-100">
                    {community.name}
                  </p>
                  {community.description ? (
                    <p className="text-xs text-slate-600 dark:text-slate-300">
                      {community.description}
                    </p>
                  ) : null}
                  <p className="text-xs uppercase tracking-wide text-slate-500 dark:text-slate-400">
                    {formatDateRange(
                      community.startedAt,
                      null,
                      t("common.presentLabel"),
                    )}
                  </p>
                </li>
              ))}
            </ul>
          </article>

          <article className="rounded-2xl border border-slate-200 bg-white/80 p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900/60">
            <h2 className="text-lg font-semibold text-slate-900 dark:text-slate-100">
              {t("profile.sections.skills.title")}
            </h2>
            <p className="mt-2 text-sm text-slate-600 dark:text-slate-300">
              {t("profile.sections.skills.description")}
            </p>
            <div className="mt-4 space-y-4">
              {techSections.map((section) => (
                <div key={section.id}>
                  <h3 className="text-sm font-semibold text-slate-900 dark:text-slate-100">
                    {section.title}
                  </h3>
                  <div className="mt-2 flex flex-wrap gap-2">
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
                </div>
              ))}
            </div>
          </article>

          <article className="rounded-2xl border border-slate-200 bg-white/80 p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900/60">
            <h2 className="text-lg font-semibold text-slate-900 dark:text-slate-100">
              {t("profile.sections.social.title")}
            </h2>
            <div className="mt-3 flex flex-wrap gap-3">
              {socialLinks.map((link) => {
                const Icon = getSocialIcon(link.provider);
                return (
                  <a
                    key={link.id}
                    href={link.url}
                    target={link.provider === "email" ? "_self" : "_blank"}
                    rel={link.provider === "email" ? undefined : "noreferrer"}
                    className="inline-flex items-center gap-2 rounded-full border border-slate-300 px-4 py-2 text-xs font-medium text-slate-700 transition hover:border-sky-400 hover:text-sky-600 dark:border-slate-700 dark:text-slate-200 dark:hover:border-sky-400 dark:hover:text-sky-300"
                  >
                    <Icon aria-hidden className="h-4 w-4" />
                    <span>{link.label}</span>
                  </a>
                );
              })}
            </div>
          </article>
        </aside>
      </div>

      <article className="rounded-2xl border border-slate-200 bg-white/80 p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900/60">
        <header className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h2 className="text-lg font-semibold text-slate-900 dark:text-slate-100">
              {t("profile.sections.lab.title")}
            </h2>
            <p className="text-sm text-slate-600 dark:text-slate-300">
              {t("profile.sections.lab.description")}
            </p>
          </div>
          {effectiveProfile.lab?.url ? (
            <a
              href={effectiveProfile.lab.url}
              target="_blank"
              rel="noreferrer"
              className="text-sm font-medium text-sky-600 underline decoration-sky-200 underline-offset-4 transition hover:text-sky-500 dark:text-sky-400 dark:hover:text-sky-300"
            >
              {t("profile.sections.lab.visit")}
            </a>
          ) : null}
        </header>
        <dl className="mt-4 grid gap-4 sm:grid-cols-2">
          {effectiveProfile.lab?.name ? (
            <div>
              <dt className="text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-400">
                {t("profile.sections.lab.name")}
              </dt>
              <dd className="mt-1 text-sm text-slate-700 dark:text-slate-200">
                {effectiveProfile.lab.name}
              </dd>
            </div>
          ) : null}
          {effectiveProfile.lab?.advisor ? (
            <div>
              <dt className="text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-400">
                {t("profile.sections.lab.advisor")}
              </dt>
              <dd className="mt-1 text-sm text-slate-700 dark:text-slate-200">
                {effectiveProfile.lab.advisor}
              </dd>
            </div>
          ) : null}
          {effectiveProfile.lab?.room ? (
            <div>
              <dt className="text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-400">
                {t("profile.sections.lab.room")}
              </dt>
              <dd className="mt-1 text-sm text-slate-700 dark:text-slate-200">
                {effectiveProfile.lab.room}
              </dd>
            </div>
          ) : null}
        </dl>
      </article>

      {error ? (
        <p
          role="alert"
          className="rounded-xl border border-rose-200 bg-rose-50 p-4 text-sm text-rose-700 dark:border-rose-900 dark:bg-rose-950/40 dark:text-rose-300"
        >
          {t("profile.error")}
        </p>
      ) : null}
    </section>
  );
}
