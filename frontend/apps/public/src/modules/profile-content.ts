import type {
  ProfileResponse,
  Project,
  ResearchEntry,
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

const canonicalProfileEn: ProfileResponse = {
  name: "Takumi Tokunaga",
  headline: "Real-world information engineering student & full-stack engineer",
  summary:
    "Undergraduate researcher at Ritsumeikan University building retrieval-augmented learning assistants, hybrid search infrastructure, and robotics-driven experiences. Combines startup leadership with large-scale web service development to deliver measurable outcomes end to end.",
  avatarUrl: undefined,
  location: "Osaka, Japan",
  affiliations: [
    {
      id: "rm2c-lab",
      organization: "Kimura Laboratory (RM²C)",
      role: "Undergraduate Researcher",
      startDate: "2025-04-01",
      endDate: null,
      location: "Osaka, Japan",
      isCurrent: true,
    },
    {
      id: "ritsumeikan-university",
      organization: "Ritsumeikan University",
      department:
        "College of Information Science and Engineering · Real-world Information Program",
      role: "B.S. Candidate",
      startDate: "2023-04-01",
      endDate: null,
      location: "Osaka, Japan",
      isCurrent: true,
    },
    {
      id: "marugame-high-school",
      organization: "Marugame High School",
      role: "Student",
      startDate: "2019-04-01",
      endDate: "2022-03-31",
      location: "Kagawa, Japan",
      isCurrent: false,
    },
  ],
  lab: {
    name: "Kimura Laboratory (RM²C Mobile Computing / Reality Media Lab)",
    advisor: "Prof. Hiroshi Kimura",
    researchFocus:
      "Human-computer interaction, multimodal XR, mobile computing, and haptics.",
    websiteUrl: "https://www.rm2c.ise.ritsumei.ac.jp/",
  },
  workHistory: [
    {
      id: "reboze-llc",
      organization: "Reboze LLC",
      role: "Business & Technology Consultant (Intern)",
      startDate: "2025-10-01",
      endDate: null,
      achievements: [
        "Designed system architecture proposals for B2B transformation projects and supported customer discovery.",
        "Coordinated cross-functional PoC initiatives spanning engineering and business development.",
      ],
      description:
        "Long-term internship within the consulting division, focusing on revenue operations optimisation and technical due diligence.",
      location: "Osaka, Japan / Remote",
    },
    {
      id: "proseeds-inc",
      organization: "Proseeds Inc.",
      role: "Software Engineer (Intern)",
      startDate: "2025-02-01",
      endDate: "2025-09-30",
      achievements: [
        "Led development of a RAG-based support chatbot for an e-learning platform.",
        "Maintained production services serving 1M+ learners with a focus on reliability and accessibility.",
      ],
      description:
        "Worked within the engineering division on SaaS learning products, covering requirements, API design, and testing.",
      location: "Osaka, Japan",
    },
    {
      id: "tier2-llc",
      organization: "Tier2 LLC",
      role: "Co-founder / Managing Partner",
      startDate: "2023-09-01",
      endDate: null,
      achievements: [
        "Automated multi-channel resale operations with in-house inventory and pricing tooling.",
        "Obtained official antique dealer certification from the Kagawa Public Safety Commission.",
      ],
      description:
        "Founded a resale and software venture, owning strategy, engineering, and financial operations.",
      location: "Kagawa, Japan / Remote",
    },
  ],
  skillGroups: [
    {
      id: "software-engineering",
      category: "Software Engineering",
      items: [
        { id: "go", name: "Go", level: "advanced" },
        { id: "typescript", name: "TypeScript", level: "advanced" },
        { id: "react", name: "React", level: "advanced" },
        { id: "grpc", name: "Gin / gRPC", level: "intermediate" },
      ],
    },
    {
      id: "infrastructure-data",
      category: "Infrastructure & Data",
      items: [
        {
          id: "gcp",
          name: "GCP (Cloud Run, Cloud Build, Artifact Registry)",
          level: "advanced",
        },
        { id: "terraform", name: "Terraform", level: "intermediate" },
        { id: "firestore", name: "Cloud Firestore", level: "advanced" },
        {
          id: "elasticsearch",
          name: "Elasticsearch / Qdrant",
          level: "intermediate",
        },
      ],
    },
    {
      id: "robotics-xr",
      category: "Robotics & XR",
      items: [
        { id: "ros", name: "ROS / ROS2", level: "intermediate" },
        { id: "unity", name: "Unity / XR", level: "intermediate" },
        { id: "python", name: "Python (Simulation, ML)", level: "advanced" },
      ],
    },
  ],
  communities: [
    "RoboCup Kansai Student Branch",
    "Ritsumeikan Robotics & Automation Team",
    "STEM Outreach Mentor (Programming Workshops)",
  ],
  socialLinks: [
    {
      id: "github",
      platform: "github",
      label: "ttokunaga-jp",
      url: "https://github.com/ttokunaga-jp",
    },
    {
      id: "lab",
      platform: "website",
      label: "RM²C Laboratory",
      url: "https://www.rm2c.ise.ritsumei.ac.jp/",
    },
    {
      id: "email",
      platform: "email",
      label: "is0732hk@ed.ritsumei.ac.jp",
      url: "mailto:is0732hk@ed.ritsumei.ac.jp",
    },
  ],
};

const canonicalProfileJa: ProfileResponse = {
  name: "徳永 拓未",
  headline: "立命館大学 情報理工学部 実世界情報コース / フルスタックエンジニア",
  summary:
    "ロボティクスと XR を核に実世界とサイバー空間をつなぐ体験設計を探究しつつ、RAG や検索基盤、教育向けプロダクトを Go / TypeScript / GCP で開発しています。起業経験と大規模サービス開発の現場経験を掛け合わせ、企画から運用まで一気通貫で価値提供することを目指しています。",
  avatarUrl: undefined,
  location: "大阪府, 日本",
  affiliations: [
    {
      id: "rm2c-lab",
      organization: "木村研究室（RM²C モバイルコンピューティング／リアリティメディア研究室）",
      role: "学部研究員",
      startDate: "2025-04-01",
      endDate: null,
      location: "大阪府, 日本",
      isCurrent: true,
    },
    {
      id: "ritsumeikan-university",
      organization: "立命館大学 情報理工学部 実世界情報コース",
      department:
        "College of Information Science and Engineering · Real-world Information Program",
      role: "学部生（B.S.）",
      startDate: "2023-04-01",
      endDate: null,
      location: "大阪府, 日本",
      isCurrent: true,
    },
    {
      id: "marugame-high-school",
      organization: "香川県立丸亀高等学校",
      role: "学生",
      startDate: "2019-04-01",
      endDate: "2022-03-31",
      location: "香川県, 日本",
      isCurrent: false,
    },
  ],
  lab: {
    name: "木村研究室（RM²C モバイルコンピューティング／リアリティメディア研究室）",
    advisor: "木村 浩嗣 教授",
    researchFocus:
      "ヒューマンコンピュータインタラクション、マルチモーダル XR、モバイルコンピューティング、ハプティクスを対象に研究。",
    websiteUrl: "https://www.rm2c.ise.ritsumei.ac.jp/",
  },
  workHistory: [
    {
      id: "reboze-llc",
      organization: "Reboze LLC",
      role: "ビジネス＆テクノロジーコンサルタント（長期インターン）",
      startDate: "2025-10-01",
      endDate: null,
      achievements: [
        "B2B 変革プロジェクト向けにシステムアーキテクチャを設計し、カスタマーディスカバリを支援。",
        "エンジニアリングと事業開発を横断する PoC を調整し、クロスファンクショナルに推進。",
      ],
      description:
        "コンサルティング部門にて収益オペレーション最適化と技術デューデリジェンスに取り組む長期インターンシップ。",
      location: "大阪 / リモート",
    },
    {
      id: "proseeds-inc",
      organization: "株式会社プロシーズ",
      role: "ソフトウェアエンジニア（インターン）",
      startDate: "2025-02-01",
      endDate: "2025-09-30",
      achievements: [
        "eラーニング向け RAG サポートチャットボットの開発をリード。",
        "100 万人超の学習者が利用する本番サービスを信頼性とアクセシビリティの観点で保守運用。",
      ],
      description:
        "SaaS 型学習プロダクトの開発部門で要件定義・API 設計・テストを担当。",
      location: "大阪府",
    },
    {
      id: "tier2-llc",
      organization: "Tier2 LLC",
      role: "共同創業者 / マネージングパートナー",
      startDate: "2023-09-01",
      endDate: null,
      achievements: [
        "在庫・価格管理ツールを内製化し、マルチチャネル転売オペレーションを自動化。",
        "香川県公安委員会より古物商許可を取得し、リユース事業を立ち上げ。",
      ],
      description:
        "戦略立案からエンジニアリング、財務までを担うソフトウェア＆リセールベンチャーを運営。",
      location: "香川県 / リモート",
    },
  ],
  skillGroups: [
    {
      id: "software-engineering",
      category: "ソフトウェアエンジニアリング",
      items: [
        { id: "go", name: "Go", level: "advanced" },
        { id: "typescript", name: "TypeScript", level: "advanced" },
        { id: "react", name: "React", level: "advanced" },
        { id: "grpc", name: "Gin / gRPC", level: "intermediate" },
      ],
    },
    {
      id: "infrastructure-data",
      category: "インフラ / データ基盤",
      items: [
        {
          id: "gcp",
          name: "GCP（Cloud Run, Cloud Build, Artifact Registry）",
          level: "advanced",
        },
        { id: "terraform", name: "Terraform", level: "intermediate" },
        { id: "firestore", name: "Cloud Firestore", level: "advanced" },
        {
          id: "elasticsearch",
          name: "Elasticsearch / Qdrant",
          level: "intermediate",
        },
      ],
    },
    {
      id: "robotics-xr",
      category: "ロボティクス / XR",
      items: [
        { id: "ros", name: "ROS / ROS2", level: "intermediate" },
        { id: "unity", name: "Unity / XR", level: "intermediate" },
        { id: "python", name: "Python（Simulation, ML）", level: "advanced" },
      ],
    },
  ],
  communities: [
    "RoboCup 関西 学生支部",
    "立命館 Robotics & Automation Team",
    "STEM Outreach Mentor（プログラミングワークショップ）",
  ],
  socialLinks: [
    {
      id: "github",
      platform: "github",
      label: "ttokunaga-jp",
      url: "https://github.com/ttokunaga-jp",
    },
    {
      id: "lab",
      platform: "website",
      label: "RM²C Laboratory",
      url: "https://www.rm2c.ise.ritsumei.ac.jp/",
    },
    {
      id: "email",
      platform: "email",
      label: "is0732hk@ed.ritsumei.ac.jp",
      url: "mailto:is0732hk@ed.ritsumei.ac.jp",
    },
  ],
};

const canonicalProjectsEn: Project[] = [
  {
    id: "project-classnav",
    title: "ClassNav: RAG-powered Learning Assistant",
    subtitle: "Adaptive study companion for LMS content",
    description:
      "Integrates with the university LMS to ingest lecture artefacts automatically, parse diverse formats via Docling, and surface cited responses backed by hybrid retrieval. Designed for reliability, latency, and educator control.",
    techStack: [
      "Go",
      "TypeScript",
      "LangChain",
      "Docling",
      "Cloud Run",
      "Firestore",
    ],
    category: "Education",
    tags: ["RAG", "Education", "GenAI"],
    period: {
      start: "2024-04-01",
      end: null,
    },
    links: [
      {
        label: "Repository",
        url: "https://github.com/ttokunaga-jp/classnav",
        type: "repo",
      },
      {
        label: "Project overview",
        url: "https://github.com/ttokunaga-jp/classnav#readme",
        type: "article",
      },
    ],
    coverImageUrl: undefined,
    highlight: true,
  },
  {
    id: "project-search-service",
    title: "searchService: Hybrid Retrieval Microservice",
    subtitle: "Elasticsearch × Qdrant microservice for RAG pipelines",
    description:
      "Implements weighted scoring across keyword and vector indices with Kafka-based ingestion, DLQ retries, and full observability (OpenTelemetry + Prometheus + Jaeger). Shared across personal products to standardise retrieval.",
    techStack: ["Go", "gRPC", "Elasticsearch", "Qdrant", "Kafka", "Docker"],
    category: "Platform",
    tags: ["Search", "Microservices", "Observability"],
    period: {
      start: "2024-06-01",
      end: "2024-12-31",
    },
    links: [
      {
        label: "Repository",
        url: "https://github.com/ttokunaga-jp/search-service",
        type: "repo",
      },
      {
        label: "Architecture notes",
        url: "https://github.com/ttokunaga-jp/search-service/blob/main/docs/architecture.md",
        type: "article",
      },
    ],
    coverImageUrl: undefined,
    highlight: false,
  },
  {
    id: "project-portfolio",
    title: "Personal Website & Admin Console",
    subtitle: "Go API + React SPA deployed on Cloud Run",
    description:
      "End-to-end portfolio stack featuring public content, admin workflows, booking automation, and CI/CD (GitHub Actions → Cloud Build). Infrastructure codified with Terraform for multi-environment rollouts.",
    techStack: ["Go", "TypeScript", "React", "Firestore", "Cloud Run", "Terraform"],
    category: "Portfolio",
    tags: ["Full-stack", "CI/CD", "Infrastructure"],
    period: {
      start: "2025-01-01",
      end: null,
    },
    links: [
      {
        label: "Repository",
        url: "https://github.com/ttokunaga-jp/personalWebsite",
        type: "repo",
      },
    ],
    coverImageUrl: undefined,
    highlight: false,
  },
];

const canonicalProjectsJa: Project[] = [
  {
    id: "project-classnav",
    title: "RAG 学習支援システム「ClassNav」",
    subtitle: "LMS 連携のアダプティブラーニングコンパニオン",
    description:
      "大学の LMS と連携し、講義資料を自動で取り込み要約・検索できる学習支援システム。Docling によるマルチフォーマット解析と RAG を組み合わせ、引用リンク提示や自動アップロード機能で NotebookLM と差別化。",
    techStack: [
      "Go",
      "TypeScript",
      "LangChain",
      "Docling",
      "Cloud Run",
      "Firestore",
    ],
    category: "Education",
    tags: ["RAG", "Education", "GenAI"],
    period: {
      start: "2024-04-01",
      end: null,
    },
    links: [
      {
        label: "リポジトリ",
        url: "https://github.com/ttokunaga-jp/classnav",
        type: "repo",
      },
      {
        label: "プロジェクト概要",
        url: "https://github.com/ttokunaga-jp/classnav#readme",
        type: "article",
      },
    ],
    coverImageUrl: undefined,
    highlight: true,
  },
  {
    id: "project-search-service",
    title: "searchService: ハイブリッド検索マイクロサービス",
    subtitle: "Elasticsearch × Qdrant を統合した RAG 基盤",
    description:
      "Elasticsearch と Qdrant を統合し、キーワード・ベクトルを加重スコアリングする検索基盤。Kafka による非同期インデックス更新や DLQ リトライ、OpenTelemetry / Prometheus / Jaeger を備え、RAG サービスの共通モジュールとして運用。",
    techStack: ["Go", "gRPC", "Elasticsearch", "Qdrant", "Kafka", "Docker"],
    category: "Platform",
    tags: ["Search", "Microservices", "Observability"],
    period: {
      start: "2024-06-01",
      end: "2024-12-31",
    },
    links: [
      {
        label: "リポジトリ",
        url: "https://github.com/ttokunaga-jp/search-service",
        type: "repo",
      },
      {
        label: "アーキテクチャノート",
        url: "https://github.com/ttokunaga-jp/search-service/blob/main/docs/architecture.md",
        type: "article",
      },
    ],
    coverImageUrl: undefined,
    highlight: false,
  },
  {
    id: "project-portfolio",
    title: "個人ポートフォリオサイト（Go + React）",
    subtitle: "Cloud Run 上の Go API + React SPA",
    description:
      "Go (Gin) と React SPA で構築した個人ポートフォリオ。公開サイトと管理 GUI を分離し、予約フォーム、Google OAuth + JWT 認証、Cloud Build → Cloud Run の CI/CD、Terraform による IaC を整備。",
    techStack: ["Go", "TypeScript", "React", "Firestore", "Cloud Run", "Terraform"],
    category: "Portfolio",
    tags: ["Full-stack", "CI/CD", "Infrastructure"],
    period: {
      start: "2025-01-01",
      end: null,
    },
    links: [
      {
        label: "リポジトリ",
        url: "https://github.com/ttokunaga-jp/personalWebsite",
        type: "repo",
      },
    ],
    coverImageUrl: undefined,
    highlight: false,
  },
];

const canonicalResearchEntriesEn: ResearchEntry[] = [
  {
    id: "rag-reliability-learning-assistants",
    title: "Improving Reliability of RAG-based Learning Assistants",
    slug: "rag-learning-assistants",
    summary:
      "Evaluates hybrid retrieval and citation UX strategies for LMS-focused learning assistants, emphasising transparency and educator oversight.",
    publishedOn: "2025-06-01",
    updatedOn: "2025-09-15",
    tags: ["Education", "RAG", "GenAI"],
    contentMarkdown: `## Overview

ClassNav synchronises lecture artefacts from the campus LMS, parses them with Docling, and maintains a hybrid Elasticsearch/Qdrant index. We compared retrieval strategies, citation formats, and educator review tools to reduce hallucinations and build learner confidence.

### Findings

- Hybrid scoring (keyword × vector) improved answer precision by 21% compared to pure dense retrieval.
- Linking to source pages with highlighted anchors decreased follow-up questions by 34%.
- Educators valued the ability to approve or redact responses before publication.

### Next steps

Extending the evaluation to timed quizzes and grading assistants while integrating telemetry for long-term monitoring.`,
    contentHtml:
      "<h2>Overview</h2><p>ClassNav synchronises lecture artefacts from the campus LMS, parses them with Docling, and maintains a hybrid Elasticsearch/Qdrant index. We compared retrieval strategies, citation formats, and educator review tools to reduce hallucinations and build learner confidence.</p><h3>Findings</h3><ul><li>Hybrid scoring (keyword × vector) improved answer precision by 21% compared to pure dense retrieval.</li><li>Linking to source pages with highlighted anchors decreased follow-up questions by 34%.</li><li>Educators valued the ability to approve or redact responses before publication.</li></ul><h3>Next steps</h3><p>Extending the evaluation to timed quizzes and grading assistants while integrating telemetry for long-term monitoring.</p>",
    assets: [
      {
        alt: "Screenshot of ClassNav retrieval diagnostics dashboard",
        url: "https://github.com/ttokunaga-jp",
        caption: "Prototype dashboard visualising hybrid retrieval scores per query.",
      },
    ],
    links: [
      {
        label: "Project repository",
        url: "https://github.com/ttokunaga-jp",
        type: "code",
      },
    ],
  },
];

const canonicalResearchEntriesJa: ResearchEntry[] = [
  {
    id: "rag-reliability-learning-assistants",
    title: "RAG による学習支援システムの信頼性向上に関する研究",
    slug: "rag-learning-assistants",
    summary:
      "大学 LMS の資料を対象に、Docling を活用した抽出とベクトル検索を組み合わせ、引用リンク提示でハルシネーションを抑制する学習支援基盤を設計。",
    publishedOn: "2025-06-01",
    updatedOn: "2025-09-15",
    tags: ["Education", "RAG", "GenAI"],
    contentMarkdown: `## 研究概要

LMS から自動取得した講義資料を Docling で構造化し、Elasticsearch と Qdrant を統合したハイブリッド検索でベクトル類似度とキーワードスコアを重み付けします。回答には参照元のリンクを必ず添付し、ユーザーが検証しやすい UI を React で実装しました。RAG パイプライン全体を可観測化し、再現性のある評価ワークフローを整備しています。

### 成果

- キーワード × ベクトルのハイブリッドスコアリングで、純粋なベクトル検索と比較して回答精度が 21% 向上。
- 参照元ページへのリンクとハイライトを提示することで、問い合わせ件数を 34% 削減。
- 教員が公開前に回答を承認・修正できるワークフローで、運用上の透明性を確保。

### 次のステップ

時間指定クイズや採点支援への適用と、長期的なモニタリングに向けたテレメトリ統合を進めています。`,
    assets: [
      {
        alt: "ClassNav の検索診断ダッシュボード",
        url: "https://github.com/ttokunaga-jp",
        caption: "ハイブリッド検索のスコアをクエリ単位で可視化したプロトタイプ。",
      },
    ],
    links: [
      {
        label: "プロジェクトリポジトリ",
        url: "https://github.com/ttokunaga-jp",
        type: "code",
      },
    ],
  },
];

const profileByLocale: Record<SupportedLocale, ProfileResponse> = {
  en: canonicalProfileEn,
  ja: canonicalProfileJa,
};

const projectsByLocale: Record<SupportedLocale, Project[]> = {
  en: canonicalProjectsEn,
  ja: canonicalProjectsJa,
};

const researchByLocale: Record<SupportedLocale, ResearchEntry[]> = {
  en: canonicalResearchEntriesEn,
  ja: canonicalResearchEntriesJa,
};

export function getCanonicalProfile(locale?: string): ProfileResponse {
  return clone(profileByLocale[resolveLocale(locale)]);
}

export function getCanonicalProjects(locale?: string): Project[] {
  return clone(projectsByLocale[resolveLocale(locale)]);
}

export function getCanonicalResearchEntries(locale?: string): ResearchEntry[] {
  return clone(researchByLocale[resolveLocale(locale)]);
}

export {
  canonicalProfileEn,
  canonicalProfileJa,
  canonicalProjectsEn,
  canonicalProjectsJa,
  canonicalResearchEntriesEn,
  canonicalResearchEntriesJa,
};
