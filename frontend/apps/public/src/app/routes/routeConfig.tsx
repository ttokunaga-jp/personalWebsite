import type { ReactElement } from "react";

import {
  AdminLandingPage,
  ContactPage,
  HomePage,
  ProfilePage,
  ProjectsPage,
  ResearchPage
} from "../../pages";
import type { NavigationItem } from "../../types/navigation";

export type RouteDefinition = {
  path: string;
  element: ReactElement;
  labelKey: string;
  showInNavigation?: boolean;
  index?: boolean;
};

export const routeDefinitions: RouteDefinition[] = [
  {
    path: "/",
    element: <HomePage />,
    labelKey: "navigation.home",
    showInNavigation: true,
    index: true
  },
  {
    path: "/profile",
    element: <ProfilePage />,
    labelKey: "navigation.profile",
    showInNavigation: true
  },
  {
    path: "/research",
    element: <ResearchPage />,
    labelKey: "navigation.research",
    showInNavigation: true
  },
  {
    path: "/projects",
    element: <ProjectsPage />,
    labelKey: "navigation.projects",
    showInNavigation: true
  },
  {
    path: "/contact",
    element: <ContactPage />,
    labelKey: "navigation.contact",
    showInNavigation: true
  },
  {
    path: "/admin",
    element: <AdminLandingPage />,
    labelKey: "navigation.admin"
  }
];

export const navigationItems: NavigationItem[] = routeDefinitions
  .filter((route) => route.showInNavigation)
  .map(({ path, labelKey }) => ({
    path,
    labelKey
  }));
