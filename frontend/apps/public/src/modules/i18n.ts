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
      common: {
        presentLabel: "Present"
      },
      home: {
        hero: {
          tagline: "Human-first innovation",
          title: "Crafting research-driven products and experiences.",
          description:
            "A modular single-page application scaffold ready for integrating portfolio content, showcasing research, projects, and contact workflows."
        },
        about: {
          title: "At a glance",
          name: "Name",
          location: "Location",
          fallbackName: "Loading profile…",
          fallbackLocation: "TBD",
          affiliation: "Primary affiliation",
          affiliationFallback: "Affiliation data will be published soon.",
          communities: "Communities",
          communitiesFallback: "Community involvement will be listed shortly.",
          error: "We were unable to load the latest profile details. Please try again later."
        },
        health: {
          title: "API status",
          caption: "Realtime health information pulled from the Go backend."
        },
        social: {
          title: "Connect",
          connectWith: "Connect via {{label}}",
          placeholder: "Social links will appear here once configured."
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
            description: "Display academic, professional, and community memberships.",
            empty: "Affiliation information will appear here soon."
          },
          skills: {
            title: "Skills and capabilities",
            description: "Expose the curated skill matrix sourced from structured data.",
            empty: "Skill groups are not yet published.",
            level: {
              beginner: "Beginner",
              intermediate: "Intermediate",
              advanced: "Advanced",
              expert: "Expert"
            }
          },
          lab: {
            title: "Laboratory",
            description: "Current research laboratory and supervision details.",
            name: "Lab",
            advisor: "Advisor",
            focus: "Focus area",
            visit: "Visit lab site",
            empty: "Laboratory information will be published when available."
          },
          work: {
            title: "Experience",
            description: "Professional experience and internships connected to ongoing research themes.",
            empty: "Work history is coming soon."
          },
          communities: {
            title: "Communities",
            description: "Open source, academic, and industry communities that inform the work.",
            empty: "Community activities will be shared in the future."
          }
        },
        error: "We could not refresh the profile data. Please reload the page."
      },
      research: {
        tagline: "Research",
        title: "Research portfolio",
        description:
          "Markdown-backed publications, abstracts, and outcomes render here with responsive layouts.",
        placeholder:
          "Content coming soon—integrate API-driven markdown rendering to showcase research output.",
        filters: {
          all: "All topics"
        },
        noEntriesForTag: "No research entries tagged “{{tag}}” yet.",
        updatedOn: "Updated {{date}}",
        error: "Research entries could not be loaded. Please retry shortly."
      },
      projects: {
        tagline: "Projects",
        title: "Project archive",
        description:
          "Highlight technical projects with tech stack chips, links, and filtering options.",
        placeholder:
          "Project cards will be rendered here once the API endpoints are connected.",
        filters: {
          all: "All stacks"
        },
        noMatchesForSelection: "No projects match the selected tech stack yet.",
        error: "Projects could not be retrieved from the API."
      },
      contact: {
        tagline: "Contact",
        title: "Get in touch",
        description:
          "Present booking availability, contact forms, and response SLAs sourced from backend services.",
        placeholder:
          "Form components and scheduling widgets will be mounted in this section.",
        availability: {
          title: "Availability",
          description: "Reserve a conversation slot that aligns with calendar availability.",
          groupLabel: "Available time slots",
          unavailable: "No open slots are currently available.",
          slotTo: "Ends {{end}}",
          timezone: "Times are displayed in {{timezone}}.",
          error: "Calendar availability could not be refreshed."
        },
        form: {
          legend: "Request details",
          name: "Your name",
          email: "Email address",
          topic: "Topic",
          topicPlaceholder: "Select a topic",
          message: "Message",
          submit: "Request booking",
          submitting: "Submitting…",
          success: "Thank you! Your request (ID: {{bookingId}}) has been received.",
          error: "We could not complete the booking. Please try again later.",
          configError: "Form configuration could not be loaded.",
          errors: {
            nameRequired: "Please provide your name.",
            emailRequired: "An email address is required.",
            emailInvalid: "Please provide a valid email address.",
            topicRequired: "Select a topic to help us route your request.",
            messageLength: "Share at least 20 characters so we can prepare effectively.",
            slotRequired: "Select an available time slot."
          }
        }
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
      common: {
        presentLabel: "現在"
      },
      home: {
        hero: {
          tagline: "人を軸にしたイノベーション",
          title: "研究を軸にしたプロダクトと体験を創出します。",
          description:
            "ポートフォリオ／研究紹介／プロジェクト情報／問い合わせ導線を統合できる SPA の骨格が整っています。"
        },
        about: {
          title: "プロフィール概要",
          name: "氏名",
          location: "所在地",
          fallbackName: "プロフィールを読込中…",
          fallbackLocation: "調整中",
          affiliation: "主たる所属",
          affiliationFallback: "所属情報は後ほど公開予定です。",
          communities: "コミュニティ",
          communitiesFallback: "コミュニティ参加情報は後日掲載します。",
          error: "最新のプロフィール情報を取得できませんでした。時間をおいて再度お試しください。"
        },
        health: {
          title: "API ステータス",
          caption: "Go バックエンドから取得した最新ステータスを表示します。"
        },
        social: {
          title: "つながる",
          connectWith: "{{label}} でつながる",
          placeholder: "ソーシャルリンクは設定後に表示されます。"
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
            description: "研究室やコミュニティ参加状況を整理して提示します。",
            empty: "所属情報は現在準備中です。"
          },
          skills: {
            title: "スキルセット",
            description: "構造化されたスキル情報からスキルマップを表示します。",
            empty: "スキル情報は後ほど公開予定です。",
            level: {
              beginner: "初級",
              intermediate: "中級",
              advanced: "上級",
              expert: "エキスパート"
            }
          },
          lab: {
            title: "研究室",
            description: "現在所属している研究室や指導教員の情報です。",
            name: "研究室",
            advisor: "指導教員",
            focus: "研究テーマ",
            visit: "研究室サイトへ",
            empty: "研究室情報は準備中です。"
          },
          work: {
            title: "職務経歴",
            description: "研究テーマと関連する職務経験やインターン実績を掲載します。",
            empty: "職務経歴は現在整理中です。"
          },
          communities: {
            title: "コミュニティ",
            description: "オープンソースや学術・産業コミュニティへの参画情報です。",
            empty: "コミュニティ参加情報は後日追加されます。"
          }
        },
        error: "プロフィール情報の取得に失敗しました。ページを再読み込みしてください。"
      },
      research: {
        tagline: "研究",
        title: "研究ポートフォリオ",
        description:
          "Markdown ベースの研究成果や発表資料をレスポンシブに表示します。",
        placeholder:
          "API 連携後に研究コンテンツがここに描画されます。",
        filters: {
          all: "すべて"
        },
        noEntriesForTag: "「{{tag}}」タグの研究コンテンツはまだありません。",
        updatedOn: "{{date}} 更新",
        error: "研究情報を取得できませんでした。時間をおいて再度お試しください。"
      },
      projects: {
        tagline: "プロジェクト",
        title: "プロジェクト一覧",
        description:
          "使用技術やリンク・フィルター付きのプロジェクトカードを表示します。",
        placeholder:
          "API 連携後にプロジェクトカードがここに表示されます。",
        filters: {
          all: "すべて"
        },
        noMatchesForSelection: "選択した技術スタックに該当するプロジェクトはありません。",
        error: "プロジェクト情報を取得できませんでした。"
      },
      contact: {
        tagline: "お問い合わせ",
        title: "コンタクト",
        description:
          "予約可能枠や問い合わせフォームをバックエンドと連携して表示します。",
        placeholder: "フォームと予約ウィジェットをこのセクションに配置します。",
        availability: {
          title: "空き状況",
          description: "カレンダーの空き状況に基づいてミーティング枠を予約できます。",
          groupLabel: "予約可能な時間枠",
          unavailable: "現在予約可能な枠はありません。",
          slotTo: "終了 {{end}}",
          timezone: "表示時刻は {{timezone}} を基準としています。",
          error: "空き状況の取得に失敗しました。"
        },
        form: {
          legend: "予約内容",
          name: "お名前",
          email: "メールアドレス",
          topic: "トピック",
          topicPlaceholder: "トピックを選択してください",
          message: "メッセージ",
          submit: "予約をリクエスト",
          submitting: "送信中…",
          success: "ありがとうございます。リクエスト (ID: {{bookingId}}) を受け付けました。",
          error: "予約リクエストを送信できませんでした。時間をおいて再度お試しください。",
          configError: "フォーム設定の取得に失敗しました。",
          errors: {
            nameRequired: "お名前を入力してください。",
            emailRequired: "メールアドレスを入力してください。",
            emailInvalid: "有効なメールアドレス形式で入力してください。",
            topicRequired: "トピックを選択してください。",
            messageLength: "具体的な内容を 20 文字以上で入力してください。",
            slotRequired: "予約する時間枠を選択してください。"
          }
        }
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
