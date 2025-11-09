import { useCallback } from "react";
import type { NavigateOptions, To } from "react-router-dom";

import { useAdminModeContext } from "../providers/AdminModeProvider";

export function useAdminAwareNavigate(): (to: To, options?: NavigateOptions) => void {
  const context = useAdminModeContext();

  return useCallback(
    (to: To, options?: NavigateOptions) => {
      context.navigateWithMode(to, options);
    },
    [context],
  );
}
