import {
  registerAuthTokenInvalidator,
  registerAuthTokenProvider,
} from "@shared/lib/api-client";
import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from "react";

export type AdminSessionInfo = {
  active: boolean;
  email?: string;
  roles?: string[];
  expiresAt?: number;
  source?: string;
  refreshed?: boolean;
};

type AuthSessionContextValue = {
  session: AdminSessionInfo | null;
  setSession: (info: AdminSessionInfo | null) => void;
  clearSession: () => void;
};

const AuthSessionContext = createContext<AuthSessionContextValue | undefined>(
  undefined,
);

let currentSession: AdminSessionInfo | null = null;

export function getSession(): AdminSessionInfo | null {
  return currentSession;
}

export function AuthSessionProvider({
  children,
}: {
  readonly children: ReactNode;
}): JSX.Element {
  const [session, setSessionState] = useState<AdminSessionInfo | null>(null);

  useEffect(() => {
    registerAuthTokenProvider(() => null);
    registerAuthTokenInvalidator(() => {
      setSessionState(null);
      currentSession = null;
    });
    return () => {
      registerAuthTokenProvider(null);
      registerAuthTokenInvalidator(null);
    };
  }, []);

  const setSession = useCallback((info: AdminSessionInfo | null) => {
    currentSession = info;
    setSessionState(info);
  }, []);

  const clearSession = useCallback(() => {
    currentSession = null;
    setSessionState(null);
  }, []);

  const value = useMemo<AuthSessionContextValue>(
    () => ({
      session,
      setSession,
      clearSession,
    }),
    [session, setSession, clearSession],
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
