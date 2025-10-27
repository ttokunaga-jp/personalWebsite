import type { ReactNode } from "react";
import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState
} from "react";

const STORAGE_KEY = "personal-website.theme";
const PREFERRED_DARK_QUERY = "(prefers-color-scheme: dark)";

type Theme = "light" | "dark";

type ThemeContextValue = {
  theme: Theme;
  toggle: () => void;
  setTheme: (theme: Theme) => void;
};

const ThemeContext = createContext<ThemeContextValue | undefined>(undefined);

type ThemeStorage = {
  getItem: (key: string) => string | null;
  setItem: (key: string, value: string) => void;
};

type ThemeMatchMedia = (query: string) => {
  matches: boolean;
  addEventListener?: (type: "change", listener: (event: MediaQueryListEvent) => void) => void;
  removeEventListener?: (type: "change", listener: (event: MediaQueryListEvent) => void) => void;
  addListener?: (listener: (event: MediaQueryListEvent) => void) => void;
  removeListener?: (listener: (event: MediaQueryListEvent) => void) => void;
} | null;

export type ThemeProviderProps = {
  children: ReactNode;
  storage?: ThemeStorage;
  matchMedia?: ThemeMatchMedia;
};

function resolveStorage(storage?: ThemeStorage): ThemeStorage | null {
  if (storage) {
    return storage;
  }

  if (typeof window === "undefined" || typeof window.localStorage === "undefined") {
    return null;
  }

  return window.localStorage;
}

function resolveMatchMedia(matchMediaFn?: ThemeMatchMedia): ThemeMatchMedia | null {
  if (matchMediaFn) {
    return matchMediaFn;
  }

  if (typeof window === "undefined" || typeof window.matchMedia !== "function") {
    return null;
  }

  return (query: string) => window.matchMedia(query);
}

function getInitialTheme(storage: ThemeStorage | null, matchMediaFn: ThemeMatchMedia | null): Theme {
  const fromStorage = storage?.getItem(STORAGE_KEY) as Theme | null;
  if (fromStorage === "light" || fromStorage === "dark") {
    return fromStorage;
  }

  const mediaQuery = matchMediaFn?.(PREFERRED_DARK_QUERY);
  if (mediaQuery?.matches) {
    return "dark";
  }

  return "light";
}

export function ThemeProvider({ children, storage, matchMedia }: ThemeProviderProps) {
  const resolvedStorage = useMemo(() => resolveStorage(storage), [storage]);
  const resolvedMatchMedia = useMemo(() => resolveMatchMedia(matchMedia), [matchMedia]);

  const [theme, setThemeState] = useState<Theme>(() =>
    getInitialTheme(resolvedStorage, resolvedMatchMedia)
  );

  useEffect(() => {
    const mediaQuery = resolvedMatchMedia?.(PREFERRED_DARK_QUERY);

    if (!mediaQuery) {
      return undefined;
    }

    const handleChange = (event: MediaQueryListEvent) => {
      setThemeState(event.matches ? "dark" : "light");
    };

    if (typeof mediaQuery.addEventListener === "function") {
      mediaQuery.addEventListener("change", handleChange);
      return () => mediaQuery.removeEventListener?.("change", handleChange);
    }

    if (typeof mediaQuery.addListener === "function") {
      mediaQuery.addListener(handleChange);
      return () => mediaQuery.removeListener?.(handleChange);
    }

    return undefined;
  }, [resolvedMatchMedia]);

  useEffect(() => {
    if (typeof document === "undefined") {
      return;
    }

    const root = document.documentElement;
    root.classList.toggle("dark", theme === "dark");
    root.style.colorScheme = theme;

    resolvedStorage?.setItem(STORAGE_KEY, theme);
  }, [theme, resolvedStorage]);

  const setTheme = useCallback((value: Theme) => {
    setThemeState(value);
  }, []);

  const toggle = useCallback(() => {
    setThemeState((prev) => (prev === "dark" ? "light" : "dark"));
  }, []);

  const value = useMemo(
    () => ({
      theme,
      toggle,
      setTheme
    }),
    [theme, setTheme, toggle]
  );

  return <ThemeContext.Provider value={value}>{children}</ThemeContext.Provider>;
}

export function useTheme() {
  const ctx = useContext(ThemeContext);
  if (!ctx) {
    throw new Error("useTheme must be used within a ThemeProvider");
  }
  return ctx;
}
