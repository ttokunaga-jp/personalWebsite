import type {
  ContactAvailabilityResponse,
  ContactConfigResponse,
  CreateBookingResponse,
  ProfileResponse,
  Project,
  ResearchEntry
} from "../modules/public-api";

const now = new Date();
const iso = (date: Date) => date.toISOString();

export const profileFixture: ProfileResponse = {
  name: "Takumi Asano",
  headline: "Lead Research Engineer",
  summary: "Building privacy-conscious, human-centered systems.",
  location: "Tokyo, Japan",
  affiliations: [
    {
      id: "aff-1",
      organization: "University of Tokyo",
      department: "Intelligent Systems Laboratory",
      role: "Graduate Researcher",
      startDate: "2022-04-01",
      endDate: null,
      location: "Tokyo, Japan",
      isCurrent: true
    },
    {
      id: "aff-2",
      organization: "Open Source Robotics Foundation",
      role: "Contributor",
      startDate: "2021-01-01",
      endDate: "2022-03-31",
      location: "Remote",
      isCurrent: false
    }
  ],
  lab: {
    name: "Intelligent Systems Lab",
    advisor: "Dr. Naomi Kobayashi",
    researchFocus: "Human-robot collaboration in manufacturing environments",
    websiteUrl: "https://example.com/lab"
  },
  workHistory: [
    {
      id: "work-1",
      organization: "Example Robotics Inc.",
      role: "Research Engineer (Intern)",
      startDate: "2023-07-01",
      endDate: "2023-09-30",
      achievements: [
        "Delivered simulation tooling that reduced iteration time by 36%.",
        "Co-authored internal whitepaper on safety testing."
      ],
      description: "Developed tooling for real-time telemetry ingestion and anomaly detection.",
      location: "Tokyo, Japan"
    }
  ],
  skillGroups: [
    {
      id: "skills-1",
      category: "Languages",
      items: [
        {
          id: "ts",
          name: "TypeScript",
          level: "expert"
        },
        {
          id: "go",
          name: "Go",
          level: "advanced"
        }
      ]
    }
  ],
  communities: ["OSS Robotics", "AI Research Society"],
  socialLinks: [
    {
      id: "link-github",
      platform: "github",
      label: "GitHub",
      url: "https://github.com/takumi"
    },
    {
      id: "link-twitter",
      platform: "x",
      label: "X",
      url: "https://twitter.com/takumi"
    },
    {
      id: "link-email",
      platform: "email",
      label: "Email",
      url: "mailto:hello@example.com"
    }
  ]
};

export const researchEntriesFixture: ResearchEntry[] = [
  {
    id: "research-1",
    title: "Embodied Agents in Crowded Warehouses",
    slug: "embodied-agents-warehouses",
    summary: "Exploring co-navigation strategies between autonomous agents and human operators.",
    publishedOn: iso(new Date(now.getTime() - 1000 * 60 * 60 * 24 * 30)),
    updatedOn: iso(new Date(now.getTime() - 1000 * 60 * 60 * 24 * 7)),
    tags: ["robots", "safety"],
    contentMarkdown: `## Abstract\n\nInvestigated shared control models for industrial cobots.`,
    contentHtml: "<h2>Abstract</h2><p>Investigated shared control models for industrial cobots.</p>",
    assets: [
      {
        alt: "Cobots navigating warehouse aisles",
        url: "https://example.com/assets/warehouse.jpg",
        caption: "Simulation environment used for trials."
      }
    ],
    links: [
      { label: "Full paper", url: "https://example.com/paper.pdf", type: "paper" },
      { label: "Code", url: "https://github.com/example/research", type: "code" }
    ]
  }
];

export const projectsFixture: Project[] = [
  {
    id: "project-1",
    title: "Edge Telemetry Kit",
    subtitle: "Unified observability for IoT fleets",
    description: "Portable toolkit for streaming telemetry from edge devices with adaptive sampling.",
    techStack: ["TypeScript", "Go", "GCP"],
    category: "Platform",
    tags: ["observability", "iot"],
    period: {
      start: "2023-01-01",
      end: null
    },
    links: [
      {
        label: "Repository",
        url: "https://github.com/example/edge-telemetry",
        type: "repo"
      }
    ],
    coverImageUrl: "https://example.com/assets/telemetry.jpg",
    highlight: true
  },
  {
    id: "project-2",
    title: "Research Portfolio CMS",
    description: "Headless CMS integration powering the public research showcase.",
    techStack: ["React", "Node.js"],
    tags: ["cms"],
    period: {
      start: "2022-05-01",
      end: "2022-12-31"
    },
    links: [],
    coverImageUrl: undefined,
    highlight: false
  }
];

export const contactAvailabilityFixture: ContactAvailabilityResponse = {
  timezone: "Asia/Tokyo",
  generatedAt: iso(now),
  slots: [
    {
      id: "slot-1",
      start: iso(new Date(now.getTime() + 1000 * 60 * 60 * 24)),
      end: iso(new Date(now.getTime() + 1000 * 60 * 60 * 24 + 30 * 60 * 1000)),
      isBookable: true
    },
    {
      id: "slot-2",
      start: iso(new Date(now.getTime() + 1000 * 60 * 60 * 48)),
      end: iso(new Date(now.getTime() + 1000 * 60 * 60 * 48 + 30 * 60 * 1000)),
      isBookable: false
    }
  ]
};

export const contactConfigFixture: ContactConfigResponse = {
  topics: ["Research collaboration", "Speaking engagement"],
  recaptchaSiteKey: "test-site-key",
  minimumLeadHours: 48,
  consentText: "We only use your information for scheduling purposes."
};

export const defaultBookingResponse: CreateBookingResponse = {
  bookingId: "bk-slot-1",
  status: "pending",
  calendarUrl: "https://calendar.example.com/bookings/bk-slot-1"
};

export function cloneFixture<T>(fixture: T): T {
  if (typeof structuredClone === "function") {
    return structuredClone(fixture);
  }

  return JSON.parse(JSON.stringify(fixture)) as T;
}
