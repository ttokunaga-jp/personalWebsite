import { beforeEach, describe, expect, it, vi } from "vitest";

import {
  FALLBACK_LANGUAGE,
  LANGUAGE_COOKIE_NAME,
  LANGUAGE_STORAGE_KEY,
} from "./config";
import {
  clearLanguagePreference,
  matchSupportedLanguage,
  persistLanguagePreference,
  resolveInitialLanguage,
} from "./preference";

function readCookieValue(name: string): string | null {
  const cookies = document.cookie
    .split(";")
    .map((entry) => entry.trim())
    .filter(Boolean);

  for (const cookie of cookies) {
    const [cookieName, value] = cookie.split("=");
    if (cookieName === name) {
      return value ?? null;
    }
  }

  return null;
}

describe("language preference helpers", () => {
  beforeEach(() => {
    clearLanguagePreference();
    window.localStorage.clear();
    vi.restoreAllMocks();
  });

  it("resolves the stored language from localStorage when available", () => {
    window.localStorage.setItem(LANGUAGE_STORAGE_KEY, "ja");

    const resolution = resolveInitialLanguage();

    expect(resolution).toEqual({
      language: "ja",
      source: "storage",
    });
  });

  it("falls back to the cookie value when localStorage is not set", () => {
    document.cookie = `${LANGUAGE_COOKIE_NAME}=ja`;

    const resolution = resolveInitialLanguage();

    expect(resolution).toEqual({
      language: "ja",
      source: "cookie",
    });
  });

  it("uses navigator languages when storage is empty", () => {
    vi.spyOn(window.navigator, "languages", "get").mockReturnValue([
      "ja-JP",
      "en-US",
    ]);
    vi.spyOn(window.navigator, "language", "get").mockReturnValue("en-US");

    const resolution = resolveInitialLanguage();

    expect(resolution).toEqual({
      language: "ja",
      source: "navigator",
    });
  });

  it("returns the fallback language when no hints are available", () => {
    vi.spyOn(window.navigator, "languages", "get").mockReturnValue([]);
    vi.spyOn(window.navigator, "language", "get").mockReturnValue("");

    const resolution = resolveInitialLanguage();

    expect(resolution).toEqual({
      language: FALLBACK_LANGUAGE,
      source: "fallback",
    });
  });

  it("persists the language in both localStorage and cookies", () => {
    persistLanguagePreference("ja");

    expect(window.localStorage.getItem(LANGUAGE_STORAGE_KEY)).toBe("ja");
    expect(readCookieValue(LANGUAGE_COOKIE_NAME)).toBe("ja");
  });

  it("matches supported languages from extended codes", () => {
    expect(matchSupportedLanguage("en-US")).toBe("en");
    expect(matchSupportedLanguage("ja-JP")).toBe("ja");
    expect(matchSupportedLanguage("de-DE")).toBeNull();
    expect(matchSupportedLanguage(undefined)).toBeNull();
  });
});
