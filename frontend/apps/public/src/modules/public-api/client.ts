import { apiClient } from "@shared/lib/api-client";

import {
  type UseApiResourceResult,
  useApiResource,
} from "../../lib/useApiResource";

import {
  transformProfile,
  transformProjects,
  transformResearchEntries,
  type RawProfileResponse,
  type RawProject,
  type RawResearchEntry,
} from "./transform";
import type {
  ContactAvailabilityResponse,
  ContactConfigResponse,
  CreateBookingPayload,
  CreateBookingResponse,
  ProfileResponse,
  Project,
  ResearchEntry,
} from "./types";

const BASE_PATH = "/v1/public";
const USE_MOCK_PUBLIC_API =
  import.meta.env.MODE !== "test" &&
  (import.meta.env.VITE_USE_MOCK_PUBLIC_API ?? "true") !== "false";

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

  const slots = Array.from({ length: 6 }).map((_, index) => {
    const start = new Date(now);
    start.setDate(start.getDate() + Math.floor(index / 2));
    start.setHours(9 + (index % 3) * 2);

    const end = new Date(start);
    end.setHours(start.getHours() + 1);

    return {
      id: `mock-slot-${index + 1}`,
      start: start.toISOString(),
      end: end.toISOString(),
      isBookable: index % 4 !== 0,
    };
  });

  return {
    timezone: "Asia/Tokyo",
    generatedAt: new Date().toISOString(),
    slots,
  };
}

function createMockContactConfig(): ContactConfigResponse {
  return {
    topics: [
      "Research collaboration",
      "Project consultation",
      "Speaking engagement",
      "Mentoring session",
    ],
    recaptchaSiteKey: import.meta.env.VITE_RECAPTCHA_SITE_KEY ?? "",
    minimumLeadHours: 48,
    consentText:
      "Provided details are used only to coordinate the requested meeting. Expect a reply within two business days.",
  };
}

function createMockBookingResponse(): CreateBookingResponse {
  return {
    bookingId: `mock-${Date.now()}`,
    status: "pending",
  };
}

export const publicApi = {
  async getProfile(signal: AbortSignal): Promise<ProfileResponse> {
    if (USE_MOCK_PUBLIC_API) {
      ensureNotAborted(signal);
      return transformProfile(undefined);
    }

    const response = await apiClient.get<ApiSuccessResponse<RawProfileResponse>>(
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
      ApiSuccessResponse<RawResearchEntry[]>
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

    const response = await apiClient.get<ApiSuccessResponse<RawProject[]>>(
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
      ApiSuccessResponse<ContactConfigResponse>
    >(`${BASE_PATH}/contact/config`, {
      ...withAbortSignal(signal),
    });
    return unwrapData(response.data);
  },
  async createBooking(
    payload: CreateBookingPayload,
  ): Promise<CreateBookingResponse> {
    if (USE_MOCK_PUBLIC_API) {
      return createMockBookingResponse();
    }

    const response = await apiClient.post<
      ApiSuccessResponse<CreateBookingResponse>
    >(`${BASE_PATH}/contact/bookings`, payload);
    return unwrapData(response.data);
  },
};

export function useProfileResource(): UseApiResourceResult<ProfileResponse> {
  return useApiResource(publicApi.getProfile, {
    initialData: () => transformProfile(undefined),
    skip: USE_MOCK_PUBLIC_API,
  });
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
