import { useMemo, useState } from "react";
import { useTranslation } from "react-i18next";

import { useProjectsResource } from "../../modules/public-api";
import type { Project } from "../../modules/public-api";
import { formatDateRange } from "../../utils/date";

const PROJECT_SKELETON_COUNT = 3;

function ProjectCardSkeleton() {
  return (
    <article className="rounded-xl border border-slate-200 bg-white/80 p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900/60">
      <header className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
        <div className="space-y-2">
          <span className="block h-6 w-56 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
          <span className="block h-4 w-40 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
          <span className="block h-5 w-24 animate-pulse rounded-full bg-slate-200 dark:bg-slate-700" />
        </div>
        <span className="block h-4 w-28 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
      </header>
      <span className="mt-4 block h-48 w-full animate-pulse rounded-lg bg-slate-200 dark:bg-slate-700" />
      <div className="mt-4 space-y-2">
        <span className="block h-4 w-full animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
        <span className="block h-4 w-5/6 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
        <span className="block h-4 w-4/6 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
      </div>
      <div className="mt-4 flex flex-wrap gap-2">
        {Array.from({ length: 4 }).map((_, index) => (
          // index is stable for skeleton placeholders
          // eslint-disable-next-line react/no-array-index-key
          <span
            key={index}
            className="inline-flex h-7 w-24 animate-pulse rounded-full bg-slate-200 dark:bg-slate-700"
          />
        ))}
      </div>
      <div className="mt-4 flex flex-wrap gap-2">
        {Array.from({ length: 3 }).map((_, index) => (
          // index is stable for skeleton placeholders
          // eslint-disable-next-line react/no-array-index-key
          <span
            key={index}
            className="inline-flex h-6 w-20 animate-pulse rounded-full border border-dashed border-slate-200 bg-transparent dark:border-slate-700"
          />
        ))}
      </div>
    </article>
  );
}

export function ProjectsPage() {
  const { t } = useTranslation();
  const { data: projects, isLoading, error } = useProjectsResource();
  const [selectedStacks, setSelectedStacks] = useState<Set<string>>(new Set());

  const techStacks = useMemo(() => {
    if (!projects?.length) {
      return [];
    }

    const stackSet = new Set<string>();
    projects.forEach((project) => {
      project.techStack.forEach((tech) => stackSet.add(tech));
    });

    return Array.from(stackSet).sort((a, b) => a.localeCompare(b));
  }, [projects]);

  const filteredProjects = useMemo(() => {
    if (!projects) {
      return [];
    }

    if (selectedStacks.size === 0) {
      return projects;
    }

    return projects.filter((project) =>
      project.techStack.some((tech) => selectedStacks.has(tech))
    );
  }, [projects, selectedStacks]);

  const toggleStack = (stack: string) => {
    setSelectedStacks((prev) => {
      const next = new Set(prev);
      if (next.has(stack)) {
        next.delete(stack);
      } else {
        next.add(stack);
      }
      return next;
    });
  };

  const renderProjectLinks = (project: Project) => {
    return project.links.map((link) => (
      <a
        key={`${project.id}-${link.url}`}
        href={link.url}
        target="_blank"
        rel="noreferrer"
        className="inline-flex items-center rounded-full border border-slate-300 px-3 py-1 text-xs font-medium text-slate-700 transition hover:border-sky-400 hover:text-sky-600 dark:border-slate-700 dark:text-slate-200 dark:hover:border-sky-400 dark:hover:text-sky-300"
      >
        {link.label}
      </a>
    ));
  };

  return (
    <section
      id="projects"
      className="mx-auto flex w-full max-w-5xl flex-col gap-6 px-4 py-12 sm:px-8"
    >
      <header className="space-y-3">
        <p className="text-sm font-semibold uppercase tracking-[0.3em] text-slate-900 dark:text-slate-100">
          {t("projects.tagline")}
        </p>
        <h1 className="text-3xl font-bold text-slate-900 dark:text-slate-50 sm:text-4xl">
          {t("projects.title")}
        </h1>
        <p className="text-base text-slate-600 dark:text-slate-300">
          {t("projects.description")}
        </p>
      </header>
      {techStacks.length ? (
        <div className="flex flex-wrap gap-2">
          <button
            type="button"
            onClick={() => setSelectedStacks(new Set())}
            className={`rounded-full border px-3 py-1 text-xs font-medium transition ${
              selectedStacks.size === 0
                ? "border-sky-500 bg-sky-500 text-white"
                : "border-slate-300 text-slate-700 hover:border-sky-400 hover:text-sky-600 dark:border-slate-700 dark:text-slate-200 dark:hover:border-sky-400 dark:hover:text-sky-300"
            }`}
          >
            {t("projects.filters.all")}
          </button>
          {techStacks.map((tech) => {
            const isActive = selectedStacks.has(tech);
            return (
              <button
                key={tech}
                type="button"
                onClick={() => toggleStack(tech)}
                className={`rounded-full border px-3 py-1 text-xs font-medium transition ${
                  isActive
                    ? "border-sky-500 bg-sky-500 text-white"
                    : "border-slate-300 text-slate-700 hover:border-sky-400 hover:text-sky-600 dark:border-slate-700 dark:text-slate-200 dark:hover:border-sky-400 dark:hover:text-sky-300"
                }`}
                aria-pressed={isActive}
              >
                {tech}
              </button>
            );
          })}
        </div>
      ) : null}

      <div className="grid gap-4">
        {isLoading
          ? Array.from({ length: PROJECT_SKELETON_COUNT }).map((_, index) => (
              // index is stable for skeleton placeholders
              // eslint-disable-next-line react/no-array-index-key
              <ProjectCardSkeleton key={index} />
            ))
          : null}
        {!isLoading && !filteredProjects.length ? (
          <article className="rounded-xl border border-slate-200 bg-white/80 p-6 text-sm text-slate-600 shadow-sm dark:border-slate-800 dark:bg-slate-900/60 dark:text-slate-300">
            {selectedStacks.size
              ? t("projects.noMatchesForSelection")
              : t("projects.placeholder")}
          </article>
        ) : null}
        {filteredProjects.map((project) => (
          <article
            key={project.id}
            className="rounded-xl border border-slate-200 bg-white/80 p-6 shadow-sm transition hover:border-sky-300 dark:border-slate-800 dark:bg-slate-900/60 dark:hover:border-sky-500"
          >
            <header className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
              <div>
                <h2 className="text-xl font-semibold text-slate-900 dark:text-slate-100">
                  {project.title}
                </h2>
                {project.subtitle ? (
                  <p className="text-sm text-slate-600 dark:text-slate-300">{project.subtitle}</p>
                ) : null}
                {project.category ? (
                  <span className="mt-2 inline-flex items-center rounded-full border border-slate-300 px-2 py-1 text-xs font-medium uppercase tracking-wide text-slate-600 dark:border-slate-700 dark:text-slate-300">
                    {project.category}
                  </span>
                ) : null}
              </div>
              {project.period ? (
                <p className="text-xs uppercase tracking-wide text-slate-500 dark:text-slate-400">
                  {formatDateRange(
                    project.period.start,
                    project.period.end,
                    t("common.presentLabel")
                  )}
                </p>
              ) : null}
            </header>
            {project.coverImageUrl ? (
              <img
                src={project.coverImageUrl}
                alt=""
                className="mt-4 h-48 w-full rounded-lg object-cover"
                loading="lazy"
              />
            ) : null}
            <p className="mt-4 text-sm text-slate-600 dark:text-slate-300">{project.description}</p>
            {project.techStack.length ? (
              <ul className="mt-4 flex flex-wrap gap-2">
                {project.techStack.map((tech) => (
                  <li
                    key={`${project.id}-${tech}`}
                    className="inline-flex items-center rounded-full border border-slate-300 px-3 py-1 text-xs font-medium text-slate-700 dark:border-slate-700 dark:text-slate-200"
                  >
                    {tech}
                  </li>
                ))}
              </ul>
            ) : null}
            {project.links.length ? (
              <div className="mt-4 flex flex-wrap gap-3">{renderProjectLinks(project)}</div>
            ) : null}
            {project.tags?.length ? (
              <div className="mt-4 flex flex-wrap gap-2">
                {project.tags.map((tag) => (
                  <span
                    key={`${project.id}-${tag}`}
                    className="inline-flex items-center rounded-full border border-dashed border-slate-300 px-2 py-1 text-[10px] uppercase tracking-wide text-slate-500 dark:border-slate-700 dark:text-slate-400"
                  >
                    {tag}
                  </span>
                ))}
              </div>
            ) : null}
          </article>
        ))}
        {error ? (
          <p role="alert" className="rounded-lg border border-rose-200 bg-rose-50 p-4 text-sm text-rose-700 dark:border-rose-700 dark:bg-rose-950/50 dark:text-rose-300">
            {t("projects.error")}
          </p>
        ) : null}
      </div>
    </section>
  );
}
