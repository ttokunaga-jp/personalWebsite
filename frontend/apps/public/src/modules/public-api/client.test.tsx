import { act, cleanup, renderHook, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import i18n from "../i18n";

import { publicApi, useProfileResourceInternal } from "./client";
import type { ProfileResponse } from "./types";

describe("useProfileResourceInternal", () => {
  const profileEn: ProfileResponse = {
    id: "profile-id",
    displayName: "Takumi (EN)",
    theme: {
      mode: "light",
    },
    affiliations: [],
    communities: [],
    workHistory: [],
    techSections: [],
    socialLinks: [],
    footerLinks: [],
    updatedAt: "2024-01-01T00:00:00.000Z",
  };

  const profileJa: ProfileResponse = {
    ...profileEn,
    displayName: "徳永 拓未",
  };

  beforeEach(async () => {
    vi.spyOn(publicApi, "getProfile")
      .mockResolvedValueOnce(profileEn)
      .mockResolvedValueOnce(profileJa);
    await i18n.changeLanguage("en");
  });

  afterEach(() => {
    cleanup();
    vi.restoreAllMocks();
  });

  it("refetches profile data when the active language changes", async () => {
    const { result } = renderHook(() => useProfileResourceInternal());

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
      expect(result.current.data?.displayName).toBe("Takumi (EN)");
    });

    await act(async () => {
      await i18n.changeLanguage("ja");
    });

    await waitFor(() => {
      expect(publicApi.getProfile).toHaveBeenCalledTimes(2);
      expect(result.current.data?.displayName).toBe("徳永 拓未");
    });
  });
});
