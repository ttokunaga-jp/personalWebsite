import axios, { AxiosHeaders } from "axios";
import type {
  AxiosError,
  AxiosPromise,
  AxiosRequestConfig,
  RawAxiosRequestHeaders,
} from "axios";

import { getCsrfToken, invalidateCsrfToken } from "./csrf";

const SAFE_METHODS = new Set(["get", "head", "options", "trace"]);

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

const apiBaseURL =
  metaEnv?.VITE_API_BASE_URL ??
  (typeof process !== "undefined" ? process.env?.VITE_API_BASE_URL : undefined) ??
  "/api";

export const apiClient = axios.create({
  baseURL: apiBaseURL,
  timeout: 10_000,
  withCredentials: true,
});

type TokenProvider = () => string | null;
type UnauthorizedHandler = () => void;

let tokenProvider: TokenProvider | null = null;
let unauthorizedHandler: UnauthorizedHandler | null = null;

export function registerAuthTokenProvider(provider: TokenProvider | null): void {
  tokenProvider = provider;
}

export function registerAuthTokenInvalidator(
  handler: UnauthorizedHandler | null,
): void {
  unauthorizedHandler = handler;
}

type CsrfAwareRequestConfig = AxiosRequestConfig & {
  _csrfRetry?: boolean;
};

function setHeader(
  config: AxiosRequestConfig,
  key: string,
  value: string,
): void {
  if (!config.headers) {
    config.headers = {};
  }

  if (config.headers instanceof AxiosHeaders) {
    config.headers.set(key, value);
    return;
  }

  (config.headers as RawAxiosRequestHeaders)[key] = value;
}

function removeHeader(config: AxiosRequestConfig, key: string): void {
  if (!config.headers) {
    return;
  }

  if (config.headers instanceof AxiosHeaders) {
    config.headers.delete(key);
    return;
  }

  delete (config.headers as RawAxiosRequestHeaders)[key];
}

apiClient.interceptors.request.use(async (config) => {
  config.withCredentials = true;

  const sessionToken = tokenProvider?.() ?? null;
  if (sessionToken) {
    setHeader(config, "Authorization", `Bearer ${sessionToken}`);
  } else {
    removeHeader(config, "Authorization");
  }

  setHeader(config, "X-Requested-With", "XMLHttpRequest");

  const method = (config.method ?? "get").toLowerCase();
  if (!SAFE_METHODS.has(method)) {
    const token = await getCsrfToken();
    if (token) {
      setHeader(config, "X-CSRF-Token", token);
    }
  }

  return config;
});

apiClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    if (error.response?.status === 401) {
      if (typeof unauthorizedHandler === "function") {
        unauthorizedHandler();
      }
    }

    const csrfConfig = error.config as CsrfAwareRequestConfig | undefined;
    const status = error.response?.status;
    const method = (csrfConfig?.method ?? "get").toLowerCase();

    if (
      status === 403 &&
      csrfConfig &&
      !SAFE_METHODS.has(method) &&
      !csrfConfig._csrfRetry
    ) {
      csrfConfig._csrfRetry = true;
      invalidateCsrfToken();
      return apiClient.request(csrfConfig);
    }

    return Promise.reject(error);
  },
);

export type ApiClientPromise<T> = AxiosPromise<T>;
