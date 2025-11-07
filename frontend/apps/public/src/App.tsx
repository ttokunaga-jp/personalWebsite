import type { Router as RemixRouter } from "@remix-run/router";
import { RouterProvider } from "react-router-dom";

import { appBrowserRouter } from "./app/router";
import { ProfileResourceProvider } from "./modules/public-api";
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
      <ProfileResourceProvider>
        <RouterProvider router={router} />
      </ProfileResourceProvider>
    </ThemeProvider>
  );
}

export default App;
