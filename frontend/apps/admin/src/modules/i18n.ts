import i18n from "i18next";
import { initReactI18next } from "react-i18next";

const resources = {
  en: {
    translation: {
      "dashboard.title": "Admin console",
      "dashboard.subtitle": "Operate the portfolio content and reservations.",
      "dashboard.systemStatus": "System health",
      "dashboard.systemStatusDescription": "Monitor the API and background jobs from here.",
      "dashboard.apiStatus": "API status"
    }
  },
  ja: {
    translation: {
      "dashboard.title": "管理コンソール",
      "dashboard.subtitle": "ポートフォリオのコンテンツと予約を管理します。",
      "dashboard.systemStatus": "システムヘルス",
      "dashboard.systemStatusDescription": "API とバックグラウンドジョブを監視します。",
      "dashboard.apiStatus": "API ステータス"
    }
  }
};

void i18n.use(initReactI18next).init({
  resources,
  fallbackLng: "ja",
  interpolation: {
    escapeValue: false
  }
});

export default i18n;
