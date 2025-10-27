import { useTranslation } from "react-i18next";

export function SiteFooter() {
  const { t } = useTranslation();
  const year = new Date().getFullYear();

  return (
    <footer className="border-t border-slate-200 bg-white/80 backdrop-blur transition-colors dark:border-slate-800 dark:bg-slate-950/80">
      <div className="mx-auto flex w-full max-w-6xl flex-col gap-4 px-4 py-6 text-sm text-slate-500 sm:flex-row sm:items-center sm:justify-between sm:px-8 dark:text-slate-400">
        <p>
          © {year} {t("footer.copyright")}
        </p>
        <div className="flex flex-wrap gap-4">
          <a
            href="https://github.com/"
            target="_blank"
            rel="noreferrer"
            className="transition hover:text-sky-600 dark:hover:text-sky-300"
          >
            {t("footer.links.github")}
          </a>
          <a
            href="https://x.com/"
            target="_blank"
            rel="noreferrer"
            className="transition hover:text-sky-600 dark:hover:text-sky-300"
          >
            {t("footer.links.twitter")}
          </a>
          <a
            href="mailto:someone@example.com"
            className="transition hover:text-sky-600 dark:hover:text-sky-300"
          >
            {t("footer.links.contact")}
          </a>
        </div>
      </div>
    </footer>
  );
}
