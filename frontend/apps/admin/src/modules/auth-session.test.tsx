import { renderHook, act } from "@testing-library/react";
import type { ReactNode } from "react";
import { describe, expect, it } from "vitest";

import {
  AuthSessionProvider,
  useAuthSession,
  type AdminSessionInfo,
} from "./auth-session";

describe("AuthSessionProvider", () => {
  it("provides session context with setters", () => {
    const wrapper = ({ children }: { children: ReactNode }) => (
      <AuthSessionProvider>{children}</AuthSessionProvider>
    );

    const { result } = renderHook(() => useAuthSession(), { wrapper });
    expect(result.current.session).toBeNull();

    const session: AdminSessionInfo = {
      active: true,
      email: "admin@example.com",
      roles: ["admin"],
    };

    act(() => {
      result.current.setSession(session);
    });

    expect(result.current.session).not.toBeNull();
    expect(result.current.session?.email).toBe("admin@example.com");

    act(() => {
      result.current.clearSession();
    });

    expect(result.current.session).toBeNull();
  });
});
