import { Fragment, useCallback, useEffect, useMemo } from "react";
import { useTranslation } from "react-i18next";
import { useLocation } from "react-router-dom";

import { useAdminMode } from "../../hooks/useAdminMode";
import { AdminConsole } from "../../modules/admin-console";

export function AdminLandingPage() {
  const { t } = useTranslation();
  const location = useLocation();
  const {
    isAdminMode,
    sessionActive,
    loading,
    toggleMode,
    sessionEmail,
  } = useAdminMode();

  const redirectTarget = useMemo(() => {
    if (location.pathname && location.pathname !== "/") {
      return `${location.pathname}${location.search ?? ""}`;
    }
    return "/admin";
  }, [location.pathname, location.search]);

  const loginHref = useMemo(() => {
    const loginBase =
      import.meta.env.VITE_ADMIN_LOGIN_URL?.trim() || "/api/admin/auth/login";

    const redirectParams = new URLSearchParams({
      redirect_uri: redirectTarget,
    }).toString();

    return loginBase.includes("?")
      ? `${loginBase}&${redirectParams}`
      : `${loginBase}?${redirectParams}`;
  }, [redirectTarget]);

  const handleLogin = useCallback(() => {
    try {
      sessionStorage.setItem("admin:postLoginRedirect", redirectTarget);
    } catch (error) {
      // Storage errors (quota/private mode) are non-fatal for login.
    }
    window.location.replace(loginHref);
  }, [loginHref, redirectTarget]);

  useEffect(() => {
    if (loading || !sessionActive) {
      return;
    }
    try {
      const target = sessionStorage.getItem("admin:postLoginRedirect");
      if (target) {
        sessionStorage.removeItem("admin:postLoginRedirect");
        const current = `${location.pathname}${location.search ?? ""}`;
        if (target !== current) {
          window.location.replace(target);
        }
      }
    } catch (error) {
      // Ignore storage access issues.
    }
  }, [loading, sessionActive, location.pathname, location.search]);

  if (loading) {
    return (
      <section className="mx-auto flex w-full max-w-3xl flex-col gap-6 px-4 py-12 sm:px-8">
        <header className="space-y-3">
          <div className="h-4 w-24 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
          <div className="h-8 w-1/2 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
          <div className="h-5 w-3/4 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
        </header>
        <div className="h-40 animate-pulse rounded-xl border border-slate-200 bg-slate-100 dark:border-slate-800 dark:bg-slate-900/60" />
      </section>
    );
  }

  if (!sessionActive) {
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
            {t("admin.modeGuard.loginRequired")}
          </p>
        </header>
        <div className="rounded-xl border border-slate-300 bg-white p-6 shadow-sm dark:border-slate-700 dark:bg-slate-900">
          <p className="text-sm text-slate-700 dark:text-slate-200">
            {t("admin.loginCallout")}
          </p>
          <button
            className="mt-4 inline-flex items-center justify-center rounded-md bg-slate-900 px-5 py-2 text-sm font-semibold text-white shadow-sm transition hover:bg-slate-800 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-slate-900 dark:bg-slate-100 dark:text-slate-900 dark:hover:bg-slate-200"
            onClick={handleLogin}
            type="button"
          >
            {t("admin.loginCta")}
          </button>
          <p className="mt-3 text-xs leading-relaxed text-slate-500 dark:text-slate-400">
            {t("admin.loginHelp")}
          </p>
        </div>
      </section>
    );
  }

  if (!isAdminMode) {
    return (
      <section className="mx-auto flex w-full max-w-4xl flex-col gap-6 px-4 py-12 sm:px-8">
        <header className="space-y-3">
          <p className="text-sm font-semibold uppercase tracking-[0.3em] text-slate-900 dark:text-slate-100">
            {t("admin.tagline")}
          </p>
          <h1 className="text-3xl font-bold text-slate-900 dark:text-slate-50 sm:text-4xl">
            {t("admin.title")}
          </h1>
          <p className="text-base text-slate-600 dark:text-slate-300">
            {t("admin.modeGuard.enableAdminMode")}
          </p>
        </header>
        <div className="rounded-xl border border-slate-300 bg-white p-6 shadow-sm dark:border-slate-700 dark:bg-slate-900">
          <p className="text-sm text-slate-700 dark:text-slate-200">
            {sessionEmail
              ? `Signed in as ${sessionEmail}.`
              : "Authenticated with administrator privileges."}
          </p>
          <button
            type="button"
            className="mt-4 inline-flex items-center justify-center rounded-md bg-sky-600 px-5 py-2 text-sm font-semibold text-white shadow-sm transition hover:bg-sky-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-sky-500 dark:bg-sky-500 dark:hover:bg-sky-400"
            onClick={() => toggleMode({ suppressPrompt: true })}
          >
            Enable admin mode
          </button>
          <p className="mt-3 text-xs leading-relaxed text-slate-500 dark:text-slate-400">
            {t("admin.modeGuard.enableAdminMode")}
          </p>
        </div>
      </section>
    );
  }

  return (
    <Fragment>
      <AdminConsole />
    </Fragment>
  );
}
