export const SUPPORTED_LANGUAGES = ["en", "ja"] as const;

export type SupportedLanguage = (typeof SUPPORTED_LANGUAGES)[number];

export const FALLBACK_LANGUAGE: SupportedLanguage = "en";

export const LANGUAGE_STORAGE_KEY = "public.language";
export const LANGUAGE_COOKIE_NAME = "public_language";
export const LANGUAGE_COOKIE_MAX_AGE_SECONDS = 60 * 60 * 24 * 365; // 1 year
