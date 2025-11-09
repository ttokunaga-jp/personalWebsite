import { useEffect, useMemo, useState } from "react";

import { useUnsavedChangesTracker } from "../../../hooks/useUnsavedChangesTracker";
import { fetchSocialLinks, replaceSocialLinks } from "../api";
import { LocalizedTextField } from "../components/LocalizedTextField";
import { SortableList } from "../components/SortableList";
import type {
  LocalizedField,
  SocialLink,
  SocialLinkInput,
  SocialProvider,
} from "../types";

type EditableSocialLink = SocialLink & { clientId: string };

const PROVIDER_OPTIONS: { value: SocialProvider; label: string }[] = [
  { value: "github", label: "GitHub" },
  { value: "zenn", label: "Zenn" },
  { value: "linkedin", label: "LinkedIn" },
  { value: "x", label: "X" },
  { value: "email", label: "Email" },
  { value: "website", label: "Website" },
  { value: "other", label: "Other" },
];

function toEditable(link: SocialLink): EditableSocialLink {
  return {
    ...link,
    clientId: `server-${link.id}`,
    label: {
      ja: link.label?.ja ?? "",
      en: link.label?.en ?? "",
    },
  };
}

function createEmptyLink(sortOrder: number): EditableSocialLink {
  const uuid =
    typeof crypto !== "undefined" && "randomUUID" in crypto
      ? crypto.randomUUID()
      : Math.random().toString(36).slice(2);
  const clientId = `draft-${uuid}`;
  return {
    id: 0,
    clientId,
    provider: "other",
    label: {
      ja: "",
      en: "",
    },
    url: "",
    isFooter: false,
    sortOrder,
  };
}

function toPayload(link: EditableSocialLink): SocialLinkInput {
  return {
    provider: link.provider,
    label: {
      ja: link.label?.ja ?? "",
      en: link.label?.en ?? "",
    },
    url: link.url,
    isFooter: link.isFooter,
    sortOrder: link.sortOrder,
  };
}

export function SocialLinksPanel() {
  const [initialLinks, setInitialLinks] = useState<EditableSocialLink[]>([]);
  const [links, setLinks] = useState<EditableSocialLink[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);
  const [formError, setFormError] = useState<string | null>(null);

  useEffect(() => {
    let active = true;
    const run = async () => {
      setLoading(true);
      setError(null);
      try {
        const data = await fetchSocialLinks();
        if (active) {
          const editable = data
            .map(toEditable)
            .sort((a, b) => a.sortOrder - b.sortOrder);
          setInitialLinks(editable);
          setLinks(editable);
        }
      } catch (err) {
        if (active) {
          setError(
            err instanceof Error
              ? err.message
              : "Unable to load social links.",
          );
        }
      } finally {
        if (active) {
          setLoading(false);
        }
      }
    };

    void run();
    return () => {
      active = false;
    };
  }, []);

  const isDirty = useMemo(() => {
    const normalize = (input: EditableSocialLink[]) =>
      input.map((link) => ({
        provider: link.provider,
        label: {
          ja: link.label?.ja ?? "",
          en: link.label?.en ?? "",
        },
        url: link.url,
        isFooter: link.isFooter,
        sortOrder: link.sortOrder,
      }));

    return (
      JSON.stringify(normalize(links)) !==
      JSON.stringify(normalize(initialLinks))
    );
  }, [initialLinks, links]);

  useUnsavedChangesTracker("admin-social-links", isDirty);

  const handleReorder = (items: EditableSocialLink[]) => {
    const next = items.map((link, index) => ({
      ...link,
      sortOrder: index + 1,
    }));
    setLinks(next);
  };

  const handleAddBlock = () => {
    const next = [
      ...links,
      createEmptyLink((links.at(-1)?.sortOrder ?? 0) + 1),
    ];
    setLinks(next);
  };

  const handleRemove = (clientId: string) => {
    setLinks((prev) => prev.filter((link) => link.clientId !== clientId));
  };

  const handleProviderChange = (clientId: string, provider: SocialProvider) => {
    setLinks((prev) =>
      prev.map((link) =>
        link.clientId === clientId ? { ...link, provider } : link,
      ),
    );
  };

  const handleLabelChange = (
    clientId: string,
    field: keyof LocalizedField,
    value: string,
  ) => {
    setLinks((prev) =>
      prev.map((link) =>
        link.clientId === clientId
          ? {
              ...link,
              label: {
                ja: field === "ja" ? value : link.label.ja,
                en: field === "en" ? value : link.label.en,
              },
            }
          : link,
      ),
    );
  };

  const handleUrlChange = (clientId: string, value: string) => {
    setLinks((prev) =>
      prev.map((link) =>
        link.clientId === clientId ? { ...link, url: value } : link,
      ),
    );
  };

  const handleFooterToggle = (clientId: string, value: boolean) => {
    setLinks((prev) =>
      prev.map((link) =>
        link.clientId === clientId ? { ...link, isFooter: value } : link,
      ),
    );
  };

  const handleSave = async () => {
    setSaving(true);
    setFormError(null);
    try {
      const payload: SocialLinkInput[] = links.map(toPayload);
      const updated = await replaceSocialLinks(payload);
      const editable = updated
        .map(toEditable)
        .sort((a, b) => a.sortOrder - b.sortOrder);
      setInitialLinks(editable);
      setLinks(editable);
    } catch (err) {
      setFormError(
        err instanceof Error ? err.message : "Failed to save social links.",
      );
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <div className="space-y-4">
        <div className="h-4 w-32 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
        <div className="h-56 animate-pulse rounded-xl bg-slate-100 dark:bg-slate-800" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="rounded-md border border-rose-300 bg-rose-50 p-4 text-sm text-rose-700 dark:border-rose-900 dark:bg-rose-950/30 dark:text-rose-300">
        {error}
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-lg font-semibold text-slate-900 dark:text-slate-100">
            Social links
          </h2>
          <p className="text-sm text-slate-600 dark:text-slate-300">
            Maintain bilingual labels, drag to reorder, and mark links for the footer.
          </p>
        </div>
        <div className="flex items-center gap-3">
          <button
            type="button"
            className="inline-flex items-center rounded-md border border-slate-300 px-3 py-2 text-sm font-medium text-slate-600 transition hover:border-sky-400 hover:text-sky-600 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-sky-500 dark:border-slate-700 dark:text-slate-300 dark:hover:border-sky-400 dark:hover:text-sky-300 dark:focus-visible:outline-sky-400"
            onClick={handleAddBlock}
          >
            + Add block
          </button>
          <button
            type="button"
            disabled={!isDirty || saving}
            className="inline-flex items-center rounded-md bg-sky-600 px-4 py-2 text-sm font-semibold text-white shadow-sm transition enabled:hover:bg-sky-500 enabled:focus-visible:outline enabled:focus-visible:outline-2 enabled:focus-visible:outline-offset-2 enabled:focus-visible:outline-sky-500 disabled:cursor-not-allowed disabled:opacity-60 dark:bg-sky-500 dark:enabled:hover:bg-sky-400"
            onClick={handleSave}
          >
            {saving ? "Saving…" : "Save links"}
          </button>
        </div>
      </div>

      {formError && (
        <div className="rounded-md border border-rose-300 bg-rose-50 p-3 text-sm text-rose-700 dark:border-rose-900 dark:bg-rose-950/30 dark:text-rose-300">
          {formError}
        </div>
      )}

      <SortableList
        items={links}
        getId={(link) => link.clientId}
        onReorder={handleReorder}
        renderItem={(link) => (
          <div className="relative flex flex-col gap-4 rounded-2xl border border-slate-200 bg-white/80 p-5 shadow-sm transition hover:border-sky-300 dark:border-slate-700 dark:bg-slate-900/60">
            <div className="flex items-start justify-between gap-3">
              <div className="flex flex-1 gap-3">
                <div className="w-40">
                  <label className="flex flex-col gap-2 text-sm text-slate-700 dark:text-slate-300">
                    <span className="font-semibold">Provider</span>
                    <select
                      value={link.provider}
                      onChange={(event) =>
                        handleProviderChange(
                          link.clientId,
                          event.target.value as SocialProvider,
                        )
                      }
                      className="h-9 rounded-md border border-slate-300 bg-white px-3 text-sm text-slate-900 shadow-sm transition focus:border-sky-500 focus:outline-none focus:ring-2 focus:ring-sky-500 dark:border-slate-700 dark:bg-slate-900 dark:text-slate-100"
                    >
                      {PROVIDER_OPTIONS.map((option) => (
                        <option key={option.value} value={option.value}>
                          {option.label}
                        </option>
                      ))}
                    </select>
                  </label>
                </div>
                <div className="flex-1">
                  <LocalizedTextField
                    id={`social-label-${link.clientId}`}
                    label="Label"
                    value={link.label}
                    onChange={(value) => {
                      handleLabelChange(link.clientId, "ja", value.ja);
                      handleLabelChange(link.clientId, "en", value.en);
                    }}
                    required
                  />
                </div>
              </div>
              <button
                type="button"
                className="rounded-full border border-slate-300 px-3 py-1 text-xs font-semibold text-slate-600 transition hover:border-rose-400 hover:text-rose-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-rose-500 dark:border-slate-700 dark:text-slate-300 dark:hover:border-rose-500 dark:hover:text-rose-300 dark:focus-visible:outline-rose-500"
                onClick={() => handleRemove(link.clientId)}
              >
                Remove
              </button>
            </div>

            <label className="flex flex-col gap-2 text-sm text-slate-700 dark:text-slate-300">
              <span className="font-semibold">URL</span>
              <input
                value={link.url}
                onChange={(event) =>
                  handleUrlChange(link.clientId, event.target.value)
                }
                placeholder="https://"
                className="rounded-md border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 shadow-sm transition focus:border-sky-500 focus:outline-none focus:ring-2 focus:ring-sky-500 dark:border-slate-700 dark:bg-slate-900 dark:text-slate-100"
              />
            </label>

            <label className="inline-flex items-center gap-2 text-sm text-slate-700 dark:text-slate-300">
              <input
                type="checkbox"
                checked={link.isFooter}
                onChange={(event) =>
                  handleFooterToggle(link.clientId, event.target.checked)
                }
                className="h-4 w-4 rounded border-slate-300 text-sky-600 focus:ring-sky-500 dark:border-slate-700 dark:bg-slate-900 dark:checked:bg-sky-400"
              />
              Display in footer
            </label>
          </div>
        )}
      />
    </div>
  );
}
