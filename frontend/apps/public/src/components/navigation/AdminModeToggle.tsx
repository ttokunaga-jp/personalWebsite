import { Fragment } from "react";
import { useTranslation } from "react-i18next";

import { useAdminMode } from "../../hooks/useAdminMode";

export function AdminModeToggle() {
  const { t } = useTranslation();
  const {
    isAdminMode,
    loading,
    sessionActive,
    setMode,
    hasUnsavedChanges,
  } = useAdminMode();

  if (loading) {
    return (
      <div className="flex h-10 items-center gap-2 rounded-full border border-slate-300 bg-white px-4 text-sm text-slate-500 shadow-sm dark:border-slate-700 dark:bg-slate-900 dark:text-slate-400">
        <span className="inline-block h-2 w-2 animate-pulse rounded-full bg-slate-400 dark:bg-slate-500" />
        {t("admin.toggle.loading")}
      </div>
    );
  }

  if (!sessionActive) {
    return null;
  }

  const handleToggle = () => {
    const nextMode = isAdminMode ? "view" : "admin";
    setMode(nextMode);
  };

  return (
    <Fragment>
      <button
        type="button"
        className="group relative inline-flex items-center gap-3 rounded-full border border-slate-300 bg-white px-4 py-2 text-sm font-medium text-slate-700 shadow-sm transition hover:border-sky-400 hover:text-sky-600 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-sky-500 dark:border-slate-700 dark:bg-slate-900 dark:text-slate-200 dark:hover:border-sky-400 dark:hover:text-sky-300 dark:focus-visible:outline-sky-400"
        aria-pressed={isAdminMode}
        onClick={handleToggle}
      >
        <span
          className={`relative inline-flex h-5 w-10 items-center rounded-full transition ${
            isAdminMode ? "bg-sky-500" : "bg-slate-300 dark:bg-slate-600"
          }`}
        >
          <span
            className={`inline-block h-4 w-4 transform rounded-full bg-white shadow transition ${
              isAdminMode ? "translate-x-5" : "translate-x-1"
            }`}
          />
        </span>
        <span className="flex flex-col text-start">
          <span className="leading-none">
            {isAdminMode
              ? t("admin.toggle.adminLabel")
              : t("admin.toggle.viewLabel")}
          </span>
          {hasUnsavedChanges && (
            <span className="text-xs font-normal text-amber-600 dark:text-amber-400">
              {t("admin.toggle.unsavedNotice")}
            </span>
          )}
        </span>
      </button>
    </Fragment>
  );
}
