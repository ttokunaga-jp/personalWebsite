import "@testing-library/jest-dom/vitest";
import { afterEach, beforeAll, vi } from "vitest";
import i18n from "./src/modules/i18n";

beforeAll(async () => {
  await i18n.changeLanguage("en");
});

afterEach(() => {
  vi.restoreAllMocks();
});
