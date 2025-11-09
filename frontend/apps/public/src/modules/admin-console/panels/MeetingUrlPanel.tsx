import { FormEvent, useEffect, useMemo, useState } from "react";

import { useUnsavedChangesTracker } from "../../../hooks/useUnsavedChangesTracker";
import { fetchMeetingUrl, updateMeetingUrl } from "../api";
import type { MeetingUrlTemplate } from "../types";

export function MeetingUrlPanel() {
  const [template, setTemplate] = useState<MeetingUrlTemplate | null>(null);
  const [draft, setDraft] = useState("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);
  const [saveError, setSaveError] = useState<string | null>(null);

  useEffect(() => {
    let mounted = true;
    const run = async () => {
      setLoading(true);
      setError(null);
      try {
        const data = await fetchMeetingUrl();
        if (mounted) {
          setTemplate(data);
          setDraft(data.template ?? "");
        }
      } catch (err) {
        if (mounted) {
          setError(
            err instanceof Error
              ? err.message
              : "Failed to load meeting URL template.",
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

  const isDirty = useMemo(() => {
    return draft !== (template?.template ?? "");
  }, [draft, template]);

  useUnsavedChangesTracker("admin-meeting-url", isDirty);

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setSaveError(null);
    setSaving(true);
    try {
      const updated = await updateMeetingUrl({ template: draft });
      setTemplate(updated);
      setDraft(updated.template ?? "");
    } catch (err) {
      setSaveError(
        err instanceof Error ? err.message : "Failed to update template.",
      );
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <div className="h-40 animate-pulse rounded-xl bg-slate-100 dark:bg-slate-800" />
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
    <form
      onSubmit={handleSubmit}
      className="space-y-4 rounded-2xl border border-slate-200 bg-white/70 p-5 shadow-sm dark:border-slate-700 dark:bg-slate-900/50"
    >
      <header className="flex flex-col gap-1">
        <h2 className="text-lg font-semibold text-slate-900 dark:text-slate-100">
          Meeting URL email template
        </h2>
        <p className="text-sm text-slate-600 dark:text-slate-300">
          This template populates the automated confirmation email sent to guests after
          accepting a reservation. Use placeholders like {"{{meeting_url}}"} and {"{{guest_name}}"}.
        </p>
      </header>

      <label className="flex flex-col gap-2 text-sm text-slate-700 dark:text-slate-300">
        <span className="font-semibold">Message body</span>
        <textarea
          value={draft}
          onChange={(event) => setDraft(event.target.value)}
          rows={12}
          className="rounded-md border border-slate-300 bg-white px-3 py-3 text-sm font-mono text-slate-900 shadow-sm transition focus:border-sky-500 focus:outline-none focus:ring-2 focus:ring-sky-500 dark:border-slate-700 dark:bg-slate-900 dark:text-slate-100"
        />
      </label>

      {saveError && (
        <p className="rounded-md border border-rose-300 bg-rose-50 px-3 py-2 text-sm text-rose-700 dark:border-rose-900 dark:bg-rose-950/30 dark:text-rose-300">
          {saveError}
        </p>
      )}

      <button
        type="submit"
        disabled={!isDirty || saving}
        className="inline-flex items-center rounded-md bg-sky-600 px-4 py-2 text-sm font-semibold text-white shadow-sm transition enabled:hover:bg-sky-500 enabled:focus-visible:outline enabled:focus-visible:outline-2 enabled:focus-visible:outline-offset-2 enabled:focus-visible:outline-sky-500 disabled:cursor-not-allowed disabled:opacity-60 dark:bg-sky-500 dark:enabled:hover:bg-sky-400"
      >
        {saving ? "Saving…" : "Save template"}
      </button>
    </form>
  );
}
