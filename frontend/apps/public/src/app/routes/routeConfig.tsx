import { lazy, Suspense, type ReactElement } from "react";

import type { NavigationItem } from "../../types/navigation";

const HomePage = lazy(() =>
  import("../../pages/Home/HomePage").then((module) => ({ default: module.HomePage }))
);
const ProfilePage = lazy(() =>
  import("../../pages/Profile/ProfilePage").then((module) => ({ default: module.ProfilePage }))
);
const ResearchPage = lazy(() =>
  import("../../pages/Research/ResearchPage").then((module) => ({ default: module.ResearchPage }))
);
const ProjectsPage = lazy(() =>
  import("../../pages/Projects/ProjectsPage").then((module) => ({ default: module.ProjectsPage }))
);
const ContactPage = lazy(() =>
  import("../../pages/Contact/ContactPage").then((module) => ({ default: module.ContactPage }))
);
const AdminLandingPage = lazy(() =>
  import("../../pages/Admin/AdminLandingPage").then((module) => ({
    default: module.AdminLandingPage
  }))
);

const routeFallback = (
  <div className="mx-auto flex w-full max-w-6xl flex-1 flex-col gap-12 px-4 py-16 sm:px-8 lg:px-12">
    <section className="space-y-4">
      <span className="block h-4 w-20 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
      <span className="block h-10 w-3/4 animate-pulse rounded bg-slate-200 dark:bg-slate-700 sm:w-2/3" />
      <span className="block h-4 w-full animate-pulse rounded bg-slate-200 dark:bg-slate-700 sm:w-5/6" />
      <span className="block h-4 w-2/3 animate-pulse rounded bg-slate-200 dark:bg-slate-700 sm:w-1/2" />
    </section>

    <section className="grid gap-6 lg:grid-cols-[2fr,1fr]">
      <div className="space-y-4 rounded-2xl border border-slate-200 bg-white/70 p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900/50">
        <span className="block h-4 w-32 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
        <span className="block h-6 w-56 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
        <div className="space-y-2">
          <span className="block h-4 w-full animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
          <span className="block h-4 w-5/6 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
          <span className="block h-4 w-2/3 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
        </div>
      </div>

      <div className="space-y-4">
        <div className="rounded-2xl border border-slate-200 bg-white/70 p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900/50">
          <span className="block h-4 w-32 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
          <span className="mt-4 block h-6 w-20 animate-pulse rounded bg-emerald-200 dark:bg-emerald-900/60" />
          <span className="mt-4 block h-4 w-full animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
        </div>
        <div className="rounded-2xl border border-slate-200 bg-white/70 p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900/50">
          <span className="block h-4 w-32 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
          <div className="mt-4 flex flex-wrap gap-3">
            {Array.from({ length: 3 }).map((_, index) => (
              <span
                // index is stable for placeholder content
                // eslint-disable-next-line react/no-array-index-key
                key={index}
                className="inline-flex h-9 w-28 animate-pulse rounded-full bg-slate-200 dark:bg-slate-700"
              />
            ))}
          </div>
        </div>
      </div>
    </section>

    <section className="grid gap-4 md:grid-cols-2">
      {Array.from({ length: 4 }).map((_, index) => (
        <article
          // index is stable for placeholder content
          // eslint-disable-next-line react/no-array-index-key
          key={index}
          className="rounded-2xl border border-slate-200 bg-white/70 p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900/50"
        >
          <span className="block h-5 w-1/2 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
          <span className="mt-4 block h-4 w-full animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
          <span className="mt-2 block h-4 w-3/4 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
          <div className="mt-4 flex flex-wrap gap-2">
            {Array.from({ length: 3 }).map((__, badgeIndex) => (
              <span
                // badgeIndex is stable within the placeholder card
                // eslint-disable-next-line react/no-array-index-key
                key={badgeIndex}
                className="inline-flex h-6 w-20 animate-pulse rounded-full bg-slate-200 dark:bg-slate-700"
              />
            ))}
          </div>
        </article>
      ))}
    </section>
  </div>
);

const withSuspense = (component: ReactElement) => (
  <Suspense fallback={routeFallback}>{component}</Suspense>
);

export type RouteDefinition = {
  path: string;
  element: ReactElement;
  labelKey: string;
  showInNavigation?: boolean;
  index?: boolean;
};

export async function preloadRouteModules(): Promise<void> {
  await Promise.all([
    import("../../pages/Home/HomePage"),
    import("../../pages/Profile/ProfilePage"),
    import("../../pages/Research/ResearchPage"),
    import("../../pages/Projects/ProjectsPage"),
    import("../../pages/Contact/ContactPage"),
    import("../../pages/Admin/AdminLandingPage")
  ]);
}

export const routeDefinitions: RouteDefinition[] = [
  {
    path: "/",
    element: withSuspense(<HomePage />),
    labelKey: "navigation.home",
    showInNavigation: true,
    index: true
  },
  {
    path: "/profile",
    element: withSuspense(<ProfilePage />),
    labelKey: "navigation.profile",
    showInNavigation: true
  },
  {
    path: "/research",
    element: withSuspense(<ResearchPage />),
    labelKey: "navigation.research",
    showInNavigation: true
  },
  {
    path: "/projects",
    element: withSuspense(<ProjectsPage />),
    labelKey: "navigation.projects",
    showInNavigation: true
  },
  {
    path: "/contact",
    element: withSuspense(<ContactPage />),
    labelKey: "navigation.contact",
    showInNavigation: true
  },
  {
    path: "/admin",
    element: withSuspense(<AdminLandingPage />),
    labelKey: "navigation.admin"
  }
];

export const navigationItems: NavigationItem[] = routeDefinitions
  .filter((route) => route.showInNavigation)
  .map(({ path, labelKey }) => ({
    path,
    labelKey
  }));
