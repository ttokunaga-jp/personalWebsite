import { useTranslation } from "react-i18next";
import { apiClient } from "@shared/lib/api-client";
import { useEffect, useState } from "react";

type AdminStatusResponse = {
  status: string;
};

function App() {
  const { t } = useTranslation();
  const [status, setStatus] = useState("unknown");

  useEffect(() => {
    apiClient
      .get<AdminStatusResponse>("/health")
      .then((response) => setStatus(response.data.status))
      .catch(() => setStatus("unreachable"));
  }, []);

  return (
    <div className="min-h-screen bg-slate-100">
      <header className="bg-slate-900 p-6 text-white">
        <h1 className="text-2xl font-bold">{t("dashboard.title")}</h1>
        <p className="text-sm text-slate-300">{t("dashboard.subtitle")}</p>
      </header>
      <main className="mx-auto max-w-4xl p-8">
        <section className="rounded-lg border border-slate-200 bg-white p-6 shadow-sm">
          <h2 className="text-lg font-semibold text-slate-800">{t("dashboard.systemStatus")}</h2>
          <p className="mt-2 text-sm text-slate-600">{t("dashboard.systemStatusDescription")}</p>
          <div className="mt-4 rounded-md bg-slate-900 p-4 text-white">
            <span className="font-mono uppercase tracking-wide text-slate-400">
              {t("dashboard.apiStatus")}
            </span>
            <p className="text-2xl font-bold text-emerald-400">{status}</p>
          </div>
        </section>
      </main>
    </div>
  );
}

export default App;
