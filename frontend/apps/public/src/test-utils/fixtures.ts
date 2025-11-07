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
          status: "available" as const,
        },
        {
          id: "slot-2",
          start: iso(new Date(now.getTime() + 1000 * 60 * 60 * 24 + 2 * 60 * 60 * 1000)),
          end: iso(
            new Date(now.getTime() + 1000 * 60 * 60 * 24 + 2 * 60 * 60 * 1000 + 30 * 60 * 1000),
          ),
          isBookable: false,
          status: "reserved" as const,
        },
      ],
    },
  ],
};

export const contactConfigFixture: ContactConfigResponse = {
  heroTitle: "Let us coordinate a session",
  heroDescription:
    "Share your context and a preferred slot to schedule a conversation with Takumi.",
  topics: [
    { id: "research", label: "Research collaboration" },
    { id: "speaking", label: "Speaking engagement" },
  ],
  minimumLeadHours: 48,
  consentText: "We only use your information for scheduling purposes.",
  recaptchaSiteKey: "test-site-key",
  supportEmail: "contact@example.com",
  calendarTimezone: "Asia/Tokyo",
  googleCalendarId: "calendar-id",
  bookingWindowDays: 14,
};

export const defaultBookingResponse: BookingResult = {
  reservation: {
    id: "bk-1",
    lookupHash: "lookup-bk-1",
    name: "Jane Doe",
    email: "jane.doe@example.com",
    topic: "consultation",
    message: "Initial consultation",
    startAt: iso(new Date(now.getTime() + 1000 * 60 * 60 * 24)),
    endAt: iso(new Date(now.getTime() + 1000 * 60 * 60 * 24 + 30 * 60 * 1000)),
    durationMinutes: 30,
    googleEventId: "event-1",
    googleCalendarStatus: "confirmed",
    status: "pending",
    confirmationSentAt: iso(
      new Date(now.getTime() + 1000 * 60 * 60 * 24 + 5 * 60 * 1000),
    ),
  },
  calendarEventId: "event-1",
  supportEmail: "contact@example.com",
  calendarTimezone: "Asia/Tokyo",
};
