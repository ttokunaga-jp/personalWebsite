import { useMemo } from "react";
import { useTranslation } from "react-i18next";

import { useProfileResource } from "../../modules/public-api";
import { formatDateRange } from "../../utils/date";

export function ProfilePage() {
  const { t } = useTranslation();
  const { data: profile, isLoading, error } = useProfileResource();

  const sortedAffiliations = useMemo(() => {
    if (!profile) {
      return [];
    }
    return [...profile.affiliations].sort((a, b) => {
      return a.isCurrent === b.isCurrent
        ? new Date(b.startDate).getTime() - new Date(a.startDate).getTime()
        : a.isCurrent
          ? -1
          : 1;
    });
  }, [profile]);

  const sortedWorkHistory = useMemo(() => {
    if (!profile) {
      return [];
    }
    return [...profile.workHistory].sort(
      (a, b) => new Date(b.startDate).getTime() - new Date(a.startDate).getTime()
    );
  }, [profile]);

  return (
    <section className="mx-auto flex w-full max-w-4xl flex-col gap-6 px-4 py-12 sm:px-8">
      <header className="space-y-3">
        <p className="text-sm font-semibold uppercase tracking-[0.3em] text-slate-900 dark:text-slate-100">
          {t("profile.tagline")}
        </p>
        <h1 className="text-3xl font-bold text-slate-900 dark:text-slate-50 sm:text-4xl">
          {t("profile.title")}
        </h1>
        <p className="text-base text-slate-600 dark:text-slate-300">
          {t("profile.description")}
        </p>
      </header>
      <div className="grid gap-6 md:grid-cols-2">
        <article className="rounded-xl border border-slate-200 bg-white/80 p-5 shadow-sm dark:border-slate-800 dark:bg-slate-900/60">
          <h2 className="text-lg font-semibold text-slate-900 dark:text-slate-100">
            {t("profile.sections.affiliations.title")}
          </h2>
          <p className="mt-2 text-sm text-slate-600 dark:text-slate-300">
            {t("profile.sections.affiliations.description")}
          </p>
          <ul className="mt-4 space-y-4" aria-live="polite">
            {isLoading && (
              <li className="space-y-2">
                <span className="block h-4 w-40 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
                <span className="block h-3 w-32 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
              </li>
            )}
            {!isLoading && !sortedAffiliations.length ? (
              <li className="text-sm text-slate-500 dark:text-slate-400">
                {t("profile.sections.affiliations.empty")}
              </li>
            ) : null}
            {sortedAffiliations.map((affiliation) => (
              <li key={affiliation.id} className="rounded-lg border border-slate-200 p-3 dark:border-slate-700">
                <p className="text-sm font-semibold text-slate-900 dark:text-slate-100">
                  {affiliation.organization}
                  {affiliation.department ? (
                    <span className="text-slate-500 dark:text-slate-400">
                      {" "}
                      · {affiliation.department}
                    </span>
                  ) : null}
                </p>
                <p className="text-sm text-slate-600 dark:text-slate-300">{affiliation.role}</p>
                <p className="text-xs uppercase tracking-wide text-slate-500 dark:text-slate-400">
                  {formatDateRange(
                    affiliation.startDate,
                    affiliation.endDate,
                    t("common.presentLabel")
                  )}
                </p>
                {affiliation.location ? (
                  <p className="text-xs text-slate-500 dark:text-slate-400">{affiliation.location}</p>
                ) : null}
              </li>
            ))}
          </ul>
        </article>
        <article className="rounded-xl border border-slate-200 bg-white/80 p-5 shadow-sm dark:border-slate-800 dark:bg-slate-900/60">
          <h2 className="text-lg font-semibold text-slate-900 dark:text-slate-100">
            {t("profile.sections.skills.title")}
          </h2>
          <p className="mt-2 text-sm text-slate-600 dark:text-slate-300">
            {t("profile.sections.skills.description")}
          </p>
          <div className="mt-4 flex flex-col gap-4">
            {isLoading && (
              <div className="space-y-2">
                <span className="inline-block h-4 w-32 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
                <div className="flex flex-wrap gap-2">
                  <span className="inline-block h-6 w-16 animate-pulse rounded-full bg-slate-200 dark:bg-slate-700" />
                  <span className="inline-block h-6 w-20 animate-pulse rounded-full bg-slate-200 dark:bg-slate-700" />
                </div>
              </div>
            )}
            {!isLoading && !profile?.skillGroups?.length ? (
              <p className="text-sm text-slate-500 dark:text-slate-400">
                {t("profile.sections.skills.empty")}
              </p>
            ) : null}
            {profile?.skillGroups?.map((group) => (
              <div key={group.id}>
                <h3 className="text-sm font-semibold text-slate-900 dark:text-slate-100">{group.category}</h3>
                <ul className="mt-2 flex flex-wrap gap-2">
                  {group.items.map((item) => (
                    <li
                      key={item.id}
                      className="inline-flex items-center gap-1 rounded-full border border-slate-300 px-3 py-1 text-xs font-medium text-slate-700 dark:border-slate-700 dark:text-slate-200"
                    >
                      <span>{item.name}</span>
                      <span className="text-[10px] uppercase text-slate-500 dark:text-slate-400">
                        {t(`profile.sections.skills.level.${item.level}` as const)}
                      </span>
                    </li>
                  ))}
                </ul>
              </div>
            ))}
          </div>
        </article>
      </div>

      <article className="rounded-xl border border-slate-200 bg-white/80 p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900/60">
        <header className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h2 className="text-lg font-semibold text-slate-900 dark:text-slate-100">
              {t("profile.sections.lab.title")}
            </h2>
            <p className="text-sm text-slate-600 dark:text-slate-300">
              {t("profile.sections.lab.description")}
            </p>
          </div>
          {profile?.lab?.websiteUrl ? (
            <a
              href={profile.lab.websiteUrl}
              target="_blank"
              rel="noreferrer"
              className="text-sm font-medium text-sky-600 underline decoration-sky-200 underline-offset-4 transition hover:text-sky-500 dark:text-sky-400 dark:hover:text-sky-300"
            >
              {t("profile.sections.lab.visit")}
            </a>
          ) : null}
        </header>
        <dl className="mt-4 grid gap-4 sm:grid-cols-2">
          {isLoading ? (
            <>
              <div>
                <span className="block h-4 w-32 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
              </div>
              <div>
                <span className="block h-4 w-24 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
              </div>
            </>
          ) : null}
          {!isLoading && profile?.lab ? (
            <>
              <div>
                <dt className="text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-400">
                  {t("profile.sections.lab.name")}
                </dt>
                <dd className="mt-1 text-sm text-slate-700 dark:text-slate-200">{profile.lab.name}</dd>
              </div>
              {profile.lab.advisor ? (
                <div>
                  <dt className="text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-400">
                    {t("profile.sections.lab.advisor")}
                  </dt>
                  <dd className="mt-1 text-sm text-slate-700 dark:text-slate-200">{profile.lab.advisor}</dd>
                </div>
              ) : null}
              {profile.lab.researchFocus ? (
                <div className="sm:col-span-2">
                  <dt className="text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-400">
                    {t("profile.sections.lab.focus")}
                  </dt>
                  <dd className="mt-1 text-sm text-slate-700 dark:text-slate-200">
                    {profile.lab.researchFocus}
                  </dd>
                </div>
              ) : null}
            </>
          ) : null}
          {!isLoading && !profile?.lab ? (
            <div className="sm:col-span-2 text-sm text-slate-500 dark:text-slate-400">
              {t("profile.sections.lab.empty")}
            </div>
          ) : null}
        </dl>
      </article>

      <article className="rounded-xl border border-slate-200 bg-white/80 p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900/60">
        <h2 className="text-lg font-semibold text-slate-900 dark:text-slate-100">
          {t("profile.sections.work.title")}
        </h2>
        <p className="text-sm text-slate-600 dark:text-slate-300">
          {t("profile.sections.work.description")}
        </p>
        <ul className="mt-4 space-y-4">
          {isLoading && (
            <li className="space-y-2">
              <span className="block h-4 w-44 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
              <span className="block h-3 w-36 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
            </li>
          )}
          {!isLoading && !sortedWorkHistory.length ? (
            <li className="text-sm text-slate-500 dark:text-slate-400">
              {t("profile.sections.work.empty")}
            </li>
          ) : null}
          {sortedWorkHistory.map((item) => (
            <li
              key={item.id}
              className="rounded-lg border border-slate-200 p-4 transition hover:border-sky-300 dark:border-slate-700 dark:hover:border-sky-500"
            >
              <div className="flex flex-col gap-2 sm:flex-row sm:justify-between">
                <div>
                  <p className="text-sm font-semibold text-slate-900 dark:text-slate-100">
                    {item.role}
                  </p>
                  <p className="text-sm text-slate-600 dark:text-slate-300">{item.organization}</p>
                </div>
                <p className="text-xs uppercase tracking-wide text-slate-500 dark:text-slate-400">
                  {formatDateRange(item.startDate, item.endDate, t("common.presentLabel"))}
                </p>
              </div>
              {item.location ? (
                <p className="text-xs text-slate-500 dark:text-slate-400">{item.location}</p>
              ) : null}
              {item.description ? (
                <p className="mt-2 text-sm text-slate-600 dark:text-slate-300">{item.description}</p>
              ) : null}
              {item.achievements?.length ? (
                <ul className="mt-2 list-disc space-y-1 pl-5 text-sm text-slate-600 dark:text-slate-300">
                  {item.achievements.map((achievement) => (
                    <li key={achievement}>{achievement}</li>
                  ))}
                </ul>
              ) : null}
            </li>
          ))}
        </ul>
      </article>

      <article className="rounded-xl border border-slate-200 bg-white/80 p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900/60">
        <h2 className="text-lg font-semibold text-slate-900 dark:text-slate-100">
          {t("profile.sections.communities.title")}
        </h2>
        <p className="text-sm text-slate-600 dark:text-slate-300">
          {t("profile.sections.communities.description")}
        </p>
        <div className="mt-4 flex flex-wrap gap-2">
          {isLoading && (
            <>
              <span className="inline-block h-6 w-16 animate-pulse rounded-full bg-slate-200 dark:bg-slate-700" />
              <span className="inline-block h-6 w-24 animate-pulse rounded-full bg-slate-200 dark:bg-slate-700" />
            </>
          )}
          {!isLoading && !profile?.communities?.length ? (
            <span className="text-sm text-slate-500 dark:text-slate-400">
              {t("profile.sections.communities.empty")}
            </span>
          ) : null}
          {profile?.communities?.map((community) => (
            <span
              key={community}
              className="inline-flex items-center rounded-full border border-slate-300 px-3 py-1 text-xs font-medium text-slate-700 dark:border-slate-700 dark:text-slate-200"
            >
              {community}
            </span>
          ))}
        </div>
      </article>

      {error ? (
        <div role="alert" className="rounded-xl border border-rose-200 bg-rose-50 p-4 text-sm text-rose-700 dark:border-rose-700 dark:bg-rose-950/50 dark:text-rose-300">
          {t("profile.error")}
        </div>
      ) : null}
    </section>
  );
}
