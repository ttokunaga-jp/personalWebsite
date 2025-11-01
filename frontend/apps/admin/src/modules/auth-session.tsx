import {
  registerAuthTokenInvalidator,
  registerAuthTokenProvider,
} from "@shared/lib/api-client";
import {
  createContext,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from "react";

const TOKEN_STORAGE_KEY = "admin.jwt";

type AuthTokenListener = (token: string | null) => void;

const listeners = new Set<AuthTokenListener>();

let cachedToken: string | null = null;
let initialized = false;

function readToken(): string | null {
  if (!initialized) {
    initialized = true;
    if (typeof window !== "undefined") {
      cachedToken = window.sessionStorage.getItem(TOKEN_STORAGE_KEY);
    } else {
      cachedToken = null;
    }
  }
  return cachedToken;
}

function persistToken(next: string | null): void {
  cachedToken = next;
  if (typeof window !== "undefined") {
    if (next) {
      window.sessionStorage.setItem(TOKEN_STORAGE_KEY, next);
    } else {
      window.sessionStorage.removeItem(TOKEN_STORAGE_KEY);
    }
  }
  listeners.forEach((listener) => listener(cachedToken));
}

function subscribe(listener: AuthTokenListener): () => void {
  listeners.add(listener);
  return () => listeners.delete(listener);
}

export function getToken(): string | null {
  return readToken();
}

export function setToken(token: string): void {
  persistToken(token.trim().length > 0 ? token : null);
}

export function clearToken(): void {
  persistToken(null);
}

export function extractTokenFromHash(hash: string): string | null {
  if (typeof hash !== "string" || hash.length === 0 || hash === "#") {
    return null;
  }

  let trimmed = hash;
  while (trimmed.startsWith("#")) {
    trimmed = trimmed.slice(1);
  }
  const params = new URLSearchParams(trimmed);
  const token = params.get("token");
  if (!token) {
    return null;
  }

  const normalized = token.trim();
  return normalized.length > 0 ? normalized : null;
}

type AuthSessionContextValue = {
  token: string | null;
  setToken: (token: string) => void;
  clearToken: () => void;
};

const AuthSessionContext = createContext<AuthSessionContextValue | undefined>(
  undefined,
);

export function AuthSessionProvider({
  children,
}: {
  readonly children: ReactNode;
}): JSX.Element {
  const [token, setTokenState] = useState<string | null>(() => getToken());

  useEffect(() => {
    setTokenState(getToken());
    return subscribe((nextToken) => setTokenState(nextToken));
  }, []);

  useEffect(() => {
    registerAuthTokenProvider(() => getToken());
    registerAuthTokenInvalidator(() => {
      clearToken();
    });
    return () => {
      registerAuthTokenProvider(null);
      registerAuthTokenInvalidator(null);
    };
  }, []);

  const value = useMemo<AuthSessionContextValue>(
    () => ({
      token,
      setToken,
      clearToken,
    }),
    [token],
  );

  return (
    <AuthSessionContext.Provider value={value}>
      {children}
    </AuthSessionContext.Provider>
  );
}

export function useAuthSession(): AuthSessionContextValue {
  const context = useContext(AuthSessionContext);
  if (!context) {
    throw new Error(
      "useAuthSession must be used within an AuthSessionProvider",
    );
  }
  return context;
}
