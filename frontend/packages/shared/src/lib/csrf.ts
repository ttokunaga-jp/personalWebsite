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

function resolveMetaEnv():
  | Record<string, string | undefined>
  | undefined {
  if (typeof import.meta === "undefined") {
    return undefined;
  }
  const candidate = import.meta as unknown;
  if (typeof candidate !== "object" || candidate === null) {
    return undefined;
  }
  const env = (candidate as { env?: unknown }).env;
  if (!env || typeof env !== "object" || env === null) {
    return undefined;
  }
  return env as Record<string, string | undefined>;
}

const metaEnv = resolveMetaEnv();

const baseURL =
  metaEnv?.VITE_API_BASE_URL ??
  (typeof process !== "undefined" ? process.env?.VITE_API_BASE_URL : undefined) ??
  "/api";

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
