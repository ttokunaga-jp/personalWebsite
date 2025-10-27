import "@testing-library/jest-dom/vitest";
import { act } from "react";
import { afterEach, beforeAll, vi } from "vitest";
import i18n from "./src/modules/i18n";

(globalThis as { IS_REACT_ACT_ENVIRONMENT?: boolean }).IS_REACT_ACT_ENVIRONMENT = true;

const listenerMap = new WeakMap<(...args: unknown[]) => void, (...args: unknown[]) => void>();
const originalOn = i18n.on.bind(i18n);
const originalOff = i18n.off.bind(i18n);

i18n.on = ((event: string, listener: (...args: unknown[]) => void, ...rest: unknown[]) => {
  if ((event === "languageChanged" || event === "initialized") && typeof listener === "function") {
    const wrapped = (...args: unknown[]) => {
      act(() => {
        listener(...args);
      });
    };
    listenerMap.set(listener, wrapped);
    return originalOn(event, wrapped, ...rest);
  }

  return originalOn(event, listener, ...rest);
}) as typeof i18n.on;

i18n.off = ((event: string, listener: (...args: unknown[]) => void, ...rest: unknown[]) => {
  const wrapped = listenerMap.get(listener);
  if (wrapped) {
    listenerMap.delete(listener);
    return originalOff(event, wrapped, ...rest);
  }

  return originalOff(event, listener, ...rest);
}) as typeof i18n.off;

if (typeof window !== "undefined" && !window.matchMedia) {
  Object.defineProperty(window, "matchMedia", {
    writable: true,
    value: (query: string) => {
      const listeners = new Set<(event: MediaQueryListEvent) => void>();
      const mediaQueryList: MediaQueryList = {
        media: query,
        matches: false,
        onchange: null,
        addEventListener: (_: string, listener: (event: MediaQueryListEvent) => void) => {
          listeners.add(listener);
        },
        removeEventListener: (_: string, listener: (event: MediaQueryListEvent) => void) => {
          listeners.delete(listener);
        },
        dispatchEvent: (event: Event) => {
          listeners.forEach((listener) => listener(event as MediaQueryListEvent));
          return true;
        },
        addListener: (listener: (event: MediaQueryListEvent) => void) => listeners.add(listener),
        removeListener: (listener: (event: MediaQueryListEvent) => void) => listeners.delete(listener)
      };
      return mediaQueryList;
    }
  });
}

beforeAll(async () => {
  await i18n.changeLanguage("en");
});

afterEach(async () => {
  vi.restoreAllMocks();
  window.localStorage.clear();
  document.documentElement.classList.remove("dark");
  await i18n.changeLanguage("en");
});
