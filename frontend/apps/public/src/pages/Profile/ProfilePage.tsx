import { useTranslation } from "react-i18next";

export function ProfilePage() {
  const { t } = useTranslation();

  return (
    <section className="mx-auto flex w-full max-w-4xl flex-col gap-6 px-4 py-12 sm:px-8">
      <header className="space-y-3">
        <p className="text-sm font-semibold uppercase tracking-[0.3em] text-sky-500 dark:text-sky-400">
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
        <div className="rounded-xl border border-slate-200 bg-white/80 p-5 shadow-sm dark:border-slate-800 dark:bg-slate-900/60">
          <h2 className="text-lg font-semibold text-slate-900 dark:text-slate-100">
            {t("profile.sections.affiliations.title")}
          </h2>
          <p className="mt-2 text-sm text-slate-600 dark:text-slate-300">
            {t("profile.sections.affiliations.description")}
          </p>
        </div>
        <div className="rounded-xl border border-slate-200 bg-white/80 p-5 shadow-sm dark:border-slate-800 dark:bg-slate-900/60">
          <h2 className="text-lg font-semibold text-slate-900 dark:text-slate-100">
            {t("profile.sections.skills.title")}
          </h2>
          <p className="mt-2 text-sm text-slate-600 dark:text-slate-300">
            {t("profile.sections.skills.description")}
          </p>
        </div>
      </div>
    </section>
  );
}
