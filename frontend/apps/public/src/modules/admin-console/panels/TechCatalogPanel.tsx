import { FormEvent, useEffect, useMemo, useState } from "react";

import { useUnsavedChangesTracker } from "../../../hooks/useUnsavedChangesTracker";
import {
  createTechCatalogEntry,
  fetchTechCatalog,
  updateTechCatalogEntry,
} from "../api";
import { SortableList } from "../components/SortableList";
import type { TechCatalogEntry, TechCatalogInput, TechLevel } from "../types";

type EditingState = {
  id: number | null;
  slug: string;
  displayName: string;
  category: string;
  level: TechLevel;
  icon: string;
  sortOrder: number;
  active: boolean;
};

const TECH_LEVEL_OPTIONS: { value: TechLevel; label: string }[] = [
  { value: "beginner", label: "Beginner" },
  { value: "intermediate", label: "Intermediate" },
  { value: "advanced", label: "Advanced" },
];

function toEditingState(entry: TechCatalogEntry): EditingState {
  return {
    id: entry.id,
    slug: entry.slug,
    displayName: entry.displayName,
    category: entry.category ?? "",
    level: entry.level,
    icon: entry.icon ?? "",
    sortOrder: entry.sortOrder,
    active: entry.active,
  };
}

function createEmptyState(sortOrder: number): EditingState {
  return {
    id: null,
    slug: "",
    displayName: "",
    category: "",
    level: "intermediate",
    icon: "",
    sortOrder,
    active: true,
  };
}

export function TechCatalogPanel() {
  const [entries, setEntries] = useState<TechCatalogEntry[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [orderedDraft, setOrderedDraft] =
    useState<TechCatalogEntry[] | null>(null);
  const [orderSaving, setOrderSaving] = useState(false);
  const [formError, setFormError] = useState<string | null>(null);
  const [editing, setEditing] = useState<EditingState>(() =>
    createEmptyState(1),
  );

  useEffect(() => {
    let mounted = true;
    const run = async () => {
      setLoading(true);
      setError(null);
      try {
        const data = await fetchTechCatalog();
        if (mounted) {
          const sorted = [...data].sort((a, b) => a.sortOrder - b.sortOrder);
          setEntries(sorted);
          setEditing(createEmptyState((sorted.at(-1)?.sortOrder ?? 0) + 1));
        }
      } catch (err) {
        if (mounted) {
          setError(
            err instanceof Error
              ? err.message
              : "Failed to load tech catalog entries.",
          );
        }
      } finally {
        if (mounted) {
          setLoading(false);
        }
      }
    };
    void run();
    return () => {
      mounted = false;
    };
  }, []);

  const activeOrderedList = orderedDraft ?? entries;

  const isOrderDirty = orderedDraft !== null;
  useUnsavedChangesTracker("admin-techcatalog-order", isOrderDirty);

  const handleReorder = (list: TechCatalogEntry[]) => {
    const normalized = list.map((entry, index) => ({
      ...entry,
      sortOrder: index + 1,
    }));
    setOrderedDraft(normalized);
  };

  const handleSaveOrder = async () => {
    if (!orderedDraft) {
      return;
    }
    setOrderSaving(true);
    setFormError(null);
    try {
      await Promise.all(
        orderedDraft.map((entry) =>
          updateTechCatalogEntry(entry.id, { sortOrder: entry.sortOrder }),
        ),
      );
      setEntries(orderedDraft);
      setOrderedDraft(null);
    } catch (err) {
      setFormError(
        err instanceof Error ? err.message : "Failed to update order.",
      );
    } finally {
      setOrderSaving(false);
    }
  };

  const originalEntry = useMemo(() => {
    if (editing.id === null) {
      return null;
    }
    return entries.find((entry) => entry.id === editing.id) ?? null;
  }, [editing.id, entries]);

  const isFormDirty = useMemo(() => {
    if (editing.id === null) {
      return Boolean(
        editing.slug.trim() &&
          editing.displayName.trim() &&
          editing.category !== "" &&
          editing.icon !== "",
      );
    }
    if (!originalEntry) return false;
    return (
      editing.slug.trim() !== originalEntry.slug.trim() ||
      editing.displayName.trim() !== originalEntry.displayName.trim() ||
      editing.category.trim() !== (originalEntry.category ?? "").trim() ||
      editing.icon.trim() !== (originalEntry.icon ?? "").trim() ||
      editing.level !== originalEntry.level ||
      editing.sortOrder !== originalEntry.sortOrder ||
      editing.active !== originalEntry.active
    );
  }, [editing, originalEntry]);

  useUnsavedChangesTracker("admin-techcatalog-form", isFormDirty);

  const resetForm = () => {
    setEditing(createEmptyState((entries.at(-1)?.sortOrder ?? 0) + 1));
  };

  const submitForm = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setFormError(null);

    if (!editing.slug.trim() || !editing.displayName.trim()) {
      setFormError("Slug and display name are required.");
      return;
    }

    const payload: TechCatalogInput = {
      slug: editing.slug.trim(),
      displayName: editing.displayName.trim(),
      category: editing.category.trim() || undefined,
      level: editing.level,
      icon: editing.icon.trim() || undefined,
      sortOrder: editing.sortOrder,
      active: editing.active,
    };

    try {
      if (editing.id === null) {
        const created = await createTechCatalogEntry(payload);
        const updated = [...entries, created].sort(
          (a, b) => a.sortOrder - b.sortOrder,
        );
        setEntries(updated);
        resetForm();
      } else {
        const updated = await updateTechCatalogEntry(editing.id, payload);
        setEntries((prev) =>
          prev
            .map((entry) => (entry.id === updated.id ? updated : entry))
            .sort((a, b) => a.sortOrder - b.sortOrder),
        );
        setEditing(toEditingState(updated));
      }
    } catch (err) {
      setFormError(
        err instanceof Error ? err.message : "Failed to persist entry.",
      );
    }
  };

  const handleEditSelect = (entry: TechCatalogEntry) => {
    setEditing(toEditingState(entry));
  };

  if (loading) {
    return (
      <div className="space-y-4">
        <div className="h-4 w-32 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
        <div className="h-48 animate-pulse rounded-xl bg-slate-100 dark:bg-slate-800" />
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
    <div className="flex flex-col gap-6 lg:flex-row">
      <form
        onSubmit={submitForm}
        className="flex w-full flex-col gap-4 rounded-2xl border border-slate-200 bg-white/70 p-5 shadow-sm dark:border-slate-700 dark:bg-slate-900/50 lg:w-96"
      >
        <header className="flex items-start justify-between">
          <div>
            <h2 className="text-lg font-semibold text-slate-900 dark:text-slate-100">
              {editing.id === null ? "Add technology" : "Edit technology"}
            </h2>
            <p className="text-sm text-slate-600 dark:text-slate-300">
              Maintain the canonical technology catalog used across projects and profiles.
            </p>
          </div>
          <button
            type="button"
            className="rounded-full border border-slate-300 px-3 py-1 text-xs font-semibold text-slate-600 transition hover:border-sky-400 hover:text-sky-600 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-sky-500 dark:border-slate-700 dark:text-slate-300 dark:hover:border-sky-400 dark:hover:text-sky-300 dark:focus-visible:outline-sky-400"
            onClick={resetForm}
          >
            New
          </button>
        </header>

        <label className="flex flex-col gap-2 text-sm text-slate-700 dark:text-slate-300">
          <span className="font-semibold">Slug</span>
          <input
            value={editing.slug}
            onChange={(event) =>
              setEditing((prev) => ({ ...prev, slug: event.target.value }))
            }
            placeholder="e.g. nextjs"
            className="rounded-md border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 shadow-sm transition focus:border-sky-500 focus:outline-none focus:ring-2 focus:ring-sky-500 dark:border-slate-700 dark:bg-slate-900 dark:text-slate-100"
            required
          />
        </label>

        <label className="flex flex-col gap-2 text-sm text-slate-700 dark:text-slate-300">
          <span className="font-semibold">Display name</span>
          <input
            value={editing.displayName}
            onChange={(event) =>
              setEditing((prev) => ({
                ...prev,
                displayName: event.target.value,
              }))
            }
            placeholder="Next.js"
            className="rounded-md border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 shadow-sm transition focus:border-sky-500 focus:outline-none focus:ring-2 focus:ring-sky-500 dark:border-slate-700 dark:bg-slate-900 dark:text-slate-100"
            required
          />
        </label>

        <label className="flex flex-col gap-2 text-sm text-slate-700 dark:text-slate-300">
          <span className="font-semibold">Category</span>
          <input
            value={editing.category}
            onChange={(event) =>
              setEditing((prev) => ({ ...prev, category: event.target.value }))
            }
            placeholder="Framework"
            className="rounded-md border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 shadow-sm transition focus:border-sky-500 focus:outline-none focus:ring-2 focus:ring-sky-500 dark:border-slate-700 dark:bg-slate-900 dark:text-slate-100"
          />
        </label>

        <label className="flex flex-col gap-2 text-sm text-slate-700 dark:text-slate-300">
          <span className="font-semibold">Level</span>
          <select
            value={editing.level}
            onChange={(event) =>
              setEditing((prev) => ({
                ...prev,
                level: event.target.value as TechLevel,
              }))
            }
            className="h-9 rounded-md border border-slate-300 bg-white px-3 text-sm text-slate-900 shadow-sm transition focus:border-sky-500 focus:outline-none focus:ring-2 focus:ring-sky-500 dark:border-slate-700 dark:bg-slate-900 dark:text-slate-100"
          >
            {TECH_LEVEL_OPTIONS.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
        </label>

        <label className="flex flex-col gap-2 text-sm text-slate-700 dark:text-slate-300">
          <span className="font-semibold">Icon</span>
          <input
            value={editing.icon}
            onChange={(event) =>
              setEditing((prev) => ({ ...prev, icon: event.target.value }))
            }
            placeholder="Emoji or icon URL"
            className="rounded-md border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 shadow-sm transition focus:border-sky-500 focus:outline-none focus:ring-2 focus:ring-sky-500 dark:border-slate-700 dark:bg-slate-900 dark:text-slate-100"
          />
        </label>

        <label className="flex flex-col gap-2 text-sm text-slate-700 dark:text-slate-300">
          <span className="font-semibold">Sort order</span>
          <input
            type="number"
            min={0}
            value={editing.sortOrder}
            onChange={(event) =>
              setEditing((prev) => ({
                ...prev,
                sortOrder: Number.parseInt(event.target.value, 10) || 0,
              }))
            }
            className="rounded-md border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 shadow-sm transition focus:border-sky-500 focus:outline-none focus:ring-2 focus:ring-sky-500 dark:border-slate-700 dark:bg-slate-900 dark:text-slate-100"
          />
        </label>

        <label className="inline-flex items-center gap-2 text-sm text-slate-700 dark:text-slate-300">
          <input
            type="checkbox"
            checked={editing.active}
            onChange={(event) =>
              setEditing((prev) => ({ ...prev, active: event.target.checked }))
            }
            className="h-4 w-4 rounded border-slate-300 text-sky-600 focus:ring-sky-500 dark:border-slate-700 dark:bg-slate-900 dark:checked:bg-sky-400"
          />
          Active
        </label>

        {formError && (
          <p className="rounded-md border border-rose-300 bg-rose-50 px-3 py-2 text-sm text-rose-700 dark:border-rose-900 dark:bg-rose-950/30 dark:text-rose-300">
            {formError}
          </p>
        )}

        <button
          type="submit"
          className="inline-flex items-center justify-center rounded-md bg-sky-600 px-4 py-2 text-sm font-semibold text-white shadow-sm transition hover:bg-sky-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-sky-500 dark:bg-sky-500 dark:hover:bg-sky-400"
        >
          {editing.id === null ? "Add technology" : "Save changes"}
        </button>
      </form>

      <section className="flex-1 space-y-4 rounded-2xl border border-slate-200 bg-white/70 p-5 shadow-sm dark:border-slate-700 dark:bg-slate-900/50">
        <header className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h2 className="text-lg font-semibold text-slate-900 dark:text-slate-100">
              Catalog entries
            </h2>
            <p className="text-sm text-slate-600 dark:text-slate-300">
              Drag to reorder. Sort order updates apply across the public pages.
            </p>
          </div>
          <button
            type="button"
            disabled={!isOrderDirty || orderSaving}
            className="inline-flex items-center rounded-md bg-slate-900 px-4 py-2 text-sm font-semibold text-white shadow-sm transition enabled:hover:bg-slate-700 enabled:focus-visible:outline enabled:focus-visible:outline-2 enabled:focus-visible:outline-offset-2 enabled:focus-visible:outline-slate-900 disabled:cursor-not-allowed disabled:opacity-60 dark:bg-slate-100 dark:text-slate-900 dark:enabled:hover:bg-slate-200"
            onClick={handleSaveOrder}
          >
            {orderSaving ? "Saving…" : "Save order"}
          </button>
        </header>

        <SortableList
          items={activeOrderedList}
          getId={(entry) => entry.id}
          onReorder={handleReorder}
          renderItem={(entry) => (
            <div className="flex flex-col gap-2 rounded-xl border border-slate-200 bg-white/90 p-4 shadow-sm transition hover:border-sky-300 dark:border-slate-700 dark:bg-slate-900/70">
              <div className="flex items-start justify-between gap-3">
                <div>
                  <p className="text-sm font-semibold text-slate-900 dark:text-slate-100">
                    {entry.displayName}{" "}
                    <span className="text-xs text-slate-500 dark:text-slate-400">
                      ({entry.slug})
                    </span>
                  </p>
                  <p className="text-xs text-slate-500 dark:text-slate-400">
                    Level: {entry.level} · Sort: {entry.sortOrder} ·{" "}
                    {entry.active ? "Active" : "Inactive"}
                  </p>
                  {entry.category && (
                    <p className="text-xs text-slate-500 dark:text-slate-400">
                      {entry.category}
                    </p>
                  )}
                </div>
                <button
                  type="button"
                  className="rounded-full border border-slate-300 px-3 py-1 text-xs font-semibold text-slate-600 transition hover:border-sky-400 hover:text-sky-600 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-sky-500 dark:border-slate-700 dark:text-slate-300 dark:hover:border-sky-400 dark:hover:text-sky-300 dark:focus-visible:outline-sky-400"
                  onClick={() => handleEditSelect(entry)}
                >
                  Edit
                </button>
              </div>
            </div>
          )}
        />
      </section>
    </div>
  );
}
