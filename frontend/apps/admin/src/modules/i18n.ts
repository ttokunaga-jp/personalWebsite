import i18n from "i18next";
import { initReactI18next } from "react-i18next";

const resources = {
  en: {
    translation: {
      "dashboard.title": "Admin console",
      "dashboard.subtitle": "Operate the portfolio content and reservations.",
      "dashboard.systemStatus": "System health",
      "dashboard.systemStatusDescription":
        "Monitor the API and background jobs from here.",
      "dashboard.apiStatus": "API status",
      "summary.title": "Overview",
      "projects.title": "Projects",
      "research.title": "Research",
      "blogs.title": "Blog posts",
      "meetings.title": "Meetings",
      "blacklist.title": "Blacklist",
      "actions.create": "Create",
      "actions.update": "Update",
      "actions.delete": "Delete",
      "status.loading": "Loading data…",
      "status.error": "Failed to load data",
      "auth.requiredTitle": "Sign in required",
      "auth.requiredDescription":
        "You need to authenticate with your administrator Google account to view this console.",
      "auth.signIn": "Continue with Google",
      "auth.signOut": "Sign out",
      "auth.supportPrompt": "Need access or facing issues?",
      "auth.contactSupport": "Contact support",
    },
  },
  ja: {
    translation: {
      "dashboard.title": "管理コンソール",
      "dashboard.subtitle": "ポートフォリオのコンテンツと予約を管理します。",
      "dashboard.systemStatus": "システムヘルス",
      "dashboard.systemStatusDescription":
        "API とバックグラウンドジョブを監視します。",
      "dashboard.apiStatus": "API ステータス",
      "summary.title": "サマリー",
      "projects.title": "プロジェクト",
      "research.title": "研究",
      "blogs.title": "ブログ投稿",
      "meetings.title": "予約",
      "blacklist.title": "ブラックリスト",
      "actions.create": "作成",
      "actions.update": "更新",
      "actions.delete": "削除",
      "status.loading": "読込中…",
      "status.error": "データの取得に失敗しました",
      "auth.requiredTitle": "サインインが必要です",
      "auth.requiredDescription":
        "管理者用の Google アカウントで認証すると、このコンソールを表示できます。",
      "auth.signIn": "Google で続行",
      "auth.signOut": "サインアウト",
      "auth.supportPrompt": "アクセス権限については、サポートまでお問い合わせください。",
      "auth.contactSupport": "サポートに連絡",
    },
  },
};

void i18n.use(initReactI18next).init({
  resources,
  fallbackLng: "ja",
  interpolation: {
    escapeValue: false,
  },
});

export default i18n;
