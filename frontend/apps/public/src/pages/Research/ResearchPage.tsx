import { useMemo, useState } from "react";
import { useTranslation } from "react-i18next";

import { useResearchResource } from "../../modules/public-api";
import { formatDate } from "../../utils/date";

type EntryKind = "all" | "research" | "blog";

const KIND_FILTERS: EntryKind[] = ["all", "research", "blog"];

export function ResearchPage() {
  const { t } = useTranslation();
  const { data: researchEntries, isLoading, error } = useResearchResource();
  const [selectedTag, setSelectedTag] = useState<string | null>(null);
  const [selectedKind, setSelectedKind] = useState<EntryKind>("all");

  const tags = useMemo(() => {
    if (!researchEntries?.length) {
      return [];
    }
    const tagSet = new Set<string>();
    for (const entry of researchEntries) {
      entry.tags.forEach((tag) => tagSet.add(tag));
    }
    return Array.from(tagSet).sort((a, b) => a.localeCompare(b));
  }, [researchEntries]);

  const filteredEntries = useMemo(() => {
    if (!researchEntries?.length) {
      return [];
    }

    return researchEntries.filter((entry) => {
      const matchesKind =
        selectedKind === "all" || entry.kind === selectedKind;
      const matchesTag =
        !selectedTag || entry.tags.includes(selectedTag);
      return matchesKind && matchesTag;
    });
  }, [researchEntries, selectedKind, selectedTag]);

  return (
    <section
      id="research"
      className="mx-auto flex w-full max-w-5xl flex-col gap-6 px-4 py-12 sm:px-8"
    >
      <header className="space-y-3">
        <p className="text-sm font-semibold uppercase tracking-[0.3em] text-slate-500 dark:text-slate-400">
          {t("research.tagline")}
        </p>
        <h1 className="text-3xl font-bold text-slate-900 dark:text-slate-50 sm:text-4xl">
          {t("research.title")}
        </h1>
        <p className="text-base text-slate-600 dark:text-slate-300">
          {t("research.description")}
        </p>
      </header>

      <div className="flex flex-wrap gap-3">
        {KIND_FILTERS.map((kind) => {
          const isActive = selectedKind === kind;
          return (
            <button
              key={kind}
              type="button"
              onClick={() => setSelectedKind(kind)}
              className={`rounded-full border px-3 py-1 text-xs font-medium transition ${
                isActive
                  ? "border-sky-500 bg-sky-500 text-white"
                  : "border-slate-300 text-slate-700 hover:border-sky-400 hover:text-sky-600 dark:border-slate-700 dark:text-slate-200 dark:hover:border-sky-400 dark:hover:text-sky-300"
              }`}
            >
              {t(`research.filters.${kind}` as const)}
            </button>
          );
        })}
        {tags.length ? (
          <>
            <span className="ml-2 inline-flex items-center text-xs uppercase tracking-wide text-slate-500 dark:text-slate-400">
              {t("research.filters.tags")}
            </span>
            <button
              type="button"
              onClick={() => setSelectedTag(null)}
              className={`rounded-full border px-3 py-1 text-xs font-medium transition ${
                selectedTag === null
                  ? "border-emerald-500 bg-emerald-500 text-white"
                  : "border-slate-300 text-slate-700 hover:border-emerald-400 hover:text-emerald-600 dark:border-slate-700 dark:text-slate-200 dark:hover:border-emerald-400 dark:hover:text-emerald-300"
              }`}
            >
              {t("research.filters.tagAll")}
            </button>
            {tags.map((tag) => {
              const isActive = selectedTag === tag;
              return (
                <button
                  key={tag}
                  type="button"
                  onClick={() => setSelectedTag(isActive ? null : tag)}
                  className={`rounded-full border px-3 py-1 text-xs font-medium transition ${
                    isActive
                      ? "border-emerald-500 bg-emerald-500 text-white"
                      : "border-slate-300 text-slate-700 hover:border-emerald-400 hover:text-emerald-600 dark:border-slate-700 dark:text-slate-200 dark:hover:border-emerald-400 dark:hover:text-emerald-300"
                  }`}
                  aria-pressed={isActive}
                >
                  {tag}
                </button>
              );
            })}
          </>
        ) : null}
      </div>

      <div className="flex flex-col gap-6">
        {isLoading ? (
          <article className="rounded-2xl border border-slate-200 bg-white/80 p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900/60">
            <div className="space-y-3">
              <span className="block h-6 w-48 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
              <span className="block h-4 w-full animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
              <span className="block h-4 w-4/5 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
            </div>
          </article>
        ) : null}

        {!isLoading && !filteredEntries.length ? (
          <article className="rounded-2xl border border-slate-200 bg-white/80 p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900/60">
            <p className="text-sm text-slate-600 dark:text-slate-300">
              {selectedTag
                ? t("research.noEntriesForTag", { tag: selectedTag })
                : t("research.placeholder")}
            </p>
          </article>
        ) : null}

        {filteredEntries.map((entry) => (
          <article
            key={entry.id}
            className="rounded-2xl border border-slate-200 bg-white/80 p-6 shadow-sm transition hover:border-sky-300 dark:border-slate-800 dark:bg-slate-900/60 dark:hover:border-sky-500"
          >
            <header className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
              <div className="space-y-2">
                <div className="flex items-center gap-3">
                  <span className="inline-flex items-center rounded-full border border-sky-300 px-2 py-1 text-[11px] font-semibold uppercase tracking-wide text-sky-600 dark:border-sky-500 dark:text-sky-300">
                    {t(`research.kind.${entry.kind}` as const)}
                  </span>
                  <p className="text-xs uppercase tracking-wide text-slate-500 dark:text-slate-400">
                    {formatDate(entry.publishedAt)}
                  </p>
                </div>
                <h2 className="text-2xl font-semibold text-slate-900 dark:text-slate-100">
                  {entry.title}
                </h2>
              </div>
              <a
                href={entry.externalUrl}
                target="_blank"
                rel="noreferrer"
                className="inline-flex items-center rounded-full border border-slate-300 px-4 py-2 text-xs font-semibold text-slate-700 transition hover:border-sky-400 hover:text-sky-600 dark:border-slate-700 dark:text-slate-200 dark:hover:border-sky-400 dark:hover:text-sky-300"
              >
                {t("research.externalLink")}
              </a>
            </header>

            <div className="mt-4 space-y-4">
              {entry.overview ? (
                <p className="text-sm text-slate-600 dark:text-slate-300">
                  {entry.overview}
                </p>
              ) : null}

              <div className="flex flex-wrap gap-2">
                {entry.tags.map((tag) => (
                  <span
                    key={tag}
                    className="inline-flex items-center rounded-full border border-slate-300 px-2 py-1 text-xs font-medium text-slate-700 dark:border-slate-700 dark:text-slate-200"
                  >
                    {tag}
                  </span>
                ))}
              </div>

              {entry.tech.length ? (
                <div className="flex flex-wrap gap-2">
                  {entry.tech.map((membership) => (
                    <span
                      key={membership.id}
                      className="inline-flex items-center gap-2 rounded-full border border-emerald-300 px-3 py-1 text-xs font-medium text-emerald-700 dark:border-emerald-700 dark:text-emerald-300"
                    >
                      <span>{membership.tech.displayName}</span>
                      <span className="text-[10px] uppercase text-emerald-600 dark:text-emerald-400">
                        {membership.tech.level}
                      </span>
                    </span>
                  ))}
                </div>
              ) : null}

              {entry.outcome ? (
                <section className="rounded-xl border border-slate-200 bg-white/60 p-4 dark:border-slate-800 dark:bg-slate-900/60">
                  <h3 className="text-sm font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-400">
                    {t("research.sections.outcome")}
                  </h3>
                  <p className="mt-1 text-sm text-slate-600 dark:text-slate-300">
                    {entry.outcome}
                  </p>
                </section>
              ) : null}

              {entry.outlook ? (
                <section className="rounded-xl border border-slate-200 bg-white/60 p-4 dark:border-slate-800 dark:bg-slate-900/60">
                  <h3 className="text-sm font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-400">
                    {t("research.sections.outlook")}
                  </h3>
                  <p className="mt-1 text-sm text-slate-600 dark:text-slate-300">
                    {entry.outlook}
                  </p>
                </section>
              ) : null}

              {entry.links.length ? (
                <div className="flex flex-wrap gap-3">
                  {entry.links.map((link) => (
                    <a
                      key={link.id}
                      href={link.url}
                      target="_blank"
                      rel="noreferrer"
                      className="inline-flex items-center rounded-full border border-slate-300 px-3 py-1 text-xs font-medium text-slate-700 transition hover:border-sky-400 hover:text-sky-600 dark:border-slate-700 dark:text-slate-200 dark:hover:border-sky-400 dark:hover:text-sky-300"
                    >
                      {link.label}
                    </a>
                  ))}
                </div>
              ) : null}
            </div>
          </article>
        ))}

        {error ? (
          <p
            role="alert"
            className="rounded-xl border border-rose-200 bg-rose-50 p-4 text-sm text-rose-700 dark:border-rose-900 dark:bg-rose-950/40 dark:text-rose-300"
          >
            {t("research.error")}
          </p>
        ) : null}
      </div>
    </section>
  );
}
