import { Outlet } from "react-router-dom";

import { SiteFooter } from "../../organisms/SiteFooter";
import { SiteHeader } from "../../organisms/SiteHeader";

export function PageShell() {
  return (
    <div className="flex min-h-screen flex-col bg-slate-50 text-slate-900 transition-colors duration-300 dark:bg-slate-950 dark:text-slate-50">
      <SiteHeader />
      <main className="flex-1">
        <Outlet />
      </main>
      <SiteFooter />
    </div>
  );
}
