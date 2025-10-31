import { useCallback, useEffect, useMemo, useRef, useState } from "react";

type UseApiResourceState<T> = {
  data: T | null;
  isLoading: boolean;
  error: Error | null;
};

export type UseApiResourceResult<T> = UseApiResourceState<T> & {
  refetch: () => void;
};

type UseApiResourceOptions<T> = {
  initialData?: T | (() => T);
  skip?: boolean;
};

export function useApiResource<T>(
  fetcher: (signal: AbortSignal) => Promise<T>,
  options?: UseApiResourceOptions<T>,
): UseApiResourceResult<T> {
  const { initialData: initialDataOption, skip = false } = options ?? {};
  const initialDataRef = useRef<T | null>();
  const abortControllerRef = useRef<AbortController | null>(null);

  const initialData = useMemo(() => {
    if (initialDataRef.current !== undefined) {
      return initialDataRef.current;
    }
    if (typeof initialDataOption === "function") {
      initialDataRef.current = (initialDataOption as () => T)();
    } else if (initialDataOption !== undefined) {
      initialDataRef.current = initialDataOption ?? null;
    } else {
      initialDataRef.current = null;
    }
    return initialDataRef.current;
  }, [initialDataOption]);

  const [state, setState] = useState<UseApiResourceState<T>>({
    data: initialData,
    isLoading: !skip && !initialData,
    error: null,
  });

  const runFetch = useCallback(async () => {
    if (skip) {
      return;
    }

    abortControllerRef.current?.abort();
    const abortController = new AbortController();
    abortControllerRef.current = abortController;

    setState((previous) => ({
      data: previous.data ?? initialData,
      isLoading: true,
      error: null,
    }));

    try {
      const response = await fetcher(abortController.signal);
      if (!abortController.signal.aborted) {
        setState({
          data: response,
          isLoading: false,
          error: null,
        });
      }
    } catch (error) {
      if (!abortController.signal.aborted) {
        setState({
          data: initialData,
          isLoading: false,
          error: error instanceof Error ? error : new Error("Unknown error"),
        });
      }
    }
  }, [fetcher, initialData, skip]);

  useEffect(() => {
    if (skip) {
      return () => {
        abortControllerRef.current?.abort();
      };
    }

    void runFetch();

    return () => {
      abortControllerRef.current?.abort();
    };
  }, [runFetch, skip]);

  const refetch = useCallback(() => {
    if (skip) {
      return;
    }
    void runFetch();
  }, [runFetch, skip]);

  return {
    ...state,
    refetch,
  };
}
