import type { Router as RemixRouter } from "@remix-run/router";
import { act, render } from "@testing-library/react";
import type { RenderResult } from "@testing-library/react";

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
  themeOverrides,
}: RenderWithRouterOptions = {}): Promise<RenderResult & { router: RemixRouter }> {
  await preloadRouteModules();

  const testRouter = router ?? createAppMemoryRouter(initialEntries);

  const renderResult = render(
    <App router={testRouter} themeOverrides={themeOverrides} />,
  );

  // Allow queued microtasks (e.g., i18n initialization) to settle before assertions.
  await act(
    () =>
      new Promise<void>((resolve) => {
        setTimeout(resolve, 0);
      }),
  );

  return { router: testRouter, ...renderResult };
}
