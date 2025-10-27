import { useTranslation } from "react-i18next";

export function ResearchPage() {
  const { t } = useTranslation();

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
      <div className="flex flex-col gap-4 rounded-xl border border-slate-200 bg-white/80 p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900/60">
        <p className="text-sm text-slate-600 dark:text-slate-300">
          {t("research.placeholder")}
        </p>
      </div>
    </section>
  );
}
