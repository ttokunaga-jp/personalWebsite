import type { SVGProps } from "react";

import type { SocialProvider } from "../modules/public-api";

type IconProps = SVGProps<SVGSVGElement>;
type IconComponent = (props: IconProps) => JSX.Element;

const GitHubIcon: IconComponent = (props) => (
  <svg viewBox="0 0 24 24" fill="currentColor" focusable="false" {...props}>
    <path d="M12 .5a12 12 0 0 0-3.79 23.4c.6.11.82-.26.82-.58l-.02-2.24c-3.34.73-4.04-1.6-4.04-1.6-.55-1.4-1.34-1.77-1.34-1.77-1.1-.76.09-.75.09-.75 1.22.09 1.87 1.27 1.87 1.27 1.08 1.86 2.84 1.32 3.53 1 .11-.78.42-1.32.77-1.62-2.67-.3-5.47-1.36-5.47-6.07 0-1.34.48-2.43 1.27-3.29-.13-.31-.55-1.56.12-3.25 0 0 1.02-.33 3.34 1.25.97-.27 2-.4 3.03-.41 1.03 0 2.06.14 3.03.41 2.31-1.58 3.33-1.25 3.33-1.25.67 1.69.25 2.94.12 3.25.79.86 1.27 1.95 1.27 3.29 0 4.72-2.8 5.77-5.48 6.07.43.37.82 1.1.82 2.22l-.01 3.29c0 .32.22.7.83.58A12 12 0 0 0 12 .5Z" />
  </svg>
);

const TwitterIcon: IconComponent = (props) => (
  <svg viewBox="0 0 24 24" fill="currentColor" focusable="false" {...props}>
    <path d="M20.5 6.03c.01.17.01.35.01.53 0 5.37-4.08 11.57-11.57 11.57-2.3 0-4.44-.67-6.24-1.84.32.04.64.06.97.06 1.9 0 3.64-.65 5.02-1.75a4.09 4.09 0 0 1-3.82-2.84c.26.05.53.08.81.08.39 0 .77-.05 1.12-.15a4.08 4.08 0 0 1-3.27-4c0-.02 0-.05.01-.07.55.3 1.18.49 1.85.51a4.07 4.07 0 0 1-1.82-3.4c0-.75.2-1.45.55-2.05a11.58 11.58 0 0 0 8.41 4.26 4.61 4.61 0 0 1-.1-.93 4.08 4.08 0 0 1 7.06-2.79 8.02 8.02 0 0 0 2.59-.99 4.08 4.08 0 0 1-1.8 2.25 8.15 8.15 0 0 0 2.34-.64 8.78 8.78 0 0 1-2.04 2.11Z" />
  </svg>
);

const LinkedInIcon: IconComponent = (props) => (
  <svg viewBox="0 0 24 24" fill="currentColor" focusable="false" {...props}>
    <path d="M19 3H5a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V5a2 2 0 0 0-2-2ZM8.34 17.34H6.16V9.75h2.18Zm-1.09-8.7a1.27 1.27 0 1 1 1.27-1.27 1.27 1.27 0 0 1-1.27 1.27Zm10.08 8.7h-2.18v-3.75c0-.89-.02-2.04-1.24-2.04-1.25 0-1.44.97-1.44 1.98v3.81h-2.18V9.75h2.09v1.03h.03a2.3 2.3 0 0 1 2.07-1.14c2.22 0 2.63 1.46 2.63 3.36Z" />
  </svg>
);

const MailIcon: IconComponent = (props) => (
  <svg viewBox="0 0 24 24" fill="currentColor" focusable="false" {...props}>
    <path d="M20 4H4a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2V6a2 2 0 0 0-2-2Zm0 2-8 5-8-5ZM4 18V8l8 5 8-5v10Z" />
  </svg>
);

const GlobeIcon: IconComponent = (props) => (
  <svg viewBox="0 0 24 24" fill="currentColor" focusable="false" {...props}>
    <path d="M12 2a10 10 0 1 0 0 20 10 10 0 0 0 0-20Zm6.32 6H16a18.19 18.19 0 0 0-1.45-3.64A8.05 8.05 0 0 1 18.32 8ZM12 4c.74 0 1.9 1.53 2.58 4H9.42C10.1 5.53 11.26 4 12 4ZM4 12a7.96 7.96 0 0 1 1.05-4h3.16a19.72 19.72 0 0 0-.21 4c0 1.41.08 2.77.21 4H5.05A7.96 7.96 0 0 1 4 12Zm8 8c-.74 0-1.9-1.53-2.58-4h5.16C13.9 18.47 12.74 20 12 20Zm2.79-6H9.21a17.32 17.32 0 0 1-.24-4c0-1.37.09-2.73.24-4h5.58c.15 1.27.24 2.63.24 4 0 1.37-.09 2.73-.24 4Zm1.66 4A18.19 18.19 0 0 0 16 13h2.32a8.05 8.05 0 0 1-1.87 5Zm1.5-7H16c.13-1.23.2-2.6.2-4 0-1.4-.07-2.77-.2-4h2.95a8 8 0 0 1 0 8Z" />
  </svg>
);

const DefaultIcon: IconComponent = (props) => (
  <svg viewBox="0 0 24 24" fill="currentColor" focusable="false" {...props}>
    <circle cx="12" cy="12" r="10" opacity="0.15" />
    <path d="M12 6.9a3.1 3.1 0 1 0 0 6.2 3.1 3.1 0 0 0 0-6.2Zm0 10.8c-2.34 0-4.46-1.2-5.68-3.02a6.89 6.89 0 0 1 11.36 0A6.78 6.78 0 0 1 12 17.7Z" />
  </svg>
);

const ZennIcon: IconComponent = (props) => (
  <svg viewBox="0 0 24 24" fill="currentColor" focusable="false" {...props}>
    <path d="M18.69 5.31A10 10 0 1 0 5.31 18.69 10 10 0 0 0 18.69 5.31ZM7.94 8.38h8.12V10h-4.8l4.8 5.62v1.62H7.94v-1.6h4.96l-4.96-5.8Z" />
  </svg>
);

const iconMap: Partial<Record<SocialProvider, IconComponent>> = {
  github: GitHubIcon,
  x: TwitterIcon,
  linkedin: LinkedInIcon,
  email: MailIcon,
  website: GlobeIcon,
  zenn: ZennIcon,
};

export function getSocialIcon(provider: SocialProvider): IconComponent {
  return iconMap[provider] ?? DefaultIcon;
}
