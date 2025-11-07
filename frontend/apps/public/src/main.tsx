import React from "react";
import ReactDOM from "react-dom/client";

import App from "./App";
import { preloadRouteModules } from "./app/routes/routeConfig";
import "./styles.css";
import "./modules/i18n";

async function waitForStyles() {
  if (typeof document === "undefined") {
    return;
  }

  const styleLinks = Array.from(
    document.querySelectorAll<HTMLLinkElement>('link[rel="stylesheet"]'),
  );

  if (!styleLinks.length) {
    return;
  }

  await Promise.all(
    styleLinks.map(
      (link) =>
        new Promise<void>((resolve) => {
          if (link.sheet) {
            resolve();
            return;
          }
          const handle = () => {
            link.removeEventListener("load", handle);
            link.removeEventListener("error", handle);
            resolve();
          };
          link.addEventListener("load", handle);
          link.addEventListener("error", handle);
        }),
    ),
  );
}

async function bootstrap() {
  const rootElement = document.getElementById("root");
  if (!rootElement) {
    throw new Error("Root element #root not found");
  }

  const root = ReactDOM.createRoot(rootElement);

  root.render(
    <React.StrictMode>
      <App />
    </React.StrictMode>,
  );

  const markAppLoaded = () => {
    rootElement.setAttribute("data-app-loaded", "true");
  };

  void waitForStyles()
    .catch(() => {})
    .finally(() => {
      if (typeof window.requestAnimationFrame === "function") {
        window.requestAnimationFrame(markAppLoaded);
      } else {
        window.setTimeout(markAppLoaded, 0);
      }
    });

  if (import.meta.env.MODE !== "test") {
    const schedulePreload = () => {
      preloadRouteModules().catch((error) => {
        if (import.meta.env.DEV) {
          console.warn("Deferred route preloading failed.", error);
        }
      });
    };

    const idleWindow = window as typeof window & {
      requestIdleCallback?: (callback: IdleRequestCallback) => number;
    };

    if (typeof idleWindow.requestIdleCallback === "function") {
      idleWindow.requestIdleCallback(schedulePreload);
    } else {
      window.setTimeout(schedulePreload, 0);
    }
  }
}

void bootstrap();
