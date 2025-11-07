import i18next from "i18next";
import { beforeEach, describe, expect, it } from "vitest";

import "../i18n";

import { getCanonicalHomeConfig } from "../profile-content";

import { transformProfile, type RawProfileDocument } from "./transform";

describe("transformProfile", () => {
  beforeEach(() => {
    i18next.changeLanguage("en");
  });

  it("maps home configuration from the API response", () => {
    const raw: RawProfileDocument = {
      id: 42,
      displayName: "Takumi",
      home: {
        heroSubtitle: { en: "Engineering the physical world" },
        quickLinks: [
          {
            id: 2,
            section: "projects",
            label: { en: "Projects" },
            description: { en: "Recent work with Go and React" },
            cta: { en: "Explore" },
            targetUrl: "/projects",
            sortOrder: 2,
          },
          {
            id: 1,
            section: "profile",
            label: { en: "Profile" },
            description: { en: "Academic background" },
            cta: { en: "View" },
            targetUrl: "/profile",
            sortOrder: 1,
          },
        ],
        chipSources: [
          {
            id: 7,
            source: "tech",
            label: { en: "Key Tech" },
            limit: 4,
            sortOrder: 1,
          },
        ],
        updatedAt: "2024-01-01T00:00:00.000Z",
      },
    };

    const profile = transformProfile(raw);

    expect(profile.home?.heroSubtitle).toBe("Engineering the physical world");
    expect(profile.home?.quickLinks[0]).toMatchObject({
      id: "1",
      label: "Profile",
      targetUrl: "/profile",
    });
    expect(profile.home?.chipSources[0]).toMatchObject({
      id: "7",
      label: "Key Tech",
      limit: 4,
    });
    expect(profile.home?.updatedAt).toBe("2024-01-01T00:00:00.000Z");
  });

  it("falls back to canonical home configuration when missing", () => {
    const raw: RawProfileDocument = {
      id: 42,
      displayName: "Takumi",
      home: null,
    };

    const profile = transformProfile(raw);
    const canonical = getCanonicalHomeConfig("en");

    expect(profile.home).toEqual(canonical);
  });
});
