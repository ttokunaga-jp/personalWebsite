import React from "react";
import ReactDOM from "react-dom/client";

import App from "./App";
import { AuthSessionProvider } from "./modules/auth-session";
import "./styles.css";
import "./modules/i18n";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <AuthSessionProvider>
      <App />
    </AuthSessionProvider>
  </React.StrictMode>,
);
