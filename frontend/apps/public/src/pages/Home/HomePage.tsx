import { apiClient } from "@shared/lib/api-client";
import { startTransition, useEffect, useState } from "react";
import { Trans, useTranslation } from "react-i18next";

type HealthResponse = {
  status: string;
};

export function HomePage() {
  const { t } = useTranslation();
  const [status, setStatus] = useState<string>("loading");

  useEffect(() => {
    let subscribed = true;

    apiClient
      .get<HealthResponse>("/health")
      .then((response) => {
        if (subscribed) {
          startTransition(() => {
            setStatus(response.data.status ?? "ok");
          });
        }
      })
      .catch(() => {
        if (subscribed) {
          startTransition(() => {
            setStatus("unreachable");
          });
        }
      });

    return () => {
      subscribed = false;
    };
  }, []);

  return (
    <section className="mx-auto flex w-full max-w-5xl flex-1 flex-col gap-10 px-4 py-16 sm:px-8 lg:px-12">
      <header className="flex flex-col gap-6 text-center md:text-left">
        <p className="text-sm font-semibold uppercase tracking-[0.3em] text-sky-500 dark:text-sky-400">
          {t("home.hero.tagline")}
        </p>
        <h1 className="text-4xl font-bold text-slate-900 dark:text-slate-50 sm:text-5xl lg:text-6xl">
          {t("home.hero.title")}
        </h1>
        <p className="text-lg leading-relaxed text-slate-600 dark:text-slate-300 md:max-w-3xl">
          <Trans i18nKey="home.hero.description" />
        </p>
      </header>

      <div className="grid gap-8 md:grid-cols-2">
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
      </div>
    </section>
  );
}
