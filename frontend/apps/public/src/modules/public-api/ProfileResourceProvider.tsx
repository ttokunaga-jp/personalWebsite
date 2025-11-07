import { type ReactNode } from "react";

import {
  ProfileResourceContext,
  useProfileResourceInternal,
} from "./client";

type ProfileResourceProviderProps = {
  children: ReactNode;
};

export function ProfileResourceProvider({
  children,
}: ProfileResourceProviderProps) {
  const resource = useProfileResourceInternal();

  return (
    <ProfileResourceContext.Provider value={resource}>
      {children}
    </ProfileResourceContext.Provider>
  );
}
