import type {
  ProfileResponse,
  Project,
  ResearchEntry,
} from "./public-api/types";

export const canonicalProfile: ProfileResponse = {
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
        { id: "mysql", name: "MySQL / Cloud SQL", level: "advanced" },
        { id: "elasticsearch", name: "Elasticsearch / Qdrant", level: "intermediate" },
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

export const canonicalProjects: Project[] = [
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
      "MySQL",
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
    techStack: ["Go", "TypeScript", "React", "MySQL", "Cloud Run", "Terraform"],
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
    highlight: true,
  },
];

export const canonicalResearchEntries: ResearchEntry[] = [
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
