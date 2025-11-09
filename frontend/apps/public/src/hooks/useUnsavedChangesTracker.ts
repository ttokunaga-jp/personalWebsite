import { useEffect } from "react";

import { useAdminMode } from "./useAdminMode";

export function useUnsavedChangesTracker(id: string, dirty: boolean): void {
  const { registerUnsavedChange, clearUnsavedChange } = useAdminMode();

  useEffect(() => {
    if (dirty) {
      registerUnsavedChange(id);
    } else {
      clearUnsavedChange(id);
    }

    return () => {
      clearUnsavedChange(id);
    };
  }, [clearUnsavedChange, dirty, id, registerUnsavedChange]);
}
