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
  ContactMessage,
  ContactStatus,
  LocalizedText,
  ResearchKind,
  ResearchLinkType,
  TechCatalogEntry,
  TechContext,
} from "./types";

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

type ProfileFormState = {
  name: LocalizedText;
  title: LocalizedText;
  affiliation: LocalizedText;
  lab: LocalizedText;
  summary: LocalizedText;
  skills: LocalizedText[];
  focusAreas: LocalizedText[];
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

type AuthState = "checking" | "authenticated" | "unauthorized";

const isUnauthorizedError = (error: unknown): boolean =>
  error instanceof DomainError && error.status === 401;

const createEmptyProfileForm = (): ProfileFormState => ({
  name: { ja: "", en: "" },
  title: { ja: "", en: "" },
  affiliation: { ja: "", en: "" },
  lab: { ja: "", en: "" },
  summary: { ja: "", en: "" },
  skills: [{ ja: "", en: "" }],
  focusAreas: [{ ja: "", en: "" }],
});

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

function profileToForm(profile: AdminProfile | null): ProfileFormState {
  if (!profile) {
    return createEmptyProfileForm();
  }
  const skills = profile.skills.length > 0 ? profile.skills : [{ ja: "", en: "" }];
  const focusAreas =
    profile.focusAreas.length > 0
      ? profile.focusAreas
      : [{ ja: "", en: "" }];
  return {
    name: { ...profile.name },
    title: { ...profile.title },
    affiliation: { ...profile.affiliation },
    lab: { ...profile.lab },
    summary: { ...profile.summary },
    skills: skills.map((item) => ({ ...item })),
    focusAreas: focusAreas.map((item) => ({ ...item })),
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
    setProfileForm(createEmptyProfileForm());
    setProjectForm({ ...emptyProjectForm });
    setProjectTechSearch("");
    setShowProjectPreview(false);
    setResearchForm(createEmptyResearchForm());
    setBlacklistForm({ ...emptyBlacklistForm });
    setEditingProjectId(null);
    setEditingResearchId(null);
    setEditingBlacklistId(null);
    setError(null);
  }, [clearSession]);

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
      ] = await Promise.all([
        adminApi.fetchSummary(),
        adminApi.getProfile(),
        adminApi.listProjects(),
        adminApi.listResearch(),
        adminApi.listContacts(),
        adminApi.listBlacklist(),
        adminApi.listTechCatalog({ includeInactive: false }),
      ]);

      setSummary(summaryRes.data);
      setProfileForm(profileToForm(profileRes.data));
      setProjects(projectRes.data);
      setResearch(researchRes.data);
      setContacts(contactRes.data);
      setContactEdits(buildContactEditMap(contactRes.data));
      setBlacklist(blacklistRes.data);
      setTechCatalog(techCatalogRes.data);
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
    const payload = {
      name: { ...profileForm.name },
      title: { ...profileForm.title },
      affiliation: { ...profileForm.affiliation },
      lab: { ...profileForm.lab },
      summary: { ...profileForm.summary },
      skills: profileForm.skills.map((item) => ({ ...item })),
      focusAreas: profileForm.focusAreas.map((item) => ({ ...item })),
    };
    await run(async () => {
      await adminApi.updateProfile(payload);
    });
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
          <form className="mt-4 space-y-4" onSubmit={handleSaveProfile}>
            <div className="grid gap-4 md:grid-cols-2">
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.nameJa")}
                </label>
                <input
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={profileForm.name.ja ?? ""}
                  onChange={(event) =>
                    setProfileForm((prev) => ({
                      ...prev,
                      name: { ...prev.name, ja: event.target.value },
                    }))
                  }
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.nameEn")}
                </label>
                <input
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={profileForm.name.en ?? ""}
                  onChange={(event) =>
                    setProfileForm((prev) => ({
                      ...prev,
                      name: { ...prev.name, en: event.target.value },
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
                  value={profileForm.title.ja ?? ""}
                  onChange={(event) =>
                    setProfileForm((prev) => ({
                      ...prev,
                      title: { ...prev.title, ja: event.target.value },
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
                  value={profileForm.title.en ?? ""}
                  onChange={(event) =>
                    setProfileForm((prev) => ({
                      ...prev,
                      title: { ...prev.title, en: event.target.value },
                    }))
                  }
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.affiliationJa")}
                </label>
                <input
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={profileForm.affiliation.ja ?? ""}
                  onChange={(event) =>
                    setProfileForm((prev) => ({
                      ...prev,
                      affiliation: {
                        ...prev.affiliation,
                        ja: event.target.value,
                      },
                    }))
                  }
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.affiliationEn")}
                </label>
                <input
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={profileForm.affiliation.en ?? ""}
                  onChange={(event) =>
                    setProfileForm((prev) => ({
                      ...prev,
                      affiliation: {
                        ...prev.affiliation,
                        en: event.target.value,
                      },
                    }))
                  }
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.labJa")}
                </label>
                <input
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={profileForm.lab.ja ?? ""}
                  onChange={(event) =>
                    setProfileForm((prev) => ({
                      ...prev,
                      lab: { ...prev.lab, ja: event.target.value },
                    }))
                  }
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.labEn")}
                </label>
                <input
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={profileForm.lab.en ?? ""}
                  onChange={(event) =>
                    setProfileForm((prev) => ({
                      ...prev,
                      lab: { ...prev.lab, en: event.target.value },
                    }))
                  }
                />
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-slate-700">
                {t("fields.summaryJa")}
              </label>
              <textarea
                className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                rows={3}
                value={profileForm.summary.ja ?? ""}
                onChange={(event) =>
                  setProfileForm((prev) => ({
                    ...prev,
                    summary: { ...prev.summary, ja: event.target.value },
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
                rows={3}
                value={profileForm.summary.en ?? ""}
                onChange={(event) =>
                  setProfileForm((prev) => ({
                    ...prev,
                    summary: { ...prev.summary, en: event.target.value },
                  }))
                }
              />
            </div>

            <div>
              <div className="flex items-center justify-between">
                <h3 className="text-sm font-semibold text-slate-700">
                  {t("profile.skills.title")}
                </h3>
                <button
                  type="button"
                  className="text-sm font-medium text-slate-600 hover:text-slate-800"
                  onClick={() =>
                    setProfileForm((prev) => ({
                      ...prev,
                      skills: [...prev.skills, { ja: "", en: "" }],
                    }))
                  }
                >
                  {t("profile.skills.add")}
                </button>
              </div>
              <div className="mt-2 space-y-2">
                {profileForm.skills.map((skill, index) => (
                  <div
                    className="flex flex-col gap-2 rounded-md border border-slate-200 p-3 md:flex-row"
                    key={`skill-${index}`}
                  >
                    <input
                      className="flex-1 rounded-md border border-slate-200 px-3 py-2 text-sm"
                      placeholder={t("fields.skillJa")}
                      value={skill.ja ?? ""}
                      onChange={(event) =>
                        setProfileForm((prev) => ({
                          ...prev,
                          skills: prev.skills.map((item, idx) =>
                            idx === index
                              ? { ...item, ja: event.target.value }
                              : item,
                          ),
                        }))
                      }
                    />
                    <input
                      className="flex-1 rounded-md border border-slate-200 px-3 py-2 text-sm"
                      placeholder={t("fields.skillEn")}
                      value={skill.en ?? ""}
                      onChange={(event) =>
                        setProfileForm((prev) => ({
                          ...prev,
                          skills: prev.skills.map((item, idx) =>
                            idx === index
                              ? { ...item, en: event.target.value }
                              : item,
                          ),
                        }))
                      }
                    />
                    <button
                      type="button"
                      className="rounded-md border border-slate-200 px-3 py-2 text-sm text-slate-600 transition hover:bg-slate-100 disabled:cursor-not-allowed disabled:opacity-50"
                      onClick={() =>
                        setProfileForm((prev) => ({
                          ...prev,
                          skills: prev.skills.filter((_, idx) => idx !== index),
                        }))
                      }
                      disabled={profileForm.skills.length === 1}
                    >
                      {t("actions.remove")}
                    </button>
                  </div>
                ))}
              </div>
            </div>

            <div>
              <div className="flex items-center justify-between">
                <h3 className="text-sm font-semibold text-slate-700">
                  {t("profile.focusAreas.title")}
                </h3>
                <button
                  type="button"
                  className="text-sm font-medium text-slate-600 hover:text-slate-800"
                  onClick={() =>
                    setProfileForm((prev) => ({
                      ...prev,
                      focusAreas: [...prev.focusAreas, { ja: "", en: "" }],
                    }))
                  }
                >
                  {t("profile.focusAreas.add")}
                </button>
              </div>
              <div className="mt-2 space-y-2">
                {profileForm.focusAreas.map((area, index) => (
                  <div
                    className="flex flex-col gap-2 rounded-md border border-slate-200 p-3 md:flex-row"
                    key={`focus-${index}`}
                  >
                    <input
                      className="flex-1 rounded-md border border-slate-200 px-3 py-2 text-sm"
                      placeholder={t("fields.focusJa")}
                      value={area.ja ?? ""}
                      onChange={(event) =>
                        setProfileForm((prev) => ({
                          ...prev,
                          focusAreas: prev.focusAreas.map((item, idx) =>
                            idx === index
                              ? { ...item, ja: event.target.value }
                              : item,
                          ),
                        }))
                      }
                    />
                    <input
                      className="flex-1 rounded-md border border-slate-200 px-3 py-2 text-sm"
                      placeholder={t("fields.focusEn")}
                      value={area.en ?? ""}
                      onChange={(event) =>
                        setProfileForm((prev) => ({
                          ...prev,
                          focusAreas: prev.focusAreas.map((item, idx) =>
                            idx === index
                              ? { ...item, en: event.target.value }
                              : item,
                          ),
                        }))
                      }
                    />
                    <button
                      type="button"
                      className="rounded-md border border-slate-200 px-3 py-2 text-sm text-slate-600 transition hover:bg-slate-100 disabled:cursor-not-allowed disabled:opacity-50"
                      onClick={() =>
                        setProfileForm((prev) => ({
                          ...prev,
                          focusAreas: prev.focusAreas.filter(
                            (_, idx) => idx !== index,
                          ),
                        }))
                      }
                      disabled={profileForm.focusAreas.length === 1}
                    >
                      {t("actions.remove")}
                    </button>
                  </div>
                ))}
              </div>
            </div>

            <div className="flex items-center justify-end gap-3">
              <button
                type="submit"
                className="inline-flex items-center justify-center rounded-md bg-slate-900 px-4 py-2 text-sm font-medium text-white shadow-sm transition hover:bg-slate-800 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-slate-900"
                disabled={loading}
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
    </div>
  );
}

export default App;
