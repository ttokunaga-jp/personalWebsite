import { http, HttpResponse } from "msw";
import { setupServer } from "msw/node";

import type {
  CreateBookingPayload,
  BookingResult,
} from "../modules/public-api";

import {
  contactAvailabilityFixture,
  contactConfigFixture,
  defaultBookingResponse,
  profileFixture,
  projectsFixture,
  researchEntriesFixture,
} from "./fixtures";

export { cloneFixture } from "./fixtures";
export {
  contactAvailabilityFixture,
  contactConfigFixture,
  defaultBookingResponse,
  profileFixture,
  projectsFixture,
  researchEntriesFixture,
};

export const defaultHandlers = [
  http.get("/api/health", () => HttpResponse.json({ status: "ok" })),
  http.get("/api/v1/public/profile", () =>
    HttpResponse.json({ data: profileFixture }),
  ),
  http.get("/api/v1/public/research", () =>
    HttpResponse.json({ data: researchEntriesFixture }),
  ),
  http.get("/api/v1/public/projects", () =>
    HttpResponse.json({ data: projectsFixture }),
  ),
  http.get("/api/v1/public/contact/availability", () =>
    HttpResponse.json({ data: contactAvailabilityFixture }),
  ),
  http.get("/api/v1/public/contact/config", () =>
    HttpResponse.json({ data: contactConfigFixture }),
  ),
  http.post("/api/v1/public/contact/bookings", async ({ request }) => {
    const payload = (await request.json()) as CreateBookingPayload;
    const result: BookingResult = {
      ...defaultBookingResponse,
      meeting: {
        ...defaultBookingResponse.meeting,
        id: String(Date.now()),
        name: payload.name,
        email: payload.email,
        datetime: payload.startTime,
        durationMinutes: payload.durationMinutes,
        notes: payload.agenda,
      },
      calendarEventId: defaultBookingResponse.calendarEventId,
    };
    return HttpResponse.json({ data: result });
  }),
  http.get("https://www.google.com/recaptcha/api.js", () =>
    HttpResponse.text("void grecaptcha;"),
  ),
];

export const server = setupServer(...defaultHandlers);
