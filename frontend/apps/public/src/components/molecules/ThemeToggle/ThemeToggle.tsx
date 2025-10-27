import { useTranslation } from "react-i18next";

import { useTheme } from "../../../providers/ThemeProvider";
import { IconButton } from "../../atoms/IconButton";

function SunIcon() {
  return (
    <svg
      aria-hidden="true"
      viewBox="0 0 24 24"
      className="h-5 w-5 fill-current"
    >
      <path d="M12 4.75a1 1 0 0 1-1-1V2a1 1 0 1 1 2 0v1.75a1 1 0 0 1-1 1Zm7.25 7.25a1 1 0 0 1 1-1H22a1 1 0 1 1 0 2h-1.75a1 1 0 0 1-1-1ZM12 19.25a1 1 0 0 1 1 1V22a1 1 0 1 1-2 0v-1.75a1 1 0 0 1 1-1ZM3.75 12a1 1 0 0 1-1 1H1a1 1 0 1 1 0-2h1.75a1 1 0 0 1 1 1ZM6.47 7.53a1 1 0 0 1-1.41-1.41L6.3 4.88a1 1 0 1 1 1.41 1.42Zm12.77 12.77a1 1 0 0 1-1.41 0l-1.06-1.06a1 1 0 1 1 1.41-1.41l1.06 1.06a1 1 0 0 1 0 1.41ZM4.88 17.7a1 1 0 0 1 1.42-1.41l1.06 1.06a1 1 0 0 1-1.42 1.41Zm12.72-12.82a1 1 0 0 1 0-1.41l1.06-1.06a1 1 0 0 1 1.41 1.41l-1.06 1.06a1 1 0 0 1-1.41 0ZM12 7a5 5 0 1 0 5 5a5 5 0 0 0-5-5Zm0 8a3 3 0 1 1 3-3a3 3 0 0 1-3 3Z" />
    </svg>
  );
}

function MoonIcon() {
  return (
    <svg
      aria-hidden="true"
      viewBox="0 0 24 24"
      className="h-5 w-5 fill-current"
    >
      <path d="M20.354 15.354a1 1 0 0 0-1.058-.21A7 7 0 0 1 8.856 5.5a7.09 7.09 0 0 1 .354-2.144a1 1 0 0 0-1.266-1.266a9 9 0 1 0 12.966 12.966a1 1 0 0 0-.198-1.702Z" />
    </svg>
  );
}

export function ThemeToggle() {
  const { theme, toggle } = useTheme();
  const { t } = useTranslation();

  const isDark = theme === "dark";

  return (
    <IconButton
      aria-label={t(isDark ? "themeToggle.setLight" : "themeToggle.setDark")}
      onClick={toggle}
      title={t(isDark ? "themeToggle.setLight" : "themeToggle.setDark")}
    >
      {isDark ? <SunIcon /> : <MoonIcon />}
    </IconButton>
  );
}
