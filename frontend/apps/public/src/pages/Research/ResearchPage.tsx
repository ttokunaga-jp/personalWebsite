import { useMemo, useState } from "react";
import { useTranslation } from "react-i18next";

import { MarkdownRenderer } from "../../components/molecules/MarkdownRenderer";
import { useResearchResource } from "../../modules/public-api";
import { formatDate } from "../../utils/date";

export function ResearchPage() {
  const { t } = useTranslation();
  const { data: researchEntries, isLoading, error } = useResearchResource();
  const [selectedTag, setSelectedTag] = useState<string | null>(null);

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
    if (!researchEntries) {
      return [];
    }

    if (!selectedTag) {
      return researchEntries;
    }

    return researchEntries.filter((entry) => entry.tags.includes(selectedTag));
  }, [researchEntries, selectedTag]);

  return (
    <section
      id="research"
      className="mx-auto flex w-full max-w-4xl flex-col gap-6 px-4 py-12 sm:px-8"
    >
      <header className="space-y-3">
        <p className="text-sm font-semibold uppercase tracking-[0.3em] text-sky-500 dark:text-sky-400">
          {t("research.tagline")}
        </p>
        <h1 className="text-3xl font-bold text-slate-900 dark:text-slate-50 sm:text-4xl">
          {t("research.title")}
        </h1>
        <p className="text-base text-slate-600 dark:text-slate-300">
          {t("research.description")}
        </p>
      </header>
      {tags.length ? (
        <div className="flex flex-wrap gap-2">
          <button
            type="button"
            onClick={() => setSelectedTag(null)}
            className={`rounded-full border px-3 py-1 text-xs font-medium transition ${
              selectedTag === null
                ? "border-sky-500 bg-sky-500 text-white"
                : "border-slate-300 text-slate-700 hover:border-sky-400 hover:text-sky-600 dark:border-slate-700 dark:text-slate-200 dark:hover:border-sky-400 dark:hover:text-sky-300"
            }`}
          >
            {t("research.filters.all")}
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
                    ? "border-sky-500 bg-sky-500 text-white"
                    : "border-slate-300 text-slate-700 hover:border-sky-400 hover:text-sky-600 dark:border-slate-700 dark:text-slate-200 dark:hover:border-sky-400 dark:hover:text-sky-300"
                }`}
                aria-pressed={isActive}
              >
                {tag}
              </button>
            );
          })}
        </div>
      ) : null}
      <div className="flex flex-col gap-6">
        {isLoading && (
          <article className="rounded-xl border border-slate-200 bg-white/80 p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900/60">
            <div className="space-y-3">
              <span className="block h-6 w-48 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
              <span className="block h-4 w-full animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
              <span className="block h-4 w-4/5 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
            </div>
          </article>
        )}
        {!isLoading && !filteredEntries.length ? (
          <article className="rounded-xl border border-slate-200 bg-white/80 p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900/60">
            <p className="text-sm text-slate-600 dark:text-slate-300">
              {selectedTag ? t("research.noEntriesForTag", { tag: selectedTag }) : t("research.placeholder")}
            </p>
          </article>
        ) : null}
        {filteredEntries.map((entry) => (
          <article
            key={entry.id}
            className="rounded-xl border border-slate-200 bg-white/80 p-6 shadow-sm transition hover:border-sky-300 dark:border-slate-800 dark:bg-slate-900/60 dark:hover:border-sky-500"
          >
            <header className="space-y-2">
              <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
                <h2 className="text-2xl font-semibold text-slate-900 dark:text-slate-100">
                  {entry.title}
                </h2>
                <p className="text-xs uppercase tracking-wide text-slate-500 dark:text-slate-400">
                  {formatDate(entry.publishedOn)}
                  {entry.updatedOn && entry.updatedOn !== entry.publishedOn
                    ? ` · ${t("research.updatedOn", { date: formatDate(entry.updatedOn) })}`
                    : null}
                </p>
              </div>
              <p className="text-sm text-slate-600 dark:text-slate-300">{entry.summary}</p>
              {entry.tags.length ? (
                <ul className="flex flex-wrap gap-2">
                  {entry.tags.map((tag) => (
                    <li
                      key={tag}
                      className="inline-flex items-center rounded-full border border-slate-300 px-2 py-1 text-xs font-medium text-slate-700 dark:border-slate-700 dark:text-slate-200"
                    >
                      {tag}
                    </li>
                  ))}
                </ul>
              ) : null}
            </header>
            <div className="mt-4 space-y-4">
              <MarkdownRenderer
                markdown={entry.contentMarkdown}
                html={entry.contentHtml}
                className="markdown-content"
              />
              {entry.assets?.length ? (
                <div className="grid gap-4 md:grid-cols-2">
                  {entry.assets.map((asset) => (
                    <figure
                      key={asset.url}
                      className="overflow-hidden rounded-lg border border-slate-200 bg-slate-100 dark:border-slate-800 dark:bg-slate-900/40"
                    >
                      <img
                        src={asset.url}
                        alt={asset.alt}
                        loading="lazy"
                        className="h-48 w-full object-cover"
                      />
                      {asset.caption ? (
                        <figcaption className="px-3 py-2 text-xs text-slate-600 dark:text-slate-400">
                          {asset.caption}
                        </figcaption>
                      ) : null}
                    </figure>
                  ))}
                </div>
              ) : null}
              {entry.links?.length ? (
                <div className="flex flex-wrap gap-3">
                  {entry.links.map((link) => (
                    <a
                      key={`${entry.id}-${link.url}`}
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
          <p role="alert" className="rounded-lg border border-rose-200 bg-rose-50 p-4 text-sm text-rose-700 dark:border-rose-700 dark:bg-rose-950/50 dark:text-rose-300">
            {t("research.error")}
          </p>
        ) : null}
      </div>
    </section>
  );
}
