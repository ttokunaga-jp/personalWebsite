import axios from "axios";

type CsrfResponse = {
  data: {
    token: string;
    expires_at: string;
  };
};

const SAFE_MARGIN_MS = 10_000;

let cachedToken: string | null = null;
let expiresAt = 0;

const baseURL = import.meta.env.VITE_API_BASE_URL ?? "/api";

async function fetchToken(): Promise<void> {
  const client = axios.create({
    baseURL,
    timeout: 5_000,
    withCredentials: true,
  });

  const response = await client.get<CsrfResponse>("/security/csrf");
  const payload = response.data.data;

  cachedToken = payload.token;
  const parsed = Date.parse(payload.expires_at ?? "");
  expiresAt = Number.isNaN(parsed)
    ? Date.now() + SAFE_MARGIN_MS
    : parsed - SAFE_MARGIN_MS;
}

export async function getCsrfToken(): Promise<string | null> {
  if (!cachedToken || Date.now() >= expiresAt) {
    await fetchToken();
  }
  return cachedToken;
}

export function invalidateCsrfToken(): void {
  cachedToken = null;
  expiresAt = 0;
}
