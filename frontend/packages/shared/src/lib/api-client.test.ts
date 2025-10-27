import type { AxiosError, AxiosRequestConfig, AxiosResponse } from "axios";
import { afterEach, describe, expect, it, vi } from "vitest";

vi.mock("./csrf", () => ({
  getCsrfToken: vi.fn().mockResolvedValue("mock-csrf-token"),
  invalidateCsrfToken: vi.fn()
}));

import { apiClient } from "./api-client";
import { getCsrfToken, invalidateCsrfToken } from "./csrf";

const getTokenMock = vi.mocked(getCsrfToken);
const invalidateTokenMock = vi.mocked(invalidateCsrfToken);

afterEach(() => {
  vi.clearAllMocks();
});

describe("apiClient", () => {
  it("uses default base url", () => {
    expect(apiClient.defaults.baseURL).toBe("/api");
  });

  it("enables credentialed requests", () => {
    expect(apiClient.defaults.withCredentials).toBe(true);
  });

  it("attaches CSRF header for unsafe methods", async () => {
    const { handlers } = apiClient.interceptors.request as unknown as {
      handlers: Array<{
        fulfilled?: (config: AxiosRequestConfig) => AxiosRequestConfig | Promise<AxiosRequestConfig>;
      }>;
    };
    const handler = handlers[0].fulfilled!;
    const config = await handler({ method: "post", headers: {} });

    expect(getTokenMock).toHaveBeenCalledTimes(1);
    expect((config.headers as Record<string, string>)["X-CSRF-Token"]).toBe("mock-csrf-token");
  });

  it("skips CSRF header for safe methods", async () => {
    const { handlers } = apiClient.interceptors.request as unknown as {
      handlers: Array<{
        fulfilled?: (config: AxiosRequestConfig) => AxiosRequestConfig | Promise<AxiosRequestConfig>;
      }>;
    };
    const handler = handlers[0].fulfilled!;
    await handler({ method: "get", headers: {} });

    expect(getTokenMock).not.toHaveBeenCalled();
  });

  it("retries once on CSRF failure", async () => {
    const { handlers } = apiClient.interceptors.response as unknown as {
      handlers: Array<{
        fulfilled?: (value: AxiosResponse) => AxiosResponse | Promise<AxiosResponse>;
        rejected?: (error: AxiosError) => unknown;
      }>;
    };
    const handler = handlers[0].rejected!;
    const requestSpy = vi.spyOn(apiClient, "request").mockResolvedValue({} as never);

    await handler({
      config: { method: "post", headers: {} },
      response: { status: 403 }
    } as AxiosError);

    expect(invalidateTokenMock).toHaveBeenCalledTimes(1);
    expect(requestSpy).toHaveBeenCalledTimes(1);
  });
});
