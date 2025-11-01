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

    const handler = apiClient.interceptors.request.handlers[0]?.fulfilled;
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

    const handler = apiClient.interceptors.response.handlers[0]?.rejected;
    expect(handler).toBeDefined();

    const response: AxiosResponse = {
      data: {},
      status: 401,
      statusText: "Unauthorized",
      headers: {},
      config: {},
    };
    const error = new AxiosError("unauthorized", undefined, {}, null, response);

    await expect(handler!(error)).rejects.toBeDefined();
    expect(wasInvalidated).toBe(true);
  });
});
