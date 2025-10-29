import type { Router as RemixRouter } from "@remix-run/router";
import { RouterProvider } from "react-router-dom";

import { appBrowserRouter } from "./app/router";
import { ThemeProvider } from "./providers/ThemeProvider";
import type { ThemeProviderProps } from "./providers/ThemeProvider";

export type AppThemeOverrides = Omit<ThemeProviderProps, "children">;

type AppProps = {
  router?: RemixRouter;
  themeOverrides?: AppThemeOverrides;
};

export function App({
  router = appBrowserRouter,
  themeOverrides,
}: AppProps = {}) {
  const providerProps = themeOverrides ?? {};

  return (
    <ThemeProvider {...providerProps}>
      <RouterProvider router={router} />
    </ThemeProvider>
  );
}

export default App;
