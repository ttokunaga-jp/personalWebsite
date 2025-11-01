import {
  canonicalProfileEn,
  canonicalProjectsEn,
  canonicalResearchEntriesEn,
} from "../modules/profile-content";
import type {
  ContactAvailabilityResponse,
  ContactConfigResponse,
  BookingResult,
  ProfileResponse,
  Project,
  ResearchEntry,
} from "../modules/public-api";

const now = new Date();
const iso = (date: Date) => date.toISOString();

export function cloneFixture<T>(fixture: T): T {
  if (typeof structuredClone === "function") {
    return structuredClone(fixture);
  }

  return JSON.parse(JSON.stringify(fixture)) as T;
}

export const profileFixture: ProfileResponse =
  cloneFixture(canonicalProfileEn);

export const researchEntriesFixture: ResearchEntry[] = cloneFixture(
  canonicalResearchEntriesEn,
);

export const projectsFixture: Project[] =
  cloneFixture(canonicalProjectsEn);

export const contactAvailabilityFixture: ContactAvailabilityResponse = {
  timezone: "Asia/Tokyo",
  generatedAt: iso(now),
  days: [
    {
      date: iso(new Date(now.getTime() + 1000 * 60 * 60 * 24)).slice(0, 10),
      slots: [
        {
          id: "slot-1",
          start: iso(new Date(now.getTime() + 1000 * 60 * 60 * 24)),
          end: iso(
            new Date(now.getTime() + 1000 * 60 * 60 * 24 + 30 * 60 * 1000),
          ),
          isBookable: true,
        },
        {
          id: "slot-2",
          start: iso(new Date(now.getTime() + 1000 * 60 * 60 * 24 + 2 * 60 * 60 * 1000)),
          end: iso(
            new Date(now.getTime() + 1000 * 60 * 60 * 24 + 2 * 60 * 60 * 1000 + 30 * 60 * 1000),
          ),
          isBookable: false,
        },
      ],
    },
  ],
};

export const contactConfigFixture: ContactConfigResponse = {
  topics: ["Research collaboration", "Speaking engagement"],
  recaptchaSiteKey: "test-site-key",
  minimumLeadHours: 48,
  consentText: "We only use your information for scheduling purposes.",
};

export const defaultBookingResponse: BookingResult = {
  meeting: {
    id: 1,
    name: "Jane Doe",
    email: "jane.doe@example.com",
    datetime: iso(new Date(now.getTime() + 1000 * 60 * 60 * 24)),
    durationMinutes: 30,
    meetUrl: "https://meet.example.com/abc",
    calendarEventId: "event-1",
    status: "pending",
    notes: "Initial consultation",
  },
  calendarEventId: "event-1",
};
