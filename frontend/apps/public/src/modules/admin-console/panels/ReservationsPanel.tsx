import { useEffect, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";

import {
  fetchReservations,
  retryReservationNotification,
  updateReservationStatus,
} from "../api";
import type { Reservation, ReservationStatus } from "../types";

const STATUS_OPTIONS: { value: ReservationStatus; label: string }[] = [
  { value: "pending", label: "Pending" },
  { value: "confirmed", label: "Confirmed" },
  { value: "cancelled", label: "Cancelled" },
];

function formatDateTime(value: string, locale: string): string {
  const parsed = Date.parse(value);
  if (Number.isNaN(parsed)) {
    return value;
  }
  return new Intl.DateTimeFormat(locale, {
    dateStyle: "medium",
    timeStyle: "short",
  }).format(parsed);
}

function formatStatus(status: ReservationStatus): string {
  switch (status) {
    case "pending":
      return "Pending";
    case "confirmed":
      return "Confirmed";
    case "cancelled":
      return "Cancelled";
    default:
      return status;
  }
}

export function ReservationsPanel() {
  const { i18n } = useTranslation();
  const [reservations, setReservations] = useState<Reservation[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [mutationError, setMutationError] = useState<string | null>(null);

  useEffect(() => {
    let mounted = true;
    const run = async () => {
      setLoading(true);
      setError(null);
      try {
        const data = await fetchReservations();
        if (mounted) {
          setReservations(data);
        }
      } catch (err) {
        if (mounted) {
          setError(
            err instanceof Error ? err.message : "Failed to load reservations.",
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

  const sortedReservations = useMemo(() => {
    return [...reservations].sort((a, b) =>
      a.startAt < b.startAt ? 1 : a.startAt > b.startAt ? -1 : 0,
    );
  }, [reservations]);

  const handleStatusChange = async (
    reservation: Reservation,
    status: ReservationStatus,
  ) => {
    if (reservation.status === status) {
      return;
    }
    setMutationError(null);
    try {
      const updated = await updateReservationStatus(reservation.id, {
        status,
        cancellationReason:
          status === "cancelled" ? reservation.cancellationReason ?? "" : "",
      });
      setReservations((prev) =>
        prev.map((entry) => (entry.id === updated.id ? updated : entry)),
      );
    } catch (err) {
      setMutationError(
        err instanceof Error ? err.message : "Failed to update reservation.",
      );
    }
  };

  const handleRetry = async (reservation: Reservation) => {
    setMutationError(null);
    try {
      const updated = await retryReservationNotification(reservation.id);
      setReservations((prev) =>
        prev.map((entry) => (entry.id === updated.id ? updated : entry)),
      );
    } catch (err) {
      setMutationError(
        err instanceof Error ? err.message : "Failed to trigger retry.",
      );
    }
  };

  if (loading) {
    return (
      <div className="space-y-4">
        <div className="h-4 w-40 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
        <div className="space-y-2">
          {Array.from({ length: 4 }).map((_, index) => (
            <div
              // eslint-disable-next-line react/no-array-index-key
              key={index}
              className="h-16 animate-pulse rounded-lg bg-slate-100 dark:bg-slate-800"
            />
          ))}
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="rounded-lg border border-rose-200 bg-rose-50 p-4 text-sm text-rose-700 dark:border-rose-900 dark:bg-rose-950/40 dark:text-rose-300">
        {error}
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {mutationError && (
        <div className="rounded-md border border-rose-300 bg-rose-50 p-3 text-sm text-rose-700 dark:border-rose-900 dark:bg-rose-950/40 dark:text-rose-300">
          {mutationError}
        </div>
      )}
      <ul className="space-y-4">
        {sortedReservations.map((reservation) => (
          <li
            key={reservation.id}
            className="rounded-2xl border border-slate-200 bg-white/70 p-5 shadow-sm transition hover:border-sky-300 dark:border-slate-700 dark:bg-slate-900/50"
          >
            <div className="flex flex-col gap-2 md:flex-row md:items-start md:justify-between">
              <div>
                <h3 className="text-base font-semibold text-slate-900 dark:text-slate-100">
                  {reservation.name}{" "}
                  <span className="text-sm font-normal text-slate-500 dark:text-slate-400">
                    ({reservation.email})
                  </span>
                </h3>
                <p className="text-sm text-slate-600 dark:text-slate-300">
                  {formatDateTime(reservation.startAt, i18n.language)} →
                  {formatDateTime(reservation.endAt, i18n.language)}
                </p>
                {reservation.topic && (
                  <p className="text-sm text-slate-500 dark:text-slate-400">
                    Topic: {reservation.topic}
                  </p>
                )}
                {reservation.message && (
                  <p className="text-sm text-slate-500 dark:text-slate-400">
                    {reservation.message}
                  </p>
                )}
              </div>
              <div className="flex items-center gap-3">
                <label className="flex flex-col gap-1 text-sm text-slate-600 dark:text-slate-300">
                  <span className="font-semibold">Status</span>
                  <select
                    value={reservation.status}
                    onChange={(event) =>
                      handleStatusChange(
                        reservation,
                        event.target.value as ReservationStatus,
                      )
                    }
                    className="h-9 rounded-md border border-slate-300 bg-white px-3 text-sm font-medium text-slate-700 shadow-sm transition hover:border-sky-400 focus:border-sky-500 focus:outline-none focus:ring-2 focus:ring-sky-500 dark:border-slate-700 dark:bg-slate-900 dark:text-slate-100"
                  >
                    {STATUS_OPTIONS.map((option) => (
                      <option key={option.value} value={option.value}>
                        {option.label}
                      </option>
                    ))}
                  </select>
                </label>
                <button
                  type="button"
                  className="inline-flex h-9 items-center rounded-md border border-slate-300 px-3 text-sm font-medium text-slate-700 transition hover:border-sky-400 hover:text-sky-600 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-sky-500 dark:border-slate-700 dark:text-slate-200 dark:hover:border-sky-400 dark:hover:text-sky-300 dark:focus-visible:outline-sky-400"
                  onClick={() => handleRetry(reservation)}
                >
                  Retry notification
                </button>
              </div>
            </div>
            <dl className="mt-3 grid gap-2 text-xs text-slate-500 dark:text-slate-400 md:grid-cols-3">
              <div>
                <dt className="font-semibold uppercase tracking-wide">
                  Calendar
                </dt>
                <dd>
                  {reservation.googleCalendarStatus
                    ? reservation.googleCalendarStatus
                    : "—"}
                </dd>
              </div>
              <div>
                <dt className="font-semibold uppercase tracking-wide">
                  Lookup hash
                </dt>
                <dd>{reservation.lookupHash}</dd>
              </div>
              <div>
                <dt className="font-semibold uppercase tracking-wide">
                  Current status
                </dt>
                <dd>{formatStatus(reservation.status)}</dd>
              </div>
            </dl>
          </li>
        ))}
      </ul>
    </div>
  );
}
