import { useState } from "react";
import { useTranslation } from "react-i18next";

import { useAdminMode } from "../../hooks/useAdminMode";

import { BlacklistPanel } from "./panels/BlacklistPanel";
import { MeetingUrlPanel } from "./panels/MeetingUrlPanel";
import { ReservationsPanel } from "./panels/ReservationsPanel";
import { SocialLinksPanel } from "./panels/SocialLinksPanel";
import { TechCatalogPanel } from "./panels/TechCatalogPanel";

type AdminView =
  | "reservations"
  | "blacklist"
  | "techCatalog"
  | "socialLinks"
  | "meetingUrl";

export function AdminConsole() {
  const { t } = useTranslation();
  const { confirmIfUnsaved } = useAdminMode();
  const [activeView, setActiveView] = useState<AdminView>("reservations");

  const navItems: { id: AdminView; label: string; description: string }[] = [
    {
      id: "reservations",
      label: t("admin.console.tabs.reservations"),
      description: "Review recent booking submissions and manage statuses.",
    },
    {
      id: "blacklist",
      label: t("admin.console.tabs.blacklist"),
      description: "Block abusive contacts and prevent future bookings.",
    },
    {
      id: "techCatalog",
      label: t("admin.console.tabs.techCatalog"),
      description: "Curate the canonical technology catalog used for tagging.",
    },
    {
      id: "socialLinks",
      label: t("admin.console.tabs.socialLinks"),
      description: "Manage bilingual labels and order for social link badges.",
    },
    {
      id: "meetingUrl",
      label: t("admin.console.tabs.meetingUrl"),
      description: "Edit the email template used when sharing meeting URLs.",
    },
  ];

  const handleSelectView = (view: AdminView) => {
    if (view === activeView) {
      return;
    }
    if (!confirmIfUnsaved()) {
      return;
    }
    setActiveView(view);
  };

  return (
    <div className="mx-auto flex w-full max-w-6xl flex-col gap-6 px-4 py-10 sm:px-8 lg:flex-row lg:gap-10">
      <aside className="w-full space-y-4 rounded-2xl border border-slate-200 bg-white/80 p-5 shadow-sm dark:border-slate-700 dark:bg-slate-900/50 lg:w-72">
        <header>
          <h1 className="text-xl font-semibold text-slate-900 dark:text-slate-100">
            {t("admin.title")}
          </h1>
          <p className="mt-1 text-sm text-slate-600 dark:text-slate-300">
            {t("admin.description")}
          </p>
        </header>
        <nav className="space-y-2">
          {navItems.map((item) => {
            const isActive = item.id === activeView;
            return (
              <button
                key={item.id}
                type="button"
                onClick={() => handleSelectView(item.id)}
                className={`w-full rounded-xl border px-4 py-3 text-left transition ${
                  isActive
                    ? "border-sky-400 bg-sky-50 text-sky-700 dark:border-sky-500 dark:bg-sky-950/30 dark:text-sky-300"
                    : "border-transparent bg-slate-50 text-slate-700 hover:border-slate-200 dark:bg-slate-900/70 dark:text-slate-200 dark:hover:border-slate-700"
                }`}
              >
                <span className="block text-sm font-semibold">{item.label}</span>
                <span className="mt-1 block text-xs text-slate-500 dark:text-slate-400">
                  {item.description}
                </span>
              </button>
            );
          })}
        </nav>
        <div className="rounded-lg border border-amber-200 bg-amber-50 px-3 py-3 text-xs text-amber-700 dark:border-amber-900 dark:bg-amber-950/40 dark:text-amber-200">
          {t("admin.console.draftWarning")}
        </div>
      </aside>

      <section className="flex-1">
        {(() => {
          switch (activeView) {
            case "reservations":
              return <ReservationsPanel />;
            case "blacklist":
              return <BlacklistPanel />;
            case "techCatalog":
              return <TechCatalogPanel />;
            case "socialLinks":
              return <SocialLinksPanel />;
            case "meetingUrl":
              return <MeetingUrlPanel />;
            default:
              return (
                <div className="rounded-xl border border-slate-200 bg-white/80 p-6 text-sm text-slate-600 dark:border-slate-700 dark:bg-slate-900/50 dark:text-slate-300">
                  {t("admin.console.emptyState")}
                </div>
              );
          }
        })()}
      </section>
    </div>
  );
}
