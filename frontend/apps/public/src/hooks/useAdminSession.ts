import { apiClient } from "@shared/lib/api-client";
import { useEffect, useState } from "react";

type AdminSessionResponse = {
  active: boolean;
  email?: string;
  roles?: string[];
  expiresAt?: number;
};

export function useAdminSession(): {
  session: AdminSessionResponse | null;
  loading: boolean;
} {
  const [session, setSession] = useState<AdminSessionResponse | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let cancelled = false;

    const fetchSession = async () => {
      try {
        const response = await apiClient.get<
          AdminSessionResponse
        >("/admin/auth/session");
        if (cancelled) {
          return;
        }
        if (response.data?.active) {
          setSession(response.data);
        } else {
          setSession(null);
        }
      } catch {
        if (!cancelled) {
          setSession(null);
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    };

    void fetchSession();
    return () => {
      cancelled = true;
    };
  }, []);

  return { session, loading };
}
