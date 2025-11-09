import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useRef,
  useState,
  type ReactNode,
} from "react";
import {
  useLocation,
  useNavigate,
  type NavigateOptions,
  type Path,
  type To,
} from "react-router-dom";

import { useAdminSession } from "../hooks/useAdminSession";

type AdminMode = "view" | "admin";

type ModeChangeOptions = {
  suppressPrompt?: boolean;
};

type AppendModeOptions = {
  targetMode?: AdminMode;
};

type AdminModeContextValue = {
  mode: AdminMode;
  isAdminMode: boolean;
  loading: boolean;
  sessionActive: boolean;
  sessionEmail?: string;
  setMode: (mode: AdminMode, options?: ModeChangeOptions) => boolean;
  toggleMode: (options?: ModeChangeOptions) => boolean;
  appendModeTo: (to: To, options?: AppendModeOptions) => To;
  navigateWithMode: (to: To, options?: NavigateOptions) => void;
  hasUnsavedChanges: boolean;
  registerUnsavedChange: (id: string) => void;
  clearUnsavedChange: (id: string) => void;
  confirmIfUnsaved: () => boolean;
};

const AdminModeContext = createContext<AdminModeContextValue | null>(null);

type AdminModeProviderProps = {
  children: ReactNode;
  unsavedPromptMessage?: string;
};

function normalizeMode(search: string | undefined | null): AdminMode {
  if (!search) {
    return "view";
  }
  const params = new URLSearchParams(search);
  return params.get("mode") === "admin" ? "admin" : "view";
}

function applyModeToSearch(
  search: string,
  mode: AdminMode,
): [string, boolean] {
  const params = new URLSearchParams(search);
  const previous = params.get("mode") === "admin" ? "admin" : "view";

  if (mode === "admin") {
    params.set("mode", "admin");
  } else {
    params.delete("mode");
  }

  const nextSearch = params.toString();
  return [nextSearch ? `?${nextSearch}` : "", previous !== mode];
}

type NormalizedTo = Partial<Path>;

function ensureToObject(to: To): NormalizedTo {
  if (typeof to === "number") {
    throw new Error("Relative navigation by delta is not supported here.");
  }

  if (typeof to === "string") {
    const [pathnameWithMaybeQuery, hash = ""] = to.split("#");
    const [pathname, query = ""] = pathnameWithMaybeQuery.split("?");
    return {
      pathname,
      search: query ? `?${query}` : "",
      hash: hash ? `#${hash}` : "",
    };
  }

  return {
    pathname: to.pathname,
    search: to.search,
    hash: to.hash,
  };
}

export function AdminModeProvider({
  children,
  unsavedPromptMessage = "You have unsaved changes. Continue and discard edits?",
}: AdminModeProviderProps) {
  const location = useLocation();
  const navigate = useNavigate();
  const { session, loading } = useAdminSession();
  const sessionActive = Boolean(session?.active);
  const sessionEmail = session?.email;

  const [mode, setModeState] = useState<AdminMode>(() =>
    normalizeMode(location.search),
  );
  const unsavedSetRef = useRef<Set<string>>(new Set());
  const [, forceUpdate] = useState(0);

  const hasUnsavedChanges = unsavedSetRef.current.size > 0;

  useEffect(() => {
    const nextMode = normalizeMode(location.search);
    if (nextMode === "admin" && !sessionActive && !loading) {
      // Strip admin mode when the session is inactive.
      const [searchWithFallback, shouldUpdate] = applyModeToSearch(
        location.search,
        "view",
      );
      if (shouldUpdate) {
        navigate(
          {
            pathname: location.pathname,
            search: searchWithFallback,
            hash: location.hash,
          },
          { replace: true },
        );
      }
      setModeState("view");
      return;
    }
    setModeState(nextMode);
  }, [
    location.hash,
    location.pathname,
    location.search,
    navigate,
    sessionActive,
    loading,
  ]);

  useEffect(() => {
    const handleBeforeUnload = (event: BeforeUnloadEvent) => {
      if (!unsavedSetRef.current.size) {
        return;
      }
      event.preventDefault();
      // Chrome requires returnValue to be set.
      // eslint-disable-next-line no-param-reassign
      event.returnValue = unsavedPromptMessage;
    };

    window.addEventListener("beforeunload", handleBeforeUnload);
    return () => {
      window.removeEventListener("beforeunload", handleBeforeUnload);
    };
  }, [unsavedPromptMessage]);

  const confirmIfUnsaved = useCallback((): boolean => {
    if (!unsavedSetRef.current.size) {
      return true;
    }
    return window.confirm(unsavedPromptMessage);
  }, [unsavedPromptMessage]);

  const updateUrlMode = useCallback(
    (nextMode: AdminMode) => {
      const [nextSearch, shouldUpdate] = applyModeToSearch(
        location.search,
        nextMode,
      );
      if (!shouldUpdate) {
        return;
      }
      navigate(
        {
          pathname: location.pathname,
          search: nextSearch,
          hash: location.hash,
        },
        { replace: true },
      );
    },
    [location.hash, location.pathname, location.search, navigate],
  );

  const setMode = useCallback(
    (nextMode: AdminMode, options: ModeChangeOptions = {}): boolean => {
      if (nextMode === "admin" && !sessionActive) {
        return false;
      }

      if (nextMode === mode) {
        return true;
      }

      if (!options.suppressPrompt && !confirmIfUnsaved()) {
        return false;
      }

      setModeState(nextMode);
      updateUrlMode(nextMode);
      return true;
    },
    [confirmIfUnsaved, mode, sessionActive, updateUrlMode],
  );

  const toggleMode = useCallback(
    (options?: ModeChangeOptions) => {
      return setMode(mode === "admin" ? "view" : "admin", options);
    },
    [mode, setMode],
  );

  const appendModeTo = useCallback(
    (to: To, options: AppendModeOptions = {}): To => {
      if (typeof to === "number") {
        return to;
      }

      const desiredMode =
        options.targetMode ?? (mode === "admin" ? "admin" : "view");

      const descriptor = ensureToObject(to);
      const [nextSearch] = applyModeToSearch(
        descriptor.search ?? "",
        desiredMode,
      );

      return {
        ...descriptor,
        search: nextSearch,
      } as Partial<Path>;
    },
    [mode],
  );

  const navigateWithMode = useCallback(
    (to: To, options?: NavigateOptions) => {
      if (!confirmIfUnsaved()) {
        return;
      }
      if (typeof to === "number") {
        navigate(to);
        return;
      }
      navigate(appendModeTo(to), options);
    },
    [appendModeTo, confirmIfUnsaved, navigate],
  );

  const registerUnsavedChange = useCallback((id: string) => {
    if (!id) {
      return;
    }
    const targetId = id.trim();
    if (!targetId) {
      return;
    }
    const next = new Set(unsavedSetRef.current);
    next.add(targetId);
    unsavedSetRef.current = next;
    forceUpdate((value) => value + 1);
  }, []);

  const clearUnsavedChange = useCallback((id: string) => {
    if (!id) {
      return;
    }
    const targetId = id.trim();
    if (!targetId) {
      return;
    }
    if (!unsavedSetRef.current.has(targetId)) {
      return;
    }
    const next = new Set(unsavedSetRef.current);
    next.delete(targetId);
    unsavedSetRef.current = next;
    forceUpdate((value) => value + 1);
  }, []);

  const value = useMemo<AdminModeContextValue>(
    () => ({
      mode,
      isAdminMode: mode === "admin",
      loading,
      sessionActive,
      sessionEmail,
      setMode,
      toggleMode,
      appendModeTo,
      navigateWithMode,
      hasUnsavedChanges,
      registerUnsavedChange,
      clearUnsavedChange,
      confirmIfUnsaved,
    }),
    [
      appendModeTo,
      clearUnsavedChange,
      confirmIfUnsaved,
      hasUnsavedChanges,
      loading,
      mode,
      navigateWithMode,
      registerUnsavedChange,
      sessionActive,
      sessionEmail,
      setMode,
      toggleMode,
    ],
  );

  return (
    <AdminModeContext.Provider value={value}>
      {children}
    </AdminModeContext.Provider>
  );
}

export function useAdminModeContext(): AdminModeContextValue {
  const context = useContext(AdminModeContext);
  if (!context) {
    throw new Error("useAdminModeContext must be used within AdminModeProvider");
  }
  return context;
}
