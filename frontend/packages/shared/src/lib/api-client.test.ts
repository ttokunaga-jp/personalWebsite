import { describe, expect, it } from "vitest";
import { apiClient } from "./api-client";

describe("apiClient", () => {
  it("uses default base url", () => {
    expect(apiClient.defaults.baseURL).toBe("/api");
  });
});
