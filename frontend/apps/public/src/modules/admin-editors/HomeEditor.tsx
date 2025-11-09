import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from "react";

import { useUnsavedChangesTracker } from "../../hooks/useUnsavedChangesTracker";
import {
  fetchHomeConfig,
  updateHomeConfig,
} from "../admin-console/api";
import { LocalizedTextField } from "../admin-console/components/LocalizedTextField";
import { SortableList } from "../admin-console/components/SortableList";
import type {
  HomeChipSourceItem,
  HomeConfigDocument,
  HomeQuickLinkItem,
  LocalizedField,
} from "../admin-console/types";

type EditableQuickLink = HomeQuickLinkItem & { clientId: string };
type EditableChipSource = HomeChipSourceItem & { clientId: string };

type EditableHomeConfig = {
  id: number;
  heroSubtitle: LocalizedField;
  quickLinks: EditableQuickLink[];
  chipSources: EditableChipSource[];
  updatedAt?: string;
};

const QUICK_LINK_SECTIONS: { value: HomeQuickLinkItem["section"]; label: string }[] =
  [
    { value: "profile", label: "Profile" },
    { value: "research_blog", label: "Research & Blog" },
    { value: "projects", label: "Projects" },
    { value: "contact", label: "Contact" },
  ];

const CHIP_SOURCE_KINDS: { value: HomeChipSourceItem["source"]; label: string }[] =
  [
    { value: "tech", label: "Tech catalog" },
    { value: "affiliation", label: "Affiliations" },
    { value: "community", label: "Communities" },
  ];

function toEditable(config: HomeConfigDocument): EditableHomeConfig {
  return {
    id: config.id,
    heroSubtitle: {
      ja: config.heroSubtitle?.ja ?? "",
      en: config.heroSubtitle?.en ?? "",
    },
    quickLinks: config.quickLinks
      .map((link) => ({
        ...link,
        clientId: `ql-${link.id}`,
        label: {
          ja: link.label?.ja ?? "",
          en: link.label?.en ?? "",
        },
        description: {
          ja: link.description?.ja ?? "",
          en: link.description?.en ?? "",
        },
        cta: {
          ja: link.cta?.ja ?? "",
          en: link.cta?.en ?? "",
        },
      }))
      .sort((a, b) => a.sortOrder - b.sortOrder),
    chipSources: config.chipSources
      .map((source) => ({
        ...source,
        clientId: `cs-${source.id}`,
        label: {
          ja: source.label?.ja ?? "",
          en: source.label?.en ?? "",
        },
      }))
      .sort((a, b) => a.sortOrder - b.sortOrder),
    updatedAt: config.updatedAt,
  };
}

function createQuickLink(sortOrder: number): EditableQuickLink {
  const uuid =
    typeof crypto !== "undefined" && "randomUUID" in crypto
      ? crypto.randomUUID()
      : Math.random().toString(36).slice(2);
  return {
    id: 0,
    clientId: `new-ql-${uuid}`,
    section: "profile",
    label: { ja: "", en: "" },
    description: { ja: "", en: "" },
    cta: { ja: "", en: "" },
    targetUrl: "",
    sortOrder,
  };
}

function createChipSource(sortOrder: number): EditableChipSource {
  const uuid =
    typeof crypto !== "undefined" && "randomUUID" in crypto
      ? crypto.randomUUID()
      : Math.random().toString(36).slice(2);
  return {
    id: 0,
    clientId: `new-cs-${uuid}`,
    source: "tech",
    label: { ja: "", en: "" },
    limit: 4,
    sortOrder,
  };
}

type HomeEditorContextValue = {
  enabled: boolean;
  loading: boolean;
  error: string | null;
  draft: EditableHomeConfig | null;
  isDirty: boolean;
  saving: boolean;
  saveError: string | null;
  setHeroSubtitle: (value: LocalizedField) => void;
  addQuickLink: () => void;
  updateQuickLink: (
    clientId: string,
    updater: (link: EditableQuickLink) => EditableQuickLink,
  ) => void;
  removeQuickLink: (clientId: string) => void;
  reorderQuickLinks: (next: EditableQuickLink[]) => void;
  addChipSource: () => void;
  updateChipSource: (
    clientId: string,
    updater: (source: EditableChipSource) => EditableChipSource,
  ) => void;
  removeChipSource: (clientId: string) => void;
  reorderChipSources: (next: EditableChipSource[]) => void;
  save: () => Promise<void>;
};

const HomeEditorContext = createContext<HomeEditorContextValue | null>(null);

type HomeEditorProviderProps = {
  children: ReactNode;
  enabled: boolean;
};

export function HomeEditorProvider({ children, enabled }: HomeEditorProviderProps) {
  const [initial, setInitial] = useState<EditableHomeConfig | null>(null);
  const [draft, setDraft] = useState<EditableHomeConfig | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);
  const [saveError, setSaveError] = useState<string | null>(null);
  const [hasLoaded, setHasLoaded] = useState(false);

  useEffect(() => {
    if (!enabled) {
      return;
    }
    if (hasLoaded) {
      return;
    }

    let mounted = true;
    const load = async () => {
      setLoading(true);
      setError(null);
      try {
        const data = await fetchHomeConfig();
        if (!mounted) {
          return;
        }
        const editable = toEditable(data);
        setInitial(editable);
        setDraft(editable);
        setHasLoaded(true);
      } catch (err) {
        if (mounted) {
          setError(
            err instanceof Error ? err.message : "Failed to load home configuration.",
          );
        }
      } finally {
        if (mounted) {
          setLoading(false);
        }
      }
    };

    void load();
    return () => {
      mounted = false;
    };
  }, [enabled, hasLoaded]);

  const isDirty = useMemo(() => {
    if (!initial || !draft) {
      return false;
    }
    return JSON.stringify(draft) !== JSON.stringify(initial);
  }, [draft, initial]);

  useUnsavedChangesTracker("admin-home-editor", enabled && isDirty);

  const setHeroSubtitle = useCallback((value: LocalizedField) => {
    if (!enabled) {
      return;
    }
    setDraft((prev) => (prev ? { ...prev, heroSubtitle: value } : prev));
  }, [enabled]);

  const updateQuickLink = useCallback(
    (
      clientId: string,
      updater: (link: EditableQuickLink) => EditableQuickLink,
    ) => {
      if (!enabled) {
        return;
      }
      setDraft((prev) => {
        if (!prev) {
          return prev;
        }
        return {
          ...prev,
          quickLinks: prev.quickLinks.map((link) =>
            link.clientId === clientId ? updater(link) : link,
          ),
        };
      });
    },
    [enabled],
  );

  const updateChipSource = useCallback(
    (
      clientId: string,
      updater: (source: EditableChipSource) => EditableChipSource,
    ) => {
      if (!enabled) {
        return;
      }
      setDraft((prev) => {
        if (!prev) {
          return prev;
        }
        return {
          ...prev,
          chipSources: prev.chipSources.map((source) =>
            source.clientId === clientId ? updater(source) : source,
          ),
        };
      });
    },
    [enabled],
  );

  const addQuickLink = useCallback(() => {
    if (!enabled) {
      return;
    }
    setDraft((prev) => {
      if (!prev) {
        return prev;
      }
      const nextLink = createQuickLink((prev.quickLinks.at(-1)?.sortOrder ?? 0) + 1);
      return {
        ...prev,
        quickLinks: [...prev.quickLinks, nextLink],
      };
    });
  }, [enabled]);

  const removeQuickLink = useCallback((clientId: string) => {
    if (!enabled) {
      return;
    }
    setDraft((prev) => {
      if (!prev) {
        return prev;
      }
      return {
        ...prev,
        quickLinks: prev.quickLinks.filter((link) => link.clientId !== clientId),
      };
    });
  }, [enabled]);

  const addChipSource = useCallback(() => {
    if (!enabled) {
      return;
    }
    setDraft((prev) => {
      if (!prev) {
        return prev;
      }
      const nextSource = createChipSource(
        (prev.chipSources.at(-1)?.sortOrder ?? 0) + 1,
      );
      return {
        ...prev,
        chipSources: [...prev.chipSources, nextSource],
      };
    });
  }, [enabled]);

  const removeChipSource = useCallback((clientId: string) => {
    if (!enabled) {
      return;
    }
    setDraft((prev) => {
      if (!prev) {
        return prev;
      }
      return {
        ...prev,
        chipSources: prev.chipSources.filter((source) => source.clientId !== clientId),
      };
    });
  }, [enabled]);

  const reorderQuickLinks = useCallback((list: EditableQuickLink[]) => {
    if (!enabled) {
      return;
    }
    setDraft((prev) => {
      if (!prev) {
        return prev;
      }
      const normalized = list.map((item, index) => ({
        ...item,
        sortOrder: index + 1,
      }));
      return { ...prev, quickLinks: normalized };
    });
  }, [enabled]);

  const reorderChipSources = useCallback((list: EditableChipSource[]) => {
    if (!enabled) {
      return;
    }
    setDraft((prev) => {
      if (!prev) {
        return prev;
      }
      const normalized = list.map((item, index) => ({
        ...item,
        sortOrder: index + 1,
      }));
      return { ...prev, chipSources: normalized };
    });
  }, [enabled]);

  const save = useCallback(async () => {
    if (!enabled || !draft) {
      return;
    }
    if (!draft) {
      return;
    }

    setSaving(true);
    setSaveError(null);
    try {
      const payload: HomeConfigDocument = {
        id: draft.id,
        heroSubtitle: draft.heroSubtitle,
        quickLinks: draft.quickLinks.map((link, index) => ({
          id: link.id,
          section: link.section,
          label: link.label,
          description: link.description,
          cta: link.cta,
          targetUrl: link.targetUrl,
          sortOrder: index + 1,
        })),
        chipSources: draft.chipSources.map((source, index) => ({
          id: source.id,
          source: source.source,
          label: source.label,
          limit: source.limit,
          sortOrder: index + 1,
        })),
        updatedAt: draft.updatedAt,
      };

      const updated = await updateHomeConfig(payload);
      const editable = toEditable(updated);
      setInitial(editable);
      setDraft(editable);
    } catch (err) {
      setSaveError(
        err instanceof Error ? err.message : "Failed to update home configuration.",
      );
      throw err;
    } finally {
      setSaving(false);
    }
  }, [draft, enabled]);

  const value = useMemo<HomeEditorContextValue>(
    () => ({
      enabled,
      loading,
      error,
      draft,
      isDirty,
      saving,
      saveError,
      setHeroSubtitle,
      addQuickLink,
      updateQuickLink,
      removeQuickLink,
      reorderQuickLinks,
      addChipSource,
      updateChipSource,
      removeChipSource,
      reorderChipSources,
      save,
    }),
    [
      enabled,
      loading,
      error,
      draft,
      isDirty,
      saving,
      saveError,
      setHeroSubtitle,
      addQuickLink,
      updateQuickLink,
      removeQuickLink,
      reorderQuickLinks,
      addChipSource,
      updateChipSource,
      removeChipSource,
      reorderChipSources,
      save,
    ],
  );

  return (
    <HomeEditorContext.Provider value={value}>
      {children}
    </HomeEditorContext.Provider>
  );
}

export function useHomeEditorContext(): HomeEditorContextValue {
  const context = useContext(HomeEditorContext);
  if (!context) {
    throw new Error("useHomeEditorContext must be used within HomeEditorProvider");
  }
  return context;
}

export function useOptionalHomeEditorContext(): HomeEditorContextValue | null {
  return useContext(HomeEditorContext);
}

export function HomeEditorPanel() {
  const {
    enabled,
    loading,
    error,
    draft,
    isDirty,
    saving,
    saveError,
    setHeroSubtitle,
    addQuickLink,
    updateQuickLink,
    removeQuickLink,
    reorderQuickLinks,
    addChipSource,
    updateChipSource,
    removeChipSource,
    reorderChipSources,
    save,
  } = useHomeEditorContext();

  if (!enabled) {
    return null;
  }

  if (loading || !draft) {
    return (
      <div className="space-y-4 rounded-2xl border border-slate-200 bg-white/70 p-5 shadow-sm dark:border-slate-700 dark:bg-slate-900/50">
        <div className="h-4 w-32 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
        <div className="h-52 animate-pulse rounded-xl bg-slate-100 dark:bg-slate-800" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="rounded-2xl border border-rose-300 bg-rose-50 p-5 text-sm text-rose-700 dark:border-rose-900 dark:bg-rose-950/40 dark:text-rose-300">
        {error}
      </div>
    );
  }

  return (
    <section className="space-y-6 rounded-2xl border border-slate-200 bg-white/80 p-6 shadow-sm dark:border-slate-700 dark:bg-slate-900/50">
      <header className="flex flex-col gap-2 md:flex-row md:items-center md:justify-between">
        <div>
          <h2 className="text-xl font-semibold text-slate-900 dark:text-slate-100">
            Home page content
          </h2>
          <p className="text-sm text-slate-600 dark:text-slate-300">
            Update hero subtitles, quick links, and chip groups directly on the public page.
          </p>
        </div>
        <button
          type="button"
          disabled={!isDirty || saving}
          onClick={() => {
            void save();
          }}
          className="inline-flex items-center rounded-md bg-sky-600 px-4 py-2 text-sm font-semibold text-white shadow-sm transition enabled:hover:bg-sky-500 enabled:focus-visible:outline enabled:focus-visible:outline-2 enabled:focus-visible:outline-offset-2 enabled:focus-visible:outline-sky-500 disabled:cursor-not-allowed disabled:opacity-60 dark:bg-sky-500 dark:enabled:hover:bg-sky-400"
        >
          {saving ? "Saving…" : "Save home content"}
        </button>
      </header>

      {saveError && (
        <p className="rounded-md border border-rose-300 bg-rose-50 px-3 py-2 text-sm text-rose-700 dark:border-rose-900 dark:bg-rose-950/40 dark:text-rose-300">
          {saveError}
        </p>
      )}

      <LocalizedTextField
        id="hero-subtitle"
        label="Hero subtitle"
        value={draft.heroSubtitle}
        onChange={setHeroSubtitle}
      />

      <QuickLinksEditor
        links={draft.quickLinks}
        onAdd={addQuickLink}
        onUpdate={updateQuickLink}
        onRemove={removeQuickLink}
        onReorder={reorderQuickLinks}
      />

      <ChipSourcesEditor
        sources={draft.chipSources}
        onAdd={addChipSource}
        onUpdate={updateChipSource}
        onRemove={removeChipSource}
        onReorder={reorderChipSources}
      />
    </section>
  );
}

type QuickLinksEditorProps = {
  links: EditableQuickLink[];
  onAdd: () => void;
  onUpdate: (
    clientId: string,
    updater: (link: EditableQuickLink) => EditableQuickLink,
  ) => void;
  onRemove: (clientId: string) => void;
  onReorder: (next: EditableQuickLink[]) => void;
};

function QuickLinksEditor({
  links,
  onAdd,
  onUpdate,
  onRemove,
  onReorder,
}: QuickLinksEditorProps) {
  return (
    <>
      <div className="flex items-center justify-between">
        <h3 className="text-lg font-semibold text-slate-900 dark:text-slate-100">
          Quick links
        </h3>
        <button
          type="button"
          className="rounded-full border border-slate-300 px-3 py-1 text-xs font-semibold text-slate-600 transition hover:border-sky-400 hover:text-sky-600 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-sky-500 dark:border-slate-700 dark:text-slate-300 dark:hover:border-sky-400 dark:hover:text-sky-300 dark:focus-visible:outline-sky-400"
          onClick={onAdd}
        >
          + Add block
        </button>
      </div>
      <SortableList
        items={links}
        getId={(link) => link.clientId}
        onReorder={onReorder}
        renderItem={(link) => (
          <div className="flex flex-col gap-4 rounded-xl border border-slate-200 bg-white/90 p-5 shadow-sm transition hover:border-sky-300 dark:border-slate-700 dark:bg-slate-900/70">
            <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
              <div className="flex flex-wrap items-center gap-3 text-xs text-slate-500 dark:text-slate-400">
                <span className="inline-flex items-center rounded-full bg-slate-100 px-3 py-1 dark:bg-slate-800">
                  Sort {link.sortOrder}
                </span>
                <label className="inline-flex items-center gap-2 text-sm text-slate-700 dark:text-slate-300">
                  <span className="font-semibold">Section</span>
                  <select
                    value={link.section}
                    onChange={(event) =>
                      onUpdate(link.clientId, (candidate) => ({
                        ...candidate,
                        section: event.target.value as HomeQuickLinkItem["section"],
                      }))
                    }
                    className="h-9 rounded-md border border-slate-300 bg-white px-3 text-sm text-slate-900 shadow-sm transition focus:border-sky-500 focus:outline-none focus:ring-2 focus:ring-sky-500 dark:border-slate-700 dark:bg-slate-900 dark:text-slate-100"
                  >
                    {QUICK_LINK_SECTIONS.map((option) => (
                      <option key={option.value} value={option.value}>
                        {option.label}
                      </option>
                    ))}
                  </select>
                </label>
              </div>
              <button
                type="button"
                className="rounded-full border border-slate-300 px-3 py-1 text-xs font-semibold text-rose-600 transition hover:border-rose-400 hover:text-rose-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-rose-500 dark:border-slate-700 dark:text-rose-300 dark:hover:border-rose-500 dark:hover:text-rose-200 dark:focus-visible:outline-rose-500"
                onClick={() => onRemove(link.clientId)}
              >
                Remove
              </button>
            </div>

            <LocalizedTextField
              id={`quicklink-label-${link.clientId}`}
              label="Label"
              value={link.label}
              onChange={(value) =>
                onUpdate(link.clientId, (candidate) => ({
                  ...candidate,
                  label: value,
                }))
              }
            />
            <LocalizedTextField
              id={`quicklink-description-${link.clientId}`}
              label="Description"
              value={link.description}
              onChange={(value) =>
                onUpdate(link.clientId, (candidate) => ({
                  ...candidate,
                  description: value,
                }))
              }
              multiline
            />
            <LocalizedTextField
              id={`quicklink-cta-${link.clientId}`}
              label="Call to action"
              value={link.cta}
              onChange={(value) =>
                onUpdate(link.clientId, (candidate) => ({
                  ...candidate,
                  cta: value,
                }))
              }
            />
            <label className="flex flex-col gap-2 text-sm text-slate-700 dark:text-slate-300">
              <span className="font-semibold">Target URL</span>
              <input
                value={link.targetUrl}
                onChange={(event) =>
                  onUpdate(link.clientId, (candidate) => ({
                    ...candidate,
                    targetUrl: event.target.value,
                  }))
                }
                placeholder="https://"
                className="rounded-md border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 shadow-sm transition focus:border-sky-500 focus:outline-none focus:ring-2 focus:ring-sky-500 dark:border-slate-700 dark:bg-slate-900 dark:text-slate-100"
              />
            </label>
          </div>
        )}
      />
    </>
  );
}

type ChipSourcesEditorProps = {
  sources: EditableChipSource[];
  onAdd: () => void;
  onUpdate: (
    clientId: string,
    updater: (source: EditableChipSource) => EditableChipSource,
  ) => void;
  onRemove: (clientId: string) => void;
  onReorder: (next: EditableChipSource[]) => void;
};

function ChipSourcesEditor({
  sources,
  onAdd,
  onUpdate,
  onRemove,
  onReorder,
}: ChipSourcesEditorProps) {
  return (
    <>
      <div className="flex items-center justify-between">
        <h3 className="text-lg font-semibold text-slate-900 dark:text-slate-100">
          Chip sources
        </h3>
        <button
          type="button"
          className="rounded-full border border-slate-300 px-3 py-1 text-xs font-semibold text-slate-600 transition hover:border-sky-400 hover:text-sky-600 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-sky-500 dark:border-slate-700 dark:text-slate-300 dark:hover:border-sky-400 dark:hover:text-sky-300 dark:focus-visible:outline-sky-400"
          onClick={onAdd}
        >
          + Add block
        </button>
      </div>

      <SortableList
        items={sources}
        getId={(source) => source.clientId}
        onReorder={onReorder}
        renderItem={(source) => (
          <div className="flex flex-col gap-4 rounded-xl border border-slate-200 bg-white/90 p-5 shadow-sm transition hover:border-sky-300 dark:border-slate-700 dark:bg-slate-900/70">
            <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
              <div className="flex flex-wrap items-center gap-3 text-xs text-slate-500 dark:text-slate-400">
                <span className="inline-flex items-center rounded-full bg-slate-100 px-3 py-1 dark:bg-slate-800">
                  Sort {source.sortOrder}
                </span>
                <label className="inline-flex items-center gap-2 text-sm text-slate-700 dark:text-slate-300">
                  <span className="font-semibold">Source</span>
                  <select
                    value={source.source}
                    onChange={(event) =>
                      onUpdate(source.clientId, (candidate) => ({
                        ...candidate,
                        source: event.target.value as HomeChipSourceItem["source"],
                      }))
                    }
                    className="h-9 rounded-md border border-slate-300 bg-white px-3 text-sm text-slate-900 shadow-sm transition focus:border-sky-500 focus:outline-none focus:ring-2 focus:ring-sky-500 dark:border-slate-700 dark:bg-slate-900 dark:text-slate-100"
                  >
                    {CHIP_SOURCE_KINDS.map((option) => (
                      <option key={option.value} value={option.value}>
                        {option.label}
                      </option>
                    ))}
                  </select>
                </label>
              </div>
              <button
                type="button"
                className="rounded-full border border-slate-300 px-3 py-1 text-xs font-semibold text-rose-600 transition hover:border-rose-400 hover:text-rose-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-rose-500 dark:border-slate-700 dark:text-rose-300 dark:hover:border-rose-500 dark:hover:text-rose-200 dark:focus-visible:outline-rose-500"
                onClick={() => onRemove(source.clientId)}
              >
                Remove
              </button>
            </div>
            <LocalizedTextField
              id={`chip-source-${source.clientId}`}
              label="Label"
              value={source.label}
              onChange={(value) =>
                onUpdate(source.clientId, (candidate) => ({
                  ...candidate,
                  label: value,
                }))
              }
            />
            <label className="flex flex-col gap-2 text-sm text-slate-700 dark:text-slate-300">
              <span className="font-semibold">Limit</span>
              <input
                type="number"
                min={1}
                value={source.limit}
                onChange={(event) =>
                  onUpdate(source.clientId, (candidate) => ({
                    ...candidate,
                    limit: Number.parseInt(event.target.value, 10) || 1,
                  }))
                }
                className="rounded-md border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 shadow-sm transition focus:border-sky-500 focus:outline-none focus:ring-2 focus:ring-sky-500 dark:border-slate-700 dark:bg-slate-900 dark:text-slate-100"
              />
            </label>
          </div>
        )}
      />
    </>
  );
}

export type { EditableHomeConfig, EditableQuickLink, EditableChipSource };
