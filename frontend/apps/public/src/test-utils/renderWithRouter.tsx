import type { Router as RemixRouter } from "@remix-run/router";
import { render } from "@testing-library/react";
import { act } from "react";

import { App, type AppThemeOverrides } from "../App";
import { createAppMemoryRouter } from "../app/router";
import { preloadRouteModules } from "../app/routes/routeConfig";

type RenderWithRouterOptions = {
  initialEntries?: string[];
  router?: RemixRouter;
  themeOverrides?: AppThemeOverrides;
};

export async function renderWithRouter({
  initialEntries,
  router,
  themeOverrides
}: RenderWithRouterOptions = {}) {
  await preloadRouteModules();

  const testRouter = router ?? createAppMemoryRouter(initialEntries);

  let renderResult: ReturnType<typeof render> | null = null;

  await act(async () => {
    renderResult = render(<App router={testRouter} themeOverrides={themeOverrides} />);
  });

  // Allow any queued microtasks (e.g., i18n initialization) to resolve before assertions.
  await act(
    () =>
      new Promise<void>((resolve) => {
        setTimeout(resolve, 0);
      })
  );

  if (!renderResult) {
    throw new Error("renderWithRouter failed to mount the component");
  }

  return Object.assign({ router: testRouter }, renderResult);
}
