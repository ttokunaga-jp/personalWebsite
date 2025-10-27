import axios, { AxiosHeaders } from "axios";
import type {
  AxiosError,
  AxiosPromise,
  AxiosRequestConfig,
  RawAxiosRequestHeaders
} from "axios";

import { getCsrfToken, invalidateCsrfToken } from "./csrf";

const SAFE_METHODS = new Set(["get", "head", "options", "trace"]);

export const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL ?? "/api",
  timeout: 10_000,
  withCredentials: true
});

type CsrfAwareRequestConfig = AxiosRequestConfig & {
  _csrfRetry?: boolean;
};

function setHeader(config: AxiosRequestConfig, key: string, value: string): void {
  if (!config.headers) {
    config.headers = {};
  }

  if (config.headers instanceof AxiosHeaders) {
    config.headers.set(key, value);
    return;
  }

  (config.headers as RawAxiosRequestHeaders)[key] = value;
}

apiClient.interceptors.request.use(async (config) => {
  config.withCredentials = true;

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
      // Placeholder for auth refresh logic.
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
  }
);

export type ApiClientPromise<T> = AxiosPromise<T>;
