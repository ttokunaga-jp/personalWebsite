import { useMemo } from "react";

import { useAdminModeContext } from "../providers/AdminModeProvider";

export function useAdminMode() {
  const context = useAdminModeContext();

  return useMemo(
    () => ({
      mode: context.mode,
      isAdminMode: context.isAdminMode,
      loading: context.loading,
      sessionActive: context.sessionActive,
      sessionEmail: context.sessionEmail,
      setMode: context.setMode,
      toggleMode: context.toggleMode,
      appendModeTo: context.appendModeTo,
      navigateWithMode: context.navigateWithMode,
      hasUnsavedChanges: context.hasUnsavedChanges,
      registerUnsavedChange: context.registerUnsavedChange,
      clearUnsavedChange: context.clearUnsavedChange,
      confirmIfUnsaved: context.confirmIfUnsaved,
    }),
    [context],
  );
}
