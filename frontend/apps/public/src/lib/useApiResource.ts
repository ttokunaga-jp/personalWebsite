import { useCallback, useEffect, useRef, useState } from "react";

type UseApiResourceState<T> = {
  data: T | null;
  isLoading: boolean;
  error: Error | null;
};

export type UseApiResourceResult<T> = UseApiResourceState<T> & {
  refetch: () => void;
};

export function useApiResource<T>(
  fetcher: (signal: AbortSignal) => Promise<T>
): UseApiResourceResult<T> {
  const [state, setState] = useState<UseApiResourceState<T>>({
    data: null,
    isLoading: true,
    error: null
  });
  const abortControllerRef = useRef<AbortController | null>(null);

  const runFetch = useCallback(async () => {
    abortControllerRef.current?.abort();
    const abortController = new AbortController();
    abortControllerRef.current = abortController;

    setState((previous) => ({
      data: previous.data,
      isLoading: true,
      error: null
    }));

    try {
      const response = await fetcher(abortController.signal);
      if (!abortController.signal.aborted) {
        setState({
          data: response,
          isLoading: false,
          error: null
        });
      }
    } catch (error) {
      if (!abortController.signal.aborted) {
        setState({
          data: null,
          isLoading: false,
          error: error instanceof Error ? error : new Error("Unknown error")
        });
      }
    }
  }, [fetcher]);

  useEffect(() => {
    runFetch();

    return () => {
      abortControllerRef.current?.abort();
    };
  }, [runFetch]);

  const refetch = useCallback(() => {
    void runFetch();
  }, [runFetch]);

  return {
    ...state,
    refetch
  };
}
