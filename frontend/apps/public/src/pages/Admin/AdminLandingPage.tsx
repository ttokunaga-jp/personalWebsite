import { useTranslation } from "react-i18next";

export function AdminLandingPage() {
  const { t } = useTranslation();

  return (
    <section className="mx-auto flex w-full max-w-3xl flex-col gap-6 px-4 py-12 sm:px-8">
      <header className="space-y-3">
        <p className="text-sm font-semibold uppercase tracking-[0.3em] text-sky-500 dark:text-sky-400">
          {t("admin.tagline")}
        </p>
        <h1 className="text-3xl font-bold text-slate-900 dark:text-slate-50 sm:text-4xl">
          {t("admin.title")}
        </h1>
        <p className="text-base text-slate-600 dark:text-slate-300">
          {t("admin.description")}
        </p>
      </header>
      <div className="rounded-xl border border-dashed border-slate-300 p-6 text-sm text-slate-600 dark:border-slate-700 dark:text-slate-300">
        {t("admin.placeholder")}
      </div>
    </section>
  );
}
