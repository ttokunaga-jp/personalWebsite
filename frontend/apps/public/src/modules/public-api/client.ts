import { apiClient } from "@shared/lib/api-client";

import { type UseApiResourceResult, useApiResource } from "../../lib/useApiResource";

import type {
  ContactAvailabilityResponse,
  ContactConfigResponse,
  CreateBookingPayload,
  CreateBookingResponse,
  ProfileResponse,
  Project,
  ResearchEntry
} from "./types";

const BASE_PATH = "/v1/public";

function withAbortSignal(signal: AbortSignal) {
  return { signal };
}

type ApiSuccessResponse<T> = {
  data: T;
};

function unwrapData<T>(response: ApiSuccessResponse<T>): T {
  return response.data;
}

export const publicApi = {
  async getProfile(signal: AbortSignal): Promise<ProfileResponse> {
    const response = await apiClient.get<ApiSuccessResponse<ProfileResponse>>(
      `${BASE_PATH}/profile`,
      {
        ...withAbortSignal(signal)
      }
    );
    return unwrapData(response.data);
  },
  async getResearch(signal: AbortSignal): Promise<ResearchEntry[]> {
    const response = await apiClient.get<ApiSuccessResponse<ResearchEntry[]>>(
      `${BASE_PATH}/research`,
      {
        ...withAbortSignal(signal)
      }
    );
    return unwrapData(response.data);
  },
  async getProjects(signal: AbortSignal): Promise<Project[]> {
    const response = await apiClient.get<ApiSuccessResponse<Project[]>>(
      `${BASE_PATH}/projects`,
      {
        ...withAbortSignal(signal)
      }
    );
    return unwrapData(response.data);
  },
  async getContactAvailability(signal: AbortSignal): Promise<ContactAvailabilityResponse> {
    const response = await apiClient.get<ApiSuccessResponse<ContactAvailabilityResponse>>(
      `${BASE_PATH}/contact/availability`,
      {
        ...withAbortSignal(signal)
      }
    );
    return unwrapData(response.data);
  },
  async getContactConfig(signal: AbortSignal): Promise<ContactConfigResponse> {
    const response = await apiClient.get<ApiSuccessResponse<ContactConfigResponse>>(
      `${BASE_PATH}/contact/config`,
      {
        ...withAbortSignal(signal)
      }
    );
    return unwrapData(response.data);
  },
  async createBooking(payload: CreateBookingPayload): Promise<CreateBookingResponse> {
    const response = await apiClient.post<ApiSuccessResponse<CreateBookingResponse>>(
      `${BASE_PATH}/contact/bookings`,
      payload
    );
    return unwrapData(response.data);
  }
};

export function useProfileResource(): UseApiResourceResult<ProfileResponse> {
  return useApiResource(publicApi.getProfile);
}

export function useResearchResource(): UseApiResourceResult<ResearchEntry[]> {
  return useApiResource(publicApi.getResearch);
}

export function useProjectsResource(): UseApiResourceResult<Project[]> {
  return useApiResource(publicApi.getProjects);
}

export function useContactAvailability(): UseApiResourceResult<ContactAvailabilityResponse> {
  return useApiResource(publicApi.getContactAvailability);
}

export function useContactConfig(): UseApiResourceResult<ContactConfigResponse> {
  return useApiResource(publicApi.getContactConfig);
}
