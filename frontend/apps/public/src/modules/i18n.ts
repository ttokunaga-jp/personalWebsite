import i18n from "i18next";
import { initReactI18next } from "react-i18next";

import {
  FALLBACK_LANGUAGE,
  SUPPORTED_LANGUAGES,
  type SupportedLanguage,
} from "./language/config";
import {
  matchSupportedLanguage,
  persistLanguagePreference,
  resolveInitialLanguage,
} from "./language/preference";

export const resources = {
  en: {
    translation: {
      branding: {
        subtitle: "Portfolio",
        title: "Takumi Tokunaga",
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
        admin: "Admin",
      },
      themeToggle: {
        setLight: "Switch to light theme",
        setDark: "Switch to dark theme",
      },
      languages: {
        en: "English",
        ja: "日本語",
      },
      common: {
        presentLabel: "Present",
      },
      home: {
        hero: {
          tagline: "Real-world innovation",
          title: "Crafting RAG-powered learning tools and robotics experiences.",
          description:
            "The public portfolio for Takumi Tokunaga, aggregating research updates, project case studies, and booking workflows backed by a Go + React stack.",
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
          error:
            "We were unable to load the latest profile details. Please try again later.",
          lastUpdated: "Profile last updated",
        },
        health: {
          title: "API status",
          caption: "Realtime health information pulled from the Go backend.",
        },
        social: {
          title: "Connect",
          connectWith: "Connect via {{label}}",
          placeholder: "Social links will appear here once configured.",
        },
        quickLinks: {
          title: "Quick links",
          projects: "Projects",
          research: "Research",
          contact: "Contact",
          supporting: "Use the navigation above for full page views.",
        },
        tech: {
          title: "Technology focus",
          description:
            "Core stacks and supporting tools powering current initiatives.",
        },
        work: {
          title: "Recent initiatives",
          description:
            "Selected engagements highlighting scope, responsibilities, and impact.",
          empty: "Work history will be published soon.",
        },
        chips: {
          empty: "Highlights will be displayed once data is available.",
        },
        affiliations: {
          title: "Affiliations overview",
        },
      },
      profile: {
        tagline: "About",
        title: "Professional profile",
        description:
          "Summaries for affiliations, work history, communities, and skill sets will be populated from the API.",
        sections: {
          affiliations: {
            title: "Affiliations",
            description:
              "Display academic, professional, and community memberships.",
            empty: "Affiliation information will appear here soon.",
          },
          skills: {
            title: "Skills and capabilities",
            description:
              "Expose the curated skill matrix sourced from structured data.",
            empty: "Skill groups are not yet published.",
            level: {
              beginner: "Beginner",
              intermediate: "Intermediate",
              advanced: "Advanced",
              expert: "Expert",
            },
          },
          lab: {
            title: "Laboratory",
            description: "Current research laboratory and supervision details.",
            name: "Lab",
            advisor: "Advisor",
            focus: "Focus area",
            room: "Room",
            visit: "Visit lab site",
            empty: "Laboratory information will be published when available.",
          },
          work: {
            title: "Experience",
            description:
              "Professional experience and internships connected to ongoing research themes.",
            empty: "Work history is coming soon.",
          },
          communities: {
            title: "Communities",
            description:
              "Open source, academic, and industry communities that inform the work.",
            empty: "Community activities will be shared in the future.",
          },
          social: {
            title: "Social links",
          },
        },
        error: "We could not refresh the profile data. Please reload the page.",
      },
      research: {
        tagline: "Research",
        title: "Research portfolio",
        description:
          "Markdown-backed publications, abstracts, and outcomes render here with responsive layouts.",
        placeholder:
          "Content coming soon—integrate API-driven entries to showcase research output.",
        filters: {
          all: "All entries",
          research: "Research",
          blog: "Blog",
          tags: "Filter by tag",
          tagAll: "All tags",
        },
        kind: {
          research: "Research",
          blog: "Blog",
        },
        sections: {
          outcome: "Outcomes",
          outlook: "Next steps",
        },
        externalLink: "Open external article",
        noEntriesForTag: "No entries tagged “{{tag}}” yet.",
        updatedOn: "Updated {{date}}",
        error: "Research entries could not be loaded. Please retry shortly.",
      },
      projects: {
        tagline: "Projects",
        title: "Project archive",
        description:
          "Highlight technical projects with tech stack chips, links, and filtering options.",
        placeholder:
          "Project cards will be rendered here once the API endpoints are connected.",
        filters: {
          all: "All stacks",
        },
        noMatchesForSelection: "No projects match the selected tech stack yet.",
        error: "Projects could not be retrieved from the API.",
        highlight: {
          title: "Highlight projects",
        },
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
          description:
            "Reserve a conversation slot that aligns with calendar availability.",
          groupLabel: "Available time slots",
          unavailable: "No open slots are currently available.",
          slotTo: "Ends {{end}}",
          timezone: "Times are displayed in {{timezone}}.",
          error: "Calendar availability could not be refreshed.",
        },
        form: {
          legend: "Request details",
          name: "Your name",
          email: "Email address",
          topic: "Topic",
          topicPlaceholder: "Select a topic",
          message: "Message",
          slot: "Time slot",
          slotPlaceholder: "Select a time slot",
          timezoneLabel: "Times display in {{timezone}}",
          slotTime: "Local time",
          view: {
            single: "Single day",
            multi: "Multi-day",
            previous: "Previous",
            next: "Next",
          },
          legendLabels: {
            available: "Available",
            reserved: "Reserved",
            blackout: "Blackout",
          },
          status: {
            available: "Available",
            reserved: "Reserved",
            blackout: "Unavailable",
          },
          noAvailability: "No availability in this range.",
          submit: "Request booking",
          submitting: "Submitting…",
          success:
            "Thank you! Your request (ID: {{id}}) was recorded. A confirmation will be sent to {{email}}.",
          error: "We could not complete the booking. Please try again later.",
          configError: "Form configuration could not be loaded.",
          consent:
            "Your details are used solely for scheduling and will be removed after the meeting concludes.",
          errors: {
            nameRequired: "Please provide your name.",
            emailRequired: "An email address is required.",
            emailInvalid: "Please provide a valid email address.",
            topicRequired: "Select a topic to help us route your request.",
            messageLength:
              "Share at least 20 characters so we can prepare effectively.",
            slotRequired: "Select an available time slot.",
            slotUnavailable: "The selected time slot is no longer available.",
          },
        },
        summary: {
          title: "Booking summary",
          timezone: "Times shown in {{timezone}}.",
          timezoneFallback: "local time",
          window: "Bookings accepted up to {{days}} days ahead.",
          supportEmail: "For urgent updates contact",
          calendarLinked: "Automatic Google Calendar syncing is enabled.",
        },
        bookingSummary: {
          title: "Reservation details",
          when: "Scheduled for {{datetime}}",
          calendarEvent: "Calendar event ID: {{id}}",
          lookup: "Confirmation code: {{hash}}",
          supportEmail: "Need help? Email {{email}}",
        },
        error: "Contact configuration could not be loaded. Please try again later.",
      },
      admin: {
        tagline: "Admin",
        title: "Management console",
        description:
          "Authenticated administrators can manage content, reservations, and blacklists from the dedicated UI.",
        loginCallout:
          "Continue to the secure admin workspace to curate content and respond to booking requests.",
        loginCta: "Sign in with Google",
        loginHelp:
          "You'll authenticate with Google SSO. Only allow-listed administrator accounts can access the console, and you will be redirected back to the dashboard after signing in.",
      },
      footer: {
        copyright: "Takumi Tokunaga. All rights reserved.",
        links: {
          github: "GitHub",
          twitter: "X (Twitter)",
          contact: "Email",
        },
      },
    },
  },
  ja: {
    translation: {
      branding: {
        subtitle: "ポートフォリオ",
        title: "徳永 拓未",
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
        admin: "管理",
      },
      themeToggle: {
        setLight: "ライトテーマに切り替え",
        setDark: "ダークテーマに切り替え",
      },
      languages: {
        en: "English",
        ja: "日本語",
      },
      common: {
        presentLabel: "現在",
      },
      home: {
        hero: {
          tagline: "実世界データと共創するイノベーション",
          title: "RAG 学習支援とロボティクスの実装に挑戦しています。",
          description:
            "Go × React による本番運用を前提としたポートフォリオ。研究活動、プロジェクト事例、インターン経験、問い合わせ導線を一体的に公開します。",
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
          error:
            "最新のプロフィール情報を取得できませんでした。時間をおいて再度お試しください。",
          lastUpdated: "最終更新日時",
        },
        health: {
          title: "API ステータス",
          caption: "Go バックエンドから取得した最新ステータスを表示します。",
        },
        social: {
          title: "つながる",
          connectWith: "{{label}} でつながる",
          placeholder: "ソーシャルリンクは設定後に表示されます。",
        },
        quickLinks: {
          title: "ショートカット",
          projects: "プロジェクト",
          research: "研究",
          contact: "お問い合わせ",
          supporting:
            "詳細は上部のナビゲーションから各ページへ移動してください。",
        },
        tech: {
          title: "技術スタック",
          description: "現在取り組む案件で活用している主要技術とサポート技術です。",
        },
        work: {
          title: "直近の取り組み",
          description: "責務と成果が明確な取り組みを抜粋して掲載しています。",
          empty: "職務経歴は後ほど公開予定です。",
        },
        chips: {
          empty: "ハイライト情報は準備が整い次第表示されます。",
        },
        affiliations: {
          title: "所属サマリー",
        },
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
            empty: "所属情報は現在準備中です。",
          },
          skills: {
            title: "スキルセット",
            description: "構造化されたスキル情報からスキルマップを表示します。",
            empty: "スキル情報は後ほど公開予定です。",
            level: {
              beginner: "初級",
              intermediate: "中級",
              advanced: "上級",
              expert: "エキスパート",
            },
          },
          lab: {
            title: "研究室",
            description: "現在所属している研究室や指導教員の情報です。",
            name: "研究室",
            advisor: "指導教員",
            focus: "研究テーマ",
            room: "居室",
            visit: "研究室サイトへ",
            empty: "研究室情報は準備中です。",
          },
          work: {
            title: "職務経歴",
            description:
              "研究テーマと関連する職務経験やインターン実績を掲載します。",
            empty: "職務経歴は現在整理中です。",
          },
          communities: {
            title: "コミュニティ",
            description:
              "オープンソースや学術・産業コミュニティへの参画情報です。",
            empty: "コミュニティ参加情報は後日追加されます。",
          },
          social: {
            title: "ソーシャルリンク",
          },
        },
        error:
          "プロフィール情報の取得に失敗しました。ページを再読み込みしてください。",
      },
      research: {
        tagline: "研究",
        title: "研究ポートフォリオ",
        description:
          "Markdown ベースの研究成果や発表資料をレスポンシブに表示します。",
        placeholder: "API 連携後に研究コンテンツがここに描画されます。",
        filters: {
          all: "すべて",
          research: "研究",
          blog: "ブログ",
          tags: "タグで絞り込む",
          tagAll: "すべてのタグ",
        },
        kind: {
          research: "研究",
          blog: "ブログ",
        },
        sections: {
          outcome: "成果",
          outlook: "次のアクション",
        },
        externalLink: "外部記事を開く",
        noEntriesForTag: "「{{tag}}」タグのエントリはまだありません。",
        updatedOn: "{{date}} 更新",
        error:
          "研究情報を取得できませんでした。時間をおいて再度お試しください。",
      },
      projects: {
        tagline: "プロジェクト",
        title: "プロジェクト一覧",
        description:
          "使用技術やリンク・フィルター付きのプロジェクトカードを表示します。",
        placeholder: "API 連携後にプロジェクトカードがここに表示されます。",
        filters: {
          all: "すべて",
        },
        noMatchesForSelection:
          "選択した技術スタックに該当するプロジェクトはありません。",
        error: "プロジェクト情報を取得できませんでした。",
        highlight: {
          title: "注目プロジェクト",
        },
      },
      contact: {
        tagline: "お問い合わせ",
        title: "コンタクト",
        description:
          "予約可能枠や問い合わせフォームをバックエンドと連携して表示します。",
        placeholder: "フォームと予約ウィジェットをこのセクションに配置します。",
        availability: {
          title: "空き状況",
          description:
            "カレンダーの空き状況に基づいてミーティング枠を予約できます。",
          groupLabel: "予約可能な時間枠",
          unavailable: "現在予約可能な枠はありません。",
          slotTo: "終了 {{end}}",
          timezone: "表示時刻は {{timezone}} を基準としています。",
          error: "空き状況の取得に失敗しました。",
        },
        form: {
          legend: "予約内容",
          name: "お名前",
          email: "メールアドレス",
          topic: "トピック",
          topicPlaceholder: "トピックを選択してください",
          message: "メッセージ",
          slot: "時間枠",
          slotPlaceholder: "時間枠を選択してください",
          timezoneLabel: "表示時刻: {{timezone}}",
          slotTime: "ローカル時刻",
          view: {
            single: "1日表示",
            multi: "複数日表示",
            previous: "前へ",
            next: "次へ",
          },
          legendLabels: {
            available: "予約可能",
            reserved: "予約済み",
            blackout: "ブロック",
          },
          status: {
            available: "予約可能",
            reserved: "予約済み",
            blackout: "予約不可",
          },
          noAvailability: "この期間に空き枠はありません。",
          submit: "予約をリクエスト",
          submitting: "送信中…",
          success:
            "ありがとうございます。リクエスト (ID: {{id}}) を受け付けました。確認メールを {{email}} に送信します。",
          error:
            "予約リクエストを送信できませんでした。時間をおいて再度お試しください。",
          configError: "フォーム設定の取得に失敗しました。",
          consent:
            "ご入力いただいた情報は日程調整のみに利用し、ミーティング終了後に削除します。",
          errors: {
            nameRequired: "お名前を入力してください。",
            emailRequired: "メールアドレスを入力してください。",
            emailInvalid: "有効なメールアドレス形式で入力してください。",
            topicRequired: "トピックを選択してください。",
            messageLength: "具体的な内容を 20 文字以上で入力してください。",
            slotRequired: "予約する時間枠を選択してください。",
            slotUnavailable: "選択した時間枠は利用できなくなりました。",
          },
        },
        summary: {
          title: "予約に関する情報",
          timezone: "表示時刻は {{timezone}} 基準です。",
          timezoneFallback: "ローカルタイム",
          window: "{{days}} 日先まで予約を受け付けています。",
          supportEmail: "お急ぎの場合は次のメールへご連絡ください：",
          calendarLinked: "Google カレンダーへの自動連携が有効です。",
        },
        bookingSummary: {
          title: "予約内容の確認",
          when: "{{datetime}} に予定されています。",
          calendarEvent: "カレンダーイベント ID: {{id}}",
          lookup: "確認コード: {{hash}}",
          supportEmail: "連絡先メールアドレス: {{email}}",
        },
        error:
          "お問い合わせ設定を取得できませんでした。時間をおいて再度お試しください。",
      },
      admin: {
        tagline: "管理",
        title: "管理者コンソール",
        description:
          "認証済み管理者がコンテンツ、予約、ブラックリストを管理する UI を提供します。",
        loginCallout:
          "管理者専用コンソールにサインインし、コンテンツ更新や予約対応を行えます。",
        loginCta: "Google でサインイン",
        loginHelp:
          "Google SSO を利用し、許可された管理者アカウントのみがアクセスできます。認証後は管理コンソールに自動的に戻ります。",
      },
      footer: {
        copyright: "Takumi Tokunaga. All rights reserved.",
        links: {
          github: "GitHub",
          twitter: "X (Twitter)",
          contact: "メール",
        },
      },
    },
  },
};

const initialLanguageResolution = resolveInitialLanguage();
const initialLanguage: SupportedLanguage = initialLanguageResolution.language;

void i18n.use(initReactI18next).init({
  resources,
  fallbackLng: FALLBACK_LANGUAGE,
  lng: initialLanguage,
  supportedLngs: SUPPORTED_LANGUAGES,
  react: {
    useSuspense: false,
  },
  initImmediate: false,
  interpolation: {
    escapeValue: false,
  },
});

persistLanguagePreference(initialLanguage);

const handleLanguageChanged = (language: string) => {
  const matched = matchSupportedLanguage(language);
  if (!matched) {
    return;
  }
  persistLanguagePreference(matched);
};

i18n.off("languageChanged", handleLanguageChanged);
i18n.on("languageChanged", handleLanguageChanged);

export default i18n;
