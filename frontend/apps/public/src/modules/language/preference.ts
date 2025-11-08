import {
  FALLBACK_LANGUAGE,
  LANGUAGE_COOKIE_MAX_AGE_SECONDS,
  LANGUAGE_COOKIE_NAME,
  LANGUAGE_STORAGE_KEY,
  SUPPORTED_LANGUAGES,
  type SupportedLanguage,
} from "./config";

const SUPPORTED_LANGUAGE_SET = new Set<SupportedLanguage>(SUPPORTED_LANGUAGES);

function normalizeLanguageCode(value: unknown): SupportedLanguage | null {
  if (typeof value !== "string") {
    return null;
  }
  const lower = value.trim().toLowerCase();
  if (!lower) {
    return null;
  }
  if (SUPPORTED_LANGUAGE_SET.has(lower as SupportedLanguage)) {
    return lower as SupportedLanguage;
  }
  const [base] = lower.split("-");
  if (base && SUPPORTED_LANGUAGE_SET.has(base as SupportedLanguage)) {
    return base as SupportedLanguage;
  }
  return null;
}

function readFromLocalStorage(): SupportedLanguage | null {
  try {
    if (typeof window === "undefined" || !("localStorage" in window)) {
      return null;
    }
    const stored = window.localStorage.getItem(LANGUAGE_STORAGE_KEY);
    return normalizeLanguageCode(stored);
  } catch {
    return null;
  }
}

function readFromCookie(): SupportedLanguage | null {
  if (typeof document === "undefined") {
    return null;
  }

  try {
    const cookies = document.cookie
      .split(";")
      .map((entry) => entry.trim())
      .filter(Boolean);

    for (const cookie of cookies) {
      const [name, value] = cookie.split("=");
      if (name === LANGUAGE_COOKIE_NAME && value) {
        return normalizeLanguageCode(decodeURIComponent(value));
      }
    }
  } catch {
    return null;
  }

  return null;
}

function readFromNavigator(): SupportedLanguage | null {
  if (typeof navigator === "undefined") {
    return null;
  }

  const candidates: string[] = [];

  if (Array.isArray(navigator.languages)) {
    candidates.push(...navigator.languages);
  }

  if (typeof navigator.language === "string") {
    candidates.push(navigator.language);
  }

  for (const candidate of candidates) {
    const normalized = normalizeLanguageCode(candidate);
    if (normalized) {
      return normalized;
    }
  }

  return null;
}

function shouldSetSecureCookie(): boolean {
  if (typeof window === "undefined") {
    return false;
  }
  try {
    return window.location?.protocol === "https:";
  } catch {
    return false;
  }
}

export type LanguageResolution = {
  language: SupportedLanguage;
  source: "storage" | "cookie" | "navigator" | "fallback";
};

export function resolveInitialLanguage(): LanguageResolution {
  const fromStorage = readFromLocalStorage();
  if (fromStorage) {
    return { language: fromStorage, source: "storage" as const };
  }

  const fromCookie = readFromCookie();
  if (fromCookie) {
    return { language: fromCookie, source: "cookie" as const };
  }

  const fromNavigator = readFromNavigator();
  if (fromNavigator) {
    return { language: fromNavigator, source: "navigator" as const };
  }

  return { language: FALLBACK_LANGUAGE, source: "fallback" as const };
}

export function matchSupportedLanguage(
  value: unknown,
): SupportedLanguage | null {
  return normalizeLanguageCode(value);
}

export function persistLanguagePreference(language: SupportedLanguage): void {
  try {
    if (typeof window !== "undefined" && "localStorage" in window) {
      window.localStorage.setItem(LANGUAGE_STORAGE_KEY, language);
    }
  } catch {
    // Errors should not block language switching UX.
  }

  if (typeof document === "undefined") {
    return;
  }

  try {
    const cookieParts = [
      `${LANGUAGE_COOKIE_NAME}=${encodeURIComponent(language)}`,
      "Path=/",
      `Max-Age=${LANGUAGE_COOKIE_MAX_AGE_SECONDS}`,
      "SameSite=Lax",
    ];
    if (shouldSetSecureCookie()) {
      cookieParts.push("Secure");
    }
    document.cookie = cookieParts.join("; ");
  } catch {
    // Ignore cookie write errors to maintain UX.
  }
}

export function clearLanguagePreference(): void {
  try {
    if (typeof window !== "undefined" && "localStorage" in window) {
      window.localStorage.removeItem(LANGUAGE_STORAGE_KEY);
    }
  } catch {
    // Ignore storage cleanup errors.
  }

  if (typeof document === "undefined") {
    return;
  }

  try {
    const cookieParts = [
      `${LANGUAGE_COOKIE_NAME}=`,
      "Path=/",
      "Max-Age=0",
      "SameSite=Lax",
    ];
    if (shouldSetSecureCookie()) {
      cookieParts.push("Secure");
    }
    document.cookie = cookieParts.join("; ");
  } catch {
    // Ignore cookie cleanup errors.
  }
}
