import i18n from "i18next";
import { initReactI18next } from "react-i18next";

const resources = {
  en: {
    translation: {
      branding: {
        subtitle: "Portfolio",
        title: "Takumi Asano"
      },
      navigation: {
        label: "Primary Navigation",
        mobileLabel: "Mobile Navigation",
        toggle: "Toggle navigation menu",
        home: "Home",
        profile: "Profile",
        research: "Research",
        projects: "Projects",
        contact: "Contact",
        admin: "Admin"
      },
      themeToggle: {
        setLight: "Switch to light theme",
        setDark: "Switch to dark theme"
      },
      languages: {
        en: "English",
        ja: "日本語"
      },
      home: {
        hero: {
          tagline: "Human-first innovation",
          title: "Crafting research-driven products and experiences.",
          description:
            "A modular single-page application scaffold ready for integrating portfolio content, showcasing research, projects, and contact workflows."
        },
        health: {
          title: "API status",
          caption: "Realtime health information pulled from the Go backend."
        },
        quickLinks: {
          title: "Quick links",
          projects: "Projects",
          research: "Research",
          contact: "Contact",
          supporting: "Use the navigation above for full page views."
        }
      },
      profile: {
        tagline: "About",
        title: "Professional profile",
        description:
          "Summaries for affiliations, work history, communities, and skill sets will be populated from the API.",
        sections: {
          affiliations: {
            title: "Affiliations",
            description: "Display academic, professional, and community memberships."
          },
          skills: {
            title: "Skills and capabilities",
            description: "Expose the curated skill matrix sourced from structured data."
          }
        }
      },
      research: {
        tagline: "Research",
        title: "Research portfolio",
        description:
          "Markdown-backed publications, abstracts, and outcomes render here with responsive layouts.",
        placeholder:
          "Content coming soon—integrate API-driven markdown rendering to showcase research output."
      },
      projects: {
        tagline: "Projects",
        title: "Project archive",
        description:
          "Highlight technical projects with tech stack chips, links, and filtering options.",
        placeholder:
          "Project cards will be rendered here once the API endpoints are connected."
      },
      contact: {
        tagline: "Contact",
        title: "Get in touch",
        description:
          "Present booking availability, contact forms, and response SLAs sourced from backend services.",
        placeholder:
          "Form components and scheduling widgets will be mounted in this section."
      },
      admin: {
        tagline: "Admin",
        title: "Management console",
        description:
          "Authenticated administrators can manage content, reservations, and blacklists from the dedicated UI.",
        placeholder:
          "Guard this route using JWT and Google SSO once the auth flow is implemented."
      },
      footer: {
        copyright: "Takumi Asano. All rights reserved.",
        links: {
          github: "GitHub",
          twitter: "X (Twitter)",
          contact: "Email"
        }
      }
    }
  },
  ja: {
    translation: {
      branding: {
        subtitle: "ポートフォリオ",
        title: "浅野 拓巳"
      },
      navigation: {
        label: "メインナビゲーション",
        mobileLabel: "モバイルナビゲーション",
        toggle: "ナビゲーションメニューを開閉",
        home: "ホーム",
        profile: "プロフィール",
        research: "研究",
        projects: "プロジェクト",
        contact: "お問い合わせ",
        admin: "管理"
      },
      themeToggle: {
        setLight: "ライトテーマに切り替え",
        setDark: "ダークテーマに切り替え"
      },
      languages: {
        en: "English",
        ja: "日本語"
      },
      home: {
        hero: {
          tagline: "人を軸にしたイノベーション",
          title: "研究を軸にしたプロダクトと体験を創出します。",
          description:
            "ポートフォリオ／研究紹介／プロジェクト情報／問い合わせ導線を統合できる SPA の骨格が整っています。"
        },
        health: {
          title: "API ステータス",
          caption: "Go バックエンドから取得した最新ステータスを表示します。"
        },
        quickLinks: {
          title: "ショートカット",
          projects: "プロジェクト",
          research: "研究",
          contact: "お問い合わせ",
          supporting: "詳細は上部のナビゲーションから各ページへ移動してください。"
        }
      },
      profile: {
        tagline: "プロフィール",
        title: "プロフィール概要",
        description:
          "所属・コミュニティ・職務経歴・スキルセットなどを API 経由で表示します。",
        sections: {
          affiliations: {
            title: "所属・コミュニティ",
            description: "研究室やコミュニティ参加状況を整理して提示します。"
          },
          skills: {
            title: "スキルセット",
            description: "構造化されたスキル情報からスキルマップを表示します。"
          }
        }
      },
      research: {
        tagline: "研究",
        title: "研究ポートフォリオ",
        description:
          "Markdown ベースの研究成果や発表資料をレスポンシブに表示します。",
        placeholder:
          "API 連携後に研究コンテンツがここに描画されます。"
      },
      projects: {
        tagline: "プロジェクト",
        title: "プロジェクト一覧",
        description:
          "使用技術やリンク・フィルター付きのプロジェクトカードを表示します。",
        placeholder:
          "API 連携後にプロジェクトカードがここに表示されます。"
      },
      contact: {
        tagline: "お問い合わせ",
        title: "コンタクト",
        description:
          "予約可能枠や問い合わせフォームをバックエンドと連携して表示します。",
        placeholder: "フォームと予約ウィジェットをこのセクションに配置します。"
      },
      admin: {
        tagline: "管理",
        title: "管理者コンソール",
        description:
          "認証済み管理者がコンテンツ、予約、ブラックリストを管理する UI を提供します。",
        placeholder: "Google SSO + JWT 連携後にアクセス制御を実装します。"
      },
      footer: {
        copyright: "Takumi Asano. All rights reserved.",
        links: {
          github: "GitHub",
          twitter: "X (Twitter)",
          contact: "メール"
        }
      }
    }
  }
};

void i18n.use(initReactI18next).init({
  resources,
  fallbackLng: "en",
  react: {
    useSuspense: false
  },
  initImmediate: false,
  interpolation: {
    escapeValue: false
  }
});

export default i18n;
