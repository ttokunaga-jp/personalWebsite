import { useTranslation } from "react-i18next";

import { getCanonicalProfile } from "../../../modules/profile-content";
import { useProfileResource } from "../../../modules/public-api";
import { getSocialIcon } from "../../../utils/icons";

export function SiteFooter() {
  const { t, i18n } = useTranslation();
  const { data: profile } = useProfileResource();
  const year = new Date().getFullYear();
  const effectiveProfile = profile ?? getCanonicalProfile(i18n.language);

  const footerLinks = effectiveProfile.footerLinks.length
    ? effectiveProfile.footerLinks
    : effectiveProfile.socialLinks.filter((link) => link.isFooter);

  return (
    <footer className="border-t border-slate-200 bg-white/80 backdrop-blur transition-colors dark:border-slate-800 dark:bg-slate-950/80">
      <div className="mx-auto flex w-full max-w-6xl flex-col gap-4 px-4 py-6 text-sm text-slate-500 sm:flex-row sm:items-center sm:justify-between sm:px-8 dark:text-slate-400">
        <p>
          © {year} {t("footer.copyright")}
        </p>
        <nav className="flex flex-wrap gap-4">
          {footerLinks.map((link) => {
            const Icon = getSocialIcon(link.provider);
            return (
              <a
                key={link.id}
                href={link.url}
                target={link.provider === "email" ? "_self" : "_blank"}
                rel={link.provider === "email" ? undefined : "noreferrer"}
                className="inline-flex items-center gap-2 transition hover:text-sky-600 dark:hover:text-sky-300"
              >
                <Icon aria-hidden className="h-4 w-4" />
                <span>{link.label}</span>
              </a>
            );
          })}
        </nav>
      </div>
    </footer>
  );
}
