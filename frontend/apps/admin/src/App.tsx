import { FormEvent, useCallback, useEffect, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";

import { adminApi, DomainError } from "./modules/admin-api";
import {
  extractTokenFromHash,
  getToken,
  useAuthSession,
} from "./modules/auth-session";
import type {
  AdminProfile,
  AdminProject,
  AdminResearch,
  AdminSummary,
  BlacklistEntry,
  ContactMessage,
  ContactStatus,
  LocalizedText,
} from "./types";

const currentYear = new Date().getFullYear();
const contactStatuses: ContactStatus[] = [
  "pending",
  "in_review",
  "resolved",
  "archived",
];

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
  techStack: string;
  linkUrl: string;
  year: string;
  published: boolean;
  sortOrder: string;
};

type ResearchFormState = {
  titleJa: string;
  titleEn: string;
  summaryJa: string;
  summaryEn: string;
  contentJa: string;
  contentEn: string;
  year: string;
  published: boolean;
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
  techStack: "",
  linkUrl: "",
  year: String(currentYear),
  published: false,
  sortOrder: "",
};

const emptyResearchForm: ResearchFormState = {
  titleJa: "",
  titleEn: "",
  summaryJa: "",
  summaryEn: "",
  contentJa: "",
  contentEn: "",
  year: String(currentYear),
  published: false,
};

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

function App() {
  const { t } = useTranslation();
  const { token, setToken: storeToken, clearToken } = useAuthSession();

  const [authState, setAuthState] = useState<AuthState>("checking");
  const [status, setStatus] = useState("unknown");
  const [summary, setSummary] = useState<AdminSummary | null>(null);
  const [projects, setProjects] = useState<AdminProject[]>([]);
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
  const [researchForm, setResearchForm] = useState<ResearchFormState>({
    ...emptyResearchForm,
  });
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

  const handleUnauthorized = useCallback(() => {
    clearToken();
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
    setResearchForm({ ...emptyResearchForm });
    setBlacklistForm({ ...emptyBlacklistForm });
    setEditingProjectId(null);
    setEditingResearchId(null);
    setEditingBlacklistId(null);
    setError(null);
  }, [clearToken]);

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
      ] = await Promise.all([
        adminApi.fetchSummary(),
        adminApi.getProfile(),
        adminApi.listProjects(),
        adminApi.listResearch(),
        adminApi.listContacts(),
        adminApi.listBlacklist(),
      ]);

      setSummary(summaryRes.data);
      setProfileForm(profileToForm(profileRes.data));
      setProjects(projectRes.data);
      setResearch(researchRes.data);
      setContacts(contactRes.data);
      setContactEdits(buildContactEditMap(contactRes.data));
      setBlacklist(blacklistRes.data);
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
        const hashToken = extractTokenFromHash(window.location.hash ?? "");
        if (hashToken) {
          storeToken(hashToken);
        }

        const sessionRes = await adminApi.session();
        const sessionToken = sessionRes.data.token?.trim();
        if (sessionToken) {
          storeToken(sessionToken);
        }
        if (sessionRes.data.active) {
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
      }
    };

    void resumeSession();

    return () => {
      cancelled = true;
    };
  }, [authState, token, refreshAll, handleUnauthorized, storeToken]);

  useEffect(() => {
    if (authState !== "authenticated") {
      return;
    }

    const currentToken = token ?? getToken();
    if (currentToken == null || currentToken === "") {
      handleUnauthorized();
    }
  }, [authState, token, handleUnauthorized]);

  useEffect(() => {
    if (authState !== "authenticated") {
      return;
    }

    const poll = async () => {
      try {
        const sessionRes = await adminApi.session();
        const sessionToken = sessionRes.data.token?.trim();
        if (sessionToken) {
          storeToken(sessionToken);
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
  }, [authState, handleUnauthorized, storeToken]);

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
    const payload = {
      title: { ja: projectForm.titleJa, en: projectForm.titleEn },
      description: {
        ja: projectForm.descriptionJa,
        en: projectForm.descriptionEn,
      },
      techStack: projectForm.techStack
        .split(",")
        .map((item) => item.trim())
        .filter(Boolean),
      linkUrl: projectForm.linkUrl.trim(),
      year: Number(projectForm.year) || currentYear,
      published: projectForm.published,
      sortOrder:
        projectForm.sortOrder === "" ? null : Number(projectForm.sortOrder),
    };

    await run(async () => {
      if (editingProjectId != null) {
        await adminApi.updateProject(editingProjectId, payload);
      } else {
        await adminApi.createProject(payload);
      }
    });
    setProjectForm({ ...emptyProjectForm });
    setEditingProjectId(null);
  };

  const handleSubmitResearch = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const payload = {
      title: { ja: researchForm.titleJa, en: researchForm.titleEn },
      summary: { ja: researchForm.summaryJa, en: researchForm.summaryEn },
      contentMd: { ja: researchForm.contentJa, en: researchForm.contentEn },
      year: Number(researchForm.year) || currentYear,
      published: researchForm.published,
    };

    await run(async () => {
      if (editingResearchId != null) {
        await adminApi.updateResearch(editingResearchId, payload);
      } else {
        await adminApi.createResearch(payload);
      }
    });
    setResearchForm({ ...emptyResearchForm });
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
      await adminApi.updateProject(project.id, {
        title: project.title,
        description: project.description,
        techStack: project.techStack,
        linkUrl: project.linkUrl,
        year: project.year,
        published: !project.published,
        sortOrder: project.sortOrder ?? null,
      });
    });

  const toggleResearchPublished = (item: AdminResearch) =>
    run(async () => {
      await adminApi.updateResearch(item.id, {
        title: item.title,
        summary: item.summary,
        contentMd: item.contentMd,
        year: item.year,
        published: !item.published,
      });
    });

  const deleteProject = (id: number) => run(() => adminApi.deleteProject(id));
  const deleteResearch = (id: number) => run(() => adminApi.deleteResearch(id));
  const deleteBlacklistEntry = (id: number) =>
    run(() => adminApi.deleteBlacklist(id));

  const handleEditProject = (project: AdminProject) => {
    setEditingProjectId(project.id);
    setProjectForm({
      titleJa: project.title.ja ?? "",
      titleEn: project.title.en ?? "",
      descriptionJa: project.description.ja ?? "",
      descriptionEn: project.description.en ?? "",
      techStack: project.techStack.join(", "),
      linkUrl: project.linkUrl,
      year: project.year.toString(),
      published: project.published,
      sortOrder: project.sortOrder != null ? project.sortOrder.toString() : "",
    });
  };

  const handleEditResearch = (item: AdminResearch) => {
    setEditingResearchId(item.id);
    setResearchForm({
      titleJa: item.title.ja ?? "",
      titleEn: item.title.en ?? "",
      summaryJa: item.summary.ja ?? "",
      summaryEn: item.summary.en ?? "",
      contentJa: item.contentMd.ja ?? "",
      contentEn: item.contentMd.en ?? "",
      year: item.year.toString(),
      published: item.published,
    });
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
  };

  const resetResearchForm = () => {
    setResearchForm({ ...emptyResearchForm });
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
            <div className="grid gap-4 md:grid-cols-2">
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.techStack")}
                </label>
                <input
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={projectForm.techStack}
                  onChange={(event) =>
                    setProjectForm((prev) => ({
                      ...prev,
                      techStack: event.target.value,
                    }))
                  }
                  placeholder={t("fields.techStackPlaceholder") ?? ""}
                />
              </div>
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
            </div>
            <div className="flex items-center justify-end gap-3">
              <button
                type="submit"
                className="inline-flex items-center justify-center rounded-md bg-slate-900 px-4 py-2 text-sm font-medium text-white shadow-sm transition hover:bg-slate-800 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-slate-900"
                disabled={loading}
              >
                {t(editingProjectId != null ? "actions.update" : "actions.create")}
              </button>
            </div>
          </form>

          <div className="mt-6 space-y-4">
            {projects.map((project) => (
              <div
                key={project.id}
                className="rounded-md border border-slate-200 p-4"
              >
                <div className="flex flex-col gap-2 sm:flex-row sm:items-start sm:justify-between">
                  <div>
                    <h3 className="text-base font-semibold text-slate-900">
                      {project.title.ja || project.title.en || t("projects.untitled")}
                    </h3>
                    <p className="text-sm text-slate-600">
                      {project.description.ja || project.description.en || t("projects.noDescription")}
                    </p>
                    <div className="mt-2 flex flex-wrap gap-2 text-xs text-slate-500">
                      <span>
                        {t("fields.year")}: {project.year}
                      </span>
                      {project.techStack.length > 0 && (
                        <span>
                          {t("fields.techStack")}: {project.techStack.join(", ")}
                        </span>
                      )}
                      <span>
                        {t("fields.published")}: {project.published ? t("status.published") : t("status.draft")}
                      </span>
                    </div>
                  </div>
                  <div className="flex gap-2">
                    <button
                      type="button"
                      className="rounded-md border border-slate-200 px-3 py-2 text-sm text-slate-600 transition hover:bg-slate-100"
                      onClick={() => toggleProjectPublished(project)}
                    >
                      {project.published ? t("actions.unpublish") : t("actions.publish")}
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
            ))}
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
          <form className="mt-4 space-y-4" onSubmit={handleSubmitResearch}>
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
                  {t("fields.summaryJa")}
                </label>
                <textarea
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  rows={3}
                  value={researchForm.summaryJa}
                  onChange={(event) =>
                    setResearchForm((prev) => ({
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
                  rows={3}
                  value={researchForm.summaryEn}
                  onChange={(event) =>
                    setResearchForm((prev) => ({
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
                  {t("fields.contentJa")}
                </label>
                <textarea
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  rows={4}
                  value={researchForm.contentJa}
                  onChange={(event) =>
                    setResearchForm((prev) => ({
                      ...prev,
                      contentJa: event.target.value,
                    }))
                  }
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.contentEn")}
                </label>
                <textarea
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  rows={4}
                  value={researchForm.contentEn}
                  onChange={(event) =>
                    setResearchForm((prev) => ({
                      ...prev,
                      contentEn: event.target.value,
                    }))
                  }
                />
              </div>
            </div>
            <div className="grid gap-4 md:grid-cols-2">
              <div>
                <label className="block text-sm font-medium text-slate-700">
                  {t("fields.year")}
                </label>
                <input
                  type="number"
                  className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 text-sm"
                  value={researchForm.year}
                  onChange={(event) =>
                    setResearchForm((prev) => ({
                      ...prev,
                      year: event.target.value,
                    }))
                  }
                />
              </div>
              <div className="flex items-center gap-2">
                <input
                  id="research-published"
                  type="checkbox"
                  className="h-4 w-4 rounded border-slate-300 text-slate-900 focus:ring-slate-900"
                  checked={researchForm.published}
                  onChange={(event) =>
                    setResearchForm((prev) => ({
                      ...prev,
                      published: event.target.checked,
                    }))
                  }
                />
                <label
                  htmlFor="research-published"
                  className="text-sm font-medium text-slate-700"
                >
                  {t("fields.published")}
                </label>
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
                      {item.summary.ja || item.summary.en || t("research.noSummary")}
                    </p>
                    <div className="mt-2 flex flex-wrap gap-2 text-xs text-slate-500">
                      <span>
                        {t("fields.year")}: {item.year}
                      </span>
                      <span>
                        {t("fields.published")}: {item.published ? t("status.published") : t("status.draft")}
                      </span>
                    </div>
                  </div>
                  <div className="flex gap-2">
                    <button
                      type="button"
                      className="rounded-md border border-slate-200 px-3 py-2 text-sm text-slate-600 transition hover:bg-slate-100"
                      onClick={() => toggleResearchPublished(item)}
                    >
                      {item.published ? t("actions.unpublish") : t("actions.publish")}
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
