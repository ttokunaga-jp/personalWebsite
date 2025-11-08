import { FormEvent, useCallback, useEffect, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";

import { adminApi, DomainError } from "./modules/admin-api";
import { useAuthSession } from "./modules/auth-session";
import type {
  AdminProfile,
  AdminProject,
  AdminResearch,
  AdminSummary,
  BlacklistEntry,
  ContactFormSettings,
  ContactMessage,
  ContactStatus,
  ResearchKind,
  ResearchLinkType,
  HomePageConfig,
  HomeQuickLinkSection,
  HomeChipSourceType,
  SocialProvider,
  TechCatalogEntry,
  TechContext,
} from "./types";

type ProfilePayload = Parameters<typeof adminApi.updateProfile>[0];
type HomeSettingsPayload = Parameters<typeof adminApi.updateHomeSettings>[0];

const currentYear = new Date().getFullYear();
const contactStatuses: ContactStatus[] = [
  "pending",
  "in_review",
  "resolved",
  "archived",
];

const researchKinds: ResearchKind[] = ["research", "blog"];
const researchLinkTypes: ResearchLinkType[] = [
  "paper",
  "slides",
  "video",
  "code",
  "external",
];
const techContexts: TechContext[] = ["primary", "supporting"];

type AffiliationForm = {
  id?: number;
  name: string;
  url: string;
  descriptionJa: string;
  descriptionEn: string;
  startedAt: string;
  sortOrder: string;
};

type WorkHistoryForm = {
  id?: number;
  organizationJa: string;
  organizationEn: string;
  roleJa: string;
  roleEn: string;
  summaryJa: string;
  summaryEn: string;
  startedAt: string;
  endedAt: string;
  externalUrl: string;
  sortOrder: string;
};

type SocialLinkForm = {
  id?: number;
  provider: SocialProvider;
  labelJa: string;
  labelEn: string;
  url: string;
  isFooter: boolean;
  sortOrder: string;
};

type HomeQuickLinkForm = {
  id?: number;
  section: HomeQuickLinkSection;
  labelJa: string;
  labelEn: string;
  descriptionJa: string;
  descriptionEn: string;
  ctaJa: string;
  ctaEn: string;
  targetUrl: string;
  sortOrder: string;
};

type HomeChipSourceForm = {
  id?: number;
  source: HomeChipSourceType;
  labelJa: string;
  labelEn: string;
  limit: string;
  sortOrder: string;
};

type HomeSettingsFormState = {
  id?: number | null;
  profileId?: number | null;
  heroSubtitleJa: string;
  heroSubtitleEn: string;
  quickLinks: HomeQuickLinkForm[];
  chipSources: HomeChipSourceForm[];
  updatedAt?: string;
};

type ProfileFormState = {
  displayName: string;
  headlineJa: string;
  headlineEn: string;
  summaryJa: string;
  summaryEn: string;
  avatarUrl: string;
  locationJa: string;
  locationEn: string;
  themeMode: "light" | "dark" | "system";
  themeAccentColor: string;
  labNameJa: string;
  labNameEn: string;
  labAdvisorJa: string;
  labAdvisorEn: string;
  labRoomJa: string;
  labRoomEn: string;
  labUrl: string;
  affiliations: AffiliationForm[];
  communities: AffiliationForm[];
  workHistory: WorkHistoryForm[];
  socialLinks: SocialLinkForm[];
};

type ProjectFormState = {
  titleJa: string;
  titleEn: string;
  descriptionJa: string;
  descriptionEn: string;
  tech: ProjectTechForm[];
  linkUrl: string;
  year: string;
  published: boolean;
  sortOrder: string;
};

type ProjectTechForm = {
  membershipId?: number;
  techId: string;
  context: TechContext;
  note: string;
  sortOrder: string;
};

type ResearchTagForm = {
  id?: number;
  value: string;
  sortOrder: string;
};

type ResearchLinkForm = {
  id?: number;
  type: ResearchLinkType;
  labelJa: string;
  labelEn: string;
  url: string;
  sortOrder: string;
};

type ResearchAssetForm = {
  id?: number;
  url: string;
  captionJa: string;
  captionEn: string;
  sortOrder: string;
};

type ResearchTechForm = {
  membershipId?: number;
  techId: string;
  context: TechContext;
  note: string;
  sortOrder: string;
};

type ResearchFormState = {
  slug: string;
  kind: ResearchKind;
  titleJa: string;
  titleEn: string;
  overviewJa: string;
  overviewEn: string;
  outcomeJa: string;
  outcomeEn: string;
  outlookJa: string;
  outlookEn: string;
  externalUrl: string;
  highlightImageUrl: string;
  imageAltJa: string;
  imageAltEn: string;
  publishedAt: string;
  isDraft: boolean;
  tags: ResearchTagForm[];
  links: ResearchLinkForm[];
  assets: ResearchAssetForm[];
  tech: ResearchTechForm[];
};

type BlacklistFormState = {
  email: string;
  reason: string;
};

type ContactEditState = {
  topic: string;
  message: string;
  status: ContactStatus;
  adminNote: string;
};

type ContactTopicForm = {
  id: string;
  labelJa: string;
  labelEn: string;
  descriptionJa: string;
  descriptionEn: string;
};

type ContactSettingsFormState = {
  heroTitleJa: string;
  heroTitleEn: string;
  heroDescriptionJa: string;
  heroDescriptionEn: string;
  consentTextJa: string;
  consentTextEn: string;
  minimumLeadHours: string;
  recaptchaSiteKey: string;
  supportEmail: string;
  calendarTimezone: string;
  googleCalendarId: string;
  bookingWindowDays: string;
  topics: ContactTopicForm[];
};

type NormalizedLocalizedText = {
  ja: string;
  en: string;
};

type ContactSettingsNormalized = {
  heroTitle: NormalizedLocalizedText;
  heroDescription: NormalizedLocalizedText;
  consentText: NormalizedLocalizedText;
  topics: {
    id: string;
    label: NormalizedLocalizedText;
    description: NormalizedLocalizedText;
  }[];
  minimumLeadHours: number;
  recaptchaSiteKey: string;
  supportEmail: string;
  calendarTimezone: string;
  googleCalendarId: string;
  bookingWindowDays: number;
};

type ContactSettingsDiffEntry = {
  key: string;
  labelKey: string;
  original: string;
  updated: string;
};

type ContactTopicDiffEntry = {
  id: string;
  change: "added" | "removed" | "updated";
  original?: {
    label: NormalizedLocalizedText;
    description: NormalizedLocalizedText;
  };
  updated?: {
    label: NormalizedLocalizedText;
    description: NormalizedLocalizedText;
  };
};

type ContactSettingsDiffSummary = {
  fields: ContactSettingsDiffEntry[];
  topics: ContactTopicDiffEntry[];
};

type AuthState = "checking" | "authenticated" | "unauthorized";

const isUnauthorizedError = (error: unknown): boolean =>
  error instanceof DomainError && error.status === 401;

const getHttpStatus = (error: unknown): number | undefined =>
  (error as { response?: { status?: number } })?.response?.status;

const isNotFoundError = (error: unknown): boolean =>
  getHttpStatus(error) === 404;

const isConflictError = (error: unknown): boolean =>
  getHttpStatus(error) === 409;

const createEmptyAffiliation = (): AffiliationForm => ({
  name: "",
  url: "",
  descriptionJa: "",
  descriptionEn: "",
  startedAt: "",
  sortOrder: "",
});

const createEmptyWorkHistory = (): WorkHistoryForm => ({
  organizationJa: "",
  organizationEn: "",
  roleJa: "",
  roleEn: "",
  summaryJa: "",
  summaryEn: "",
  startedAt: "",
  endedAt: "",
  externalUrl: "",
  sortOrder: "",
});

const createSocialLink = (provider: SocialProvider): SocialLinkForm => ({
  provider,
  labelJa: "",
  labelEn: "",
  url: "",
  isFooter: true,
  sortOrder: "",
});

const createEmptyProfileForm = (): ProfileFormState => ({
  displayName: "",
  headlineJa: "",
  headlineEn: "",
  summaryJa: "",
  summaryEn: "",
  avatarUrl: "",
  locationJa: "",
  locationEn: "",
  themeMode: "system",
  themeAccentColor: "",
  labNameJa: "",
  labNameEn: "",
  labAdvisorJa: "",
  labAdvisorEn: "",
  labRoomJa: "",
  labRoomEn: "",
  labUrl: "",
  affiliations: [createEmptyAffiliation()],
  communities: [createEmptyAffiliation()],
  workHistory: [createEmptyWorkHistory()],
  socialLinks: [
    createSocialLink("github"),
    createSocialLink("zenn"),
    createSocialLink("linkedin"),
  ],
});

const createEmptyQuickLink = (): HomeQuickLinkForm => ({
  section: "profile",
  labelJa: "",
  labelEn: "",
  descriptionJa: "",
  descriptionEn: "",
  ctaJa: "",
  ctaEn: "",
  targetUrl: "",
  sortOrder: "",
});

const createEmptyChipSource = (): HomeChipSourceForm => ({
  source: "affiliation",
  labelJa: "",
  labelEn: "",
  limit: "",
  sortOrder: "",
});

const createEmptyHomeSettingsForm = (
  profileId?: number | null,
): HomeSettingsFormState => ({
  id: null,
  profileId: profileId ?? null,
  heroSubtitleJa: "",
  heroSubtitleEn: "",
  quickLinks: [createEmptyQuickLink()],
  chipSources: [createEmptyChipSource()],
  updatedAt: undefined,
});

const themeOptions: ProfileFormState["themeMode"][] = [
  "system",
  "light",
  "dark",
];

const socialProviderOptions: SocialProvider[] = [
  "github",
  "zenn",
  "linkedin",
  "x",
  "email",
  "other",
];

const requiredSocialProviders: SocialProvider[] = [
  "github",
  "zenn",
  "linkedin",
];

const homeQuickLinkSections: HomeQuickLinkSection[] = [
  "profile",
  "research_blog",
  "projects",
  "contact",
];

const homeChipSourceOptions: HomeChipSourceType[] = [
  "affiliation",
  "community",
  "skill",
];

type ListUpdater<T> = (prev: T[]) => T[];

type ProfileAffiliationListProps = {
  title: string;
  items: AffiliationForm[];
  onChange: (updater: ListUpdater<AffiliationForm>) => void;
};

type HomeQuickLinksEditorProps = {
  items: HomeQuickLinkForm[];
  onChange: (updater: ListUpdater<HomeQuickLinkForm>) => void;
};

type HomeChipSourcesEditorProps = {
  items: HomeChipSourceForm[];
  onChange: (updater: ListUpdater<HomeChipSourceForm>) => void;
};

const ProfileAffiliationList = ({
  title,
  items,
  onChange,
}: ProfileAffiliationListProps) => {
  const { t } = useTranslation();

  return (
    <div>
      <h3 className="text-sm font-semibold text-slate-700">{title}</h3>
      <div className="mt-2 space-y-4">
        {items.map((item, index) => (
          <div
            key={item.id ?? index}
            className="space-y-3 rounded-lg border border-slate-200 p-4"
          >
            <div className="grid gap-4 md:grid-cols-2">
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.name", { defaultValue: "Name" })}
                </label>
                <input
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={item.name}
                  onChange={(event) =>
                    onChange((prev) => {
                      const next = [...prev];
                      next[index] = { ...next[index], name: event.target.value };
                      return next;
                    })
                  }
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.url", { defaultValue: "URL" })}
                </label>
                <input
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={item.url}
                  onChange={(event) =>
                    onChange((prev) => {
                      const next = [...prev];
                      next[index] = { ...next[index], url: event.target.value };
                      return next;
                    })
                  }
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.startedAt", { defaultValue: "Started at" })}
                </label>
                <input
                  type="datetime-local"
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={toDateTimeLocal(item.startedAt)}
                  onChange={(event) =>
                    onChange((prev) => {
                      const next = [...prev];
                      next[index] = {
                        ...next[index],
                        startedAt: toISOStringWithFallback(event.target.value),
                      };
                      return next;
                    })
                  }
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.sortOrder", { defaultValue: "Sort order" })}
                </label>
                <input
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={item.sortOrder}
                  onChange={(event) =>
                    onChange((prev) => {
                      const next = [...prev];
                      next[index] = {
                        ...next[index],
                        sortOrder: event.target.value,
                      };
                      return next;
                    })
                  }
                />
              </div>
            </div>
            <div className="grid gap-4 md:grid-cols-2">
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.descriptionJa")}
                </label>
                <input
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={item.descriptionJa}
                  onChange={(event) =>
                    onChange((prev) => {
                      const next = [...prev];
                      next[index] = {
                        ...next[index],
                        descriptionJa: event.target.value,
                      };
                      return next;
                    })
                  }
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.descriptionEn")}
                </label>
                <input
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={item.descriptionEn}
                  onChange={(event) =>
                    onChange((prev) => {
                      const next = [...prev];
                      next[index] = {
                        ...next[index],
                        descriptionEn: event.target.value,
                      };
                      return next;
                    })
                  }
                />
              </div>
            </div>
            <div className="flex justify-end">
              <button
                className="rounded-md border border-red-200 px-2 py-1 text-xs text-red-600"
                type="button"
                onClick={() =>
                  onChange((prev) => prev.filter((_, i) => i !== index))
                }
                disabled={items.length === 1}
              >
                {t("actions.remove")}
              </button>
            </div>
          </div>
        ))}
      </div>
      <button
        className="mt-3 rounded-md border border-slate-200 px-3 py-1 text-sm text-slate-700"
        type="button"
        onClick={() =>
          onChange((prev) => [...prev, createEmptyAffiliation()])
        }
      >
        {t("actions.addAffiliation", { defaultValue: "Add affiliation" })}
      </button>
    </div>
  );
};

type ProfileWorkHistoryListProps = {
  items: WorkHistoryForm[];
  onChange: (updater: ListUpdater<WorkHistoryForm>) => void;
};

const ProfileWorkHistoryList = ({
  items,
  onChange,
}: ProfileWorkHistoryListProps) => {
  const { t } = useTranslation();

  return (
    <div>
      <h3 className="text-sm font-semibold text-slate-700">
        {t("profile.sections.workHistory", { defaultValue: "Work history" })}
      </h3>
      <div className="mt-2 space-y-4">
        {items.map((item, index) => (
          <div
            key={item.id ?? index}
            className="space-y-3 rounded-lg border border-slate-200 p-4"
          >
            <div className="grid gap-4 md:grid-cols-2">
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("profile.fields.organizationJa", {
                    defaultValue: "Organization (JA)",
                  })}
                </label>
                <input
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={item.organizationJa}
                  onChange={(event) =>
                    onChange((prev) => {
                      const next = [...prev];
                      next[index] = {
                        ...next[index],
                        organizationJa: event.target.value,
                      };
                      return next;
                    })
                  }
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("profile.fields.organizationEn", {
                    defaultValue: "Organization (EN)",
                  })}
                </label>
                <input
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={item.organizationEn}
                  onChange={(event) =>
                    onChange((prev) => {
                      const next = [...prev];
                      next[index] = {
                        ...next[index],
                        organizationEn: event.target.value,
                      };
                      return next;
                    })
                  }
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("profile.fields.roleJa", { defaultValue: "Role (JA)" })}
                </label>
                <input
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={item.roleJa}
                  onChange={(event) =>
                    onChange((prev) => {
                      const next = [...prev];
                      next[index] = {
                        ...next[index],
                        roleJa: event.target.value,
                      };
                      return next;
                    })
                  }
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("profile.fields.roleEn", { defaultValue: "Role (EN)" })}
                </label>
                <input
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={item.roleEn}
                  onChange={(event) =>
                    onChange((prev) => {
                      const next = [...prev];
                      next[index] = {
                        ...next[index],
                        roleEn: event.target.value,
                      };
                      return next;
                    })
                  }
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.startedAt", { defaultValue: "Started at" })}
                </label>
                <input
                  type="datetime-local"
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={toDateTimeLocal(item.startedAt)}
                  onChange={(event) =>
                    onChange((prev) => {
                      const next = [...prev];
                      next[index] = {
                        ...next[index],
                        startedAt: toISOStringWithFallback(event.target.value),
                      };
                      return next;
                    })
                  }
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.endedAt", { defaultValue: "Ended at" })}
                </label>
                <input
                  type="datetime-local"
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={toDateTimeLocal(item.endedAt)}
                  onChange={(event) =>
                    onChange((prev) => {
                      const next = [...prev];
                      next[index] = {
                        ...next[index],
                        endedAt: event.target.value
                          ? toISOStringWithFallback(event.target.value)
                          : "",
                      };
                      return next;
                    })
                  }
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.sortOrder", { defaultValue: "Sort order" })}
                </label>
                <input
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={item.sortOrder}
                  onChange={(event) =>
                    onChange((prev) => {
                      const next = [...prev];
                      next[index] = {
                        ...next[index],
                        sortOrder: event.target.value,
                      };
                      return next;
                    })
                  }
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.url", { defaultValue: "External URL" })}
                </label>
                <input
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={item.externalUrl}
                  onChange={(event) =>
                    onChange((prev) => {
                      const next = [...prev];
                      next[index] = {
                        ...next[index],
                        externalUrl: event.target.value,
                      };
                      return next;
                    })
                  }
                />
              </div>
            </div>
            <div className="grid gap-4 md:grid-cols-2">
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.descriptionJa")}
                </label>
                <textarea
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  rows={3}
                  value={item.summaryJa}
                  onChange={(event) =>
                    onChange((prev) => {
                      const next = [...prev];
                      next[index] = {
                        ...next[index],
                        summaryJa: event.target.value,
                      };
                      return next;
                    })
                  }
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.descriptionEn")}
                </label>
                <textarea
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  rows={3}
                  value={item.summaryEn}
                  onChange={(event) =>
                    onChange((prev) => {
                      const next = [...prev];
                      next[index] = {
                        ...next[index],
                        summaryEn: event.target.value,
                      };
                      return next;
                    })
                  }
                />
              </div>
            </div>
            <div className="flex justify-end">
              <button
                className="rounded-md border border-red-200 px-2 py-1 text-xs text-red-600"
                type="button"
                onClick={() =>
                  onChange((prev) => prev.filter((_, i) => i !== index))
                }
                disabled={items.length === 1}
              >
                {t("actions.remove")}
              </button>
            </div>
          </div>
        ))}
      </div>
      <button
        className="mt-3 rounded-md border border-slate-200 px-3 py-1 text-sm text-slate-700"
        type="button"
        onClick={() =>
          onChange((prev) => [...prev, createEmptyWorkHistory()])
        }
      >
        {t("actions.addWorkHistory", { defaultValue: "Add work history" })}
      </button>
    </div>
  );
};

type ProfileSocialLinkListProps = {
  items: SocialLinkForm[];
  onChange: (updater: ListUpdater<SocialLinkForm>) => void;
};

const ProfileSocialLinkList = ({
  items,
  onChange,
}: ProfileSocialLinkListProps) => {
  const { t } = useTranslation();

  return (
    <div>
      <h3 className="text-sm font-semibold text-slate-700">
        {t("profile.sections.socialLinks", {
          defaultValue: "Social links",
        })}
      </h3>
      <div className="mt-2 space-y-4">
        {items.map((item, index) => {
          const disableRemove = requiredSocialProviders.includes(item.provider);
          return (
            <div
              key={item.id ?? `${item.provider}-${index}`}
              className="space-y-3 rounded-lg border border-slate-200 p-4"
            >
              <div className="grid gap-4 md:grid-cols-2">
                <div>
                  <label className="block text-sm font-medium text-slate-700">
                    {t("fields.provider", { defaultValue: "Provider" })}
                  </label>
                  <select
                    className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                    value={item.provider}
                    onChange={(event) =>
                      onChange((prev) => {
                        const next = [...prev];
                        next[index] = {
                          ...next[index],
                          provider: event.target.value as SocialProvider,
                        };
                        return next;
                      })
                    }
                  >
                    {socialProviderOptions.map((provider) => (
                      <option key={provider} value={provider}>
                        {provider}
                      </option>
                    ))}
                  </select>
                </div>
                <div>
                  <label className="block text-sm font-medium text-slate-700">
                    {t("fields.url", { defaultValue: "URL" })}
                  </label>
                  <input
                    className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                    value={item.url}
                    onChange={(event) =>
                      onChange((prev) => {
                        const next = [...prev];
                        next[index] = { ...next[index], url: event.target.value };
                        return next;
                      })
                    }
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-slate-700">
                    {t("fields.labelJa")}
                  </label>
                  <input
                    className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                    value={item.labelJa}
                    onChange={(event) =>
                      onChange((prev) => {
                        const next = [...prev];
                        next[index] = {
                          ...next[index],
                          labelJa: event.target.value,
                        };
                        return next;
                      })
                    }
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-slate-700">
                    {t("fields.labelEn")}
                  </label>
                  <input
                    className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                    value={item.labelEn}
                    onChange={(event) =>
                      onChange((prev) => {
                        const next = [...prev];
                        next[index] = {
                          ...next[index],
                          labelEn: event.target.value,
                        };
                        return next;
                      })
                    }
                  />
                </div>
                <div className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    className="h-4 w-4 rounded border-slate-300 text-slate-900"
                    checked={item.isFooter}
                    onChange={(event) =>
                      onChange((prev) => {
                        const next = [...prev];
                        next[index] = {
                          ...next[index],
                          isFooter: event.target.checked,
                        };
                        return next;
                      })
                    }
                  />
                  <span className="text-sm text-slate-700">
                    {t("profile.fields.showInFooter", {
                      defaultValue: "Show in footer",
                    })}
                  </span>
                </div>
                <div>
                  <label className="block text-sm font-medium text-slate-700">
                    {t("fields.sortOrder", { defaultValue: "Sort order" })}
                  </label>
                  <input
                    className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                    value={item.sortOrder}
                    onChange={(event) =>
                      onChange((prev) => {
                        const next = [...prev];
                        next[index] = {
                          ...next[index],
                          sortOrder: event.target.value,
                        };
                        return next;
                      })
                    }
                  />
                </div>
              </div>
              <div className="flex justify-end">
                <button
                  className="rounded-md border border-red-200 px-2 py-1 text-xs text-red-600"
                  type="button"
                  onClick={() =>
                    onChange((prev) => prev.filter((_, i) => i !== index))
                  }
                  disabled={disableRemove}
                >
                  {t("actions.remove")}
                </button>
              </div>
            </div>
          );
        })}
      </div>
      <button
        className="mt-3 rounded-md border border-slate-200 px-3 py-1 text-sm text-slate-700"
        type="button"
        onClick={() =>
          onChange((prev) => [...prev, createSocialLink("other")])
        }
      >
        {t("actions.addSocialLink", { defaultValue: "Add social link" })}
      </button>
    </div>
  );
};

const HomeQuickLinksEditor = ({ items, onChange }: HomeQuickLinksEditorProps) => {
  const { t } = useTranslation();

  return (
    <div className="space-y-4">
      {items.map((item, index) => (
        <div
          key={item.id ?? index}
          className="space-y-3 rounded-lg border border-slate-200 p-4"
        >
          <div className="grid gap-4 md:grid-cols-2">
            <div>
              <label className="block text-sm font-medium text-slate-700">
                {t("homeSettings.fields.section", { defaultValue: "Section" })}
              </label>
              <select
                className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                value={item.section}
                onChange={(event) =>
                  onChange((prev) => {
                    const next = [...prev];
                    next[index] = {
                      ...next[index],
                      section: event.target.value as HomeQuickLinkSection,
                    };
                    return next;
                  })
                }
              >
                {homeQuickLinkSections.map((section) => (
                  <option key={section} value={section}>
                    {section}
                  </option>
                ))}
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-700">
                {t("fields.sortOrder", { defaultValue: "Sort order" })}
              </label>
              <input
                className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                value={item.sortOrder}
                onChange={(event) =>
                  onChange((prev) => {
                    const next = [...prev];
                    next[index] = {
                      ...next[index],
                      sortOrder: event.target.value,
                    };
                    return next;
                  })
                }
              />
            </div>
          </div>
          <div className="grid gap-4 md:grid-cols-2">
            <div>
              <label className="block text-sm font-medium text-slate-700">
                {t("fields.labelJa")}
              </label>
              <input
                className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                value={item.labelJa}
                onChange={(event) =>
                  onChange((prev) => {
                    const next = [...prev];
                    next[index] = {
                      ...next[index],
                      labelJa: event.target.value,
                    };
                    return next;
                  })
                }
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-700">
                {t("fields.labelEn")}
              </label>
              <input
                className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                value={item.labelEn}
                onChange={(event) =>
                  onChange((prev) => {
                    const next = [...prev];
                    next[index] = {
                      ...next[index],
                      labelEn: event.target.value,
                    };
                    return next;
                  })
                }
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-700">
                {t("fields.descriptionJa")}
              </label>
              <textarea
                className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                rows={3}
                value={item.descriptionJa}
                onChange={(event) =>
                  onChange((prev) => {
                    const next = [...prev];
                    next[index] = {
                      ...next[index],
                      descriptionJa: event.target.value,
                    };
                    return next;
                  })
                }
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-700">
                {t("fields.descriptionEn")}
              </label>
              <textarea
                className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                rows={3}
                value={item.descriptionEn}
                onChange={(event) =>
                  onChange((prev) => {
                    const next = [...prev];
                    next[index] = {
                      ...next[index],
                      descriptionEn: event.target.value,
                    };
                    return next;
                  })
                }
              />
            </div>
          </div>
          <div className="grid gap-4 md:grid-cols-2">
            <div>
              <label className="block text-sm font-medium text-slate-700">
                {t("homeSettings.fields.ctaJa", { defaultValue: "CTA (JA)" })}
              </label>
              <input
                className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                value={item.ctaJa}
                onChange={(event) =>
                  onChange((prev) => {
                    const next = [...prev];
                    next[index] = {
                      ...next[index],
                      ctaJa: event.target.value,
                    };
                    return next;
                  })
                }
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-700">
                {t("homeSettings.fields.ctaEn", { defaultValue: "CTA (EN)" })}
              </label>
              <input
                className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                value={item.ctaEn}
                onChange={(event) =>
                  onChange((prev) => {
                    const next = [...prev];
                    next[index] = {
                      ...next[index],
                      ctaEn: event.target.value,
                    };
                    return next;
                  })
                }
              />
            </div>
            <div className="md:col-span-2">
              <label className="block text-sm font-medium text-slate-700">
                {t("fields.url", { defaultValue: "Target URL" })}
              </label>
              <input
                className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                value={item.targetUrl}
                onChange={(event) =>
                  onChange((prev) => {
                    const next = [...prev];
                    next[index] = {
                      ...next[index],
                      targetUrl: event.target.value,
                    };
                    return next;
                  })
                }
              />
            </div>
          </div>
          <div className="flex justify-end">
            <button
              className="rounded-md border border-red-200 px-2 py-1 text-xs text-red-600"
              type="button"
              onClick={() =>
                onChange((prev) => prev.filter((_, i) => i !== index))
              }
              disabled={items.length === 1}
            >
              {t("actions.remove")}
            </button>
          </div>
        </div>
      ))}
      <button
        className="rounded-md border border-slate-200 px-3 py-1 text-sm text-slate-700"
        type="button"
        onClick={() => onChange((prev) => [...prev, createEmptyQuickLink()])}
      >
        {t("homeSettings.actions.addQuickLink", {
          defaultValue: "Add quick link",
        })}
      </button>
    </div>
  );
};

const HomeChipSourcesEditor = ({ items, onChange }: HomeChipSourcesEditorProps) => {
  const { t } = useTranslation();

  return (
    <div className="space-y-4">
      {items.map((item, index) => (
        <div
          key={item.id ?? index}
          className="space-y-3 rounded-lg border border-slate-200 p-4"
        >
          <div className="grid gap-4 md:grid-cols-2">
            <div>
              <label className="block text-sm font-medium text-slate-700">
                {t("homeSettings.fields.source", { defaultValue: "Source" })}
              </label>
              <select
                className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                value={item.source}
                onChange={(event) =>
                  onChange((prev) => {
                    const next = [...prev];
                    next[index] = {
                      ...next[index],
                      source: event.target.value as HomeChipSourceType,
                    };
                    return next;
                  })
                }
              >
                {homeChipSourceOptions.map((source) => (
                  <option key={source} value={source}>
                    {source}
                  </option>
                ))}
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-700">
                {t("fields.sortOrder", { defaultValue: "Sort order" })}
              </label>
              <input
                className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                value={item.sortOrder}
                onChange={(event) =>
                  onChange((prev) => {
                    const next = [...prev];
                    next[index] = {
                      ...next[index],
                      sortOrder: event.target.value,
                    };
                    return next;
                  })
                }
              />
            </div>
          </div>
          <div className="grid gap-4 md:grid-cols-2">
            <div>
              <label className="block text-sm font-medium text-slate-700">
                {t("fields.labelJa")}
              </label>
              <input
                className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                value={item.labelJa}
                onChange={(event) =>
                  onChange((prev) => {
                    const next = [...prev];
                    next[index] = {
                      ...next[index],
                      labelJa: event.target.value,
                    };
                    return next;
                  })
                }
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-700">
                {t("fields.labelEn")}
              </label>
              <input
                className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                value={item.labelEn}
                onChange={(event) =>
                  onChange((prev) => {
                    const next = [...prev];
                    next[index] = {
                      ...next[index],
                      labelEn: event.target.value,
                    };
                    return next;
                  })
                }
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-700">
                {t("homeSettings.fields.limit", { defaultValue: "Limit" })}
              </label>
              <input
                className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                value={item.limit}
                onChange={(event) =>
                  onChange((prev) => {
                    const next = [...prev];
                    next[index] = {
                      ...next[index],
                      limit: event.target.value,
                    };
                    return next;
                  })
                }
              />
            </div>
          </div>
          <div className="flex justify-end">
            <button
              className="rounded-md border border-red-200 px-2 py-1 text-xs text-red-600"
              type="button"
              onClick={() =>
                onChange((prev) => prev.filter((_, i) => i !== index))
              }
              disabled={items.length === 1}
            >
              {t("actions.remove")}
            </button>
          </div>
        </div>
      ))}
      <button
        className="rounded-md border border-slate-200 px-3 py-1 text-sm text-slate-700"
        type="button"
        onClick={() => onChange((prev) => [...prev, createEmptyChipSource()])}
      >
        {t("homeSettings.actions.addChipSource", {
          defaultValue: "Add chip source",
        })}
      </button>
    </div>
  );
};


const emptyProjectForm: ProjectFormState = {
  titleJa: "",
  titleEn: "",
  descriptionJa: "",
  descriptionEn: "",
  tech: [],
  linkUrl: "",
  year: String(currentYear),
  published: false,
  sortOrder: "",
};

const createEmptyResearchForm = (): ResearchFormState => ({
  slug: "",
  kind: "research",
  titleJa: "",
  titleEn: "",
  overviewJa: "",
  overviewEn: "",
  outcomeJa: "",
  outcomeEn: "",
  outlookJa: "",
  outlookEn: "",
  externalUrl: "",
  highlightImageUrl: "",
  imageAltJa: "",
  imageAltEn: "",
  publishedAt: "",
  isDraft: true,
  tags: [],
  links: [],
  assets: [],
  tech: [],
});

const emptyBlacklistForm: BlacklistFormState = {
  email: "",
  reason: "",
};

const createEmptyContactSettingsForm = (): ContactSettingsFormState => ({
  heroTitleJa: "",
  heroTitleEn: "",
  heroDescriptionJa: "",
  heroDescriptionEn: "",
  consentTextJa: "",
  consentTextEn: "",
  minimumLeadHours: "",
  recaptchaSiteKey: "",
  supportEmail: "",
  calendarTimezone: "",
  googleCalendarId: "",
  bookingWindowDays: "",
  topics: [],
});

const contactSettingsToForm = (
  settings: ContactFormSettings | null,
): ContactSettingsFormState => {
  if (!settings) {
    return createEmptyContactSettingsForm();
  }
  return {
    heroTitleJa: settings.heroTitle.ja ?? "",
    heroTitleEn: settings.heroTitle.en ?? "",
    heroDescriptionJa: settings.heroDescription.ja ?? "",
    heroDescriptionEn: settings.heroDescription.en ?? "",
    consentTextJa: settings.consentText.ja ?? "",
    consentTextEn: settings.consentText.en ?? "",
    minimumLeadHours: String(settings.minimumLeadHours),
    recaptchaSiteKey: settings.recaptchaSiteKey ?? "",
    supportEmail: settings.supportEmail ?? "",
    calendarTimezone: settings.calendarTimezone ?? "",
    googleCalendarId: settings.googleCalendarId ?? "",
    bookingWindowDays: String(settings.bookingWindowDays),
    topics: settings.topics.map((topic) => ({
      id: topic.id,
      labelJa: topic.label.ja ?? "",
      labelEn: topic.label.en ?? "",
      descriptionJa: topic.description.ja ?? "",
      descriptionEn: topic.description.en ?? "",
    })),
  };
};

const trimValue = (value?: string | null): string => value?.trim() ?? "";

const normalizeLocalizedField = (
  value: { ja?: string; en?: string },
): NormalizedLocalizedText => ({
  ja: trimValue(value.ja),
  en: trimValue(value.en),
});

const parseNumber = (value: string, fallback: number): number => {
  const parsed = Number.parseInt(value, 10);
  if (Number.isFinite(parsed)) {
    return parsed;
  }
  return fallback;
};

const normalizeContactSettingsForm = (
  form: ContactSettingsFormState,
): ContactSettingsNormalized => ({
  heroTitle: normalizeLocalizedField({
    ja: form.heroTitleJa,
    en: form.heroTitleEn,
  }),
  heroDescription: normalizeLocalizedField({
    ja: form.heroDescriptionJa,
    en: form.heroDescriptionEn,
  }),
  consentText: normalizeLocalizedField({
    ja: form.consentTextJa,
    en: form.consentTextEn,
  }),
  topics: form.topics.map((topic) => ({
    id: trimValue(topic.id),
    label: normalizeLocalizedField({
      ja: topic.labelJa,
      en: topic.labelEn,
    }),
    description: normalizeLocalizedField({
      ja: topic.descriptionJa,
      en: topic.descriptionEn,
    }),
  })),
  minimumLeadHours: parseNumber(form.minimumLeadHours, 0),
  recaptchaSiteKey: trimValue(form.recaptchaSiteKey),
  supportEmail: trimValue(form.supportEmail),
  calendarTimezone: trimValue(form.calendarTimezone),
  googleCalendarId: trimValue(form.googleCalendarId),
  bookingWindowDays: parseNumber(form.bookingWindowDays, 0),
});

const normalizeProfileFormState = (
  form: ProfileFormState,
): ProfilePayload => ({
  displayName: trimValue(form.displayName),
  headline: normalizeLocalizedField({
    ja: form.headlineJa,
    en: form.headlineEn,
  }),
  summary: normalizeLocalizedField({
    ja: form.summaryJa,
    en: form.summaryEn,
  }),
  avatarUrl: trimValue(form.avatarUrl),
  location: normalizeLocalizedField({
    ja: form.locationJa,
    en: form.locationEn,
  }),
  theme: {
    mode: form.themeMode,
    accentColor: trimValue(form.themeAccentColor) || undefined,
  },
  lab: {
    name: normalizeLocalizedField({
      ja: form.labNameJa,
      en: form.labNameEn,
    }),
    advisor: normalizeLocalizedField({
      ja: form.labAdvisorJa,
      en: form.labAdvisorEn,
    }),
    room: normalizeLocalizedField({
      ja: form.labRoomJa,
      en: form.labRoomEn,
    }),
    url: trimValue(form.labUrl) || undefined,
  },
  affiliations: form.affiliations
    .filter((item) => trimValue(item.name))
    .map((item) => ({
      id: item.id,
      name: trimValue(item.name),
      url: trimValue(item.url) || undefined,
      description: normalizeLocalizedField({
        ja: item.descriptionJa,
        en: item.descriptionEn,
      }),
      startedAt: item.startedAt,
      sortOrder: normalizeSortOrder(item.sortOrder),
    })),
  communities: form.communities
    .filter((item) => trimValue(item.name))
    .map((item) => ({
      id: item.id,
      name: trimValue(item.name),
      url: trimValue(item.url) || undefined,
      description: normalizeLocalizedField({
        ja: item.descriptionJa,
        en: item.descriptionEn,
      }),
      startedAt: item.startedAt,
      sortOrder: normalizeSortOrder(item.sortOrder),
    })),
  workHistory: form.workHistory
    .filter(
      (item) =>
        trimValue(item.organizationJa) ||
        trimValue(item.organizationEn) ||
        trimValue(item.roleJa) ||
        trimValue(item.roleEn),
    )
    .map((item) => ({
      id: item.id,
      organization: normalizeLocalizedField({
        ja: item.organizationJa,
        en: item.organizationEn,
      }),
      role: normalizeLocalizedField({
        ja: item.roleJa,
        en: item.roleEn,
      }),
      summary: normalizeLocalizedField({
        ja: item.summaryJa,
        en: item.summaryEn,
      }),
      startedAt: item.startedAt,
      endedAt: trimValue(item.endedAt) || undefined,
      externalUrl: trimValue(item.externalUrl) || undefined,
      sortOrder: normalizeSortOrder(item.sortOrder),
    })),
  socialLinks: form.socialLinks.map((item) => ({
    id: item.id,
    provider: item.provider,
    label: normalizeLocalizedField({
      ja: item.labelJa,
      en: item.labelEn,
    }),
    url: trimValue(item.url),
    isFooter: item.isFooter,
    sortOrder: normalizeSortOrder(item.sortOrder),
  })),
});

const normalizeContactSettingsOriginal = (
  settings: ContactFormSettings | null,
): ContactSettingsNormalized | null => {
  if (!settings) {
    return null;
  }
  return {
    heroTitle: normalizeLocalizedField(settings.heroTitle),
    heroDescription: normalizeLocalizedField(settings.heroDescription),
    consentText: normalizeLocalizedField(settings.consentText),
    topics: settings.topics.map((topic) => ({
      id: topic.id,
      label: normalizeLocalizedField(topic.label),
      description: normalizeLocalizedField(topic.description),
    })),
    minimumLeadHours: settings.minimumLeadHours,
    recaptchaSiteKey: trimValue(settings.recaptchaSiteKey),
    supportEmail: trimValue(settings.supportEmail),
    calendarTimezone: trimValue(settings.calendarTimezone),
    googleCalendarId: trimValue(settings.googleCalendarId),
    bookingWindowDays: settings.bookingWindowDays,
  };
};

const buildContactSettingsDiff = (
  original: ContactFormSettings | null,
  form: ContactSettingsFormState,
): ContactSettingsDiffSummary => {
  const normalizedDraft = normalizeContactSettingsForm(form);
  const normalizedOriginal =
    normalizeContactSettingsOriginal(original) ??
    normalizeContactSettingsForm(createEmptyContactSettingsForm());

  const addLocalizedDiff = (
    entries: ContactSettingsDiffEntry[],
    keyPrefix: string,
    labelPrefix: string,
    originalValue: NormalizedLocalizedText,
    updatedValue: NormalizedLocalizedText,
  ) => {
    if (originalValue.ja !== updatedValue.ja) {
      entries.push({
        key: `${keyPrefix}.ja`,
        labelKey: `${labelPrefix}Ja`,
        original: originalValue.ja,
        updated: updatedValue.ja,
      });
    }
    if (originalValue.en !== updatedValue.en) {
      entries.push({
        key: `${keyPrefix}.en`,
        labelKey: `${labelPrefix}En`,
        original: originalValue.en,
        updated: updatedValue.en,
      });
    }
  };

  const fieldDiffs: ContactSettingsDiffEntry[] = [];
  addLocalizedDiff(
    fieldDiffs,
    "heroTitle",
    "contactSettings.fields.heroTitle",
    normalizedOriginal.heroTitle,
    normalizedDraft.heroTitle,
  );
  addLocalizedDiff(
    fieldDiffs,
    "heroDescription",
    "contactSettings.fields.heroDescription",
    normalizedOriginal.heroDescription,
    normalizedDraft.heroDescription,
  );
  addLocalizedDiff(
    fieldDiffs,
    "consentText",
    "contactSettings.fields.consentText",
    normalizedOriginal.consentText,
    normalizedDraft.consentText,
  );

  const addSimpleDiff = (
    key: string,
    labelKey: string,
    originalValue: string | number,
    updatedValue: string | number,
  ) => {
    if (`${originalValue}` !== `${updatedValue}`) {
      fieldDiffs.push({
        key,
        labelKey,
        original: String(originalValue),
        updated: String(updatedValue),
      });
    }
  };

  addSimpleDiff(
    "minimumLeadHours",
    "contactSettings.fields.minimumLeadHours",
    normalizedOriginal.minimumLeadHours,
    normalizedDraft.minimumLeadHours,
  );
  addSimpleDiff(
    "bookingWindowDays",
    "contactSettings.fields.bookingWindowDays",
    normalizedOriginal.bookingWindowDays,
    normalizedDraft.bookingWindowDays,
  );
  addSimpleDiff(
    "supportEmail",
    "contactSettings.fields.supportEmail",
    normalizedOriginal.supportEmail,
    normalizedDraft.supportEmail,
  );
  addSimpleDiff(
    "recaptchaSiteKey",
    "contactSettings.fields.recaptchaSiteKey",
    normalizedOriginal.recaptchaSiteKey,
    normalizedDraft.recaptchaSiteKey,
  );
  addSimpleDiff(
    "calendarTimezone",
    "contactSettings.fields.calendarTimezone",
    normalizedOriginal.calendarTimezone,
    normalizedDraft.calendarTimezone,
  );
  addSimpleDiff(
    "googleCalendarId",
    "contactSettings.fields.googleCalendarId",
    normalizedOriginal.googleCalendarId,
    normalizedDraft.googleCalendarId,
  );

  const topicDiffs: ContactTopicDiffEntry[] = [];
  const originalTopics = new Map<string, ContactSettingsNormalized["topics"][number]>();
  normalizedOriginal.topics.forEach((topic) => {
    originalTopics.set(topic.id, topic);
  });

  normalizedDraft.topics.forEach((topic) => {
    const originalTopic = originalTopics.get(topic.id);
    if (!originalTopic) {
      topicDiffs.push({
        id: topic.id,
        change: "added",
        updated: {
          label: topic.label,
          description: topic.description,
        },
      });
      return;
    }

    const labelChanged =
      originalTopic.label.ja !== topic.label.ja ||
      originalTopic.label.en !== topic.label.en;
    const descriptionChanged =
      originalTopic.description.ja !== topic.description.ja ||
      originalTopic.description.en !== topic.description.en;
    if (labelChanged || descriptionChanged) {
      topicDiffs.push({
        id: topic.id,
        change: "updated",
        original: {
          label: originalTopic.label,
          description: originalTopic.description,
        },
        updated: {
          label: topic.label,
          description: topic.description,
        },
      });
    }
    originalTopics.delete(topic.id);
  });

  originalTopics.forEach((topic) => {
    topicDiffs.push({
      id: topic.id,
      change: "removed",
      original: {
        label: topic.label,
        description: topic.description,
      },
    });
  });

  return {
    fields: fieldDiffs,
    topics: topicDiffs,
  };
};

const homeSettingsToForm = (
  settings: HomePageConfig | null,
  fallbackProfileId?: number,
): HomeSettingsFormState => {
  if (!settings) {
    return createEmptyHomeSettingsForm(fallbackProfileId);
  }
  return {
    id: settings.id,
    profileId: settings.profileId ?? fallbackProfileId ?? null,
    heroSubtitleJa: settings.heroSubtitle.ja ?? "",
    heroSubtitleEn: settings.heroSubtitle.en ?? "",
    quickLinks:
      settings.quickLinks.length > 0
        ? settings.quickLinks.map((link) => ({
            id: link.id,
            section: link.section,
            labelJa: link.label.ja ?? "",
            labelEn: link.label.en ?? "",
            descriptionJa: link.description.ja ?? "",
            descriptionEn: link.description.en ?? "",
            ctaJa: link.cta.ja ?? "",
            ctaEn: link.cta.en ?? "",
            targetUrl: link.targetUrl,
            sortOrder: String(link.sortOrder ?? 0),
          }))
        : [createEmptyQuickLink()],
    chipSources:
      settings.chipSources.length > 0
        ? settings.chipSources.map((chip) => ({
            id: chip.id,
            source: chip.source as HomeChipSourceType,
            labelJa: chip.label.ja ?? "",
            labelEn: chip.label.en ?? "",
            limit: String(chip.limit ?? 0),
            sortOrder: String(chip.sortOrder ?? 0),
          }))
        : [createEmptyChipSource()],
    updatedAt: settings.updatedAt,
  };
};

const normalizeHomeSettingsForm = (
  form: HomeSettingsFormState,
): HomeSettingsPayload => {
  if (form.id == null || form.profileId == null || !form.updatedAt) {
    throw new Error("home settings missing required identifiers");
  }

  return {
    id: form.id,
    profileId: form.profileId,
    heroSubtitle: normalizeLocalizedField({
      ja: form.heroSubtitleJa,
      en: form.heroSubtitleEn,
    }),
    quickLinks: form.quickLinks
      .filter((link) => trimValue(link.labelJa) || trimValue(link.labelEn))
      .map((link) => ({
        id: link.id,
        section: link.section,
        label: normalizeLocalizedField({
          ja: link.labelJa,
          en: link.labelEn,
        }),
        description: normalizeLocalizedField({
          ja: link.descriptionJa,
          en: link.descriptionEn,
        }),
        cta: normalizeLocalizedField({
          ja: link.ctaJa,
          en: link.ctaEn,
        }),
        targetUrl: trimValue(link.targetUrl),
        sortOrder: normalizeSortOrder(link.sortOrder),
      })),
    chipSources: form.chipSources.map((chip) => ({
      id: chip.id,
      source: chip.source,
      label: normalizeLocalizedField({
        ja: chip.labelJa,
        en: chip.labelEn,
      }),
      limit: parseNumber(chip.limit, 0),
      sortOrder: normalizeSortOrder(chip.sortOrder),
    })),
    updatedAt: form.updatedAt,
  };
};

const validateHomeSettingsForm = (
  form: HomeSettingsFormState,
): string | null => {
  if (form.id == null || form.profileId == null || !form.updatedAt) {
    return "homeSettings.validation.missingRecord";
  }

  const heroSubtitleJa = trimValue(form.heroSubtitleJa);
  const heroSubtitleEn = trimValue(form.heroSubtitleEn);
  if (!heroSubtitleJa && !heroSubtitleEn) {
    return "homeSettings.validation.heroSubtitle";
  }

  if (form.quickLinks.length === 0) {
    return "homeSettings.validation.quickLinks";
  }
  for (const link of form.quickLinks) {
    const hasLabel = trimValue(link.labelJa) || trimValue(link.labelEn);
    if (!hasLabel) {
      return "homeSettings.validation.quickLinkLabel";
    }
    if (!trimValue(link.targetUrl)) {
      return "homeSettings.validation.quickLinkUrl";
    }
  }

  if (form.chipSources.length === 0) {
    return "homeSettings.validation.chipSources";
  }

  return null;
};

const isValidEmail = (value: string): boolean =>
  /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value);

const validateContactSettingsForm = (
  form: ContactSettingsFormState,
): string | null => {
  const heroTitleJa = trimValue(form.heroTitleJa);
  const heroTitleEn = trimValue(form.heroTitleEn);
  if (!heroTitleJa && !heroTitleEn) {
    return "contactSettings.validation.heroTitle";
  }

  const consentJa = trimValue(form.consentTextJa);
  const consentEn = trimValue(form.consentTextEn);
  if (!consentJa && !consentEn) {
    return "contactSettings.validation.consentText";
  }

  const minimumLeadHours = parseNumber(form.minimumLeadHours, Number.NaN);
  if (!Number.isFinite(minimumLeadHours) || minimumLeadHours < 0) {
    return "contactSettings.validation.minimumLeadHours";
  }

  const bookingWindowDays = parseNumber(form.bookingWindowDays, Number.NaN);
  if (!Number.isFinite(bookingWindowDays) || bookingWindowDays < 1) {
    return "contactSettings.validation.bookingWindowDays";
  }

  if (!trimValue(form.supportEmail)) {
    return "contactSettings.validation.supportEmail";
  }
  if (!isValidEmail(trimValue(form.supportEmail))) {
    return "contactSettings.validation.supportEmailFormat";
  }

  if (!trimValue(form.calendarTimezone)) {
    return "contactSettings.validation.calendarTimezone";
  }

  if (form.topics.length === 0) {
    return "contactSettings.validation.topicsRequired";
  }

  const seenIds = new Set<string>();
  for (const topic of form.topics) {
    const topicId = trimValue(topic.id);
    if (!topicId) {
      return "contactSettings.validation.topicId";
    }
    if (seenIds.has(topicId)) {
      return "contactSettings.validation.topicIdUnique";
    }
    seenIds.add(topicId);
    const labelJa = trimValue(topic.labelJa);
    const labelEn = trimValue(topic.labelEn);
    if (!labelJa && !labelEn) {
      return "contactSettings.validation.topicLabel";
    }
  }

  return null;
};

const buildContactSettingsPayload = (
  form: ContactSettingsFormState,
  original: ContactFormSettings,
) => {
  const normalized = normalizeContactSettingsForm(form);
  const payload: Parameters<
    typeof adminApi.updateContactSettings
  >[0] = {
    id: original.id,
    heroTitle: normalized.heroTitle,
    heroDescription: normalized.heroDescription,
    topics: normalized.topics,
    consentText: normalized.consentText,
    minimumLeadHours: normalized.minimumLeadHours,
    recaptchaSiteKey: normalized.recaptchaSiteKey,
    supportEmail: normalized.supportEmail,
    calendarTimezone: normalized.calendarTimezone,
    googleCalendarId: normalized.googleCalendarId,
    bookingWindowDays: normalized.bookingWindowDays,
    updatedAt: original.updatedAt,
  };
  return payload;
};

function profileToForm(profile: AdminProfile | null): ProfileFormState {
  if (!profile) {
    return createEmptyProfileForm();
  }

  const affiliations =
    profile.affiliations.length > 0
      ? profile.affiliations.map((item) => ({
          id: item.id,
          name: item.name,
          url: item.url ?? "",
          descriptionJa: item.description.ja ?? "",
          descriptionEn: item.description.en ?? "",
          startedAt: item.startedAt,
          sortOrder: String(item.sortOrder ?? 0),
        }))
      : [createEmptyAffiliation()];

  const communities =
    profile.communities.length > 0
      ? profile.communities.map((item) => ({
          id: item.id,
          name: item.name,
          url: item.url ?? "",
          descriptionJa: item.description.ja ?? "",
          descriptionEn: item.description.en ?? "",
          startedAt: item.startedAt,
          sortOrder: String(item.sortOrder ?? 0),
        }))
      : [createEmptyAffiliation()];

  const workHistory =
    profile.workHistory.length > 0
      ? profile.workHistory.map((item) => ({
          id: item.id,
          organizationJa: item.organization.ja ?? "",
          organizationEn: item.organization.en ?? "",
          roleJa: item.role.ja ?? "",
          roleEn: item.role.en ?? "",
          summaryJa: item.summary.ja ?? "",
          summaryEn: item.summary.en ?? "",
          startedAt: item.startedAt,
          endedAt: item.endedAt ?? "",
          externalUrl: item.externalUrl ?? "",
          sortOrder: String(item.sortOrder ?? 0),
        }))
      : [createEmptyWorkHistory()];

  const socialLinksMap = new Map<SocialProvider, SocialLinkForm>();
  profile.socialLinks.forEach((link) => {
    const provider = link.provider as SocialProvider;
    socialLinksMap.set(provider, {
      id: link.id,
      provider,
      labelJa: link.label.ja ?? "",
      labelEn: link.label.en ?? "",
      url: link.url,
      isFooter: link.isFooter,
      sortOrder: String(link.sortOrder ?? 0),
    });
  });

  const requiredProviders: SocialProvider[] = [
    "github",
    "zenn",
    "linkedin",
  ];
  requiredProviders.forEach((provider) => {
    if (!socialLinksMap.has(provider)) {
      socialLinksMap.set(provider, createSocialLink(provider));
    }
  });

  const socialLinks = Array.from(socialLinksMap.values());

  return {
    displayName: profile.displayName,
    headlineJa: profile.headline.ja ?? "",
    headlineEn: profile.headline.en ?? "",
    summaryJa: profile.summary.ja ?? "",
    summaryEn: profile.summary.en ?? "",
    avatarUrl: profile.avatarUrl ?? "",
    locationJa: profile.location.ja ?? "",
    locationEn: profile.location.en ?? "",
    themeMode: profile.theme.mode,
    themeAccentColor: profile.theme.accentColor ?? "",
    labNameJa: profile.lab.name.ja ?? "",
    labNameEn: profile.lab.name.en ?? "",
    labAdvisorJa: profile.lab.advisor.ja ?? "",
    labAdvisorEn: profile.lab.advisor.en ?? "",
    labRoomJa: profile.lab.room.ja ?? "",
    labRoomEn: profile.lab.room.en ?? "",
    labUrl: profile.lab.url ?? "",
    affiliations,
    communities,
    workHistory,
    socialLinks,
  };
}

function buildContactEditMap(
  contacts: ContactMessage[],
): Record<string, ContactEditState> {
  return contacts.reduce<Record<string, ContactEditState>>((acc, contact) => {
    acc[contact.id] = {
      topic: contact.topic,
      message: contact.message,
      status: contact.status,
      adminNote: contact.adminNote,
    };
    return acc;
  }, {});
}

const toDateTimeLocal = (iso: string): string => {
  if (!iso) {
    return "";
  }
  const date = new Date(iso);
  if (Number.isNaN(date.getTime())) {
    return "";
  }
  const offset = date.getTimezoneOffset();
  const local = new Date(date.getTime() - offset * 60 * 1000);
  return local.toISOString().slice(0, 16);
};

const toISOStringWithFallback = (value: string): string => {
  if (!value) {
    return new Date().toISOString();
  }
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return new Date().toISOString();
  }
  return date.toISOString();
};

const normalizeSortOrder = (value: string): number =>
  Number.parseInt(value, 10) || 0;

const projectToForm = (project: AdminProject): ProjectFormState => ({
  titleJa: project.title.ja ?? "",
  titleEn: project.title.en ?? "",
  descriptionJa: project.description.ja ?? "",
  descriptionEn: project.description.en ?? "",
  tech: project.tech.map((membership) => ({
    membershipId: membership.membershipId,
    techId: membership.tech?.id ? String(membership.tech.id) : "",
    context: membership.context,
    note: membership.note,
    sortOrder: String(membership.sortOrder ?? 0),
  })),
  linkUrl: project.linkUrl,
  year: String(project.year),
  published: project.published,
  sortOrder: project.sortOrder != null ? String(project.sortOrder) : "",
});

const projectFormToPayload = (form: ProjectFormState) => {
  const techMembers: {
    membershipId?: number;
    techId: number;
    context: TechContext;
    note: string;
    sortOrder: number;
  }[] = [];

  form.tech.forEach((membership) => {
    const parsed = Number.parseInt(membership.techId, 10);
    if (!Number.isFinite(parsed) || parsed <= 0) {
      return;
    }
    techMembers.push({
      membershipId: membership.membershipId,
      techId: parsed,
      context: membership.context,
      note: membership.note.trim(),
      sortOrder: normalizeSortOrder(membership.sortOrder),
    });
  });

  return {
    title: { ja: form.titleJa.trim(), en: form.titleEn.trim() },
    description: {
      ja: form.descriptionJa.trim(),
      en: form.descriptionEn.trim(),
    },
    tech: techMembers,
    linkUrl: form.linkUrl.trim(),
    year: Number.parseInt(form.year, 10) || currentYear,
    published: form.published,
    sortOrder: form.sortOrder === "" ? null : normalizeSortOrder(form.sortOrder),
  };
};

const projectToPayload = (project: AdminProject) =>
  projectFormToPayload(projectToForm(project));

const researchToForm = (item: AdminResearch): ResearchFormState => ({
  slug: item.slug,
  kind: item.kind,
  titleJa: item.title.ja ?? "",
  titleEn: item.title.en ?? "",
  overviewJa: item.overview.ja ?? "",
  overviewEn: item.overview.en ?? "",
  outcomeJa: item.outcome.ja ?? "",
  outcomeEn: item.outcome.en ?? "",
  outlookJa: item.outlook.ja ?? "",
  outlookEn: item.outlook.en ?? "",
  externalUrl: item.externalUrl,
  highlightImageUrl: item.highlightImageUrl,
  imageAltJa: item.imageAlt.ja ?? "",
  imageAltEn: item.imageAlt.en ?? "",
  publishedAt: toDateTimeLocal(item.publishedAt),
  isDraft: item.isDraft,
  tags: item.tags.map((tag) => ({
    id: tag.id,
    value: tag.value,
    sortOrder: String(tag.sortOrder),
  })),
  links: item.links.map((link) => ({
    id: link.id,
    type: link.type,
    labelJa: link.label.ja ?? "",
    labelEn: link.label.en ?? "",
    url: link.url,
    sortOrder: String(link.sortOrder),
  })),
  assets: item.assets.map((asset) => ({
    id: asset.id,
    url: asset.url,
    captionJa: asset.caption.ja ?? "",
    captionEn: asset.caption.en ?? "",
    sortOrder: String(asset.sortOrder),
  })),
  tech: item.tech.map((membership) => ({
    membershipId: membership.membershipId,
    techId: membership.tech?.id ? String(membership.tech.id) : "",
    context: membership.context,
    note: membership.note,
    sortOrder: String(membership.sortOrder),
  })),
});

const researchFormToPayload = (form: ResearchFormState) => {
  const techMembers: {
    membershipId?: number;
    techId: number;
    context: TechContext;
    note: string;
    sortOrder: number;
  }[] = [];

  form.tech.forEach((membership) => {
    const parsed = Number.parseInt(membership.techId, 10);
    if (!Number.isFinite(parsed) || parsed <= 0) {
      return;
    }
    techMembers.push({
      membershipId: membership.membershipId,
      techId: parsed,
      context: membership.context,
      note: membership.note.trim(),
      sortOrder: normalizeSortOrder(membership.sortOrder),
    });
  });

  return {
    slug: form.slug.trim(),
    kind: form.kind,
    title: { ja: form.titleJa.trim(), en: form.titleEn.trim() },
    overview: { ja: form.overviewJa.trim(), en: form.overviewEn.trim() },
    outcome: { ja: form.outcomeJa.trim(), en: form.outcomeEn.trim() },
    outlook: { ja: form.outlookJa.trim(), en: form.outlookEn.trim() },
    externalUrl: form.externalUrl.trim(),
    highlightImageUrl: form.highlightImageUrl.trim(),
    imageAlt: { ja: form.imageAltJa.trim(), en: form.imageAltEn.trim() },
    publishedAt: toISOStringWithFallback(form.publishedAt),
    isDraft: form.isDraft,
    tags: form.tags
      .filter((tag) => tag.value.trim())
      .map((tag) => ({
        id: tag.id,
        value: tag.value.trim(),
        sortOrder: normalizeSortOrder(tag.sortOrder),
      })),
    links: form.links
      .filter((link) => link.url.trim())
      .map((link) => ({
        id: link.id,
        type: link.type,
        label: { ja: link.labelJa.trim(), en: link.labelEn.trim() },
        url: link.url.trim(),
        sortOrder: normalizeSortOrder(link.sortOrder),
      })),
    assets: form.assets
      .filter((asset) => asset.url.trim())
      .map((asset) => ({
        id: asset.id,
        url: asset.url.trim(),
        caption: { ja: asset.captionJa.trim(), en: asset.captionEn.trim() },
        sortOrder: normalizeSortOrder(asset.sortOrder),
      })),
    tech: techMembers,
  };
};

const researchToPayload = (item: AdminResearch) =>
  researchFormToPayload(researchToForm(item));

function App() {
  const { t } = useTranslation();
  const { setSession, clearSession } = useAuthSession();

  const [authState, setAuthState] = useState<AuthState>("checking");
  const [status, setStatus] = useState("unknown");
  const [summary, setSummary] = useState<AdminSummary | null>(null);
  const [projects, setProjects] = useState<AdminProject[]>([]);
  const [techCatalog, setTechCatalog] = useState<TechCatalogEntry[]>([]);
  const [research, setResearch] = useState<AdminResearch[]>([]);
  const [contacts, setContacts] = useState<ContactMessage[]>([]);
  const [blacklist, setBlacklist] = useState<BlacklistEntry[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [profileForm, setProfileForm] = useState<ProfileFormState>(
    createEmptyProfileForm(),
  );
  const [projectForm, setProjectForm] = useState<ProjectFormState>({
    ...emptyProjectForm,
  });
  const [projectTechSearch, setProjectTechSearch] = useState("");
  const [showProjectPreview, setShowProjectPreview] = useState(false);
  const [researchForm, setResearchForm] = useState<ResearchFormState>(
    createEmptyResearchForm(),
  );
  const [blacklistForm, setBlacklistForm] = useState<BlacklistFormState>({
    ...emptyBlacklistForm,
  });
  const [contactEdits, setContactEdits] = useState<
    Record<string, ContactEditState>
  >({});
  const [contactSettings, setContactSettings] =
    useState<ContactFormSettings | null>(null);
  const [contactSettingsForm, setContactSettingsForm] =
    useState<ContactSettingsFormState>(createEmptyContactSettingsForm());
  const [contactSettingsError, setContactSettingsError] = useState<string | null>(
    null,
  );
  const [contactSettingsNotice, setContactSettingsNotice] = useState<
    string | null
  >(null);
  const [contactSettingsSaving, setContactSettingsSaving] = useState(false);
  const [showContactSettingsDiff, setShowContactSettingsDiff] = useState(false);
  const [contactPreviewLocale, setContactPreviewLocale] = useState<"ja" | "en">(
    "ja",
  );
  const [homeSettingsForm, setHomeSettingsForm] =
    useState<HomeSettingsFormState>(createEmptyHomeSettingsForm());
  const [homeSettingsSaving, setHomeSettingsSaving] = useState(false);
  const [homeSettingsNotice, setHomeSettingsNotice] = useState<string | null>(
    null,
  );
  const [homeSettingsError, setHomeSettingsError] = useState<string | null>(null);

  const [editingProjectId, setEditingProjectId] = useState<number | null>(null);
  const [editingResearchId, setEditingResearchId] = useState<number | null>(
    null,
  );
  const [editingBlacklistId, setEditingBlacklistId] =
    useState<number | null>(null);

  const selectedProjectTechIds = useMemo(() => {
    const ids = new Set<number>();
    projectForm.tech.forEach((membership) => {
      const parsed = Number.parseInt(membership.techId, 10);
      if (Number.isFinite(parsed) && parsed > 0) {
        ids.add(parsed);
      }
    });
    return ids;
  }, [projectForm.tech]);

  const filteredProjectTech = useMemo(() => {
    const query = projectTechSearch.trim().toLowerCase();
    const matches = techCatalog
      .filter((entry) => entry.active && !selectedProjectTechIds.has(entry.id))
      .filter((entry) => {
        if (!query) {
          return true;
        }
        const haystack = `${entry.displayName} ${entry.slug} ${
          entry.category ?? ""
        }`.toLowerCase();
        return haystack.includes(query);
      })
      .sort((a, b) => a.displayName.localeCompare(b.displayName));
    return matches.slice(0, 10);
  }, [projectTechSearch, techCatalog, selectedProjectTechIds]);

  const contactSettingsDiffSummary = useMemo(
    () => buildContactSettingsDiff(contactSettings, contactSettingsForm),
    [contactSettings, contactSettingsForm],
  );

  const hasContactSettingsChanges =
    contactSettingsDiffSummary.fields.length > 0 ||
    contactSettingsDiffSummary.topics.length > 0;

  const selectLocalizedValue = useCallback(
    (ja: string, en: string) =>
      contactPreviewLocale === "ja"
        ? trimValue(ja) || trimValue(en)
        : trimValue(en) || trimValue(ja),
    [contactPreviewLocale],
  );

  const previewHeroTitle = selectLocalizedValue(
    contactSettingsForm.heroTitleJa,
    contactSettingsForm.heroTitleEn,
  );
  const previewHeroDescription = selectLocalizedValue(
    contactSettingsForm.heroDescriptionJa,
    contactSettingsForm.heroDescriptionEn,
  );
  const previewConsent = selectLocalizedValue(
    contactSettingsForm.consentTextJa,
    contactSettingsForm.consentTextEn,
  );
  const previewTopics = useMemo(
    () =>
      contactSettingsForm.topics.map((topic) => ({
        id: trimValue(topic.id),
        label: selectLocalizedValue(topic.labelJa, topic.labelEn),
        description: selectLocalizedValue(
          topic.descriptionJa,
          topic.descriptionEn,
        ),
      })),
    [contactSettingsForm.topics, selectLocalizedValue],
  );
  const previewSupportEmail = trimValue(contactSettingsForm.supportEmail);
  const previewCalendarTimezone = trimValue(
    contactSettingsForm.calendarTimezone,
  );
  const previewRecaptchaSiteKey = trimValue(
    contactSettingsForm.recaptchaSiteKey,
  );
  const previewMinimumLeadHours = parseNumber(
    contactSettingsForm.minimumLeadHours,
    0,
  );
  const previewBookingWindowDays = parseNumber(
    contactSettingsForm.bookingWindowDays,
    0,
  );

  const handleUnauthorized = useCallback(() => {
    clearSession();
    setAuthState("unauthorized");
    setLoading(false);
    setSummary(null);
    setProjects([]);
    setResearch([]);
    setContacts([]);
    setBlacklist([]);
    setContactEdits({});
    setContactSettings(null);
    setContactSettingsForm(createEmptyContactSettingsForm());
    setContactSettingsError(null);
    setContactSettingsNotice(null);
    setContactSettingsSaving(false);
    setShowContactSettingsDiff(false);
    setContactPreviewLocale("ja");
    setProfileForm(createEmptyProfileForm());
    setProjectForm({ ...emptyProjectForm });
    setProjectTechSearch("");
    setShowProjectPreview(false);
    setResearchForm(createEmptyResearchForm());
    setBlacklistForm({ ...emptyBlacklistForm });
    setEditingProjectId(null);
    setEditingResearchId(null);
    setEditingBlacklistId(null);
    setHomeSettingsForm(createEmptyHomeSettingsForm());
    setHomeSettingsError(null);
    setHomeSettingsNotice(null);
    setHomeSettingsSaving(false);
    setError(null);
  }, [clearSession]);

  const handleContactSettingsFieldChange = useCallback(
    (
      field: Exclude<keyof ContactSettingsFormState, "topics">,
      value: string,
    ) => {
      setContactSettingsForm((prev) => ({
        ...prev,
        [field]: value,
      }));
      setContactSettingsError(null);
      setContactSettingsNotice(null);
    },
    [],
  );

  const handleContactTopicChange = useCallback(
    (index: number, field: keyof ContactTopicForm, value: string) => {
      setContactSettingsForm((prev) => {
        const nextTopics = [...prev.topics];
        nextTopics[index] = {
          ...nextTopics[index],
          [field]: value,
        };
        return {
          ...prev,
          topics: nextTopics,
        };
      });
      setContactSettingsError(null);
      setContactSettingsNotice(null);
    },
    [],
  );

  const handleAddContactTopic = useCallback(() => {
    setContactSettingsForm((prev) => ({
      ...prev,
      topics: [
        ...prev.topics,
        {
          id: `topic-${Math.random().toString(36).slice(2, 8)}`,
          labelJa: "",
          labelEn: "",
          descriptionJa: "",
          descriptionEn: "",
        },
      ],
    }));
    setContactSettingsError(null);
    setContactSettingsNotice(null);
  }, []);

  const handleHomeSettingsFieldChange = useCallback(
    (field: "heroSubtitleJa" | "heroSubtitleEn", value: string) => {
      setHomeSettingsForm((prev) => ({
        ...prev,
        [field]: value,
      }));
      setHomeSettingsError(null);
      setHomeSettingsNotice(null);
    },
    [],
  );

  const handleHomeQuickLinksChange = useCallback(
    (updater: ListUpdater<HomeQuickLinkForm>) => {
      setHomeSettingsForm((prev) => ({
        ...prev,
        quickLinks: updater(prev.quickLinks),
      }));
      setHomeSettingsError(null);
      setHomeSettingsNotice(null);
    },
    [],
  );

  const handleHomeChipSourcesChange = useCallback(
    (updater: ListUpdater<HomeChipSourceForm>) => {
      setHomeSettingsForm((prev) => ({
        ...prev,
        chipSources: updater(prev.chipSources),
      }));
      setHomeSettingsError(null);
      setHomeSettingsNotice(null);
    },
    [],
  );

  const handleRemoveContactTopic = useCallback((index: number) => {
    setContactSettingsForm((prev) => {
      const nextTopics = prev.topics.filter((_, i) => i !== index);
      return {
        ...prev,
        topics: nextTopics,
      };
    });
    setContactSettingsError(null);
    setContactSettingsNotice(null);
  }, []);

  const handleResetContactSettings = useCallback(() => {
    setContactSettingsForm(contactSettingsToForm(contactSettings));
    setContactSettingsError(null);
    setContactSettingsNotice(null);
    setShowContactSettingsDiff(false);
  }, [contactSettings]);

  const handleSaveContactSettings = useCallback(async () => {
    if (!contactSettings) {
      setContactSettingsError("contactSettings.validation.missingRecord");
      return;
    }
    const validationError = validateContactSettingsForm(contactSettingsForm);
    if (validationError) {
      setContactSettingsError(validationError);
      return;
    }
    if (!hasContactSettingsChanges) {
      setContactSettingsNotice("contactSettings.noChanges");
      return;
    }
    setContactSettingsSaving(true);
    setContactSettingsError(null);
    setContactSettingsNotice(null);
    try {
      const payload = buildContactSettingsPayload(
        contactSettingsForm,
        contactSettings,
      );
      const response = await adminApi.updateContactSettings(payload);
      setContactSettings(response.data);
      setContactSettingsForm(contactSettingsToForm(response.data));
      setContactSettingsNotice("contactSettings.saveSuccess");
      setShowContactSettingsDiff(false);
    } catch (err) {
      if (isUnauthorizedError(err)) {
        handleUnauthorized();
      } else if (isConflictError(err)) {
        try {
          const latest = await adminApi.getContactSettings();
          setContactSettings(latest.data);
        } catch (refreshErr) {
          if (isUnauthorizedError(refreshErr)) {
            handleUnauthorized();
          } else if (isNotFoundError(refreshErr)) {
            setContactSettings(null);
          } else {
            console.error(refreshErr);
            setContactSettingsError("contactSettings.saveError");
          }
        }
        setContactSettingsNotice(null);
        setContactSettingsError("contactSettings.conflict");
        setShowContactSettingsDiff(true);
      } else {
        console.error(err);
        setContactSettingsError("contactSettings.saveError");
      }
    } finally {
      setContactSettingsSaving(false);
    }
  }, [
    contactSettings,
    contactSettingsForm,
    hasContactSettingsChanges,
    handleUnauthorized,
  ]);

  const handleReviewContactSettingsDiff = useCallback(() => {
    setShowContactSettingsDiff(true);
    setContactSettingsError(null);
    setContactSettingsNotice(null);
  }, []);

  const handleContactPreviewLocaleChange = useCallback((locale: "ja" | "en") => {
    setContactPreviewLocale(locale);
  }, []);

  const refreshAll = useCallback(async () => {
    setLoading(true);
    try {
      const statusRes = await adminApi.health();
      setStatus(statusRes.data.status);
      setAuthState("authenticated");
    } catch (err) {
      if (isUnauthorizedError(err)) {
        handleUnauthorized();
      } else {
        console.error(err);
        setError("status.error");
      }
      setLoading(false);
      return;
    }

    try {
      const [
        summaryRes,
        profileRes,
        projectRes,
        researchRes,
        contactRes,
        blacklistRes,
        techCatalogRes,
        homeRes,
      ] = await Promise.all([
        adminApi.fetchSummary(),
        adminApi.getProfile(),
        adminApi.listProjects(),
        adminApi.listResearch(),
        adminApi.listContacts(),
        adminApi.listBlacklist(),
        adminApi.listTechCatalog({ includeInactive: false }),
        adminApi
          .getHomeSettings()
          .then((res) => res)
          .catch((err) => {
            if (isNotFoundError(err)) {
              return { data: null };
            }
            throw err;
          }),
      ]);

      let contactSettingsData: ContactFormSettings | null = null;
      try {
        const contactSettingsRes = await adminApi.getContactSettings();
        contactSettingsData = contactSettingsRes.data;
      } catch (err) {
        if (isUnauthorizedError(err)) {
          throw err;
        }
        if (!isNotFoundError(err)) {
          throw err;
        }
      }

      setSummary(summaryRes.data);
      const profileData = profileRes.data;
      setProfileForm(profileToForm(profileData));
      setProjects(projectRes.data);
      setResearch(researchRes.data);
      setContacts(contactRes.data);
      setContactEdits(buildContactEditMap(contactRes.data));
      setBlacklist(blacklistRes.data);
      setTechCatalog(techCatalogRes.data);
      setHomeSettingsForm(
        homeSettingsToForm(homeRes.data ?? null, profileData.id),
      );
      setContactSettings(contactSettingsData);
      setContactSettingsForm(contactSettingsToForm(contactSettingsData));
      setContactSettingsError(null);
      setContactSettingsNotice(null);
      setContactSettingsSaving(false);
      setShowContactSettingsDiff(false);
      setContactPreviewLocale("ja");
      setHomeSettingsError(null);
      setHomeSettingsNotice(null);
      setHomeSettingsSaving(false);
      setError(null);
    } catch (err) {
      if (isUnauthorizedError(err)) {
        handleUnauthorized();
      } else {
        console.error(err);
        setError("status.error");
      }
    } finally {
      setLoading(false);
    }
  }, [handleUnauthorized]);

  useEffect(() => {
    if (authState === "authenticated") {
      return;
    }

    let cancelled = false;
  const resumeSession = async () => {
    try {
      const sessionRes = await adminApi.session();
      if (cancelled) {
        return;
      }
      if (!sessionRes?.data) {
        console.error("session response missing data");
        handleUnauthorized();
        return;
      }
      if (sessionRes.data.active) {
        setSession(sessionRes.data);
        setAuthState("authenticated");
        void refreshAll();
      } else {
        handleUnauthorized();
        }
      } catch (err) {
        if (cancelled) {
          return;
        }
        if (isUnauthorizedError(err)) {
          handleUnauthorized();
        } else {
          console.error(err);
          handleUnauthorized();
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    };

    void resumeSession();

    return () => {
      cancelled = true;
    };
  }, [authState, refreshAll, handleUnauthorized, setSession]);

  useEffect(() => {
    if (authState !== "authenticated") {
      return;
    }

  const poll = async () => {
    try {
      const sessionRes = await adminApi.session();
      if (!sessionRes?.data) {
        return;
      }
      if (sessionRes.data.active) {
        setSession(sessionRes.data);
      } else {
        handleUnauthorized();
      }
      } catch (err) {
        if (isUnauthorizedError(err)) {
          handleUnauthorized();
        } else {
          console.error(err);
        }
      }
    };

    void poll();
    const intervalId = window.setInterval(poll, 60_000);
    return () => window.clearInterval(intervalId);
  }, [authState, handleUnauthorized, setSession]);

  const run = useCallback(
    async (operation: () => Promise<unknown>) => {
      try {
        await operation();
        await refreshAll();
      } catch (err) {
        if (isUnauthorizedError(err)) {
          handleUnauthorized();
        } else {
          console.error(err);
          setError("status.error");
        }
      }
    },
    [refreshAll, handleUnauthorized],
  );

  const logout = useCallback(() => {
    if (typeof window !== "undefined") {
      const cleanUrl = `${window.location.pathname}${window.location.search}`;
      window.history.replaceState(null, "", cleanUrl);
    }
    handleUnauthorized();
  }, [handleUnauthorized]);

  const handleSaveProfile = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const missingProviders = requiredSocialProviders.filter(
      (provider) =>
        !profileForm.socialLinks.some(
          (link) =>
            link.provider === provider && trimValue(link.url).length > 0,
        ),
    );
    if (missingProviders.length > 0) {
      setError("profile.validation.requiredSocialLinks");
      return;
    }
    setError(null);
    const payload = normalizeProfileFormState(profileForm);
    await run(async () => {
      await adminApi.updateProfile(payload);
    });
  };

  const handleSaveHomeSettings = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setHomeSettingsNotice(null);
    setHomeSettingsError(null);

    const validationMessage = validateHomeSettingsForm(homeSettingsForm);
    if (validationMessage) {
      setHomeSettingsError(validationMessage);
      return;
    }

    let payload: Parameters<
      typeof adminApi.updateHomeSettings
    >[0];
    try {
      payload = normalizeHomeSettingsForm(homeSettingsForm);
    } catch {
      setHomeSettingsError("homeSettings.validation.missingRecord");
      return;
    }

    setHomeSettingsSaving(true);
    try {
      const response = await adminApi.updateHomeSettings(payload);
      setHomeSettingsForm(
        homeSettingsToForm(response.data, payload.profileId),
      );
      setHomeSettingsNotice("homeSettings.saveSuccess");
    } catch (err) {
      if (isUnauthorizedError(err)) {
        handleUnauthorized();
      } else if (isConflictError(err)) {
        setHomeSettingsError("homeSettings.conflict");
        await refreshAll();
      } else {
        console.error(err);
        setHomeSettingsError("homeSettings.saveError");
      }
    } finally {
      setHomeSettingsSaving(false);
    }
  };

  const handleSubmitProject = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const payload = projectFormToPayload(projectForm);

    await run(async () => {
      if (editingProjectId != null) {
        await adminApi.updateProject(editingProjectId, payload);
      } else {
        await adminApi.createProject(payload);
      }
    });
    setProjectForm({ ...emptyProjectForm });
    setEditingProjectId(null);
    setShowProjectPreview(false);
  };

  const handleSubmitResearch = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const payload = researchFormToPayload(researchForm);

    await run(async () => {
      if (editingResearchId != null) {
        await adminApi.updateResearch(editingResearchId, payload);
      } else {
        await adminApi.createResearch(payload);
      }
    });
    setResearchForm(createEmptyResearchForm());
    setEditingResearchId(null);
  };

  const handleSubmitBlacklist = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const payload = {
      email: blacklistForm.email.trim(),
      reason: blacklistForm.reason.trim(),
    };

    await run(async () => {
      if (editingBlacklistId != null) {
        await adminApi.updateBlacklist(editingBlacklistId, payload);
      } else {
        await adminApi.createBlacklist(payload);
      }
    });
    setBlacklistForm({ ...emptyBlacklistForm });
    setEditingBlacklistId(null);
  };

  const toggleProjectPublished = (project: AdminProject) =>
    run(async () => {
      const payload = projectToPayload({
        ...project,
        published: !project.published,
      });
      await adminApi.updateProject(project.id, payload);
    });

  const toggleResearchDraft = (item: AdminResearch) =>
    run(async () => {
      const payload = researchToPayload(item);
      payload.isDraft = !item.isDraft;
      await adminApi.updateResearch(item.id, payload);
    });

  const handleAddResearchTag = () =>
    setResearchForm((prev) => ({
      ...prev,
      tags: [
        ...prev.tags,
        { value: "", sortOrder: String(prev.tags.length + 1) },
      ],
    }));

  const handleUpdateResearchTag = (
    index: number,
    payload: Partial<ResearchTagForm>,
  ) =>
    setResearchForm((prev) => {
      const tags = prev.tags.map((tag, idx) =>
        idx === index ? { ...tag, ...payload } : tag,
      );
      return { ...prev, tags };
    });

  const handleRemoveResearchTag = (index: number) =>
    setResearchForm((prev) => ({
      ...prev,
      tags: prev.tags.filter((_, idx) => idx !== index),
    }));

  const handleAddResearchLink = () =>
    setResearchForm((prev) => ({
      ...prev,
      links: [
        ...prev.links,
        {
          type: "paper",
          labelJa: "",
          labelEn: "",
          url: "",
          sortOrder: String(prev.links.length + 1),
        },
      ],
    }));

  const handleUpdateResearchLink = (
    index: number,
    payload: Partial<ResearchLinkForm>,
  ) =>
    setResearchForm((prev) => {
      const links = prev.links.map((link, idx) =>
        idx === index ? { ...link, ...payload } : link,
      );
      return { ...prev, links };
    });

  const handleRemoveResearchLink = (index: number) =>
    setResearchForm((prev) => ({
      ...prev,
      links: prev.links.filter((_, idx) => idx !== index),
    }));

  const handleAddResearchAsset = () =>
    setResearchForm((prev) => ({
      ...prev,
      assets: [
        ...prev.assets,
        {
          url: "",
          captionJa: "",
          captionEn: "",
          sortOrder: String(prev.assets.length + 1),
        },
      ],
    }));

  const handleUpdateResearchAsset = (
    index: number,
    payload: Partial<ResearchAssetForm>,
  ) =>
    setResearchForm((prev) => {
      const assets = prev.assets.map((asset, idx) =>
        idx === index ? { ...asset, ...payload } : asset,
      );
      return { ...prev, assets };
    });

  const handleRemoveResearchAsset = (index: number) =>
    setResearchForm((prev) => ({
      ...prev,
      assets: prev.assets.filter((_, idx) => idx !== index),
    }));

  const handleAddResearchTech = () =>
    setResearchForm((prev) => ({
      ...prev,
      tech: [
        ...prev.tech,
        {
          techId: "",
          context: "primary",
          note: "",
          sortOrder: String(prev.tech.length + 1),
        },
      ],
    }));

  const handleUpdateResearchTech = (
    index: number,
    payload: Partial<ResearchTechForm>,
  ) =>
    setResearchForm((prev) => {
      const tech = prev.tech.map((membership, idx) =>
        idx === index ? { ...membership, ...payload } : membership,
      );
      return { ...prev, tech };
    });

  const handleRemoveResearchTech = (index: number) =>
    setResearchForm((prev) => ({
      ...prev,
      tech: prev.tech.filter((_, idx) => idx !== index),
    }));

  const handleAddProjectTech = useCallback(
    (techId: number) =>
      setProjectForm((prev) => ({
        ...prev,
        tech: [
          ...prev.tech,
          {
            techId: String(techId),
            context: "primary",
            note: "",
            sortOrder: String(prev.tech.length + 1),
          },
        ],
      })),
    [],
  );

  const handleUpdateProjectTech = (
    index: number,
    payload: Partial<ProjectTechForm>,
  ) =>
    setProjectForm((prev) => {
      const tech = prev.tech.map((membership, idx) =>
        idx === index ? { ...membership, ...payload } : membership,
      );
      return { ...prev, tech };
    });

  const handleRemoveProjectTech = (index: number) =>
    setProjectForm((prev) => ({
      ...prev,
      tech: prev.tech.filter((_, idx) => idx !== index),
    }));

  const renderProjectTechSection = () => (
    <div className="md:col-span-2">
      <div className="space-y-3 rounded-md border border-slate-200 p-4">
        <div className="flex flex-col gap-2 md:flex-row md:items-center md:justify-between">
          <div>
            <span className="block text-sm font-medium text-slate-700">
              {t("fields.tech")}
            </span>
            <p className="text-xs text-slate-500">
              {t("fields.techSearchDescription") ??
                "Search the catalog to add technologies. Each entry can be marked as primary or supporting."}
            </p>
          </div>
          <div className="flex w-full flex-col gap-2 md:w-96">
            <input
              className="w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
              value={projectTechSearch}
              onChange={(event) => setProjectTechSearch(event.target.value)}
              placeholder={t("fields.techSearchPlaceholder") ?? ""}
            />
            <div className="flex flex-wrap gap-2">
              {filteredProjectTech.length > 0 ? (
                filteredProjectTech.map((entry) => (
                  <button
                    type="button"
                    key={`catalog-${entry.id}`}
                    className="inline-flex items-center rounded-full border border-slate-200 px-3 py-1 text-xs text-slate-600 transition hover:border-sky-400 hover:text-sky-600"
                    onClick={() => {
                      handleAddProjectTech(entry.id);
                      setProjectTechSearch("");
                    }}
                  >
                    {entry.displayName}
                  </button>
                ))
              ) : (
                <span className="text-xs text-slate-500">
                  {t("fields.techSearchEmpty")}
                </span>
              )}
            </div>
          </div>
        </div>
        <div className="space-y-2">
          {projectForm.tech.length === 0 && (
            <p className="text-sm text-slate-500">
              {t("projects.noTech") ??
                "No technologies added yet. Use the search above to add from the catalog."}
            </p>
          )}
          {projectForm.tech.map((membership, index) => {
            const selected = techCatalog.find(
              (entry) => String(entry.id) === membership.techId,
            );
            return (
              <div
                key={`project-tech-${index}`}
                className="grid gap-2 md:grid-cols-5"
              >
                <select
                  className="rounded-md border border-slate-200 px-3 py-2 text-sm md:col-span-2"
                  value={membership.techId}
                  onChange={(event) =>
                    handleUpdateProjectTech(index, {
                      techId: event.target.value,
                    })
                  }
                >
                  <option value="">
                    {t("fields.selectTech") ?? "Select technology"}
                  </option>
                  {techCatalog.map((entry) => (
                    <option key={entry.id} value={entry.id}>
                      {entry.displayName}
                    </option>
                  ))}
                  {membership.techId && !selected && (
                    <option value={membership.techId}>
                      {t("fields.unknownTech", {
                        id: membership.techId,
                      }) ??
                        `Unknown (#${membership.techId})`}
                    </option>
                  )}
                </select>
                <select
                  className="rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={membership.context}
                  onChange={(event) =>
                    handleUpdateProjectTech(index, {
                      context: event.target.value as TechContext,
                    })
                  }
                >
                  {techContexts.map((context) => (
                    <option key={context} value={context}>
                      {context}
                    </option>
                  ))}
                </select>
                <input
                  className="rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={membership.sortOrder}
                  onChange={(event) =>
                    handleUpdateProjectTech(index, {
                      sortOrder: event.target.value,
                    })
                  }
                  placeholder={t("fields.sortOrder") ?? "Sort"}
                />
                <div className="flex items-center gap-2">
                  <input
                    className="flex-1 rounded-md border border-slate-200 px-3 py-2 text-sm"
                    value={membership.note}
                    onChange={(event) =>
                      handleUpdateProjectTech(index, {
                        note: event.target.value,
                      })
                    }
                    placeholder={t("fields.note") ?? "Note"}
                  />
                  <button
                    type="button"
                    className="rounded-md border border-slate-200 px-3 py-2 text-sm text-slate-600 transition hover:bg-slate-100"
                    onClick={() => handleRemoveProjectTech(index)}
                  >
                    {t("actions.remove")}
                  </button>
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );

  const deleteProject = (id: number) => run(() => adminApi.deleteProject(id));
  const deleteResearch = (id: number) => run(() => adminApi.deleteResearch(id));
  const deleteBlacklistEntry = (id: number) =>
    run(() => adminApi.deleteBlacklist(id));

const handleEditProject = (project: AdminProject) => {
  setEditingProjectId(project.id);
  setProjectForm(projectToForm(project));
  setProjectTechSearch("");
  setShowProjectPreview(false);
};

  const handleEditResearch = (item: AdminResearch) => {
    setEditingResearchId(item.id);
    setResearchForm(researchToForm(item));
  };

  const handleEditBlacklist = (entry: BlacklistEntry) => {
    setEditingBlacklistId(entry.id);
    setBlacklistForm({
      email: entry.email,
      reason: entry.reason,
    });
  };

const resetProjectForm = () => {
  setProjectForm({ ...emptyProjectForm });
  setEditingProjectId(null);
  setProjectTechSearch("");
};

  const resetResearchForm = () => {
    setResearchForm(createEmptyResearchForm());
    setEditingResearchId(null);
  };

  const resetBlacklistForm = () => {
    setBlacklistForm({ ...emptyBlacklistForm });
    setEditingBlacklistId(null);
  };

  const handleContactEditChange = (
    id: string,
    field: keyof ContactEditState,
    value: string,
  ) => {
    setContactEdits((prev) => ({
      ...prev,
      [id]: {
        ...prev[id],
        [field]: field === "status" ? (value as ContactStatus) : value,
      },
    }));
  };

  const handleSaveContact = async (id: string) => {
    const edit = contactEdits[id];
    if (!edit) {
      return;
    }
    await run(async () => {
      await adminApi.updateContact(id, {
        topic: edit.topic,
        message: edit.message,
        status: edit.status,
        adminNote: edit.adminNote,
      });
    });
  };

  const handleResetContact = (contact: ContactMessage) => {
    setContactEdits((prev) => ({
      ...prev,
      [contact.id]: {
        topic: contact.topic,
        message: contact.message,
        status: contact.status,
        adminNote: contact.adminNote,
      },
    }));
  };

  const handleDeleteContact = (id: string) =>
    run(() => adminApi.deleteContact(id));

  const profileUpdatedDisplay = useMemo(() => {
    if (!summary?.profileUpdatedAt) {
      return t("summary.notUpdated");
    }
    const date = new Date(summary.profileUpdatedAt);
    if (Number.isNaN(date.getTime())) {
      return t("summary.notUpdated");
    }
    return new Intl.DateTimeFormat(undefined, {
      dateStyle: "medium",
      timeStyle: "short",
    }).format(date);
  }, [summary?.profileUpdatedAt, t]);

  if (authState === "checking") {
    return (
      <div className="flex min-h-screen items-center justify-center bg-slate-100 p-6">
        <p className="text-sm text-slate-600">{t("status.loading")}</p>
      </div>
    );
  }

  if (authState === "unauthorized") {
    const loginUrl =
      import.meta.env.VITE_ADMIN_LOGIN_URL ?? "/api/admin/auth/login";
    const supportEmail =
      import.meta.env.VITE_ADMIN_SUPPORT_EMAIL ?? "support@example.com";

    return (
      <div className="flex min-h-screen items-center justify-center bg-slate-100 p-6">
        <div className="w-full max-w-md rounded-lg border border-slate-200 bg-white p-6 text-center shadow-sm">
          <h1 className="text-xl font-semibold text-slate-900">
            {t("auth.requiredTitle")}
          </h1>
          <p className="mt-2 text-sm text-slate-600">
            {t("auth.requiredDescription")}
          </p>
          <button
            className="mt-4 inline-flex items-center justify-center rounded-md bg-slate-900 px-4 py-2 text-sm font-medium text-white shadow-sm transition hover:bg-slate-800 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-slate-900"
            type="button"
            onClick={() => window.location.assign(loginUrl)}
          >
            {t("auth.signIn")}
          </button>
          <p className="mt-4 text-xs text-slate-500">
            {t("auth.supportPrompt")}{" "}
            <a
              className="font-medium text-slate-700 underline hover:text-slate-900"
              href={`mailto:${supportEmail}`}
              rel="noreferrer"
            >
              {t("auth.contactSupport")}
            </a>
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-slate-100">
      <header className="bg-slate-900 p-6 text-white">
        <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h1 className="text-2xl font-bold">{t("dashboard.title")}</h1>
            <p className="text-sm text-slate-300">{t("dashboard.subtitle")}</p>
          </div>
          <button
            type="button"
            onClick={logout}
            className="inline-flex items-center justify-center rounded-md bg-white/10 px-4 py-2 text-sm font-medium text-white transition hover:bg-white/20 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white"
          >
            {t("auth.signOut")}
          </button>
        </div>
      </header>
      <main className="mx-auto flex max-w-6xl flex-col gap-6 p-6">
        <section className="rounded-lg border border-slate-200 bg-white p-6 shadow-sm">
          <div className="flex flex-col gap-3 md:flex-row md:items-start md:justify-between">
            <div>
              <h2 className="text-lg font-semibold text-slate-800">
                {t("contactSettings.title")}
              </h2>
              <p className="mt-1 text-sm text-slate-600">
                {t("contactSettings.description")}
              </p>
            </div>
            <div className="flex flex-wrap gap-2">
              <button
                type="button"
                className="rounded-md border border-slate-200 px-3 py-2 text-sm text-slate-600 transition hover:bg-slate-100 disabled:cursor-not-allowed disabled:opacity-60"
                onClick={handleReviewContactSettingsDiff}
                disabled={!hasContactSettingsChanges}
              >
                {t("contactSettings.reviewDiff")}
              </button>
              <button
                type="button"
                className="rounded-md border border-slate-200 px-3 py-2 text-sm text-slate-600 transition hover:bg-slate-100 disabled:cursor-not-allowed disabled:opacity-60"
                onClick={handleResetContactSettings}
                disabled={!hasContactSettingsChanges}
              >
                {t("actions.reset")}
              </button>
            </div>
          </div>
          {contactSettingsError && (
            <p className="mt-3 rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700">
              {t(contactSettingsError)}
            </p>
          )}
          {contactSettingsNotice && (
            <p className="mt-3 rounded-md border border-emerald-200 bg-emerald-50 px-3 py-2 text-sm text-emerald-700">
              {t(contactSettingsNotice)}
            </p>
          )}
          <div className="mt-6 grid gap-4 md:grid-cols-2">
            <div>
              <label className="block text-sm font-medium text-slate-700">
                {t("contactSettings.fields.heroTitleJa")}
              </label>
              <input
                className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                value={contactSettingsForm.heroTitleJa}
                onChange={(event) =>
                  handleContactSettingsFieldChange(
                    "heroTitleJa",
                    event.target.value,
                  )
                }
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-700">
                {t("contactSettings.fields.heroTitleEn")}
              </label>
              <input
                className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                value={contactSettingsForm.heroTitleEn}
                onChange={(event) =>
                  handleContactSettingsFieldChange(
                    "heroTitleEn",
                    event.target.value,
                  )
                }
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-700">
                {t("contactSettings.fields.heroDescriptionJa")}
              </label>
              <textarea
                className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                rows={3}
                value={contactSettingsForm.heroDescriptionJa}
                onChange={(event) =>
                  handleContactSettingsFieldChange(
                    "heroDescriptionJa",
                    event.target.value,
                  )
                }
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-700">
                {t("contactSettings.fields.heroDescriptionEn")}
              </label>
              <textarea
                className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                rows={3}
                value={contactSettingsForm.heroDescriptionEn}
                onChange={(event) =>
                  handleContactSettingsFieldChange(
                    "heroDescriptionEn",
                    event.target.value,
                  )
                }
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-700">
                {t("contactSettings.fields.consentTextJa")}
              </label>
              <textarea
                className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                rows={3}
                value={contactSettingsForm.consentTextJa}
                onChange={(event) =>
                  handleContactSettingsFieldChange(
                    "consentTextJa",
                    event.target.value,
                  )
                }
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-700">
                {t("contactSettings.fields.consentTextEn")}
              </label>
              <textarea
                className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                rows={3}
                value={contactSettingsForm.consentTextEn}
                onChange={(event) =>
                  handleContactSettingsFieldChange(
                    "consentTextEn",
                    event.target.value,
                  )
                }
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-700">
                {t("contactSettings.fields.minimumLeadHours")}
              </label>
              <input
                type="number"
                min={0}
                className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                placeholder="24"
                value={contactSettingsForm.minimumLeadHours}
                onChange={(event) =>
                  handleContactSettingsFieldChange(
                    "minimumLeadHours",
                    event.target.value,
                  )
                }
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-700">
                {t("contactSettings.fields.bookingWindowDays")}
              </label>
              <input
                type="number"
                min={1}
                className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                placeholder="30"
                value={contactSettingsForm.bookingWindowDays}
                onChange={(event) =>
                  handleContactSettingsFieldChange(
                    "bookingWindowDays",
                    event.target.value,
                  )
                }
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-700">
                {t("contactSettings.fields.supportEmail")}
              </label>
              <input
                className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                value={contactSettingsForm.supportEmail}
                onChange={(event) =>
                  handleContactSettingsFieldChange(
                    "supportEmail",
                    event.target.value,
                  )
                }
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-700">
                {t("contactSettings.fields.recaptchaSiteKey")}
              </label>
              <input
                className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                value={contactSettingsForm.recaptchaSiteKey}
                onChange={(event) =>
                  handleContactSettingsFieldChange(
                    "recaptchaSiteKey",
                    event.target.value,
                  )
                }
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-700">
                {t("contactSettings.fields.calendarTimezone")}
              </label>
              <input
                className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                placeholder="Asia/Tokyo"
                value={contactSettingsForm.calendarTimezone}
                onChange={(event) =>
                  handleContactSettingsFieldChange(
                    "calendarTimezone",
                    event.target.value,
                  )
                }
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-700">
                {t("contactSettings.fields.googleCalendarId")}
              </label>
              <input
                className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                value={contactSettingsForm.googleCalendarId}
                onChange={(event) =>
                  handleContactSettingsFieldChange(
                    "googleCalendarId",
                    event.target.value,
                  )
                }
              />
            </div>
          </div>
          <div className="mt-6">
            <div className="flex flex-wrap items-center justify-between gap-2">
              <h3 className="text-base font-semibold text-slate-800">
                {t("contactSettings.topics.title")}
              </h3>
              <button
                type="button"
                className="rounded-md border border-slate-200 px-3 py-2 text-sm text-slate-600 transition hover:bg-slate-100"
                onClick={handleAddContactTopic}
              >
                {t("contactSettings.addTopic")}
              </button>
            </div>
            <div className="mt-3 space-y-4">
              {contactSettingsForm.topics.map((topic, index) => (
                <div
                  key={`${topic.id}-${index}`}
                  className="rounded-md border border-slate-200 p-4"
                >
                  <div className="flex flex-wrap items-center justify-between gap-2">
                    <h4 className="text-sm font-semibold text-slate-800">
                      {topic.id || t("contactSettings.topics.placeholderId")}
                    </h4>
                    <button
                      type="button"
                      className="rounded-md border border-red-200 px-2 py-1 text-xs text-red-600 transition hover:bg-red-50"
                      onClick={() => handleRemoveContactTopic(index)}
                    >
                      {t("contactSettings.removeTopic")}
                    </button>
                  </div>
                  <div className="mt-3 grid gap-3 md:grid-cols-2">
                    <div>
                      <label className="block text-xs font-medium text-slate-600">
                        {t("contactSettings.fields.topicId")}
                      </label>
                      <input
                        className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                        value={topic.id}
                        onChange={(event) =>
                          handleContactTopicChange(
                            index,
                            "id",
                            event.target.value,
                          )
                        }
                      />
                    </div>
                    <div>
                      <label className="block text-xs font-medium text-slate-600">
                        {t("contactSettings.fields.topicLabelJa")}
                      </label>
                      <input
                        className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                        value={topic.labelJa}
                        onChange={(event) =>
                          handleContactTopicChange(
                            index,
                            "labelJa",
                            event.target.value,
                          )
                        }
                      />
                    </div>
                    <div>
                      <label className="block text-xs font-medium text-slate-600">
                        {t("contactSettings.fields.topicLabelEn")}
                      </label>
                      <input
                        className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                        value={topic.labelEn}
                        onChange={(event) =>
                          handleContactTopicChange(
                            index,
                            "labelEn",
                            event.target.value,
                          )
                        }
                      />
                    </div>
                    <div>
                      <label className="block text-xs font-medium text-slate-600">
                        {t("contactSettings.fields.topicDescriptionJa")}
                      </label>
                      <textarea
                        className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                        rows={2}
                        value={topic.descriptionJa}
                        onChange={(event) =>
                          handleContactTopicChange(
                            index,
                            "descriptionJa",
                            event.target.value,
                          )
                        }
                      />
                    </div>
                    <div>
                      <label className="block text-xs font-medium text-slate-600">
                        {t("contactSettings.fields.topicDescriptionEn")}
                      </label>
                      <textarea
                        className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                        rows={2}
                        value={topic.descriptionEn}
                        onChange={(event) =>
                          handleContactTopicChange(
                            index,
                            "descriptionEn",
                            event.target.value,
                          )
                        }
                      />
                    </div>
                  </div>
                </div>
              ))}
              {contactSettingsForm.topics.length === 0 && (
                <p className="text-sm text-slate-500">
                  {t("contactSettings.topics.empty")}
                </p>
              )}
            </div>
          </div>
          <div className="mt-6">
            <div className="flex flex-wrap items-center justify-between gap-2">
              <h3 className="text-base font-semibold text-slate-800">
                {t("contactSettings.preview.title")}
              </h3>
              <div className="inline-flex overflow-hidden rounded-md border border-slate-200">
                <button
                  type="button"
                  className={`px-3 py-1 text-xs font-medium ${
                    contactPreviewLocale === "ja"
                      ? "bg-slate-900 text-white"
                      : "bg-white text-slate-700"
                  }`}
                  onClick={() => handleContactPreviewLocaleChange("ja")}
                >
                  JA
                </button>
                <button
                  type="button"
                  className={`px-3 py-1 text-xs font-medium ${
                    contactPreviewLocale === "en"
                      ? "bg-slate-900 text-white"
                      : "bg-white text-slate-700"
                  }`}
                  onClick={() => handleContactPreviewLocaleChange("en")}
                >
                  EN
                </button>
              </div>
            </div>
            <div className="mt-3 overflow-hidden rounded-lg border border-slate-200">
              <div className="bg-slate-900 px-6 py-6 text-white">
                <h4 className="text-xl font-semibold">
                  {previewHeroTitle || t("contactSettings.preview.emptyTitle")}
                </h4>
                <p className="mt-2 text-sm opacity-90">
                  {previewHeroDescription ||
                    t("contactSettings.preview.emptyDescription")}
                </p>
              </div>
              <div className="space-y-4 px-6 py-6">
                <div className="flex flex-wrap gap-3 text-xs text-slate-500">
                  <span>
                    {t("contactSettings.preview.minimumLeadHours", {
                      count: previewMinimumLeadHours,
                    })}
                  </span>
                  <span>
                    {t("contactSettings.preview.bookingWindowDays", {
                      count: previewBookingWindowDays,
                    })}
                  </span>
                </div>
                <div>
                  <h4 className="text-sm font-semibold text-slate-700">
                    {t("contactSettings.preview.topicsHeading")}
                  </h4>
                  <div className="mt-2 space-y-3">
                    {previewTopics.length > 0 ? (
                      previewTopics.map((topic) => (
                        <div
                          key={`${topic.id}-${topic.label}`}
                          className="rounded-md border border-slate-200 p-3"
                        >
                          <p className="text-sm font-medium text-slate-800">
                            {topic.label ||
                              t("contactSettings.preview.topicPlaceholder")}
                          </p>
                          <p className="mt-1 text-sm text-slate-600">
                            {topic.description ||
                              t(
                                "contactSettings.preview.topicDescriptionPlaceholder",
                              )}
                          </p>
                        </div>
                      ))
                    ) : (
                      <p className="text-sm text-slate-500">
                        {t("contactSettings.preview.topicsEmpty")}
                      </p>
                    )}
                  </div>
                </div>
                <div className="grid gap-3 md:grid-cols-2">
                  <div>
                    <p className="text-xs uppercase tracking-wide text-slate-500">
                      {t("contactSettings.fields.supportEmail")}
                    </p>
                    <p className="text-sm text-slate-700">
                      {previewSupportEmail ||
                        t("contactSettings.preview.supportEmailPlaceholder")}
                    </p>
                  </div>
                  <div>
                    <p className="text-xs uppercase tracking-wide text-slate-500">
                      {t("contactSettings.fields.calendarTimezone")}
                    </p>
                    <p className="text-sm text-slate-700">
                      {previewCalendarTimezone ||
                        t(
                          "contactSettings.preview.calendarTimezonePlaceholder",
                        )}
                    </p>
                  </div>
                </div>
                <div>
                  <p className="text-xs uppercase tracking-wide text-slate-500">
                    {t("contactSettings.fields.consentText")}
                  </p>
                  <p className="mt-1 text-sm text-slate-700">
                    {previewConsent ||
                      t("contactSettings.preview.consentPlaceholder")}
                  </p>
                </div>
                <div>
                  <p className="text-xs uppercase tracking-wide text-slate-500">
                    {t("contactSettings.fields.recaptchaSiteKey")}
                  </p>
                  <p className="text-sm text-slate-700">
                    {previewRecaptchaSiteKey ||
                      t("contactSettings.preview.recaptchaPlaceholder")}
                  </p>
                </div>
              </div>
            </div>
          </div>
          <div className="mt-6 flex flex-wrap gap-2">
            <button
              type="button"
              className="rounded-md bg-slate-900 px-4 py-2 text-sm font-medium text-white transition hover:bg-slate-800 disabled:cursor-not-allowed disabled:bg-slate-400 disabled:text-white/80"
              onClick={handleSaveContactSettings}
              disabled={
                contactSettingsSaving ||
                !contactSettings ||
                !hasContactSettingsChanges
              }
            >
              {contactSettingsSaving
                ? t("contactSettings.saving")
                : t("actions.save")}
            </button>
          </div>
        </section>

        <section className="rounded-lg border border-slate-200 bg-white p-6 shadow-sm">
          <h2 className="text-lg font-semibold text-slate-800">
            {t("homeSettings.title")}
          </h2>
          <p className="mt-1 text-sm text-slate-600">
            {t("homeSettings.description")}
          </p>
          {homeSettingsError && (
            <div className="mt-4 rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-700">
              {t(homeSettingsError)}
            </div>
          )}
          {homeSettingsNotice && (
            <div className="mt-4 rounded-md border border-emerald-200 bg-emerald-50 p-3 text-sm text-emerald-700">
              {t(homeSettingsNotice)}
            </div>
          )}
          <form className="mt-4 space-y-6" onSubmit={handleSaveHomeSettings}>
            <div className="grid gap-4 md:grid-cols-2">
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("homeSettings.fields.heroSubtitleJa")}
                </label>
                <textarea
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  rows={3}
                  value={homeSettingsForm.heroSubtitleJa}
                  onChange={(event) =>
                    handleHomeSettingsFieldChange(
                      "heroSubtitleJa",
                      event.target.value,
                    )
                  }
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("homeSettings.fields.heroSubtitleEn")}
                </label>
                <textarea
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  rows={3}
                  value={homeSettingsForm.heroSubtitleEn}
                  onChange={(event) =>
                    handleHomeSettingsFieldChange(
                      "heroSubtitleEn",
                      event.target.value,
                    )
                  }
                />
              </div>
            </div>

            <div>
              <div className="flex items-center justify-between">
                <h3 className="text-sm font-semibold text-slate-700">
                  {t("homeSettings.sections.quickLinks")}
                </h3>
                <span className="text-xs text-slate-500">
                  {t("homeSettings.hints.quickLinks")}
                </span>
              </div>
              <div className="mt-3">
                <HomeQuickLinksEditor
                  items={homeSettingsForm.quickLinks}
                  onChange={handleHomeQuickLinksChange}
                />
              </div>
            </div>

            <div>
              <div className="flex items-center justify-between">
                <h3 className="text-sm font-semibold text-slate-700">
                  {t("homeSettings.sections.chipSources")}
                </h3>
                <span className="text-xs text-slate-500">
                  {t("homeSettings.hints.chipSources")}
                </span>
              </div>
              <div className="mt-3">
                <HomeChipSourcesEditor
                  items={homeSettingsForm.chipSources}
                  onChange={handleHomeChipSourcesChange}
                />
              </div>
            </div>

            <div className="flex justify-end">
              <button
                className="inline-flex items-center justify-center rounded-md bg-slate-900 px-4 py-2 text-sm font-medium text-white transition hover:bg-slate-800 disabled:cursor-not-allowed disabled:bg-slate-400"
                type="submit"
                disabled={homeSettingsSaving}
              >
                {homeSettingsSaving
                  ? t("homeSettings.saving", { defaultValue: "Saving…" })
                  : t("actions.save")}
              </button>
            </div>
          </form>
        </section>

        <section className="rounded-lg border border-slate-200 bg-white p-6 shadow-sm">
          <h2 className="text-lg font-semibold text-slate-800">
            {t("dashboard.systemStatus")}
          </h2>
          <div className="mt-4 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            <div className="rounded-md bg-slate-900 p-4 text-white">
              <span className="font-mono uppercase tracking-wide text-slate-400">
                {t("dashboard.apiStatus")}
              </span>
              <p className="text-2xl font-bold text-emerald-400">{status}</p>
            </div>
            <div className="rounded-md border border-slate-200 bg-white p-4">
              <span className="text-xs font-semibold uppercase text-slate-500">
                {t("summary.profileUpdated")}
              </span>
              <p className="mt-1 text-lg font-semibold text-slate-800">
                {profileUpdatedDisplay}
              </p>
            </div>
            <div className="rounded-md border border-slate-200 bg-white p-4">
              <span className="text-xs font-semibold uppercase text-slate-500">
                {t("summary.skillCount")}
              </span>
              <p className="mt-1 text-lg font-semibold text-slate-800">
                {summary?.skillCount ?? 0}
              </p>
            </div>
            <div className="rounded-md border border-slate-200 bg-white p-4">
              <span className="text-xs font-semibold uppercase text-slate-500">
                {t("summary.focusAreaCount")}
              </span>
              <p className="mt-1 text-lg font-semibold text-slate-800">
                {summary?.focusAreaCount ?? 0}
              </p>
            </div>
            <div className="rounded-md border border-slate-200 bg-white p-4">
              <span className="text-xs font-semibold uppercase text-slate-500">
                {t("summary.projects")}
              </span>
              <p className="mt-1 text-lg font-semibold text-slate-800">
                {summary
                  ? `${summary.publishedProjects} / ${summary.draftProjects}`
                  : "0 / 0"}
              </p>
            </div>
            <div className="rounded-md border border-slate-200 bg-white p-4">
              <span className="text-xs font-semibold uppercase text-slate-500">
                {t("summary.research")}
              </span>
              <p className="mt-1 text-lg font-semibold text-slate-800">
                {summary
                  ? `${summary.publishedResearch} / ${summary.draftResearch}`
                  : "0 / 0"}
              </p>
            </div>
            <div className="rounded-md border border-slate-200 bg-white p-4">
              <span className="text-xs font-semibold uppercase text-slate-500">
                {t("summary.pendingContacts")}
              </span>
              <p className="mt-1 text-lg font-semibold text-slate-800">
                {summary?.pendingContacts ?? 0}
              </p>
            </div>
            <div className="rounded-md border border-slate-200 bg-white p-4">
              <span className="text-xs font-semibold uppercase text-slate-500">
                {t("summary.blacklist")}
              </span>
              <p className="mt-1 text-lg font-semibold text-slate-800">
                {summary?.blacklistEntries ?? 0}
              </p>
            </div>
          </div>
        </section>

        {error && (
          <div className="rounded-md border border-red-200 bg-red-50 p-4 text-sm text-red-700">
            {t(error)}
          </div>
        )}

        <section className="rounded-lg border border-slate-200 bg-white p-6 shadow-sm">
          <h2 className="text-lg font-semibold text-slate-800">
            {t("profile.title")}
          </h2>
          <p className="mt-1 text-sm text-slate-600">{t("profile.description")}</p>
          <form className="mt-4 space-y-6" onSubmit={handleSaveProfile}>
            <div className="grid gap-4 md:grid-cols-2">
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("profile.fields.displayName", {
                    defaultValue: "Display name",
                  })}
                </label>
                <input
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={profileForm.displayName}
                  onChange={(event) =>
                    setProfileForm((prev) => ({
                      ...prev,
                      displayName: event.target.value,
                    }))
                  }
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("profile.fields.avatarUrl", {
                    defaultValue: "Avatar URL",
                  })}
                </label>
                <input
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={profileForm.avatarUrl}
                  onChange={(event) =>
                    setProfileForm((prev) => ({
                      ...prev,
                      avatarUrl: event.target.value,
                    }))
                  }
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.titleJa")}
                </label>
                <input
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={profileForm.headlineJa}
                  onChange={(event) =>
                    setProfileForm((prev) => ({
                      ...prev,
                      headlineJa: event.target.value,
                    }))
                  }
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.titleEn")}
                </label>
                <input
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={profileForm.headlineEn}
                  onChange={(event) =>
                    setProfileForm((prev) => ({
                      ...prev,
                      headlineEn: event.target.value,
                    }))
                  }
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("profile.fields.themeMode", {
                    defaultValue: "Theme mode",
                  })}
                </label>
                <select
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={profileForm.themeMode}
                  onChange={(event) =>
                    setProfileForm((prev) => ({
                      ...prev,
                      themeMode: event.target.value as ProfileFormState["themeMode"],
                    }))
                  }
                >
                  {themeOptions.map((mode) => (
                    <option key={mode} value={mode}>
                      {mode}
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("profile.fields.themeAccent", {
                    defaultValue: "Accent color",
                  })}
                </label>
                <input
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={profileForm.themeAccentColor}
                  onChange={(event) =>
                    setProfileForm((prev) => ({
                      ...prev,
                      themeAccentColor: event.target.value,
                    }))
                  }
                />
              </div>
            </div>

            <div className="grid gap-4 md:grid-cols-2">
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.summaryJa")}
                </label>
                <textarea
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  rows={4}
                  value={profileForm.summaryJa}
                  onChange={(event) =>
                    setProfileForm((prev) => ({
                      ...prev,
                      summaryJa: event.target.value,
                    }))
                  }
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.summaryEn")}
                </label>
                <textarea
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  rows={4}
                  value={profileForm.summaryEn}
                  onChange={(event) =>
                    setProfileForm((prev) => ({
                      ...prev,
                      summaryEn: event.target.value,
                    }))
                  }
                />
              </div>
            </div>

            <div className="grid gap-4 md:grid-cols-2">
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("profile.fields.locationJa", {
                    defaultValue: "Location (JA)",
                  })}
                </label>
                <input
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={profileForm.locationJa}
                  onChange={(event) =>
                    setProfileForm((prev) => ({
                      ...prev,
                      locationJa: event.target.value,
                    }))
                  }
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("profile.fields.locationEn", {
                    defaultValue: "Location (EN)",
                  })}
                </label>
                <input
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={profileForm.locationEn}
                  onChange={(event) =>
                    setProfileForm((prev) => ({
                      ...prev,
                      locationEn: event.target.value,
                    }))
                  }
                />
              </div>
            </div>

            <div>
              <h3 className="text-sm font-semibold text-slate-700">
                {t("profile.sections.lab", { defaultValue: "Lab information" })}
              </h3>
              <div className="mt-2 grid gap-4 md:grid-cols-2">
                <div>
                  <label className="block text-sm font-medium text-slate-700">
                    {t("profile.fields.labNameJa", {
                      defaultValue: "Lab name (JA)",
                    })}
                  </label>
                  <input
                    className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                    value={profileForm.labNameJa}
                    onChange={(event) =>
                      setProfileForm((prev) => ({
                        ...prev,
                        labNameJa: event.target.value,
                      }))
                    }
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-slate-700">
                    {t("profile.fields.labNameEn", {
                      defaultValue: "Lab name (EN)",
                    })}
                  </label>
                  <input
                    className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                    value={profileForm.labNameEn}
                    onChange={(event) =>
                      setProfileForm((prev) => ({
                        ...prev,
                        labNameEn: event.target.value,
                      }))
                    }
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-slate-700">
                    {t("profile.fields.labAdvisorJa", {
                      defaultValue: "Advisor (JA)",
                    })}
                  </label>
                  <input
                    className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                    value={profileForm.labAdvisorJa}
                    onChange={(event) =>
                      setProfileForm((prev) => ({
                        ...prev,
                        labAdvisorJa: event.target.value,
                      }))
                    }
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-slate-700">
                    {t("profile.fields.labAdvisorEn", {
                      defaultValue: "Advisor (EN)",
                    })}
                  </label>
                  <input
                    className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                    value={profileForm.labAdvisorEn}
                    onChange={(event) =>
                      setProfileForm((prev) => ({
                        ...prev,
                        labAdvisorEn: event.target.value,
                      }))
                    }
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-slate-700">
                    {t("profile.fields.labRoomJa", {
                      defaultValue: "Room (JA)",
                    })}
                  </label>
                  <input
                    className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                    value={profileForm.labRoomJa}
                    onChange={(event) =>
                      setProfileForm((prev) => ({
                        ...prev,
                        labRoomJa: event.target.value,
                      }))
                    }
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-slate-700">
                    {t("profile.fields.labRoomEn", {
                      defaultValue: "Room (EN)",
                    })}
                  </label>
                  <input
                    className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                    value={profileForm.labRoomEn}
                    onChange={(event) =>
                      setProfileForm((prev) => ({
                        ...prev,
                        labRoomEn: event.target.value,
                      }))
                    }
                  />
                </div>
                <div className="md:col-span-2">
                  <label className="block text-sm font-medium text-slate-700">
                    {t("profile.fields.labUrl", { defaultValue: "Lab URL" })}
                  </label>
                  <input
                    className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                    value={profileForm.labUrl}
                    onChange={(event) =>
                      setProfileForm((prev) => ({
                        ...prev,
                        labUrl: event.target.value,
                      }))
                    }
                  />
                </div>
              </div>
            </div>

            <ProfileAffiliationList
              title={t("profile.sections.affiliations", { defaultValue: "Affiliations" })}
              items={profileForm.affiliations}
              onChange={(updater) =>
                setProfileForm((prev) => ({
                  ...prev,
                  affiliations: updater(prev.affiliations),
                }))
              }
            />

            <ProfileAffiliationList
              title={t("profile.sections.communities", { defaultValue: "Communities" })}
              items={profileForm.communities}
              onChange={(updater) =>
                setProfileForm((prev) => ({
                  ...prev,
                  communities: updater(prev.communities),
                }))
              }
            />

            <ProfileWorkHistoryList
              items={profileForm.workHistory}
              onChange={(updater) =>
                setProfileForm((prev) => ({
                  ...prev,
                  workHistory: updater(prev.workHistory),
                }))
              }
            />

            <ProfileSocialLinkList
              items={profileForm.socialLinks}
              onChange={(updater) =>
                setProfileForm((prev) => ({
                  ...prev,
                  socialLinks: updater(prev.socialLinks),
                }))
              }
            />

            <div className="flex justify-end">
              <button
                className="rounded-md bg-slate-900 px-4 py-2 text-sm font-medium text-white hover:bg-slate-800"
                type="submit"
              >
                {t("actions.save")}
              </button>
            </div>
          </form>
        </section>

        <section className="rounded-lg border border-slate-200 bg-white p-6 shadow-sm">
          <div className="flex flex-col gap-3 border-b border-slate-200 pb-4 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <h2 className="text-lg font-semibold text-slate-800">
                {t("projects.title")}
              </h2>
              <p className="text-sm text-slate-600">
                {t("projects.description")}
              </p>
            </div>
            {editingProjectId != null && (
              <button
                type="button"
                className="text-sm font-medium text-slate-600 hover:text-slate-800"
                onClick={resetProjectForm}
              >
                {t("actions.cancel")}
              </button>
            )}
          </div>
          <form className="mt-4 space-y-4" onSubmit={handleSubmitProject}>
            <div className="grid gap-4 md:grid-cols-2">
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.titleJa")}
                </label>
                <input
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={projectForm.titleJa}
                  onChange={(event) =>
                    setProjectForm((prev) => ({
                      ...prev,
                      titleJa: event.target.value,
                    }))
                  }
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.titleEn")}
                </label>
                <input
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={projectForm.titleEn}
                  onChange={(event) =>
                    setProjectForm((prev) => ({
                      ...prev,
                      titleEn: event.target.value,
                    }))
                  }
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.descriptionJa")}
                </label>
                <textarea
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  rows={3}
                  value={projectForm.descriptionJa}
                  onChange={(event) =>
                    setProjectForm((prev) => ({
                      ...prev,
                      descriptionJa: event.target.value,
                    }))
                  }
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.descriptionEn")}
                </label>
                <textarea
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  rows={3}
                  value={projectForm.descriptionEn}
                  onChange={(event) =>
                    setProjectForm((prev) => ({
                      ...prev,
                      descriptionEn: event.target.value,
                    }))
                  }
                />
              </div>
            </div>
            {renderProjectTechSection()}
            <div>
              <label className="block text-sm font-medium text-slate-700">
                {t("fields.linkUrl")}
              </label>
              <input
                className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                value={projectForm.linkUrl}
                onChange={(event) =>
                  setProjectForm((prev) => ({
                    ...prev,
                    linkUrl: event.target.value,
                  }))
                }
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-700">
                {t("fields.year")}
              </label>
              <input
                type="number"
                className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                value={projectForm.year}
                onChange={(event) =>
                  setProjectForm((prev) => ({
                    ...prev,
                    year: event.target.value,
                  }))
                }
              />
            </div>
            <div className="flex items-center gap-2">
              <input
                id="project-published"
                type="checkbox"
                className="h-4 w-4 rounded border-slate-300 text-slate-900 focus:ring-slate-900"
                checked={projectForm.published}
                onChange={(event) =>
                  setProjectForm((prev) => ({
                    ...prev,
                    published: event.target.checked,
                  }))
                }
              />
              <label
                htmlFor="project-published"
                className="text-sm font-medium text-slate-700"
              >
                {t("fields.published")}
              </label>
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-700">
                {t("fields.sortOrder")}
              </label>
              <input
                type="number"
                className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                value={projectForm.sortOrder}
                onChange={(event) =>
                  setProjectForm((prev) => ({
                    ...prev,
                    sortOrder: event.target.value,
                  }))
                }
              />
            </div>
            <div className="flex items-center justify-end gap-3">
              <button
                type="button"
                className="rounded-md border border-slate-200 px-3 py-2 text-sm text-slate-600 transition hover:bg-slate-100"
                onClick={() => setShowProjectPreview((prev) => !prev)}
              >
                {showProjectPreview
                  ? t("actions.hidePreview")
                  : t("actions.runPreview")}
              </button>
              <button
                type="submit"
                className="inline-flex items-center justify-center rounded-md bg-slate-900 px-4 py-2 text-sm font-medium text-white shadow-sm transition hover:bg-slate-800 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-slate-900"
                disabled={loading}
              >
                {t(editingProjectId != null ? "actions.update" : "actions.create")}
              </button>
            </div>
            {showProjectPreview ? (
              <pre className="mt-4 overflow-x-auto rounded-md border border-slate-200 bg-slate-900/80 p-4 text-xs text-slate-100 dark:border-slate-700">
                {JSON.stringify(projectFormToPayload(projectForm), null, 2)}
              </pre>
            ) : null}
          </form>

          <div className="mt-6 space-y-4">
            {projects.map((project) => {
              const techLabels = project.tech
                .map((membership) => {
                  if (membership.tech?.displayName) {
                    return membership.tech.displayName;
                  }
                  const fallback = techCatalog.find(
                    (entry) => entry.id === membership.tech?.id,
                  );
                  return fallback?.displayName ?? `#${membership.tech?.id ?? ""}`;
                })
                .filter((label) => label != null && label !== "");

              return (
                <div
                  key={project.id}
                  className="rounded-md border border-slate-200 p-4"
                >
                  <div className="flex flex-col gap-2 sm:flex-row sm:items-start sm:justify-between">
                    <div>
                      <h3 className="text-base font-semibold text-slate-900">
                        {project.title.ja ||
                          project.title.en ||
                          t("projects.untitled")}
                      </h3>
                      <p className="text-sm text-slate-600">
                        {project.description.ja ||
                          project.description.en ||
                          t("projects.noDescription")}
                      </p>
                      <div className="mt-2 flex flex-wrap gap-2 text-xs text-slate-500">
                        <span>
                          {t("fields.year")}: {project.year}
                        </span>
                        {techLabels.length > 0 && (
                          <span>
                            {t("fields.tech")}: {techLabels.join(", ")}
                          </span>
                        )}
                        <span>
                          {t("fields.published")}:{" "}
                          {project.published
                            ? t("status.published")
                            : t("status.draft")}
                        </span>
                      </div>
                    </div>
                    <div className="flex gap-2">
                      <button
                        type="button"
                        className="rounded-md border border-slate-200 px-3 py-2 text-sm text-slate-600 transition hover:bg-slate-100"
                        onClick={() => toggleProjectPublished(project)}
                      >
                        {project.published
                          ? t("actions.unpublish")
                          : t("actions.publish")}
                      </button>
                      <button
                        type="button"
                        className="rounded-md border border-slate-200 px-3 py-2 text-sm text-slate-600 transition hover:bg-slate-100"
                        onClick={() => handleEditProject(project)}
                      >
                        {t("actions.edit")}
                      </button>
                      <button
                        type="button"
                        className="rounded-md border border-red-200 px-3 py-2 text-sm text-red-600 transition hover:bg-red-50"
                        onClick={() => deleteProject(project.id)}
                      >
                        {t("actions.delete")}
                      </button>
                    </div>
                  </div>
                  {project.linkUrl && (
                    <a
                      className="mt-3 inline-block text-sm font-medium text-slate-700 underline hover:text-slate-900"
                      href={project.linkUrl}
                      target="_blank"
                      rel="noreferrer"
                    >
                      {project.linkUrl}
                    </a>
                  )}
                </div>
              );
            })}
            {projects.length === 0 && (
              <p className="text-sm text-slate-500">{t("projects.empty")}</p>
            )}
          </div>
        </section>

        <section className="rounded-lg border border-slate-200 bg-white p-6 shadow-sm">
          <div className="flex flex-col gap-3 border-b border-slate-200 pb-4 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <h2 className="text-lg font-semibold text-slate-800">
                {t("research.title")}
              </h2>
              <p className="text-sm text-slate-600">
                {t("research.description")}
              </p>
            </div>
            {editingResearchId != null && (
              <button
                type="button"
                className="text-sm font-medium text-slate-600 hover:text-slate-800"
                onClick={resetResearchForm}
              >
                {t("actions.cancel")}
              </button>
            )}
          </div>
          <form className="mt-4 space-y-5" onSubmit={handleSubmitResearch}>
            <div className="grid gap-4 md:grid-cols-2">
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.slug") ?? "Slug"}
                </label>
                <input
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={researchForm.slug}
                  onChange={(event) =>
                    setResearchForm((prev) => ({
                      ...prev,
                      slug: event.target.value,
                    }))
                  }
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.kind") ?? "Kind"}
                </label>
                <select
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={researchForm.kind}
                  onChange={(event) =>
                    setResearchForm((prev) => ({
                      ...prev,
                      kind: event.target.value as ResearchKind,
                    }))
                  }
                >
                  {researchKinds.map((kind) => (
                    <option key={kind} value={kind}>
                      {kind}
                    </option>
                  ))}
                </select>
              </div>
            </div>

            <div className="grid gap-4 md:grid-cols-2">
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.publishedAt") ?? "Published at"}
                </label>
                <input
                  type="datetime-local"
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={researchForm.publishedAt}
                  onChange={(event) =>
                    setResearchForm((prev) => ({
                      ...prev,
                      publishedAt: event.target.value,
                    }))
                  }
                />
              </div>
              <div className="flex items-center gap-2">
                <input
                  id="research-draft"
                  type="checkbox"
                  className="h-4 w-4 rounded border-slate-300 text-slate-900 focus:ring-slate-900"
                  checked={researchForm.isDraft}
                  onChange={(event) =>
                    setResearchForm((prev) => ({
                      ...prev,
                      isDraft: event.target.checked,
                    }))
                  }
                />
                <label
                  htmlFor="research-draft"
                  className="text-sm font-medium text-slate-700"
                >
                  {t("fields.draft") ?? "Save as draft"}
                </label>
              </div>
            </div>

            <div className="grid gap-4 md:grid-cols-2">
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.titleJa")}
                </label>
                <input
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={researchForm.titleJa}
                  onChange={(event) =>
                    setResearchForm((prev) => ({
                      ...prev,
                      titleJa: event.target.value,
                    }))
                  }
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.titleEn")}
                </label>
                <input
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={researchForm.titleEn}
                  onChange={(event) =>
                    setResearchForm((prev) => ({
                      ...prev,
                      titleEn: event.target.value,
                    }))
                  }
                />
              </div>
            </div>

            <div className="grid gap-4 md:grid-cols-2">
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.overviewJa") ?? "Overview (JA)"}
                </label>
                <textarea
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  rows={3}
                  value={researchForm.overviewJa}
                  onChange={(event) =>
                    setResearchForm((prev) => ({
                      ...prev,
                      overviewJa: event.target.value,
                    }))
                  }
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.overviewEn") ?? "Overview (EN)"}
                </label>
                <textarea
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  rows={3}
                  value={researchForm.overviewEn}
                  onChange={(event) =>
                    setResearchForm((prev) => ({
                      ...prev,
                      overviewEn: event.target.value,
                    }))
                  }
                />
              </div>
            </div>

            <div className="grid gap-4 md:grid-cols-2">
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.outcomeJa") ?? "Outcome (JA)"}
                </label>
                <textarea
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  rows={3}
                  value={researchForm.outcomeJa}
                  onChange={(event) =>
                    setResearchForm((prev) => ({
                      ...prev,
                      outcomeJa: event.target.value,
                    }))
                  }
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.outcomeEn") ?? "Outcome (EN)"}
                </label>
                <textarea
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  rows={3}
                  value={researchForm.outcomeEn}
                  onChange={(event) =>
                    setResearchForm((prev) => ({
                      ...prev,
                      outcomeEn: event.target.value,
                    }))
                  }
                />
              </div>
            </div>

            <div className="grid gap-4 md:grid-cols-2">
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.outlookJa") ?? "Outlook (JA)"}
                </label>
                <textarea
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  rows={3}
                  value={researchForm.outlookJa}
                  onChange={(event) =>
                    setResearchForm((prev) => ({
                      ...prev,
                      outlookJa: event.target.value,
                    }))
                  }
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.outlookEn") ?? "Outlook (EN)"}
                </label>
                <textarea
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  rows={3}
                  value={researchForm.outlookEn}
                  onChange={(event) =>
                    setResearchForm((prev) => ({
                      ...prev,
                      outlookEn: event.target.value,
                    }))
                  }
                />
              </div>
            </div>

            <div className="grid gap-4 md:grid-cols-2">
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.externalUrl") ?? "External URL"}
                </label>
                <input
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={researchForm.externalUrl}
                  onChange={(event) =>
                    setResearchForm((prev) => ({
                      ...prev,
                      externalUrl: event.target.value,
                    }))
                  }
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.highlightImageUrl") ?? "Highlight image URL"}
                </label>
                <input
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={researchForm.highlightImageUrl}
                  onChange={(event) =>
                    setResearchForm((prev) => ({
                      ...prev,
                      highlightImageUrl: event.target.value,
                    }))
                  }
                />
              </div>
            </div>

            <div className="grid gap-4 md:grid-cols-2">
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.imageAltJa") ?? "Image alt (JA)"}
                </label>
                <input
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={researchForm.imageAltJa}
                  onChange={(event) =>
                    setResearchForm((prev) => ({
                      ...prev,
                      imageAltJa: event.target.value,
                    }))
                  }
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.imageAltEn") ?? "Image alt (EN)"}
                </label>
                <input
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={researchForm.imageAltEn}
                  onChange={(event) =>
                    setResearchForm((prev) => ({
                      ...prev,
                      imageAltEn: event.target.value,
                    }))
                  }
                />
              </div>
            </div>

            <div>
              <div className="flex items-center justify-between">
                <label className="text-sm font-medium text-slate-700">
                  {t("fields.tags") ?? "Tags"}
                </label>
                <button
                  type="button"
                  className="text-sm font-medium text-slate-600 hover:text-slate-800"
                  onClick={handleAddResearchTag}
                >
                  {t("actions.addTag") ?? "Add tag"}
                </button>
              </div>
              <div className="mt-2 space-y-3">
                {researchForm.tags.length === 0 ? (
                  <p className="text-sm text-slate-500">
                    {t("research.noTags") ?? "No tags yet."}
                  </p>
                ) : (
                  researchForm.tags.map((tag, index) => (
                    <div
                      key={`tag-${index}`}
                      className="grid gap-3 md:grid-cols-[1fr,120px,auto]"
                    >
                      <input
                        className="rounded-md border border-slate-200 px-3 py-2 text-sm"
                        placeholder="value"
                        value={tag.value}
                        onChange={(event) =>
                          handleUpdateResearchTag(index, {
                            value: event.target.value,
                          })
                        }
                      />
                      <input
                        type="number"
                        className="rounded-md border border-slate-200 px-3 py-2 text-sm"
                        value={tag.sortOrder}
                        onChange={(event) =>
                          handleUpdateResearchTag(index, {
                            sortOrder: event.target.value,
                          })
                        }
                      />
                      <button
                        type="button"
                        className="rounded-md border border-red-200 px-3 py-2 text-sm text-red-600 transition hover:bg-red-50"
                        onClick={() => handleRemoveResearchTag(index)}
                      >
                        {t("actions.remove") ?? "Remove"}
                      </button>
                    </div>
                  ))
                )}
              </div>
            </div>

            <div>
              <div className="flex items-center justify-between">
                <label className="text-sm font-medium text-slate-700">
                  {t("fields.links") ?? "Links"}
                </label>
                <button
                  type="button"
                  className="text-sm font-medium text-slate-600 hover:text-slate-800"
                  onClick={handleAddResearchLink}
                >
                  {t("actions.addLink") ?? "Add link"}
                </button>
              </div>
              <div className="mt-2 space-y-3">
                {researchForm.links.length === 0 ? (
                  <p className="text-sm text-slate-500">
                    {t("research.noLinks") ?? "No links yet."}
                  </p>
                ) : (
                  researchForm.links.map((link, index) => (
                    <div
                      key={`link-${index}`}
                      className="space-y-3 rounded-md border border-slate-200 p-3"
                    >
                      <div className="grid gap-3 md:grid-cols-[160px,1fr,auto]">
                        <select
                          className="rounded-md border border-slate-200 px-3 py-2 text-sm"
                          value={link.type}
                          onChange={(event) =>
                            handleUpdateResearchLink(index, {
                              type: event.target.value as ResearchLinkType,
                            })
                          }
                        >
                          {researchLinkTypes.map((type) => (
                            <option key={type} value={type}>
                              {type}
                            </option>
                          ))}
                        </select>
                        <input
                          className="rounded-md border border-slate-200 px-3 py-2 text-sm"
                          placeholder="https://"
                          value={link.url}
                          onChange={(event) =>
                            handleUpdateResearchLink(index, {
                              url: event.target.value,
                            })
                          }
                        />
                        <div className="flex items-center justify-end gap-2">
                          <input
                            type="number"
                            className="w-24 rounded-md border border-slate-200 px-3 py-2 text-sm"
                            value={link.sortOrder}
                            onChange={(event) =>
                              handleUpdateResearchLink(index, {
                                sortOrder: event.target.value,
                              })
                            }
                          />
                          <button
                            type="button"
                            className="rounded-md border border-red-200 px-3 py-2 text-sm text-red-600 transition hover:bg-red-50"
                            onClick={() => handleRemoveResearchLink(index)}
                          >
                            {t("actions.remove") ?? "Remove"}
                          </button>
                        </div>
                      </div>
                      <div className="grid gap-3 md:grid-cols-2">
                        <input
                          className="rounded-md border border-slate-200 px-3 py-2 text-sm"
                          placeholder="Label (JA)"
                          value={link.labelJa}
                          onChange={(event) =>
                            handleUpdateResearchLink(index, {
                              labelJa: event.target.value,
                            })
                          }
                        />
                        <input
                          className="rounded-md border border-slate-200 px-3 py-2 text-sm"
                          placeholder="Label (EN)"
                          value={link.labelEn}
                          onChange={(event) =>
                            handleUpdateResearchLink(index, {
                              labelEn: event.target.value,
                            })
                          }
                        />
                      </div>
                    </div>
                  ))
                )}
              </div>
            </div>

            <div>
              <div className="flex items-center justify-between">
                <label className="text-sm font-medium text-slate-700">
                  {t("fields.assets") ?? "Assets"}
                </label>
                <button
                  type="button"
                  className="text-sm font-medium text-slate-600 hover:text-slate-800"
                  onClick={handleAddResearchAsset}
                >
                  {t("actions.addAsset") ?? "Add asset"}
                </button>
              </div>
              <div className="mt-2 space-y-3">
                {researchForm.assets.length === 0 ? (
                  <p className="text-sm text-slate-500">
                    {t("research.noAssets") ?? "No assets yet."}
                  </p>
                ) : (
                  researchForm.assets.map((asset, index) => (
                    <div
                      key={`asset-${index}`}
                      className="space-y-3 rounded-md border border-slate-200 p-3"
                    >
                      <div className="grid gap-3 md:grid-cols-[1fr,120px,auto]">
                        <input
                          className="rounded-md border border-slate-200 px-3 py-2 text-sm"
                          placeholder="https://"
                          value={asset.url}
                          onChange={(event) =>
                            handleUpdateResearchAsset(index, {
                              url: event.target.value,
                            })
                          }
                        />
                        <input
                          type="number"
                          className="rounded-md border border-slate-200 px-3 py-2 text-sm"
                          value={asset.sortOrder}
                          onChange={(event) =>
                            handleUpdateResearchAsset(index, {
                              sortOrder: event.target.value,
                            })
                          }
                        />
                        <button
                          type="button"
                          className="rounded-md border border-red-200 px-3 py-2 text-sm text-red-600 transition hover:bg-red-50"
                          onClick={() => handleRemoveResearchAsset(index)}
                        >
                          {t("actions.remove") ?? "Remove"}
                        </button>
                      </div>
                      <div className="grid gap-3 md:grid-cols-2">
                        <input
                          className="rounded-md border border-slate-200 px-3 py-2 text-sm"
                          placeholder="Caption (JA)"
                          value={asset.captionJa}
                          onChange={(event) =>
                            handleUpdateResearchAsset(index, {
                              captionJa: event.target.value,
                            })
                          }
                        />
                        <input
                          className="rounded-md border border-slate-200 px-3 py-2 text-sm"
                          placeholder="Caption (EN)"
                          value={asset.captionEn}
                          onChange={(event) =>
                            handleUpdateResearchAsset(index, {
                              captionEn: event.target.value,
                            })
                          }
                        />
                      </div>
                    </div>
                  ))
                )}
              </div>
            </div>

            <div>
              <div className="flex items-center justify-between">
                <label className="text-sm font-medium text-slate-700">
                  {t("fields.tech") ?? "Tech relationships"}
                </label>
                <button
                  type="button"
                  className="text-sm font-medium text-slate-600 hover:text-slate-800"
                  onClick={handleAddResearchTech}
                >
                  {t("actions.addTech") ?? "Add tech"}
                </button>
              </div>
              <div className="mt-2 space-y-3">
                {researchForm.tech.length === 0 ? (
                  <p className="text-sm text-slate-500">
                    {t("research.noTech") ?? "No technology relationships yet."}
                  </p>
                ) : (
                  researchForm.tech.map((membership, index) => (
                    <div
                      key={`tech-${index}`}
                      className="space-y-3 rounded-md border border-slate-200 p-3"
                    >
                      <div className="grid gap-3 md:grid-cols-[160px,1fr,auto]">
                        <input
                          className="rounded-md border border-slate-200 px-3 py-2 text-sm"
                          placeholder="Tech ID"
                          value={membership.techId}
                          onChange={(event) =>
                            handleUpdateResearchTech(index, {
                              techId: event.target.value,
                            })
                          }
                        />
                        <select
                          className="rounded-md border border-slate-200 px-3 py-2 text-sm"
                          value={membership.context}
                          onChange={(event) =>
                            handleUpdateResearchTech(index, {
                              context: event.target.value as TechContext,
                            })
                          }
                        >
                          {techContexts.map((context) => (
                            <option key={context} value={context}>
                              {context}
                            </option>
                          ))}
                        </select>
                        <div className="flex items-center justify-end gap-2">
                          <input
                            type="number"
                            className="w-24 rounded-md border border-slate-200 px-3 py-2 text-sm"
                            value={membership.sortOrder}
                            onChange={(event) =>
                              handleUpdateResearchTech(index, {
                                sortOrder: event.target.value,
                              })
                            }
                          />
                          <button
                            type="button"
                            className="rounded-md border border-red-200 px-3 py-2 text-sm text-red-600 transition hover:bg-red-50"
                            onClick={() => handleRemoveResearchTech(index)}
                          >
                            {t("actions.remove") ?? "Remove"}
                          </button>
                        </div>
                      </div>
                      <textarea
                        className="w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                        rows={2}
                        placeholder={t("fields.note") ?? "Note"}
                        value={membership.note}
                        onChange={(event) =>
                          handleUpdateResearchTech(index, {
                            note: event.target.value,
                          })
                        }
                      />
                    </div>
                  ))
                )}
              </div>
            </div>

            <div className="flex items-center justify-end gap-3">
              <button
                type="submit"
                className="inline-flex items-center justify-center rounded-md bg-slate-900 px-4 py-2 text-sm font-medium text-white shadow-sm transition hover:bg-slate-800 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-slate-900"
                disabled={loading}
              >
                {t(editingResearchId != null ? "actions.update" : "actions.create")}
              </button>
            </div>
          </form>

          <div className="mt-6 space-y-4">
            {research.map((item) => (
              <div key={item.id} className="rounded-md border border-slate-200 p-4">
                <div className="flex flex-col gap-2 sm:flex-row sm:items-start sm:justify-between">
                  <div>
                    <h3 className="text-base font-semibold text-slate-900">
                      {item.title.ja || item.title.en || t("research.untitled")}
                    </h3>
                    <p className="text-sm text-slate-600">
                      {item.overview.ja || item.overview.en || t("research.noSummary")}
                    </p>
                    <div className="mt-2 flex flex-wrap gap-2 text-xs text-slate-500">
                      <span>slug: {item.slug}</span>
                      <span>kind: {item.kind}</span>
                      <span>
                        {t("fields.publishedAt") ?? "Published"}:{" "}
                        {new Date(item.publishedAt).toLocaleString()}
                      </span>
                      <span>
                        {item.isDraft
                          ? t("status.draft") ?? "Draft"
                          : t("status.published") ?? "Published"}
                      </span>
                      {item.tags.length > 0 && (
                        <span>{t("fields.tags") ?? "Tags"}: {item.tags.length}</span>
                      )}
                      {item.links.length > 0 && (
                        <span>{t("fields.links") ?? "Links"}: {item.links.length}</span>
                      )}
                      {item.tech.length > 0 && (
                        <span>{t("fields.tech") ?? "Tech"}: {item.tech.length}</span>
                      )}
                    </div>
                    {item.externalUrl && (
                      <a
                        className="mt-2 inline-block text-sm font-medium text-slate-700 underline hover:text-slate-900"
                        href={item.externalUrl}
                        target="_blank"
                        rel="noreferrer"
                      >
                        {item.externalUrl}
                      </a>
                    )}
                  </div>
                  <div className="flex gap-2">
                    <button
                      type="button"
                      className="rounded-md border border-slate-200 px-3 py-2 text-sm text-slate-600 transition hover:bg-slate-100"
                      onClick={() => toggleResearchDraft(item)}
                    >
                      {item.isDraft
                        ? t("actions.publish") ?? "Publish"
                        : t("actions.markDraft") ?? "Mark as draft"}
                    </button>
                    <button
                      type="button"
                      className="rounded-md border border-slate-200 px-3 py-2 text-sm text-slate-600 transition hover:bg-slate-100"
                      onClick={() => handleEditResearch(item)}
                    >
                      {t("actions.edit")}
                    </button>
                    <button
                      type="button"
                      className="rounded-md border border-red-200 px-3 py-2 text-sm text-red-600 transition hover:bg-red-50"
                      onClick={() => deleteResearch(item.id)}
                    >
                      {t("actions.delete")}
                    </button>
                  </div>
                </div>
              </div>
            ))}
            {research.length === 0 && (
              <p className="text-sm text-slate-500">{t("research.empty")}</p>
            )}
          </div>
        </section>

        <section className="rounded-lg border border-slate-200 bg-white p-6 shadow-sm">
          <h2 className="text-lg font-semibold text-slate-800">
            {t("contacts.title")}
          </h2>
          <p className="mt-1 text-sm text-slate-600">
            {t("contacts.description")}
          </p>
          <div className="mt-4 space-y-4">
            {contacts.map((contact) => {
              const edit = contactEdits[contact.id] ?? {
                topic: contact.topic,
                message: contact.message,
                status: contact.status,
                adminNote: contact.adminNote,
              };
              return (
                <div
                  key={contact.id}
                  className="rounded-md border border-slate-200 p-4"
                >
                  <div className="flex flex-col gap-2 sm:flex-row sm:items-start sm:justify-between">
                    <div className="text-sm text-slate-700">
                      <p className="font-semibold text-slate-900">
                        {contact.name} · {contact.email}
                      </p>
                      {contact.topic && (
                        <p className="text-slate-600">{contact.topic}</p>
                      )}
                      <p className="mt-2 whitespace-pre-wrap text-slate-600">
                        {contact.message}
                      </p>
                      <p className="mt-2 text-xs text-slate-500">
                        {t("fields.createdAt")}:
                        {" "}
                        {new Date(contact.createdAt).toLocaleString()}
                      </p>
                      <p className="text-xs text-slate-500">
                        {t("fields.updatedAt")}:
                        {" "}
                        {new Date(contact.updatedAt).toLocaleString()}
                      </p>
                    </div>
                    <div className="flex gap-2">
                      <button
                        type="button"
                        className="rounded-md border border-red-200 px-3 py-2 text-sm text-red-600 transition hover:bg-red-50"
                        onClick={() => handleDeleteContact(contact.id)}
                      >
                        {t("actions.delete")}
                      </button>
                    </div>
                  </div>
                  <div className="mt-4 grid gap-3 md:grid-cols-2">
                    <div>
                      <label className="block text-xs font-medium uppercase text-slate-500">
                        {t("fields.topic")}
                      </label>
                      <input
                        className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                        value={edit.topic}
                        onChange={(event) =>
                          handleContactEditChange(
                            contact.id,
                            "topic",
                            event.target.value,
                          )
                        }
                      />
                    </div>
                    <div>
                      <label className="block text-xs font-medium uppercase text-slate-500">
                        {t("fields.status")}
                      </label>
                      <select
                        className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                        value={edit.status}
                        onChange={(event) =>
                          handleContactEditChange(
                            contact.id,
                            "status",
                            event.target.value,
                          )
                        }
                      >
                        {contactStatuses.map((statusValue) => (
                          <option key={statusValue} value={statusValue}>
                            {t(`contacts.status.${statusValue}`)}
                          </option>
                        ))}
                      </select>
                    </div>
                    <div className="md:col-span-2">
                      <label className="block text-xs font-medium uppercase text-slate-500">
                        {t("fields.adminNote")}
                      </label>
                      <textarea
                        className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                        rows={2}
                        value={edit.adminNote}
                        onChange={(event) =>
                          handleContactEditChange(
                            contact.id,
                            "adminNote",
                            event.target.value,
                          )
                        }
                      />
                    </div>
                    <div className="md:col-span-2">
                      <label className="block text-xs font-medium uppercase text-slate-500">
                        {t("fields.message")}
                      </label>
                      <textarea
                        className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                        rows={3}
                        value={edit.message}
                        onChange={(event) =>
                          handleContactEditChange(
                            contact.id,
                            "message",
                            event.target.value,
                          )
                        }
                      />
                    </div>
                  </div>
                  <div className="mt-4 flex items-center justify-end gap-3">
                    <button
                      type="button"
                      className="rounded-md border border-slate-200 px-3 py-2 text-sm text-slate-600 transition hover:bg-slate-100"
                      onClick={() => handleResetContact(contact)}
                    >
                      {t("actions.reset")}
                    </button>
                    <button
                      type="button"
                      className="inline-flex items-center justify-center rounded-md bg-slate-900 px-4 py-2 text-sm font-medium text-white shadow-sm transition hover:bg-slate-800 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-slate-900"
                      onClick={() => handleSaveContact(contact.id)}
                      disabled={loading}
                    >
                      {t("actions.save")}
                    </button>
                  </div>
                </div>
              );
            })}
            {contacts.length === 0 && (
              <p className="text-sm text-slate-500">{t("contacts.empty")}</p>
            )}
          </div>
        </section>

        <section className="rounded-lg border border-slate-200 bg-white p-6 shadow-sm">
          <div className="flex flex-col gap-3 border-b border-slate-200 pb-4 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <h2 className="text-lg font-semibold text-slate-800">
                {t("blacklist.title")}
              </h2>
              <p className="text-sm text-slate-600">
                {t("blacklist.description")}
              </p>
            </div>
            {editingBlacklistId != null && (
              <button
                type="button"
                className="text-sm font-medium text-slate-600 hover:text-slate-800"
                onClick={resetBlacklistForm}
              >
                {t("actions.cancel")}
              </button>
            )}
          </div>
          <form className="mt-4 space-y-4" onSubmit={handleSubmitBlacklist}>
            <div className="grid gap-4 md:grid-cols-2">
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.email")}
                </label>
                <input
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={blacklistForm.email}
                  onChange={(event) =>
                    setBlacklistForm((prev) => ({
                      ...prev,
                      email: event.target.value,
                    }))
                  }
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.reason")}
                </label>
                <input
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={blacklistForm.reason}
                  onChange={(event) =>
                    setBlacklistForm((prev) => ({
                      ...prev,
                      reason: event.target.value,
                    }))
                  }
                />
              </div>
            </div>
            <div className="flex items-center justify-end gap-3">
              <button
                type="submit"
                className="inline-flex items-center justify-center rounded-md bg-slate-900 px-4 py-2 text-sm font-medium text-white shadow-sm transition hover:bg-slate-800 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-slate-900"
                disabled={loading}
              >
                {t(editingBlacklistId != null ? "actions.update" : "actions.create")}
              </button>
            </div>
          </form>

          <div className="mt-6 space-y-4">
            {blacklist.map((entry) => (
              <div key={entry.id} className="rounded-md border border-slate-200 p-4">
                <div className="flex flex-col gap-2 sm:flex-row sm:items-start sm:justify-between">
                  <div>
                    <p className="font-semibold text-slate-900">{entry.email}</p>
                    <p className="text-sm text-slate-600">{entry.reason}</p>
                    <p className="mt-2 text-xs text-slate-500">
                      {t("fields.createdAt")}:
                      {" "}
                      {new Date(entry.createdAt).toLocaleString()}
                    </p>
                  </div>
                  <div className="flex gap-2">
                    <button
                      type="button"
                      className="rounded-md border border-slate-200 px-3 py-2 text-sm text-slate-600 transition hover:bg-slate-100"
                      onClick={() => handleEditBlacklist(entry)}
                    >
                      {t("actions.edit")}
                    </button>
                    <button
                      type="button"
                      className="rounded-md border border-red-200 px-3 py-2 text-sm text-red-600 transition hover:bg-red-50"
                      onClick={() => deleteBlacklistEntry(entry.id)}
                    >
                      {t("actions.delete")}
                    </button>
                  </div>
                </div>
              </div>
            ))}
            {blacklist.length === 0 && (
              <p className="text-sm text-slate-500">{t("blacklist.empty")}</p>
            )}
          </div>
        </section>
      </main>

      {showContactSettingsDiff && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/40 p-4">
          <div className="max-h-full w-full max-w-4xl overflow-y-auto rounded-lg bg-white p-6 shadow-xl">
            <div className="flex items-start justify-between gap-4">
              <h3 className="text-lg font-semibold text-slate-900">
                {t("contactSettings.diff.title")}
              </h3>
              <button
                type="button"
                className="rounded-md border border-slate-200 px-3 py-1 text-sm text-slate-600 transition hover:bg-slate-100"
                onClick={() => setShowContactSettingsDiff(false)}
              >
                {t("actions.close")}
              </button>
            </div>
            {contactSettingsDiffSummary.fields.length === 0 &&
            contactSettingsDiffSummary.topics.length === 0 ? (
              <p className="mt-4 text-sm text-slate-600">
                {t("contactSettings.diff.noChanges")}
              </p>
            ) : (
              <div className="mt-4 space-y-6">
                {contactSettingsDiffSummary.fields.length > 0 && (
                  <div>
                    <h4 className="text-sm font-semibold text-slate-700">
                      {t("contactSettings.diff.section.fields")}
                    </h4>
                    <div className="mt-3 divide-y rounded-md border border-slate-200">
                      {contactSettingsDiffSummary.fields.map((entry) => (
                        <div
                          key={entry.key}
                          className="grid gap-4 px-4 py-3 md:grid-cols-2"
                        >
                          <div>
                            <p className="text-xs uppercase tracking-wide text-slate-500">
                              {t(entry.labelKey)}
                            </p>
                            <p className="mt-1 whitespace-pre-wrap text-sm text-slate-600">
                              {entry.original ||
                                t("contactSettings.diff.emptyValue")}
                            </p>
                          </div>
                          <div>
                            <p className="text-xs uppercase tracking-wide text-slate-500">
                              {t("contactSettings.diff.updated")}
                            </p>
                            <p className="mt-1 whitespace-pre-wrap text-sm text-slate-800">
                              {entry.updated ||
                                t("contactSettings.diff.emptyValue")}
                            </p>
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                )}
                {contactSettingsDiffSummary.topics.length > 0 && (
                  <div>
                    <h4 className="text-sm font-semibold text-slate-700">
                      {t("contactSettings.diff.section.topics")}
                    </h4>
                    <div className="mt-3 space-y-4">
                      {contactSettingsDiffSummary.topics.map((topic) => (
                        <div
                          key={`${topic.change}-${topic.id}`}
                          className="rounded-md border border-slate-200 p-4"
                        >
                          <div className="flex flex-wrap items-center justify-between gap-2">
                            <span className="text-sm font-semibold text-slate-800">
                              {topic.id ||
                                t("contactSettings.topics.placeholderId")}
                            </span>
                            <span
                              className={`rounded-full px-2 py-1 text-xs font-medium ${
                                topic.change === "added"
                                  ? "bg-emerald-50 text-emerald-600"
                                  : topic.change === "removed"
                                    ? "bg-rose-50 text-rose-600"
                                    : "bg-amber-50 text-amber-600"
                              }`}
                            >
                              {t(
                                `contactSettings.diff.topic.${topic.change}` as const,
                              )}
                            </span>
                          </div>
                          {topic.original && (
                            <div className="mt-3 grid gap-2 md:grid-cols-2">
                              <div>
                                <p className="text-xs uppercase tracking-wide text-slate-500">
                                  {t("contactSettings.diff.previousLabel")}
                                </p>
                                <p className="text-sm text-slate-700">
                                  {topic.original.label.ja ||
                                    t("contactSettings.diff.emptyValue")}
                                </p>
                                <p className="text-xs text-slate-500">
                                  EN:{" "}
                                  {topic.original.label.en ||
                                    t("contactSettings.diff.emptyValue")}
                                </p>
                              </div>
                              <div>
                                <p className="text-xs uppercase tracking-wide text-slate-500">
                                  {t("contactSettings.diff.previousDescription")}
                                </p>
                                <p className="text-sm text-slate-700">
                                  {topic.original.description.ja ||
                                    t("contactSettings.diff.emptyValue")}
                                </p>
                                <p className="text-xs text-slate-500">
                                  EN:{" "}
                                  {topic.original.description.en ||
                                    t("contactSettings.diff.emptyValue")}
                                </p>
                              </div>
                            </div>
                          )}
                          {topic.updated && (
                            <div className="mt-3 grid gap-2 md:grid-cols-2">
                              <div>
                                <p className="text-xs uppercase tracking-wide text-slate-500">
                                  {t("contactSettings.diff.updatedLabel")}
                                </p>
                                <p className="text-sm text-slate-700">
                                  {topic.updated.label.ja ||
                                    t("contactSettings.diff.emptyValue")}
                                </p>
                                <p className="text-xs text-slate-500">
                                  EN:{" "}
                                  {topic.updated.label.en ||
                                    t("contactSettings.diff.emptyValue")}
                                </p>
                              </div>
                              <div>
                                <p className="text-xs uppercase tracking-wide text-slate-500">
                                  {t("contactSettings.diff.updatedDescription")}
                                </p>
                                <p className="text-sm text-slate-700">
                                  {topic.updated.description.ja ||
                                    t("contactSettings.diff.emptyValue")}
                                </p>
                                <p className="text-xs text-slate-500">
                                  EN:{" "}
                                  {topic.updated.description.en ||
                                    t("contactSettings.diff.emptyValue")}
                                </p>
                              </div>
                            </div>
                          )}
                        </div>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
}

export default App;
