import { FormEvent, useEffect, useMemo, useState } from "react";

import { useUnsavedChangesTracker } from "../../../hooks/useUnsavedChangesTracker";
import {
  createBlacklistEntry,
  deleteBlacklistEntry,
  fetchBlacklist,
  updateBlacklistEntry,
} from "../api";
import type { BlacklistEntry } from "../types";

type EditingState = {
  id: number | null;
  email: string;
  reason: string;
};

function createEmptyEditingState(): EditingState {
  return {
    id: null,
    email: "",
    reason: "",
  };
}

export function BlacklistPanel() {
  const [entries, setEntries] = useState<BlacklistEntry[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [formError, setFormError] = useState<string | null>(null);
  const [editing, setEditing] = useState<EditingState>(() =>
    createEmptyEditingState(),
  );

  useEffect(() => {
    let active = true;
    const run = async () => {
      setLoading(true);
      setError(null);
      try {
        const data = await fetchBlacklist();
        if (active) {
          setEntries(data);
        }
      } catch (err) {
        if (active) {
          setError(
            err instanceof Error ? err.message : "Failed to load blacklist entries.",
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

  const originalEntry = useMemo(() => {
    if (editing.id === null) {
      return null;
    }
    return entries.find((entry) => entry.id === editing.id) ?? null;
  }, [editing.id, entries]);

  const isDirty = useMemo(() => {
    if (editing.id === null) {
      return editing.email.trim().length > 0 || editing.reason.trim().length > 0;
    }
    if (!originalEntry) {
      return editing.email.trim().length > 0 || editing.reason.trim().length > 0;
    }
    return (
      editing.email.trim() !== originalEntry.email.trim() ||
      editing.reason.trim() !== (originalEntry.reason ?? "").trim()
    );
  }, [editing, originalEntry]);

  useUnsavedChangesTracker("admin-blacklist-form", isDirty);

  const resetForm = () => {
    setEditing(createEmptyEditingState());
  };

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setFormError(null);

    if (!editing.email.trim()) {
      setFormError("Email is required.");
      return;
    }

    try {
      if (editing.id === null) {
        const created = await createBlacklistEntry({
          email: editing.email.trim(),
          reason: editing.reason.trim(),
        });
        setEntries((prev) => [created, ...prev]);
        resetForm();
      } else {
        const updated = await updateBlacklistEntry(editing.id, {
          email: editing.email.trim(),
          reason: editing.reason.trim(),
        });
        setEntries((prev) =>
          prev.map((entry) => (entry.id === updated.id ? updated : entry)),
        );
        setEditing({
          id: updated.id,
          email: updated.email,
          reason: updated.reason ?? "",
        });
      }
    } catch (err) {
      setFormError(
        err instanceof Error
          ? err.message
          : "Unable to persist blacklist entry. Try again later.",
      );
    }
  };

  const handleEditSelect = (entry: BlacklistEntry) => {
    setEditing({
      id: entry.id,
      email: entry.email,
      reason: entry.reason ?? "",
    });
  };

  const handleDelete = async (entry: BlacklistEntry) => {
    if (!window.confirm(`Remove ${entry.email} from the blacklist?`)) {
      return;
    }
    try {
      await deleteBlacklistEntry(entry.id);
      setEntries((prev) => prev.filter((candidate) => candidate.id !== entry.id));
      if (editing.id === entry.id) {
        resetForm();
      }
    } catch (err) {
      setFormError(
        err instanceof Error
          ? err.message
          : "Failed to delete blacklist entry.",
      );
    }
  };

  const handleNewEntry = () => {
    resetForm();
  };

  if (loading) {
    return (
      <div className="space-y-4">
        <div className="h-4 w-32 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
        <div className="h-32 animate-pulse rounded-lg bg-slate-100 dark:bg-slate-800" />
      </div>
    );
  }

  return (
    <div className="grid gap-6 lg:grid-cols-[minmax(0,2fr),minmax(0,3fr)]">
      <form
        onSubmit={handleSubmit}
        className="flex flex-col gap-4 rounded-2xl border border-slate-200 bg-white/70 p-5 shadow-sm dark:border-slate-700 dark:bg-slate-900/50"
      >
        <header className="flex items-start justify-between gap-3">
          <div>
            <h2 className="text-lg font-semibold text-slate-900 dark:text-slate-100">
              {editing.id === null ? "Add blacklist entry" : "Edit blacklist entry"}
            </h2>
            <p className="text-sm text-slate-600 dark:text-slate-300">
              Prevent specific email addresses from booking meetings.
            </p>
          </div>
          <button
            type="button"
            className="rounded-full border border-slate-300 px-3 py-1 text-xs font-medium text-slate-600 transition hover:border-sky-400 hover:text-sky-600 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-sky-500 dark:border-slate-700 dark:text-slate-300 dark:hover:border-sky-400 dark:hover:text-sky-300 dark:focus-visible:outline-sky-400"
            onClick={handleNewEntry}
          >
            New entry
          </button>
        </header>

        <label className="flex flex-col gap-2 text-sm text-slate-700 dark:text-slate-300">
          <span className="font-semibold">Email</span>
          <input
            value={editing.email}
            onChange={(event) =>
              setEditing((prev) => ({ ...prev, email: event.target.value }))
            }
            placeholder="user@example.com"
            className="rounded-md border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 shadow-sm transition focus:border-sky-500 focus:outline-none focus:ring-2 focus:ring-sky-500 dark:border-slate-700 dark:bg-slate-900 dark:text-slate-100"
            required
            type="email"
          />
        </label>

        <label className="flex flex-col gap-2 text-sm text-slate-700 dark:text-slate-300">
          <span className="font-semibold">Reason</span>
          <textarea
            value={editing.reason}
            onChange={(event) =>
              setEditing((prev) => ({ ...prev, reason: event.target.value }))
            }
            rows={4}
            placeholder="Optional context for the block"
            className="rounded-md border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 shadow-sm transition focus:border-sky-500 focus:outline-none focus:ring-2 focus:ring-sky-500 dark:border-slate-700 dark:bg-slate-900 dark:text-slate-100"
          />
        </label>

        {formError && (
          <p className="rounded-md border border-rose-300 bg-rose-50 px-3 py-2 text-sm text-rose-700 dark:border-rose-900 dark:bg-rose-950/40 dark:text-rose-300">
            {formError}
          </p>
        )}

        <div className="flex items-center justify-between">
          <button
            type="submit"
            disabled={!isDirty}
            className="inline-flex items-center rounded-md bg-sky-600 px-4 py-2 text-sm font-semibold text-white shadow-sm transition enabled:hover:bg-sky-500 enabled:focus-visible:outline enabled:focus-visible:outline-2 enabled:focus-visible:outline-offset-2 enabled:focus-visible:outline-sky-500 disabled:cursor-not-allowed disabled:opacity-60 dark:bg-sky-500 dark:enabled:hover:bg-sky-400"
          >
            {editing.id === null ? "Add entry" : "Save changes"}
          </button>
          {editing.id !== null && (
            <button
              type="button"
              className="rounded-md border border-rose-400 px-4 py-2 text-sm font-semibold text-rose-600 transition hover:border-rose-500 hover:text-rose-700 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-rose-500 dark:border-rose-800 dark:text-rose-300 dark:hover:border-rose-600 dark:hover:text-rose-200"
              onClick={() => {
                const entry = entries.find((candidate) => candidate.id === editing.id);
                if (entry) {
                  handleDelete(entry);
                }
              }}
            >
              Delete
            </button>
          )}
        </div>
      </form>

      <section className="space-y-3 rounded-2xl border border-slate-200 bg-white/70 p-5 shadow-sm dark:border-slate-700 dark:bg-slate-900/50">
        <header>
          <h2 className="text-lg font-semibold text-slate-900 dark:text-slate-100">
            Blocked emails
          </h2>
          <p className="text-sm text-slate-600 dark:text-slate-300">
            Entries are evaluated server-side when a booking request is submitted.
          </p>
        </header>

        {error && (
          <p className="rounded-md border border-rose-300 bg-rose-50 px-3 py-2 text-sm text-rose-700 dark:border-rose-900 dark:bg-rose-950/40 dark:text-rose-300">
            {error}
          </p>
        )}

        <ul className="space-y-3">
          {entries.map((entry) => (
            <li
              key={entry.id}
              className="flex flex-col gap-2 rounded-xl border border-slate-200 bg-white/80 p-4 shadow-sm transition hover:border-sky-300 dark:border-slate-700 dark:bg-slate-900/60"
            >
              <div className="flex items-start justify-between gap-3">
                <div>
                  <p className="text-sm font-semibold text-slate-900 dark:text-slate-100">
                    {entry.email}
                  </p>
                  {entry.reason && (
                    <p className="text-sm text-slate-600 dark:text-slate-300">
                      {entry.reason}
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
              <p className="text-xs uppercase tracking-wide text-slate-500 dark:text-slate-400">
                Added {new Date(entry.createdAt).toLocaleString()}
              </p>
            </li>
          ))}
        </ul>
      </section>
    </div>
  );
}
