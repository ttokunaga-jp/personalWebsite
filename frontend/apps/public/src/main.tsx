import React from "react";
import ReactDOM from "react-dom/client";

import App from "./App";
import { preloadRouteModules } from "./app/routes/routeConfig";
import "./styles.css";
import "./modules/i18n";

async function bootstrap() {
  ReactDOM.createRoot(document.getElementById("root")!).render(
    <React.StrictMode>
      <App />
    </React.StrictMode>,
  );

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
