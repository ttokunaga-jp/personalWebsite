import { apiClient } from "@shared/lib/api-client";

import type {
  BlacklistEntry,
  BlacklistInput,
  HomeConfigDocument,
  MeetingUrlPayload,
  MeetingUrlTemplate,
  ProfileDocument,
  Reservation,
  ReservationUpdatePayload,
  SocialLink,
  SocialLinkInput,
  TechCatalogEntry,
  TechCatalogInput,
} from "./types";

const ADMIN_MODE_PARAMS = {
  mode: "admin",
};

type ApiListResponse<T> = {
  data: T[];
};

type ApiItemResponse<T> = {
  data: T;
};

export async function fetchReservations(): Promise<Reservation[]> {
  const response = await apiClient.get<ApiListResponse<Reservation>>(
    "/admin/reservations",
    {
      params: ADMIN_MODE_PARAMS,
    },
  );
  return response.data?.data ?? [];
}

export async function updateReservationStatus(
  id: number,
  payload: ReservationUpdatePayload,
): Promise<Reservation> {
  const response = await apiClient.put<ApiItemResponse<Reservation>>(
    `/admin/reservations/${id}`,
    payload,
    {
      params: ADMIN_MODE_PARAMS,
    },
  );
  return response.data.data;
}

export async function retryReservationNotification(
  id: number,
): Promise<Reservation> {
  const response = await apiClient.post<ApiItemResponse<Reservation>>(
    `/admin/reservations/${id}/retry`,
    {},
    {
      params: ADMIN_MODE_PARAMS,
    },
  );
  return response.data.data;
}

export async function fetchBlacklist(): Promise<BlacklistEntry[]> {
  const response = await apiClient.get<ApiListResponse<BlacklistEntry>>(
    "/admin/blacklist",
    {
      params: ADMIN_MODE_PARAMS,
    },
  );
  return response.data?.data ?? [];
}

export async function createBlacklistEntry(
  payload: BlacklistInput,
): Promise<BlacklistEntry> {
  const response = await apiClient.post<ApiItemResponse<BlacklistEntry>>(
    "/admin/blacklist",
    payload,
    {
      params: ADMIN_MODE_PARAMS,
    },
  );
  return response.data.data;
}

export async function updateBlacklistEntry(
  id: number,
  payload: BlacklistInput,
): Promise<BlacklistEntry> {
  const response = await apiClient.put<ApiItemResponse<BlacklistEntry>>(
    `/admin/blacklist/${id}`,
    payload,
    {
      params: ADMIN_MODE_PARAMS,
    },
  );
  return response.data.data;
}

export async function deleteBlacklistEntry(id: number): Promise<void> {
  await apiClient.delete(`/admin/blacklist/${id}`, {
    params: ADMIN_MODE_PARAMS,
  });
}

export async function fetchTechCatalog(): Promise<TechCatalogEntry[]> {
  const response = await apiClient.get<ApiListResponse<TechCatalogEntry>>(
    "/admin/tech-catalog",
    {
      params: ADMIN_MODE_PARAMS,
    },
  );
  return response.data?.data ?? [];
}

export async function createTechCatalogEntry(
  payload: TechCatalogInput,
): Promise<TechCatalogEntry> {
  const response = await apiClient.post<ApiItemResponse<TechCatalogEntry>>(
    "/admin/tech-catalog",
    payload,
    {
      params: ADMIN_MODE_PARAMS,
    },
  );
  return response.data.data;
}

export async function updateTechCatalogEntry(
  id: number,
  payload: Partial<TechCatalogInput>,
): Promise<TechCatalogEntry> {
  const response = await apiClient.put<ApiItemResponse<TechCatalogEntry>>(
    `/admin/tech-catalog/${id}`,
    payload,
    {
      params: ADMIN_MODE_PARAMS,
    },
  );
  return response.data.data;
}

export async function fetchSocialLinks(): Promise<SocialLink[]> {
  const response = await apiClient.get<ApiListResponse<SocialLink>>(
    "/admin/social-links",
    {
      params: ADMIN_MODE_PARAMS,
    },
  );
  return response.data?.data ?? [];
}

export async function replaceSocialLinks(
  links: SocialLinkInput[],
): Promise<SocialLink[]> {
  const response = await apiClient.put<ApiListResponse<SocialLink>>(
    "/admin/social-links",
    { links },
    { params: ADMIN_MODE_PARAMS },
  );
  return response.data?.data ?? [];
}

export async function fetchMeetingUrl(): Promise<MeetingUrlTemplate> {
  const response = await apiClient.get<ApiItemResponse<MeetingUrlTemplate>>(
    "/admin/meeting-url",
    {
      params: ADMIN_MODE_PARAMS,
    },
  );
  return response.data.data;
}

export async function updateMeetingUrl(
  payload: MeetingUrlPayload,
): Promise<MeetingUrlTemplate> {
  const response = await apiClient.put<ApiItemResponse<MeetingUrlTemplate>>(
    "/admin/meeting-url",
    payload,
    {
      params: ADMIN_MODE_PARAMS,
    },
  );
  return response.data.data;
}

export async function fetchHomeConfig(): Promise<HomeConfigDocument> {
  const response = await apiClient.get<ApiItemResponse<HomeConfigDocument>>(
    "/admin/home",
    {
      params: ADMIN_MODE_PARAMS,
    },
  );
  return response.data.data;
}

export async function updateHomeConfig(
  payload: HomeConfigDocument,
): Promise<HomeConfigDocument> {
  const response = await apiClient.put<ApiItemResponse<HomeConfigDocument>>(
    "/admin/home",
    payload,
    { params: ADMIN_MODE_PARAMS },
  );
  return response.data.data;
}

export async function fetchProfileDocument(): Promise<ProfileDocument> {
  const response = await apiClient.get<ApiItemResponse<ProfileDocument>>(
    "/admin/profile",
    {
      params: ADMIN_MODE_PARAMS,
    },
  );
  return response.data.data;
}

export async function updateProfileDocument(
  payload: ProfileDocument,
): Promise<ProfileDocument> {
  const response = await apiClient.put<ApiItemResponse<ProfileDocument>>(
    "/admin/profile",
    payload,
    {
      params: ADMIN_MODE_PARAMS,
    },
  );
  return response.data.data;
}
