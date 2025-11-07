import { apiClient } from "@shared/lib/api-client";
import { createContext, useContext } from "react";
import { useTranslation } from "react-i18next";

import {
  type UseApiResourceResult,
  useApiResource,
} from "../../lib/useApiResource";

import {
  transformContactConfig,
  transformProfile,
  transformProjects,
  transformResearchEntries,
  type RawContactConfig,
  type RawProfileDocument,
  type RawProjectDocument,
  type RawResearchDocument,
} from "./transform";
import type {
  ContactAvailabilityResponse,
  ContactConfigResponse,
  CreateBookingPayload,
  BookingResult,
  MeetingReservationStatus,
  ProfileResponse,
  Project,
  ResearchEntry,
} from "./types";

const BASE_PATH = "/v1/public";
const USE_MOCK_PUBLIC_API =
  import.meta.env.MODE !== "test" &&
  (import.meta.env.VITE_USE_MOCK_PUBLIC_API ?? "false") === "true";

function withAbortSignal(signal: AbortSignal) {
  return { signal };
}

type ApiSuccessResponse<T> = {
  data: T;
};

function unwrapData<T>(response: ApiSuccessResponse<T>): T {
  return response.data;
}

function ensureNotAborted(signal: AbortSignal) {
  if (signal.aborted) {
    throw new DOMException("The operation was aborted.", "AbortError");
  }
}

function createMockAvailability(): ContactAvailabilityResponse {
  const now = new Date();
  now.setMinutes(0, 0, 0);

  const days: ContactAvailabilityResponse["days"] = [];
  for (let offset = 0; offset < 3; offset += 1) {
    const base = new Date(now);
    base.setDate(base.getDate() + offset);
    const slots = Array.from({ length: 3 }).map((_, index) => {
      const start = new Date(base);
      start.setHours(9 + index * 2, 0, 0, 0);
      const end = new Date(start);
      end.setHours(start.getHours() + 1);
      return {
        id: start.toISOString(),
        start: start.toISOString(),
        end: end.toISOString(),
        isBookable: true,
        status: "available" as const,
      };
    });
    days.push({
      date: base.toISOString().slice(0, 10),
      slots,
    });
  }

  return {
    timezone: "Asia/Tokyo",
    generatedAt: new Date().toISOString(),
    days,
  };
}

function createMockContactConfig(): ContactConfigResponse {
  return {
    heroTitle: "Schedule a conversation",
    heroDescription:
      "Pick a topic and a preferred timeslot to coordinate with Takumi. A confirmation email will follow after review.",
    topics: [
      { id: "research", label: "Research collaboration" },
      { id: "consultation", label: "Project consultation" },
      { id: "speaking", label: "Speaking engagement" },
      { id: "mentoring", label: "Mentoring session" },
    ],
    consentText:
      "Provided details are used only to coordinate the requested meeting. Expect a reply within two business days.",
    minimumLeadHours: 48,
    recaptchaSiteKey: import.meta.env.VITE_RECAPTCHA_SITE_KEY ?? "",
    supportEmail: "contact@example.com",
    calendarTimezone: "Asia/Tokyo",
    googleCalendarId: "mock-calendar",
    bookingWindowDays: 14,
  };
}

type RawMeetingReservation = {
  id?: number | string | null;
  lookupHash?: string | null;
  name?: string | null;
  email?: string | null;
  topic?: string | null;
  message?: string | null;
  startAt?: string | null;
  endAt?: string | null;
  durationMinutes?: number | null;
  googleEventId?: string | null;
  googleCalendarStatus?: string | null;
  status?: string | null;
  confirmationSentAt?: string | null;
  lastNotificationSentAt?: string | null;
  cancellationReason?: string | null;
  createdAt?: string | null;
  updatedAt?: string | null;
};

type RawBookingResult = {
  reservation?: RawMeetingReservation;
  calendarEventId?: string | null;
  supportEmail?: string | null;
  calendarTimezone?: string | null;
};

function normalizeString(value?: string | null): string | undefined {
  const trimmed = value?.trim();
  return trimmed ? trimmed : undefined;
}

function transformBookingResult(raw: RawBookingResult): BookingResult {
  const reservation = raw.reservation ?? {};
  const id =
    reservation.id !== undefined && reservation.id !== null
      ? String(reservation.id)
      : "";

  return {
    reservation: {
      id,
      lookupHash: normalizeString(reservation.lookupHash),
      name: normalizeString(reservation.name) ?? "",
      email: normalizeString(reservation.email) ?? "",
      topic: normalizeString(reservation.topic),
      message: normalizeString(reservation.message),
      startAt: reservation.startAt ?? "",
      endAt: reservation.endAt ?? "",
      durationMinutes: reservation.durationMinutes ?? 0,
      googleEventId: normalizeString(reservation.googleEventId),
      googleCalendarStatus: normalizeString(reservation.googleCalendarStatus),
      status: (reservation.status ?? "pending") as MeetingReservationStatus,
      confirmationSentAt: reservation.confirmationSentAt ?? undefined,
      lastNotificationSentAt: reservation.lastNotificationSentAt ?? undefined,
      cancellationReason: normalizeString(reservation.cancellationReason),
      createdAt: reservation.createdAt ?? undefined,
      updatedAt: reservation.updatedAt ?? undefined,
    },
    calendarEventId: normalizeString(raw.calendarEventId),
    supportEmail: normalizeString(raw.supportEmail),
    calendarTimezone: normalizeString(raw.calendarTimezone),
  };
}

export const ProfileResourceContext =
  createContext<UseApiResourceResult<ProfileResponse> | null>(null);

export const publicApi = {
  async getProfile(signal: AbortSignal): Promise<ProfileResponse> {
    if (USE_MOCK_PUBLIC_API) {
      ensureNotAborted(signal);
      return transformProfile(undefined);
    }

    const response = await apiClient.get<ApiSuccessResponse<RawProfileDocument>>(
      `${BASE_PATH}/profile`,
      {
        ...withAbortSignal(signal),
      },
    );
    return transformProfile(unwrapData(response.data));
  },
  async getResearch(signal: AbortSignal): Promise<ResearchEntry[]> {
    if (USE_MOCK_PUBLIC_API) {
      ensureNotAborted(signal);
      return transformResearchEntries(undefined);
    }

    const response = await apiClient.get<
      ApiSuccessResponse<RawResearchDocument[]>
    >(`${BASE_PATH}/research`, {
      ...withAbortSignal(signal),
    });
    return transformResearchEntries(unwrapData(response.data));
  },
  async getProjects(signal: AbortSignal): Promise<Project[]> {
    if (USE_MOCK_PUBLIC_API) {
      ensureNotAborted(signal);
      return transformProjects(undefined);
    }

    const response = await apiClient.get<ApiSuccessResponse<RawProjectDocument[]>>(
      `${BASE_PATH}/projects`,
      {
        ...withAbortSignal(signal),
      },
    );
    return transformProjects(unwrapData(response.data));
  },
  async getContactAvailability(
    signal: AbortSignal,
  ): Promise<ContactAvailabilityResponse> {
    if (USE_MOCK_PUBLIC_API) {
      ensureNotAborted(signal);
      return createMockAvailability();
    }

    const response = await apiClient.get<
      ApiSuccessResponse<ContactAvailabilityResponse>
    >(`${BASE_PATH}/contact/availability`, {
      ...withAbortSignal(signal),
    });
    return unwrapData(response.data);
  },
  async getContactConfig(signal: AbortSignal): Promise<ContactConfigResponse> {
    if (USE_MOCK_PUBLIC_API) {
      ensureNotAborted(signal);
      return createMockContactConfig();
    }

    const response = await apiClient.get<
      ApiSuccessResponse<RawContactConfig>
    >(`${BASE_PATH}/contact/config`, {
      ...withAbortSignal(signal),
    });
    return transformContactConfig(unwrapData(response.data));
  },
  async createBooking(
    payload: CreateBookingPayload,
  ): Promise<BookingResult> {
    if (USE_MOCK_PUBLIC_API) {
      return {
        reservation: {
          id: String(Date.now()),
          lookupHash: `mock-${Date.now()}`,
          name: payload.name,
          email: payload.email,
          topic: payload.topic,
          message: payload.agenda,
          startAt: payload.startTime,
          endAt: payload.startTime,
          durationMinutes: payload.durationMinutes,
          status: "pending",
          googleCalendarStatus: "confirmed",
        },
        calendarEventId: undefined,
        supportEmail: "contact@example.com",
        calendarTimezone: "Asia/Tokyo",
      };
    }

    const response = await apiClient.post<
      ApiSuccessResponse<RawBookingResult>
    >(`${BASE_PATH}/contact/bookings`, payload);
    return transformBookingResult(unwrapData(response.data));
  },
};

type UseProfileResourceOptions = {
  skip?: boolean;
};

export function useProfileResourceInternal(
  options?: UseProfileResourceOptions,
): UseApiResourceResult<ProfileResponse> {
  const { i18n } = useTranslation();
  const skip = options?.skip ?? USE_MOCK_PUBLIC_API;

  return useApiResource(publicApi.getProfile, {
    initialData: () => transformProfile(undefined),
    skip,
    dependencies: [i18n.language],
  });
}

export function useProfileResource(): UseApiResourceResult<ProfileResponse> {
  const contextValue = useContext(ProfileResourceContext);
  const fallbackValue = useProfileResourceInternal({
    skip: contextValue !== null,
  });

  return contextValue ?? fallbackValue;
}

export function useResearchResource(): UseApiResourceResult<ResearchEntry[]> {
  return useApiResource(publicApi.getResearch, {
    initialData: () => transformResearchEntries(undefined),
    skip: USE_MOCK_PUBLIC_API,
  });
}

export function useProjectsResource(): UseApiResourceResult<Project[]> {
  return useApiResource(publicApi.getProjects, {
    initialData: () => transformProjects(undefined),
    skip: USE_MOCK_PUBLIC_API,
  });
}

export function useContactAvailability(): UseApiResourceResult<ContactAvailabilityResponse> {
  return useApiResource(publicApi.getContactAvailability, {
    initialData: () => createMockAvailability(),
    skip: USE_MOCK_PUBLIC_API,
  });
}

export function useContactConfig(): UseApiResourceResult<ContactConfigResponse> {
  return useApiResource(publicApi.getContactConfig, {
    initialData: () => createMockContactConfig(),
    skip: USE_MOCK_PUBLIC_API,
  });
}
