import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { NavLink } from "react-router-dom";

import { navigationItems } from "../../../app/routes/routeConfig";
import { classNames } from "../../../lib/classNames";
import { LanguageSwitcher } from "../../molecules/LanguageSwitcher";
import { ThemeToggle } from "../../molecules/ThemeToggle";

export function SiteHeader() {
  const { t } = useTranslation();
  const [isMobileOpen, setIsMobileOpen] = useState(false);

  useEffect(() => {
    const handleResize = () => {
      if (window.innerWidth >= 768) {
        setIsMobileOpen(false);
      }
    };
    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  }, []);

  return (
    <header className="border-b border-slate-200 bg-white/80 backdrop-blur-md transition-colors dark:border-slate-800 dark:bg-slate-950/80">
      <div className="mx-auto flex w-full max-w-6xl items-center justify-between px-4 py-4 sm:px-8">
        <div className="flex items-center gap-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-full bg-sky-600 font-bold text-white dark:bg-sky-500">
            TA
          </div>
          <div className="flex flex-col">
            <span className="text-sm font-semibold uppercase tracking-[0.3em] text-slate-400 dark:text-slate-500">
              {t("branding.subtitle")}
            </span>
            <span className="text-base font-semibold text-slate-900 dark:text-slate-100">
              {t("branding.title")}
            </span>
          </div>
        </div>

        <div className="flex items-center gap-3 md:hidden">
          <ThemeToggle />
          <button
            type="button"
            className="inline-flex h-10 w-10 items-center justify-center rounded-full border border-slate-300 text-slate-600 transition hover:border-sky-400 hover:text-sky-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-sky-500 focus-visible:ring-offset-2 dark:border-slate-700 dark:text-slate-200 dark:hover:border-sky-400 dark:hover:text-sky-300 dark:focus-visible:ring-sky-400"
            aria-expanded={isMobileOpen}
            aria-label={t("navigation.toggle")}
            onClick={() => setIsMobileOpen((prev) => !prev)}
          >
            <span className="sr-only">{t("navigation.toggle")}</span>
            <div className="flex flex-col gap-1.5">
              <span className="block h-0.5 w-5 rounded-full bg-current transition" />
              <span className="block h-0.5 w-4 rounded-full bg-current transition" />
              <span className="block h-0.5 w-5 rounded-full bg-current transition" />
            </div>
          </button>
        </div>

        <nav
          className="hidden items-center gap-6 md:flex"
          aria-label={t("navigation.label")}
        >
          {navigationItems.map((item) => (
            <NavLink
              key={item.path}
              to={item.path}
              className={({ isActive }) =>
                classNames(
                  "text-sm font-medium text-slate-600 transition hover:text-sky-600 dark:text-slate-300 dark:hover:text-sky-300",
                  isActive && "text-slate-900 dark:text-white",
                )
              }
            >
              {t(item.labelKey)}
            </NavLink>
          ))}
          <LanguageSwitcher />
          <ThemeToggle />
        </nav>
      </div>

      <div
        className={classNames("md:hidden", isMobileOpen ? "block" : "hidden")}
      >
        <nav
          className="flex flex-col gap-2 border-t border-slate-200 bg-white px-4 py-4 shadow-sm dark:border-slate-800 dark:bg-slate-950"
          aria-label={t("navigation.mobileLabel")}
        >
          {navigationItems.map((item) => (
            <NavLink
              key={item.path}
              to={item.path}
              className={({ isActive }) =>
                classNames(
                  "rounded-lg px-3 py-2 text-sm font-medium text-slate-700 transition hover:bg-slate-100 dark:text-slate-200 dark:hover:bg-slate-800",
                  isActive &&
                    "bg-slate-100 text-slate-900 dark:bg-slate-800 dark:text-white",
                )
              }
              onClick={() => setIsMobileOpen(false)}
            >
              {t(item.labelKey)}
            </NavLink>
          ))}
          <div className="flex items-center justify-between rounded-lg bg-slate-100 px-3 py-2 dark:bg-slate-800">
            <LanguageSwitcher />
            <ThemeToggle />
          </div>
        </nav>
      </div>
    </header>
  );
}
