import { useMemo, useState } from "react";
import { useTranslation } from "react-i18next";

import { getCanonicalProjects } from "../../modules/profile-content";
import { useProjectsResource } from "../../modules/public-api";
import type { Project } from "../../modules/public-api";
import { formatDateRange } from "../../utils/date";

function extractTechLabel(project: Project): string[] {
  return project.tech.map((membership) => membership.tech.displayName).filter(Boolean);
}

export function ProjectsPage() {
  const { t, i18n } = useTranslation();
  const { data: projects, isLoading, error } = useProjectsResource();
  const canonicalProjects = useMemo(
    () => getCanonicalProjects(i18n.language),
    [i18n.language],
  );
  const effectiveProjects = projects ?? canonicalProjects;

  const techFilters = useMemo(() => {
    const set = new Set<string>();
    effectiveProjects.forEach((project) => {
      extractTechLabel(project).forEach((label) => set.add(label));
    });
    return Array.from(set).sort((a, b) => a.localeCompare(b));
  }, [effectiveProjects]);

  const [selectedTech, setSelectedTech] = useState<Set<string>>(new Set());

  const toggleTech = (label: string) => {
    setSelectedTech((prev) => {
      const next = new Set(prev);
      if (next.has(label)) {
        next.delete(label);
      } else {
        next.add(label);
      }
      return next;
    });
  };

  const filteredProjects = useMemo(() => {
    if (!effectiveProjects.length) {
      return [];
    }
    if (selectedTech.size === 0) {
      return effectiveProjects;
    }
    return effectiveProjects.filter((project) => {
      const techLabels = extractTechLabel(project);
      return techLabels.some((label) => selectedTech.has(label));
    });
  }, [effectiveProjects, selectedTech]);

  const highlightProjects = useMemo(
    () =>
      filteredProjects.filter((project) => project.highlight).sort(
        (a, b) => a.sortOrder - b.sortOrder,
      ),
    [filteredProjects],
  );

  const standardProjects = useMemo(
    () =>
      filteredProjects
        .filter((project) => !project.highlight)
        .sort((a, b) => a.sortOrder - b.sortOrder),
    [filteredProjects],
  );

  return (
    <section
      id="projects"
      className="mx-auto flex w-full max-w-5xl flex-col gap-6 px-4 py-12 sm:px-8"
    >
      <header className="space-y-3">
        <p className="text-sm font-semibold uppercase tracking-[0.3em] text-slate-500 dark:text-slate-400">
          {t("projects.tagline")}
        </p>
        <h1 className="text-3xl font-bold text-slate-900 dark:text-slate-50 sm:text-4xl">
          {t("projects.title")}
        </h1>
        <p className="text-base text-slate-600 dark:text-slate-300">
          {t("projects.description")}
        </p>
      </header>

  <div className="flex flex-wrap items-center gap-2">
        <button
          type="button"
          onClick={() => setSelectedTech(new Set())}
          className={`rounded-full border px-3 py-1 text-xs font-medium transition ${
            selectedTech.size === 0
              ? "border-sky-500 bg-sky-500 text-white"
              : "border-slate-300 text-slate-700 hover:border-sky-400 hover:text-sky-600 dark:border-slate-700 dark:text-slate-200 dark:hover:border-sky-400 dark:hover:text-sky-300"
          }`}
        >
          {t("projects.filters.all")}
        </button>
        {isLoading && !techFilters.length ? (
          <>
            <span className="inline-flex h-6 w-16 animate-pulse rounded-full bg-slate-200 dark:bg-slate-700" />
            <span className="inline-flex h-6 w-20 animate-pulse rounded-full bg-slate-200 dark:bg-slate-700" />
          </>
        ) : null}
        {techFilters.map((label) => {
          const isActive = selectedTech.has(label);
          return (
            <button
              key={label}
              type="button"
              onClick={() => toggleTech(label)}
              className={`rounded-full border px-3 py-1 text-xs font-medium transition ${
                isActive
                  ? "border-sky-500 bg-sky-500 text-white"
                  : "border-slate-300 text-slate-700 hover:border-sky-400 hover:text-sky-600 dark:border-slate-700 dark:text-slate-200 dark:hover:border-sky-400 dark:hover:text-sky-300"
              }`}
              aria-pressed={isActive}
            >
              {label}
            </button>
          );
        })}
      </div>

      {highlightProjects.length ? (
        <section className="space-y-4">
          <h2 className="text-lg font-semibold text-slate-900 dark:text-slate-100">
            {t("projects.highlight.title")}
          </h2>
          <div className="grid gap-4 lg:grid-cols-2">
            {highlightProjects.map((project) => (
              <ProjectCard
                key={project.id}
                project={project}
                emphasis
                presentLabel={t("common.presentLabel")}
              />
            ))}
          </div>
        </section>
      ) : null}

      <div className="grid gap-4">
        {!highlightProjects.length && !standardProjects.length && !isLoading ? (
          <article className="rounded-2xl border border-slate-200 bg-white/80 p-6 text-sm text-slate-600 shadow-sm dark:border-slate-800 dark:bg-slate-900/60 dark:text-slate-300">
            {selectedTech.size
              ? t("projects.noMatchesForSelection")
              : t("projects.placeholder")}
          </article>
        ) : null}
        {standardProjects.map((project) => (
          <ProjectCard
            key={project.id}
            project={project}
            presentLabel={t("common.presentLabel")}
          />
        ))}
        {error ? (
          <p
            role="alert"
            className="rounded-xl border border-rose-200 bg-rose-50 p-4 text-sm text-rose-700 dark:border-rose-900 dark:bg-rose-950/40 dark:text-rose-300"
          >
            {t("projects.error")}
          </p>
        ) : null}
      </div>
    </section>
  );
}

type ProjectCardProps = {
  project: Project;
  emphasis?: boolean;
  presentLabel: string;
};

function ProjectCard({ project, emphasis = false, presentLabel }: ProjectCardProps) {
  const techLabels = extractTechLabel(project);

  return (
    <article
      className={`rounded-2xl border p-6 shadow-sm transition ${
        emphasis
          ? "border-sky-300 bg-white dark:border-sky-500 dark:bg-slate-900/60"
          : "border-slate-200 bg-white/80 hover:border-sky-300 dark:border-slate-800 dark:bg-slate-900/60 dark:hover:border-sky-500"
      }`}
    >
      <header className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
        <div>
          <h3 className="text-xl font-semibold text-slate-900 dark:text-slate-100">
            {project.title}
          </h3>
          {project.summary ? (
            <p className="text-sm text-slate-600 dark:text-slate-300">
              {project.summary}
            </p>
          ) : null}
        </div>
        <p className="text-xs uppercase tracking-wide text-slate-500 dark:text-slate-400">
          {formatDateRange(
            project.period.start ?? "",
            project.period.end ?? undefined,
            presentLabel,
          )}
        </p>
      </header>

      {project.coverImageUrl ? (
        <div className="mt-4 aspect-video w-full overflow-hidden rounded-lg">
          <img
            src={project.coverImageUrl}
            alt=""
            loading="lazy"
            className="h-full w-full object-cover"
          />
        </div>
      ) : null}

      {project.description ? (
        <p className="mt-4 text-sm text-slate-600 dark:text-slate-300">
          {project.description}
        </p>
      ) : null}

      {techLabels.length ? (
        <div className="mt-4 flex flex-wrap gap-2">
          {project.tech.map((membership) => (
            <span
              key={membership.id}
              className="inline-flex items-center gap-2 rounded-full border border-slate-300 px-3 py-1 text-xs font-medium text-slate-700 dark:border-slate-700 dark:text-slate-200"
            >
              <span>{membership.tech.displayName}</span>
              <span className="text-[10px] uppercase text-slate-500 dark:text-slate-400">
                {membership.tech.level}
              </span>
            </span>
          ))}
        </div>
      ) : null}

      {project.links.length ? (
        <div className="mt-4 flex flex-wrap gap-3">
          {project.links.map((link) => (
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

      {project.primaryLink ? (
        <a
          href={project.primaryLink}
          target="_blank"
          rel="noreferrer"
          className="mt-4 inline-flex items-center text-sm font-medium text-sky-600 underline decoration-sky-200 underline-offset-4 transition hover:text-sky-500 dark:text-sky-400 dark:hover:text-sky-300"
        >
          {project.primaryLink}
        </a>
      ) : null}
    </article>
  );
}
