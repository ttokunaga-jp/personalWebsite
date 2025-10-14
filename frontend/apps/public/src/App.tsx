import { apiClient } from "@shared/lib/api-client";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";

type HealthResponse = {
  status: string;
};

function App() {
  const { t } = useTranslation();
  const [status, setStatus] = useState<string>("unknown");

  useEffect(() => {
    apiClient
      .get<HealthResponse>("/health")
      .then((response) => setStatus(response.data.status))
      .catch(() => setStatus("unreachable"));
  }, []);

  return (
    <div className="min-h-screen bg-slate-900 text-white">
      <div className="mx-auto flex max-w-3xl flex-col gap-6 p-10 text-center">
        <h1 className="text-4xl font-bold">{t("welcome")}</h1>
        <p className="text-lg">{t("intro")}</p>
        <div className="rounded-md bg-slate-800 p-4 shadow">
          <p className="font-mono text-sm uppercase tracking-wide text-slate-400">
            API status
          </p>
          <p className="text-2xl font-semibold text-emerald-400">{status}</p>
        </div>
      </div>
    </div>
  );
}

export default App;
