import type { AxiosRequestConfig, AxiosResponse } from "axios";
import { AxiosError } from "axios";
import { afterEach, describe, expect, it } from "vitest";

import {
  apiClient,
  registerAuthTokenInvalidator,
  registerAuthTokenProvider,
} from "./api-client";

describe("apiClient authentication interceptors", () => {
  afterEach(() => {
    registerAuthTokenProvider(null);
    registerAuthTokenInvalidator(null);
  });

  it("attaches a bearer token from the registered provider", async () => {
    registerAuthTokenProvider(() => "test-bearer-token");

    const requestInterceptors =
      apiClient.interceptors.request as unknown as {
        handlers: Array<{
          fulfilled?: (config: AxiosRequestConfig) =>
            | AxiosRequestConfig
            | Promise<AxiosRequestConfig>;
        }>;
      };
    const handler = requestInterceptors.handlers[0]?.fulfilled;
    expect(handler).toBeDefined();

    const config = await handler!({ headers: {} } as AxiosRequestConfig);
    expect(config.headers).toBeDefined();
    expect((config.headers as Record<string, string>).Authorization).toBe(
      "Bearer test-bearer-token",
    );
  });

  it("invokes the unauthorized callback on 401 responses", async () => {
    let wasInvalidated = false;
    registerAuthTokenInvalidator(() => {
      wasInvalidated = true;
    });

    const responseInterceptors =
      apiClient.interceptors.response as unknown as {
        handlers: Array<{
          rejected?: (error: unknown) => unknown | Promise<unknown>;
        }>;
      };
    const handler = responseInterceptors.handlers[0]?.rejected;
    expect(handler).toBeDefined();

    const response = {
      data: {},
      status: 401,
      statusText: "Unauthorized",
      headers: {},
      config: { headers: {} } as AxiosRequestConfig,
    } as AxiosResponse;
    const error = new AxiosError(
      "unauthorized",
      undefined,
      response.config,
      null,
      response,
    );

    await expect(handler!(error)).rejects.toBeDefined();
    expect(wasInvalidated).toBe(true);
  });
});
