import type {
  HomePageConfig,
  ProfileResponse,
  Project,
  ResearchEntry,
  TechCatalogEntry,
  TechMembership,
} from "./public-api/types";

type SupportedLocale = "en" | "ja";

function resolveLocale(locale?: string): SupportedLocale {
  if (locale && locale.toLowerCase().startsWith("ja")) {
    return "ja";
  }
  return "en";
}

function clone<T>(value: T): T {
  if (typeof structuredClone === "function") {
    return structuredClone(value);
  }
  return JSON.parse(JSON.stringify(value)) as T;
}

const techCatalog: Record<
  string,
  TechCatalogEntry
> = {
  typescript: {
    id: "tech-typescript",
    slug: "typescript",
    displayName: "TypeScript",
    category: "frontend",
    level: "advanced",
    icon: undefined,
    sortOrder: 1,
    active: true,
  },
  react: {
    id: "tech-react",
    slug: "react",
    displayName: "React",
    category: "frontend",
    level: "advanced",
    icon: undefined,
    sortOrder: 2,
    active: true,
  },
  go: {
    id: "tech-go",
    slug: "go",
    displayName: "Go",
    category: "backend",
    level: "advanced",
    icon: undefined,
    sortOrder: 3,
    active: true,
  },
  gcp: {
    id: "tech-gcp",
    slug: "google-cloud",
    displayName: "Google Cloud",
    category: "infrastructure",
    level: "advanced",
    icon: undefined,
    sortOrder: 4,
    active: true,
  },
  firestore: {
    id: "tech-firestore",
    slug: "cloud-firestore",
    displayName: "Cloud Firestore",
    category: "data",
    level: "advanced",
    icon: undefined,
    sortOrder: 5,
    active: true,
  },
  terraform: {
    id: "tech-terraform",
    slug: "terraform",
    displayName: "Terraform",
    category: "infrastructure",
    level: "intermediate",
    icon: undefined,
    sortOrder: 6,
    active: true,
  },
  python: {
    id: "tech-python",
    slug: "python",
    displayName: "Python",
    category: "ml",
    level: "advanced",
    icon: undefined,
    sortOrder: 7,
    active: true,
  },
  ros: {
    id: "tech-ros",
    slug: "ros",
    displayName: "ROS / ROS2",
    category: "robotics",
    level: "intermediate",
    icon: undefined,
    sortOrder: 8,
    active: true,
  },
  unity: {
    id: "tech-unity",
    slug: "unity",
    displayName: "Unity",
    category: "xr",
    level: "intermediate",
    icon: undefined,
    sortOrder: 9,
    active: true,
  },
};

function createMembership(
  tech: TechCatalogEntry,
  sortOrder: number,
  context: "primary" | "supporting" = "primary",
  note?: string,
): TechMembership {
  return {
    id: `${tech.id}-${context}-${sortOrder}`,
    context,
    note,
    sortOrder,
    tech: clone(tech),
  };
}

const canonicalHomeConfigs: Record<SupportedLocale, HomePageConfig> = {
  en: {
    heroSubtitle: "Software engineer and real-world information student",
    quickLinks: [
      {
        id: "quick-profile",
        section: "profile",
        label: "Academic profile",
        description: "Affiliations, lab details, and study focus areas.",
        cta: "View profile",
        targetUrl: "/profile",
        sortOrder: 1,
      },
      {
        id: "quick-research",
        section: "research_blog",
        label: "Research & blog",
        description: "Integrated updates on papers, talks, and field notes.",
        cta: "Read updates",
        targetUrl: "/research",
        sortOrder: 2,
      },
      {
        id: "quick-projects",
        section: "projects",
        label: "Project archive",
        description: "Selected engineering work with technology highlights.",
        cta: "Browse projects",
        targetUrl: "/projects",
        sortOrder: 3,
      },
      {
        id: "quick-contact",
        section: "contact",
        label: "Book a time",
        description: "Coordinate conversations with calendar integration.",
        cta: "Request booking",
        targetUrl: "/contact",
        sortOrder: 4,
      },
    ],
    chipSources: [
      {
        id: "chips-tech",
        source: "tech",
        label: "Focus technologies",
        limit: 6,
        sortOrder: 1,
      },
      {
        id: "chips-affiliation",
        source: "affiliation",
        label: "Active affiliations",
        limit: 3,
        sortOrder: 2,
      },
      {
        id: "chips-community",
        source: "community",
        label: "Community work",
        limit: 3,
        sortOrder: 3,
      },
    ],
    updatedAt: "2024-11-01T00:00:00Z",
  },
  ja: {
    heroSubtitle: "実世界情報コース所属のソフトウェアエンジニア",
    quickLinks: [
      {
        id: "quick-profile",
        section: "profile",
        label: "プロフィール",
        description: "所属・研究室・学習テーマの概要。",
        cta: "プロフィールを見る",
        targetUrl: "/profile",
        sortOrder: 1,
      },
      {
        id: "quick-research",
        section: "research_blog",
        label: "研究・ブログ",
        description: "論文・登壇・フィールドノートを統合した更新情報。",
        cta: "記事を読む",
        targetUrl: "/research",
        sortOrder: 2,
      },
      {
        id: "quick-projects",
        section: "projects",
        label: "プロジェクト一覧",
        description: "技術ハイライト付きの開発実績。",
        cta: "プロジェクトを見る",
        targetUrl: "/projects",
        sortOrder: 3,
      },
      {
        id: "quick-contact",
        section: "contact",
        label: "お問い合わせ / 予約",
        description: "Google カレンダー連携による面談調整。",
        cta: "日程をリクエスト",
        targetUrl: "/contact",
        sortOrder: 4,
      },
    ],
    chipSources: [
      {
        id: "chips-tech",
        source: "tech",
        label: "得意な技術領域",
        limit: 6,
        sortOrder: 1,
      },
      {
        id: "chips-affiliation",
        source: "affiliation",
        label: "在籍・所属",
        limit: 3,
        sortOrder: 2,
      },
      {
        id: "chips-community",
        source: "community",
        label: "コミュニティ活動",
        limit: 3,
        sortOrder: 3,
      },
    ],
    updatedAt: "2024-11-01T00:00:00Z",
  },
};

function buildProfile(
  locale: SupportedLocale,
  base: Omit<ProfileResponse, "footerLinks" | "home">,
): ProfileResponse {
  const socialLinks = base.socialLinks.map((link) => clone(link));
  const footerLinks = socialLinks.filter((link) => link.isFooter);
  return {
    ...base,
    socialLinks,
    footerLinks,
    home: clone(canonicalHomeConfigs[locale]),
  };
}

const canonicalProfiles: Record<SupportedLocale, ProfileResponse> = {
  en: buildProfile("en", {
    id: "profile-canonical",
    displayName: "Takumi Tokunaga",
    headline:
      "Real-world information engineering student and full-stack engineer building reliable systems.",
    summary:
      "Undergraduate researcher at Ritsumeikan University working on retrieval-augmented assistants, robotics-driven experiences, and resilient cloud services. Combines startup leadership with production engineering to deliver measurable outcomes end to end.",
    avatarUrl: undefined,
    location: "Osaka, Japan",
    theme: {
      mode: "system",
      accentColor: "#0ea5e9",
    },
    lab: {
      name: "Kimura Laboratory (RM²C)",
      advisor: "Professor Hiroshi Kimura",
      room: "AN236",
      url: "https://www.rm2c.ise.ritsumei.ac.jp/",
    },
    affiliations: [
      {
        id: "affiliation-rm2c",
        name: "Kimura Laboratory (RM²C)",
        url: "https://www.rm2c.ise.ritsumei.ac.jp/",
        description: "Undergraduate researcher focusing on human-computer interaction and XR.",
        startedAt: "2025-04-01",
        sortOrder: 1,
      },
      {
        id: "affiliation-ritsumeikan",
        name: "Ritsumeikan University · Real-world Information Program",
        url: "https://www.ritsumei.ac.jp/",
        description: "B.S. candidate in Information Science and Engineering.",
        startedAt: "2023-04-01",
        sortOrder: 2,
      },
    ],
    communities: [
      {
        id: "community-robocup",
        name: "RoboCup Kansai Student Branch",
        description: "Regional robotics outreach and competition programming.",
        startedAt: "2024-04-01",
        sortOrder: 1,
      },
      {
        id: "community-rrat",
        name: "Ritsumeikan Robotics & Automation Team",
        description: "Cross-disciplinary robotics project-based community.",
        startedAt: "2023-04-01",
        sortOrder: 2,
      },
      {
        id: "community-stem",
        name: "STEM Outreach Workshops",
        description: "Instructor for programming and engineering bootcamps.",
        startedAt: "2022-04-01",
        sortOrder: 3,
      },
    ],
    workHistory: [
      {
        id: "work-reboze",
        organization: "Reboze LLC",
        role: "Business & Technology Consultant (Intern)",
        summary:
          "Designed architecture proposals for enterprise transformation projects and coordinated discovery workshops.",
        startedAt: "2025-10-01",
        endedAt: null,
        externalUrl: "https://reboze.com/",
        sortOrder: 1,
      },
      {
        id: "work-proseeds",
        organization: "Proseeds Inc.",
        role: "Software Engineer (Intern)",
        summary:
          "Maintained e-learning services serving 1M+ learners and prototyped retrieval-augmented support experiences.",
        startedAt: "2025-02-01",
        endedAt: "2025-09-30",
        externalUrl: "https://www.pro-seeds.com/",
        sortOrder: 2,
      },
      {
        id: "work-tier2",
        organization: "Tier2 LLC",
        role: "Co-founder / Managing Partner",
        summary:
          "Built inventory automation tooling and analytics for a multi-channel resale business with 20% margin lift.",
        startedAt: "2023-09-01",
        endedAt: null,
        externalUrl: undefined,
        sortOrder: 3,
      },
    ],
    techSections: [
      {
        id: "tech-core",
        title: "Core engineering",
        layout: "grid",
        breakpoint: "md",
        sortOrder: 1,
        members: [
          createMembership(techCatalog.typescript, 1),
          createMembership(techCatalog.react, 2),
          createMembership(techCatalog.go, 3),
        ],
      },
      {
        id: "tech-infra",
        title: "Infrastructure & data",
        layout: "grid",
        breakpoint: "md",
        sortOrder: 2,
        members: [
          createMembership(techCatalog.gcp, 1),
          createMembership(techCatalog.terraform, 2, "supporting"),
          createMembership(techCatalog.firestore, 3, "supporting"),
        ],
      },
      {
        id: "tech-robotics",
        title: "Robotics & XR",
        layout: "grid",
        breakpoint: "md",
        sortOrder: 3,
        members: [
          createMembership(techCatalog.ros, 1, "supporting"),
          createMembership(techCatalog.unity, 2, "supporting"),
          createMembership(techCatalog.python, 3),
        ],
      },
    ],
    socialLinks: [
      {
        id: "github",
        provider: "github",
        label: "ttokunaga-jp",
        url: "https://github.com/ttokunaga-jp",
        isFooter: true,
        sortOrder: 1,
      },
      {
        id: "zenn",
        provider: "zenn",
        label: "Zenn",
        url: "https://zenn.dev/ttokunaga",
        isFooter: true,
        sortOrder: 2,
      },
      {
        id: "linkedin",
        provider: "linkedin",
        label: "LinkedIn",
        url: "https://www.linkedin.com/in/takumi-tokunaga/",
        isFooter: true,
        sortOrder: 3,
      },
      {
        id: "email",
        provider: "email",
        label: "is0732hk@ed.ritsumei.ac.jp",
        url: "mailto:is0732hk@ed.ritsumei.ac.jp",
        isFooter: false,
        sortOrder: 4,
      },
    ],
    updatedAt: "2024-11-01T00:00:00Z",
  }),
  ja: buildProfile("ja", {
    id: "profile-canonical-ja",
    displayName: "徳永 拓未",
    headline:
      "立命館大学 実世界情報コース / フルスタックエンジニアとして信頼性ある体験を設計。",
    summary:
      "立命館大学にて RAG 支援エージェントやロボティクスを活用した体験設計、堅牢なクラウドサービス開発に取り組む学部研究員。起業と現場開発の経験を生かし、要件定義から運用まで一気通貫で価値創出を行う。",
    avatarUrl: undefined,
    location: "日本 大阪府",
    theme: {
      mode: "system",
      accentColor: "#0ea5e9",
    },
    lab: {
      name: "木村研究室（RM²C）",
      advisor: "木村 浩 先生",
      room: "AN236",
      url: "https://www.rm2c.ise.ritsumei.ac.jp/",
    },
    affiliations: [
      {
        id: "affiliation-rm2c",
        name: "木村研究室（RM²C モバイルコンピューティング／リアリティメディア研究室）",
        url: "https://www.rm2c.ise.ritsumei.ac.jp/",
        description: "ヒューマンコンピュータインタラクションと XR の研究に従事。",
        startedAt: "2025-04-01",
        sortOrder: 1,
      },
      {
        id: "affiliation-ritsumeikan",
        name: "立命館大学 情報理工学部 実世界情報コース",
        url: "https://www.ritsumei.ac.jp/",
        description: "実世界情報プログラム B3 / 研究室配属前期。",
        startedAt: "2023-04-01",
        sortOrder: 2,
      },
    ],
    communities: [
      {
        id: "community-robocup",
        name: "RoboCup Kansai Student Branch",
        description: "高校生向けロボット競技の普及と教育支援を担当。",
        startedAt: "2024-04-01",
        sortOrder: 1,
      },
      {
        id: "community-rrat",
        name: "Ritsumeikan Robotics & Automation Team",
        description: "学内外のロボティクス共同開発コミュニティ。",
        startedAt: "2023-04-01",
        sortOrder: 2,
      },
      {
        id: "community-stem",
        name: "STEM Outreach ワークショップ",
        description: "小中高校生向けプログラミング講座の運営・講師を担当。",
        startedAt: "2022-04-01",
        sortOrder: 3,
      },
    ],
    workHistory: [
      {
        id: "work-reboze",
        organization: "Reboze LLC",
        role: "ビジネス / テクノロジーコンサルタント（長期インターン）",
        summary:
          "企業の業務変革に向けたシステムアーキテクチャ検討、PoC 調整、顧客課題のヒアリングを推進。",
        startedAt: "2025-10-01",
        endedAt: null,
        externalUrl: "https://reboze.com/",
        sortOrder: 1,
      },
      {
        id: "work-proseeds",
        organization: "株式会社プロシーズ",
        role: "ソフトウェアエンジニア（長期インターン）",
        summary:
          "100 万人以上が利用する e ラーニング基盤の開発・保守、RAG を活用したサポート自動化を担当。",
        startedAt: "2025-02-01",
        endedAt: "2025-09-30",
        externalUrl: "https://www.pro-seeds.com/",
        sortOrder: 2,
      },
      {
        id: "work-tier2",
        organization: "合同会社 Tier2",
        role: "共同創業者 / マネージングパートナー",
        summary:
          "物販自動化ツールの内製化・在庫分析による粗利改善、古物商許可の取得と多拠点運営をリード。",
        startedAt: "2023-09-01",
        endedAt: null,
        externalUrl: undefined,
        sortOrder: 3,
      },
    ],
    techSections: [
      {
        id: "tech-core",
        title: "ソフトウェアエンジニアリング",
        layout: "grid",
        breakpoint: "md",
        sortOrder: 1,
        members: [
          createMembership(techCatalog.typescript, 1),
          createMembership(techCatalog.react, 2),
          createMembership(techCatalog.go, 3),
        ],
      },
      {
        id: "tech-infra",
        title: "インフラ・データ基盤",
        layout: "grid",
        breakpoint: "md",
        sortOrder: 2,
        members: [
          createMembership(techCatalog.gcp, 1),
          createMembership(techCatalog.terraform, 2, "supporting"),
          createMembership(techCatalog.firestore, 3, "supporting"),
        ],
      },
      {
        id: "tech-robotics",
        title: "ロボティクス・XR",
        layout: "grid",
        breakpoint: "md",
        sortOrder: 3,
        members: [
          createMembership(techCatalog.ros, 1, "supporting"),
          createMembership(techCatalog.unity, 2, "supporting"),
          createMembership(techCatalog.python, 3),
        ],
      },
    ],
    socialLinks: [
      {
        id: "github",
        provider: "github",
        label: "ttokunaga-jp",
        url: "https://github.com/ttokunaga-jp",
        isFooter: true,
        sortOrder: 1,
      },
      {
        id: "zenn",
        provider: "zenn",
        label: "Zenn",
        url: "https://zenn.dev/ttokunaga",
        isFooter: true,
        sortOrder: 2,
      },
      {
        id: "linkedin",
        provider: "linkedin",
        label: "LinkedIn",
        url: "https://www.linkedin.com/in/takumi-tokunaga/",
        isFooter: true,
        sortOrder: 3,
      },
      {
        id: "email",
        provider: "email",
        label: "is0732hk@ed.ritsumei.ac.jp",
        url: "mailto:is0732hk@ed.ritsumei.ac.jp",
        isFooter: false,
        sortOrder: 4,
      },
    ],
    updatedAt: "2024-11-01T00:00:00Z",
  }),
};

const canonicalProjectsMap: Record<SupportedLocale, Project[]> = {
  en: [
    {
      id: "project-knowledge-guide",
      slug: "knowledge-guide",
      title: "Knowledge Guide",
      summary: "Retrieval-augmented learning assistant for robotics clubs.",
      description:
        "Built a retrieval-augmented assistant that consolidates lab notes, CAD documents, and mentor Q&A. Deployed on Google Cloud Run with robust observability, GDPR compliant data retention, and multilingual content search.",
      coverImageUrl: "https://images.takumi.dev/projects/knowledge-guide.jpg",
      primaryLink: "https://github.com/ttokunaga-jp/knowledge-guide",
      links: [
        {
          id: "knowledge-guide-repo",
          type: "repo",
          label: "Source repository",
          url: "https://github.com/ttokunaga-jp/knowledge-guide",
          sortOrder: 1,
        },
        {
          id: "knowledge-guide-demo",
          type: "demo",
          label: "Interactive demo",
          url: "https://demo.takumi.dev/knowledge-guide",
          sortOrder: 2,
        },
      ],
      period: {
        start: "2024-05-01",
        end: null,
      },
      tech: [
        createMembership(techCatalog.go, 1),
        createMembership(techCatalog.typescript, 2, "supporting"),
        createMembership(techCatalog.gcp, 3, "primary"),
      ],
      highlight: true,
      published: true,
      sortOrder: 1,
      createdAt: "2024-05-01T00:00:00Z",
      updatedAt: "2024-10-12T00:00:00Z",
    },
    {
      id: "project-robotics-insight",
      slug: "robotics-insight",
      title: "Robotics Insight Platform",
      summary: "Telemetry pipeline and dashboards for field robotics trials.",
      description:
        "Designed a streaming telemetry pipeline capturing ROS2 topics, edge perception, and operator annotations. Provided low-latency dashboards for competition organizers and exported data to BigQuery for long-term analysis.",
      coverImageUrl: "https://images.takumi.dev/projects/robotics-insight.jpg",
      primaryLink: "https://github.com/ttokunaga-jp/robotics-insight",
      links: [
        {
          id: "robotics-insight-repo",
          type: "repo",
          label: "Source repository",
          url: "https://github.com/ttokunaga-jp/robotics-insight",
          sortOrder: 1,
        },
        {
          id: "robotics-insight-paper",
          type: "article",
          label: "Field report",
          url: "https://takumi.dev/blog/robotics-insight",
          sortOrder: 2,
        },
      ],
      period: {
        start: "2024-01-01",
        end: "2024-08-31",
      },
      tech: [
        createMembership(techCatalog.ros, 1),
        createMembership(techCatalog.python, 2),
        createMembership(techCatalog.gcp, 3, "supporting"),
      ],
      highlight: false,
      published: true,
      sortOrder: 2,
      createdAt: "2024-01-15T00:00:00Z",
      updatedAt: "2024-08-31T00:00:00Z",
    },
    {
      id: "project-campus-nav",
      slug: "campus-navigation-xr",
      title: "Campus Navigation XR",
      summary: "Mixed reality indoor navigation prototype for campus tours.",
      description:
        "Implemented a Unity-based XR navigation prototype with spatial anchors, route optimisation, and haptic cues for accessibility. Evaluated in collaboration with university admissions and mobility support teams.",
      coverImageUrl: "https://images.takumi.dev/projects/campus-navigation.jpg",
      primaryLink: "https://takumi.dev/projects/campus-navigation-xr",
      links: [
        {
          id: "campus-navigation-slides",
          type: "slides",
          label: "Presentation slides",
          url: "https://speakerdeck.com/takumi/campus-navigation-xr",
          sortOrder: 1,
        },
      ],
      period: {
        start: "2023-09-01",
        end: "2024-03-31",
      },
      tech: [
        createMembership(techCatalog.unity, 1),
        createMembership(techCatalog.react, 2, "supporting"),
        createMembership(techCatalog.gcp, 3, "supporting"),
      ],
      highlight: false,
      published: true,
      sortOrder: 3,
      createdAt: "2023-09-01T00:00:00Z",
      updatedAt: "2024-04-01T00:00:00Z",
    },
  ],
  ja: [
    {
      id: "project-knowledge-guide",
      slug: "knowledge-guide",
      title: "Knowledge Guide",
      summary: "ロボット競技部向けの RAG 学習支援アシスタント。",
      description:
        "部内ナレッジ・CAD・メンター Q&A を統合する検索アシスタント。Cloud Run 上で動作し、可観測性と多言語検索、GDPR 準拠のデータ保持設計を実装した。",
      coverImageUrl: "https://images.takumi.dev/projects/knowledge-guide.jpg",
      primaryLink: "https://github.com/ttokunaga-jp/knowledge-guide",
      links: [
        {
          id: "knowledge-guide-repo",
          type: "repo",
          label: "ソースコード",
          url: "https://github.com/ttokunaga-jp/knowledge-guide",
          sortOrder: 1,
        },
        {
          id: "knowledge-guide-demo",
          type: "demo",
          label: "デモサイト",
          url: "https://demo.takumi.dev/knowledge-guide",
          sortOrder: 2,
        },
      ],
      period: {
        start: "2024-05-01",
        end: null,
      },
      tech: [
        createMembership(techCatalog.go, 1),
        createMembership(techCatalog.typescript, 2, "supporting"),
        createMembership(techCatalog.gcp, 3, "primary"),
      ],
      highlight: true,
      published: true,
      sortOrder: 1,
      createdAt: "2024-05-01T00:00:00Z",
      updatedAt: "2024-10-12T00:00:00Z",
    },
    {
      id: "project-robotics-insight",
      slug: "robotics-insight",
      title: "Robotics Insight Platform",
      summary: "ロボット実験のテレメトリ取得と可視化基盤。",
      description:
        "ROS2 トピックやエッジ推論、オペレーターの注釈を統合するストリーミングパイプラインを構築。Cloud Run と BigQuery で低レイテンシなダッシュボードと長期分析を実現した。",
      coverImageUrl: "https://images.takumi.dev/projects/robotics-insight.jpg",
      primaryLink: "https://github.com/ttokunaga-jp/robotics-insight",
      links: [
        {
          id: "robotics-insight-repo",
          type: "repo",
          label: "ソースコード",
          url: "https://github.com/ttokunaga-jp/robotics-insight",
          sortOrder: 1,
        },
        {
          id: "robotics-insight-paper",
          type: "article",
          label: "フィールドレポート",
          url: "https://takumi.dev/blog/robotics-insight",
          sortOrder: 2,
        },
      ],
      period: {
        start: "2024-01-01",
        end: "2024-08-31",
      },
      tech: [
        createMembership(techCatalog.ros, 1),
        createMembership(techCatalog.python, 2),
        createMembership(techCatalog.gcp, 3, "supporting"),
      ],
      highlight: false,
      published: true,
      sortOrder: 2,
      createdAt: "2024-01-15T00:00:00Z",
      updatedAt: "2024-08-31T00:00:00Z",
    },
    {
      id: "project-campus-nav",
      slug: "campus-navigation-xr",
      title: "Campus Navigation XR",
      summary: "キャンパスツアー向けの MR 屋内ナビゲーション試作。",
      description:
        "空間アンカーとハプティクスを活用した Unity ベースのナビゲーションプロトタイプ。入試広報やバリアフリー支援チームと共同で運用検証を実施。",
      coverImageUrl: "https://images.takumi.dev/projects/campus-navigation.jpg",
      primaryLink: "https://takumi.dev/projects/campus-navigation-xr",
      links: [
        {
          id: "campus-navigation-slides",
          type: "slides",
          label: "発表資料",
          url: "https://speakerdeck.com/takumi/campus-navigation-xr",
          sortOrder: 1,
        },
      ],
      period: {
        start: "2023-09-01",
        end: "2024-03-31",
      },
      tech: [
        createMembership(techCatalog.unity, 1),
        createMembership(techCatalog.react, 2, "supporting"),
        createMembership(techCatalog.gcp, 3, "supporting"),
      ],
      highlight: false,
      published: true,
      sortOrder: 3,
      createdAt: "2023-09-01T00:00:00Z",
      updatedAt: "2024-04-01T00:00:00Z",
    },
  ],
};

const canonicalResearchMap: Record<SupportedLocale, ResearchEntry[]> = {
  en: [
    {
      id: "research-rag-evaluation",
      slug: "rag-evaluation-lab-notes",
      kind: "research",
      title: "Evaluating RAG pipelines for robotics field guides",
      overview:
        "Measured retrieval quality and operator satisfaction when providing robotics troubleshooting assistance via RAG.",
      outcome:
        "Achieved 27% reduction in time-to-resolution for common hardware faults and increased mentor coverage by 3x.",
      outlook:
        "Next iteration focuses on multimodal documentation ingestion and alignment with safety runbooks.",
      externalUrl: "https://takumi.dev/research/rag-evaluation",
      publishedAt: "2024-09-18T00:00:00Z",
      updatedAt: "2024-10-05T00:00:00Z",
      highlightImageUrl: "https://images.takumi.dev/research/rag-evaluation.jpg",
      imageAlt: "Evaluation dashboard visualising retrieval metrics.",
      isDraft: false,
      tags: ["retrieval", "robotics", "evaluation"],
      links: [
        {
          id: "rag-evaluation-paper",
          type: "paper",
          label: "Field study report",
          url: "https://takumi.dev/files/rag-evaluation-report.pdf",
          sortOrder: 1,
        },
        {
          id: "rag-evaluation-slides",
          type: "slides",
          label: "Talk slides",
          url: "https://speakerdeck.com/takumi/rag-evaluation",
          sortOrder: 2,
        },
      ],
      assets: [],
      tech: [
        createMembership(techCatalog.go, 1),
        createMembership(techCatalog.python, 2),
        createMembership(techCatalog.firestore, 3, "supporting"),
      ],
    },
    {
      id: "blog-lab-onboarding",
      slug: "lab-onboarding-guide",
      kind: "blog",
      title: "Lab onboarding guide: aligning real-world information practice",
      overview:
        "Documented onboarding practices for RM²C lab members covering experiment logging, asset management, and mentorship loops.",
      outcome:
        "Standardised onboarding enabled new members to ship production contributions within their first month.",
      outlook:
        "Plan to evolve into an interactive handbook integrated with Knowledge Guide and calendar automations.",
      externalUrl: "https://takumi.dev/blog/lab-onboarding-guide",
      publishedAt: "2024-06-01T00:00:00Z",
      updatedAt: "2024-07-15T00:00:00Z",
      highlightImageUrl: "https://images.takumi.dev/research/lab-onboarding.jpg",
      imageAlt: "Students collaborating with laptops and lab equipment.",
      isDraft: false,
      tags: ["lab", "operations", "handbook"],
      links: [
        {
          id: "lab-onboarding-notion",
          type: "external",
          label: "Notion playbook",
          url: "https://takumi.dev/notion/lab-handbook",
          sortOrder: 1,
        },
      ],
      assets: [],
      tech: [
        createMembership(techCatalog.typescript, 1, "supporting"),
        createMembership(techCatalog.gcp, 2, "supporting"),
      ],
    },
  ],
  ja: [
    {
      id: "research-rag-evaluation",
      slug: "rag-evaluation-lab-notes",
      kind: "research",
      title: "RAG を活用したロボット現場支援の評価",
      overview:
        "ロボット競技チーム向けに構築した RAG パイプラインの検索品質と現場満足度を測定。",
      outcome:
        "ハードウェア故障対応の平均時間を 27% 短縮し、メンター対応可能な案件を 3 倍に拡大。",
      outlook:
        "次フェーズではマルチモーダル資料の取り込みと安全運用 Runbook 連携を予定。",
      externalUrl: "https://takumi.dev/research/rag-evaluation",
      publishedAt: "2024-09-18T00:00:00Z",
      updatedAt: "2024-10-05T00:00:00Z",
      highlightImageUrl: "https://images.takumi.dev/research/rag-evaluation.jpg",
      imageAlt: "検索評価ダッシュボードのスクリーンショット。",
      isDraft: false,
      tags: ["RAG", "ロボティクス", "評価"],
      links: [
        {
          id: "rag-evaluation-paper",
          type: "paper",
          label: "フィールドスタディ報告書",
          url: "https://takumi.dev/files/rag-evaluation-report.pdf",
          sortOrder: 1,
        },
        {
          id: "rag-evaluation-slides",
          type: "slides",
          label: "発表資料",
          url: "https://speakerdeck.com/takumi/rag-evaluation",
          sortOrder: 2,
        },
      ],
      assets: [],
      tech: [
        createMembership(techCatalog.go, 1),
        createMembership(techCatalog.python, 2),
        createMembership(techCatalog.firestore, 3, "supporting"),
      ],
    },
    {
      id: "blog-lab-onboarding",
      slug: "lab-onboarding-guide",
      kind: "blog",
      title: "研究室オンボーディングガイド：実世界情報の実践をそろえる",
      overview:
        "RM²C 研究室の新規メンバー向けに、実験ログ、資産管理、メンタリングのループを整理。",
      outcome:
        "標準化されたオンボーディングにより、参加 1 か月以内での本番貢献が可能に。",
      outlook:
        "今後は Knowledge Guide やカレンダー自動化と統合したインタラクティブハンドブックへ発展予定。",
      externalUrl: "https://takumi.dev/blog/lab-onboarding-guide",
      publishedAt: "2024-06-01T00:00:00Z",
      updatedAt: "2024-07-15T00:00:00Z",
      highlightImageUrl: "https://images.takumi.dev/research/lab-onboarding.jpg",
      imageAlt: "研究室メンバーがノート PC と機材で議論する様子。",
      isDraft: false,
      tags: ["ラボ運営", "ドキュメント", "ナレッジ管理"],
      links: [
        {
          id: "lab-onboarding-notion",
          type: "external",
          label: "Notion プレイブック",
          url: "https://takumi.dev/notion/lab-handbook",
          sortOrder: 1,
        },
      ],
      assets: [],
      tech: [
        createMembership(techCatalog.typescript, 1, "supporting"),
        createMembership(techCatalog.gcp, 2, "supporting"),
      ],
    },
  ],
};

export const canonicalProfileEn = canonicalProfiles.en;
export const canonicalProfileJa = canonicalProfiles.ja;

export const canonicalProjectsEn = canonicalProjectsMap.en;
export const canonicalProjectsJa = canonicalProjectsMap.ja;

export const canonicalResearchEntriesEn = canonicalResearchMap.en;
export const canonicalResearchEntriesJa = canonicalResearchMap.ja;

export function getCanonicalProfile(locale?: string): ProfileResponse {
  const resolved = resolveLocale(locale);
  return clone(canonicalProfiles[resolved]);
}

export function getCanonicalProjects(locale?: string): Project[] {
  const resolved = resolveLocale(locale);
  return canonicalProjectsMap[resolved].map((project) => clone(project));
}

export function getCanonicalResearchEntries(locale?: string): ResearchEntry[] {
  const resolved = resolveLocale(locale);
  return canonicalResearchMap[resolved].map((entry) => clone(entry));
}

export function getCanonicalHomeConfig(locale?: string): HomePageConfig {
  const resolved = resolveLocale(locale);
  return clone(canonicalHomeConfigs[resolved]);
}
