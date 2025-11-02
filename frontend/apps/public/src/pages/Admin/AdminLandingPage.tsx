import { useTranslation } from "react-i18next";

export function AdminLandingPage() {
  const { t } = useTranslation();
  const loginBase =
    import.meta.env.VITE_ADMIN_LOGIN_URL?.trim() || "/api/admin/auth/login";
  const redirectParams = new URLSearchParams({ redirect_uri: "/admin" }).toString();
  const loginHref = loginBase.includes("?")
    ? `${loginBase}&${redirectParams}`
    : `${loginBase}?${redirectParams}`;

  return (
    <section className="mx-auto flex w-full max-w-3xl flex-col gap-6 px-4 py-12 sm:px-8">
      <header className="space-y-3">
        <p className="text-sm font-semibold uppercase tracking-[0.3em] text-slate-900 dark:text-slate-100">
          {t("admin.tagline")}
        </p>
        <h1 className="text-3xl font-bold text-slate-900 dark:text-slate-50 sm:text-4xl">
          {t("admin.title")}
        </h1>
        <p className="text-base text-slate-600 dark:text-slate-300">
          {t("admin.description")}
        </p>
      </header>
      <div className="rounded-xl border border-slate-300 bg-white p-6 shadow-sm dark:border-slate-700 dark:bg-slate-900">
        <p className="text-sm text-slate-700 dark:text-slate-200">
          {t("admin.loginCallout")}
        </p>
        <a
          href={loginHref}
          className="mt-4 inline-flex items-center justify-center rounded-md bg-slate-900 px-5 py-2 text-sm font-semibold text-white shadow-sm transition hover:bg-slate-800 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-slate-900 dark:bg-slate-100 dark:text-slate-900 dark:hover:bg-slate-200"
        >
          {t("admin.loginCta")}
        </a>
        <p className="mt-3 text-xs leading-relaxed text-slate-500 dark:text-slate-400">
          {t("admin.loginHelp")}
        </p>
      </div>
    </section>
  );
}
