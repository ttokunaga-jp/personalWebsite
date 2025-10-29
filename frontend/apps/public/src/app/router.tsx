import type { Router as RemixRouter } from "@remix-run/router";
import {
  Navigate,
  createBrowserRouter,
  createMemoryRouter,
  type RouteObject,
} from "react-router-dom";

import { PageShell } from "../components/templates/PageShell";

import { routeDefinitions } from "./routes/routeConfig";

const FUTURE_FLAGS = {
  v7_startTransition: true,
  v7_relativeSplatPath: true,
} as const;

const appRoutes: RouteObject[] = [
  {
    path: "/",
    element: <PageShell />,
    children: [
      ...routeDefinitions.map((definition) => {
        if (definition.index) {
          return {
            index: true,
            element: definition.element,
          };
        }

        return {
          path: definition.path.replace(/^\//, ""),
          element: definition.element,
        };
      }),
      {
        path: "*",
        element: <Navigate to="/" replace />,
      },
    ],
  },
];

export function createAppBrowserRouter(): RemixRouter {
  return createBrowserRouter(appRoutes, {
    future: FUTURE_FLAGS,
  });
}

export function createAppMemoryRouter(
  initialEntries: string[] = ["/"],
): RemixRouter {
  return createMemoryRouter(appRoutes, {
    initialEntries,
    future: FUTURE_FLAGS,
  });
}

export const appBrowserRouter = createAppBrowserRouter();
