import { Fragment } from "react";
import { useTranslation } from "react-i18next";

import { Button } from "../../atoms/Button";

const languages = [
  { code: "en", labelKey: "languages.en" },
  { code: "ja", labelKey: "languages.ja" }
] as const;

export function LanguageSwitcher() {
  const { i18n, t } = useTranslation();
  const current = i18n.language;

  return (
    <div className="flex items-center gap-1 rounded-full border border-slate-200 bg-white/70 p-1 dark:border-slate-700 dark:bg-slate-900/40">
      {languages.map((language, index) => {
        const isActive = current === language.code;
        return (
          <Fragment key={language.code}>
            <Button
              size="sm"
              variant={isActive ? "primary" : "ghost"}
              aria-pressed={isActive}
              onClick={() => {
                if (current !== language.code) {
                  void i18n.changeLanguage(language.code);
                }
              }}
            >
              {t(language.labelKey)}
            </Button>
            {index < languages.length - 1 ? (
              <span className="text-xs text-slate-400 dark:text-slate-600">/</span>
            ) : null}
          </Fragment>
        );
      })}
    </div>
  );
}
