import {
  FormEvent,
  useCallback,
  useEffect,
  useId,
  useMemo,
  useState,
} from "react";
import { useTranslation } from "react-i18next";

import {
  publicApi,
  useContactAvailability,
  useContactConfig,
} from "../../modules/public-api";
import type {
  ContactAvailabilityDay,
  ContactAvailabilitySlot,
} from "../../modules/public-api";
import { formatDateTime, formatTime } from "../../utils/date";

const EMPTY_AVAILABILITY_DAYS: ContactAvailabilityDay[] = [];

type FormState = {
  name: string;
  email: string;
  topic: string;
  agenda: string;
  slotId: string;
};

type FormErrors = Partial<Record<keyof FormState, string>>;

declare global {
  interface Window {
    grecaptcha?: {
      ready: (callback: () => void) => void;
      execute: (
        siteKey: string,
        options: { action: string },
      ) => Promise<string>;
    };
  }
}

const initialFormState: FormState = {
  name: "",
  email: "",
  topic: "",
  agenda: "",
  slotId: "",
};

export function ContactPage() {
  const { t, i18n } = useTranslation();
  const agendaFieldId = useId();
  const agendaErrorId = useId();
  const consentId = useId();

  const {
    data: availability,
    isLoading: isAvailabilityLoading,
    error: availabilityError,
  } = useContactAvailability();
  const {
    data: config,
    isLoading: isConfigLoading,
    error: configError,
  } = useContactConfig();

  type SlotWithContext = ContactAvailabilitySlot & { day: string };

  const days = availability?.days ?? EMPTY_AVAILABILITY_DAYS;

  const { dayOrder, timeKeys, slotByDay, slotById } = useMemo(() => {
    const dayOrder = days.map((day) => day.date);
    const timeSet = new Set<number>();
    const slotByDay = new Map<string, Map<number, SlotWithContext>>();
    const slotById = new Map<string, SlotWithContext>();

    days.forEach((day) => {
      const slotsForDay = new Map<number, SlotWithContext>();
      day.slots.forEach((slot) => {
        const slotId = slot.id || slot.start;
        const startKey = new Date(slot.start).getTime();
        const enriched: SlotWithContext = {
          ...slot,
          id: slotId,
          day: day.date,
        };
        slotsForDay.set(startKey, enriched);
        slotById.set(slotId, enriched);
        timeSet.add(startKey);
      });
      slotByDay.set(day.date, slotsForDay);
    });

    const timeKeys = Array.from(timeSet).sort((a, b) => a - b);

    return { dayOrder, timeKeys, slotByDay, slotById };
  }, [days]);

  const [viewMode, setViewMode] = useState<"single" | "multi">("single");
  const [activeDayIndex, setActiveDayIndex] = useState(0);
  const maxMultiColumns = 3;

  useEffect(() => {
    setActiveDayIndex(0);
  }, [dayOrder.length]);

  const displayedDayIndices = useMemo(() => {
    if (dayOrder.length === 0) {
      return [];
    }
    if (viewMode === "single") {
      return [Math.min(activeDayIndex, dayOrder.length - 1)];
    }
    const start = Math.min(activeDayIndex, Math.max(0, dayOrder.length - 1));
    const end = Math.min(start + maxMultiColumns, dayOrder.length);
    return Array.from({ length: end - start }, (_value, index) => start + index);
  }, [activeDayIndex, dayOrder.length, viewMode, maxMultiColumns]);

  const displayedDays = displayedDayIndices
    .map((index) => dayOrder[index])
    .filter((value): value is string => Boolean(value));

  const canGoPrev = activeDayIndex > 0;
  const canGoNext =
    viewMode === "single"
      ? activeDayIndex < dayOrder.length - 1
      : activeDayIndex + maxMultiColumns < dayOrder.length;

  const goPrev = useCallback(() => {
    setActiveDayIndex((index) => Math.max(0, index - 1));
  }, []);

  const goNext = useCallback(() => {
    setActiveDayIndex((index) =>
      Math.min(dayOrder.length > 0 ? dayOrder.length - 1 : 0, index + 1),
    );
  }, [dayOrder.length]);

  const topics = useMemo(
    () => config?.topics ?? [],
    [config?.topics],
  );

  const [formState, setFormState] = useState<FormState>(initialFormState);
  const [formErrors, setFormErrors] = useState<FormErrors>({});
  const [statusMessage, setStatusMessage] = useState<string | null>(null);
  const [bookingResult, setBookingResult] = useState<
    Awaited<ReturnType<typeof publicApi.createBooking>> | null
  >(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const timezoneLabel = useMemo(
    () => availability?.timezone ?? config?.calendarTimezone ?? "",
    [availability?.timezone, config?.calendarTimezone],
  );

  const selectedTopic = useMemo(
    () => topics.find((topic) => topic.id === formState.topic) ?? null,
    [topics, formState.topic],
  );

  const validate = useCallback(
    (state: FormState): FormErrors => {
      const errors: FormErrors = {};

      if (!state.name.trim()) {
        errors.name = t("contact.form.errors.nameRequired");
      }

      if (!state.email.trim()) {
        errors.email = t("contact.form.errors.emailRequired");
      } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(state.email.trim())) {
        errors.email = t("contact.form.errors.emailInvalid");
      }

      if (!state.topic) {
        errors.topic = t("contact.form.errors.topicRequired");
      }

      if (!state.agenda.trim() || state.agenda.trim().length < 20) {
        errors.agenda = t("contact.form.errors.messageLength");
      }

      if (!state.slotId) {
        errors.slotId = t("contact.form.errors.slotRequired");
      } else {
        const slot = slotById.get(state.slotId);
        if (!slot || slot.status !== "available" || !slot.isBookable) {
          errors.slotId = t("contact.form.errors.slotUnavailable");
        }
      }

      return errors;
    },
    [slotById, t],
  );

  const handleInputChange =
    (field: keyof FormState) =>
    (
      event: React.ChangeEvent<
        HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement
      >,
    ) => {
      setFormState((previous) => ({
        ...previous,
        [field]: event.target.value,
      }));
    };

  const loadRecaptchaToken = useCallback(async (): Promise<string> => {
    if (!config?.recaptchaSiteKey) {
      return "";
    }

    await new Promise<void>((resolve) => {
      if (window.grecaptcha) {
        window.grecaptcha.ready(resolve);
        return;
      }

      const scriptId = "google-recaptcha-script";
      if (!document.getElementById(scriptId)) {
        const script = document.createElement("script");
        script.id = scriptId;
        script.src = `https://www.google.com/recaptcha/api.js?render=${config.recaptchaSiteKey}`;
        script.async = true;
        script.defer = true;
        script.onload = () => resolve();
        script.onerror = () => resolve();
        document.body.append(script);
      } else {
        resolve();
      }
    });

    if (!window.grecaptcha) {
      return "";
    }

    return window.grecaptcha.execute(config.recaptchaSiteKey, {
      action: "submit",
    });
  }, [config?.recaptchaSiteKey]);

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setStatusMessage(null);
    setBookingResult(null);

    const errors = validate(formState);
    setFormErrors(errors);

    if (Object.keys(errors).length > 0) {
      return;
    }

    try {
      setIsSubmitting(true);
      const recaptchaToken = await loadRecaptchaToken();
      const selectedSlot = slotById.get(formState.slotId);
      if (!selectedSlot) {
        setFormErrors({ slotId: t("contact.form.errors.slotRequired") });
        return;
      }
      if (selectedSlot.status !== "available" || !selectedSlot.isBookable) {
        setFormErrors({ slotId: t("contact.form.errors.slotUnavailable") });
        return;
      }

      const slotStart = new Date(selectedSlot.start);
      const slotEnd = new Date(selectedSlot.end);
      const startTime = slotStart.toISOString();
      const durationMinutes = Math.max(
        1,
        Math.round((slotEnd.getTime() - slotStart.getTime()) / 60000),
      );

      const response = await publicApi.createBooking({
        name: formState.name.trim(),
        email: formState.email.trim(),
        topic: formState.topic,
        agenda: formState.agenda.trim(),
        startTime,
        durationMinutes,
        recaptchaToken,
      });

      setBookingResult(response);
      setStatusMessage(
        t("contact.form.success", {
          id: response.reservation.id,
          email: response.supportEmail ?? response.reservation.email,
        }),
      );
      setFormState(initialFormState);
    } catch {
      setStatusMessage(t("contact.form.error"));
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleSlotPick = useCallback(
    (slotId: string) => {
      const slot = slotById.get(slotId);
      if (!slot || slot.status !== "available" || !slot.isBookable) {
        return;
      }
      setFormState((previous) => ({
        ...previous,
        slotId,
      }));
      setFormErrors((previous) => {
        if (!previous.slotId) {
          return previous;
        }
        const next = { ...previous };
        delete next.slotId;
        return next;
      });
    },
    [slotById],
  );

  return (
    <section className="mx-auto flex w-full max-w-5xl flex-col gap-8 px-4 py-12 sm:px-8">
      <header className="space-y-3 text-center md:text-left">
        <p className="text-sm font-semibold uppercase tracking-[0.3em] text-slate-500 dark:text-slate-400">
          {t("contact.tagline")}
        </p>
        <h1 className="text-3xl font-bold text-slate-900 dark:text-slate-50 sm:text-4xl">
          {config?.heroTitle ?? t("contact.title")}
        </h1>
        <p className="text-base text-slate-600 dark:text-slate-300">
          {config?.heroDescription ?? t("contact.description")}
        </p>
      </header>

      <div className="grid gap-8 lg:grid-cols-[3fr,2fr]">
        <form
          onSubmit={handleSubmit}
          noValidate
          className="space-y-6 rounded-2xl border border-slate-200 bg-white/80 p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900/60"
        >
          <div className="grid gap-4 sm:grid-cols-2">
            <label className="flex flex-col gap-2">
              <span className="text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-400">
                {t("contact.form.name")}
              </span>
              <input
                type="text"
                name="name"
                value={formState.name}
                onChange={handleInputChange("name")}
                required
                className="rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 focus:border-sky-400 focus:outline-none focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-900 dark:text-slate-100 dark:focus:border-sky-500 dark:focus:ring-sky-900/30"
              />
              {formErrors.name ? (
                <span className="text-xs text-rose-500 dark:text-rose-400">
                  {formErrors.name}
                </span>
              ) : null}
            </label>
            <label className="flex flex-col gap-2">
              <span className="text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-400">
                {t("contact.form.email")}
              </span>
              <input
                type="email"
                name="email"
                value={formState.email}
                onChange={handleInputChange("email")}
                required
                className="rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 focus:border-sky-400 focus:outline-none focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-900 dark:text-slate-100 dark:focus:border-sky-500 dark:focus:ring-sky-900/30"
              />
              {formErrors.email ? (
                <span className="text-xs text-rose-500 dark:text-rose-400">
                  {formErrors.email}
                </span>
              ) : null}
            </label>
          </div>

          <label className="flex flex-col gap-2">
            <span className="text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-400">
              {t("contact.form.topic")}
            </span>
            <select
              name="topic"
              value={formState.topic}
              onChange={handleInputChange("topic")}
              required
              className="rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 focus:border-sky-400 focus:outline-none focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-900 dark:text-slate-100 dark:focus:border-sky-500 dark:focus:ring-sky-900/30"
            >
              <option value="">{t("contact.form.topicPlaceholder")}</option>
              {topics.map((topic) => (
                <option key={topic.id} value={topic.id}>
                  {topic.label}
                </option>
              ))}
            </select>
            {selectedTopic?.description ? (
              <span className="text-xs text-slate-500 dark:text-slate-400">
                {selectedTopic.description}
              </span>
            ) : null}
            {formErrors.topic ? (
              <span className="text-xs text-rose-500 dark:text-rose-400">
                {formErrors.topic}
              </span>
            ) : null}
          </label>

          <label className="flex flex-col gap-2">
            <span className="text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-400">
              {t("contact.form.message")}
            </span>
            <textarea
              id={agendaFieldId}
              name="agenda"
              value={formState.agenda}
              onChange={handleInputChange("agenda")}
              required
              aria-describedby={
                formErrors.agenda ? `${consentId} ${agendaErrorId}` : consentId
              }
              rows={6}
              className="rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 focus:border-sky-400 focus:outline-none focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-900 dark:text-slate-100 dark:focus:border-sky-500 dark:focus:ring-sky-900/30"
            />
          </label>
          <span
            id={consentId}
            className="text-xs text-slate-500 dark:text-slate-400"
          >
            {config?.consentText ?? t("contact.form.consent")}
          </span>
          {formErrors.agenda ? (
            <span
              id={agendaErrorId}
              className="text-xs text-rose-500 dark:text-rose-400"
            >
              {formErrors.agenda}
            </span>
          ) : null}

          <div className="flex flex-col gap-3">
            <div className="flex flex-wrap items-center justify-between gap-3">
              <div className="flex flex-col">
                <span className="text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-400">
                  {t("contact.form.slot")}
                </span>
                <span className="text-xs text-slate-500 dark:text-slate-400">
                  {t("contact.form.timezoneLabel", {
                    timezone: timezoneLabel,
                  })}
                </span>
              </div>
              <div className="flex flex-wrap items-center gap-2">
                <div className="inline-flex overflow-hidden rounded-full border border-slate-200 bg-white shadow-sm dark:border-slate-700 dark:bg-slate-900/60">
                  <button
                    type="button"
                    onClick={() => setViewMode("single")}
                    className={`px-3 py-1 text-xs font-semibold ${
                      viewMode === "single"
                        ? "bg-sky-600 text-white dark:bg-sky-500"
                        : "text-slate-600 hover:bg-slate-100 dark:text-slate-300 dark:hover:bg-slate-800/80"
                    }`}
                  >
                    {t("contact.form.view.single")}
                  </button>
                  <button
                    type="button"
                    onClick={() => setViewMode("multi")}
                    className={`px-3 py-1 text-xs font-semibold ${
                      viewMode === "multi"
                        ? "bg-sky-600 text-white dark:bg-sky-500"
                        : "text-slate-600 hover:bg-slate-100 dark:text-slate-300 dark:hover:bg-slate-800/80"
                    }`}
                  >
                    {t("contact.form.view.multi")}
                  </button>
                </div>
                <div className="inline-flex overflow-hidden rounded-full border border-slate-200 bg-white shadow-sm dark:border-slate-700 dark:bg-slate-900/60">
                  <button
                    type="button"
                    onClick={goPrev}
                    disabled={!canGoPrev}
                    className="px-3 py-1 text-xs font-semibold text-slate-600 transition hover:bg-slate-100 disabled:cursor-not-allowed disabled:opacity-40 dark:text-slate-300 dark:hover:bg-slate-800/80"
                  >
                    {t("contact.form.view.previous")}
                  </button>
                  <button
                    type="button"
                    onClick={goNext}
                    disabled={!canGoNext}
                    className="px-3 py-1 text-xs font-semibold text-slate-600 transition hover:bg-slate-100 disabled:cursor-not-allowed disabled:opacity-40 dark:text-slate-300 dark:hover:bg-slate-800/80"
                  >
                    {t("contact.form.view.next")}
                  </button>
                </div>
              </div>
            </div>

            <div className="overflow-x-auto">
              {displayedDays.length > 0 && timeKeys.length > 0 ? (
                <div
                  className="grid gap-px rounded-2xl border border-slate-200 bg-slate-200 dark:border-slate-700 dark:bg-slate-800"
                  style={{
                    gridTemplateColumns: `120px repeat(${displayedDays.length}, minmax(140px, 1fr))`,
                  }}
                >
                  <div className="sticky left-0 top-0 flex h-14 items-center justify-center bg-slate-100 text-xs font-semibold uppercase text-slate-500 dark:bg-slate-900/80 dark:text-slate-300">
                    {t("contact.form.slotTime")}
                  </div>
                  {displayedDays.map((date) => (
                    <div
                      key={`header-${date}`}
                      className="flex h-14 flex-col items-center justify-center bg-white px-3 text-center text-xs font-semibold text-slate-600 dark:bg-slate-900/80 dark:text-slate-200"
                    >
                      <span>
                        {new Date(`${date}T00:00:00`).toLocaleDateString(undefined, {
                          month: "short",
                          day: "numeric",
                        })}
                      </span>
                      <span className="text-[11px] font-normal text-slate-400 dark:text-slate-500">
                        {new Date(`${date}T00:00:00`).toLocaleDateString(undefined, {
                          weekday: "short",
                        })}
                      </span>
                    </div>
                  ))}

                  {timeKeys.map((timeKey) => {
                    const timeLabel = formatTime(
                      new Date(timeKey).toISOString(),
                      timezoneLabel,
                    );
                    return (
                      <div key={`row-${timeKey}`} className="contents">
                        <div className="flex h-14 items-center justify-center bg-white px-3 text-xs font-semibold text-slate-500 dark:bg-slate-900/70 dark:text-slate-300">
                          {timeLabel}
                        </div>
                        {displayedDays.map((date) => {
                          const slot = slotByDay.get(date)?.get(timeKey);
                          if (!slot) {
                            return (
                              <div
                                key={`${date}-${timeKey}`}
                                className="h-14 bg-white dark:bg-slate-900/40"
                              />
                            );
                          }
                          const isSelected = slot.id === formState.slotId;
                          const isDisabled =
                            slot.status !== "available" || !slot.isBookable;
                          let statusClass = "";
                          if (slot.status === "available") {
                            statusClass = isSelected
                              ? "bg-sky-600 text-white shadow-sm dark:bg-sky-500"
                              : "bg-sky-100 text-sky-700 hover:bg-sky-200 dark:bg-sky-900/40 dark:text-sky-200 dark:hover:bg-sky-900/60";
                          } else if (slot.status === "reserved") {
                            statusClass =
                              "bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-200";
                          } else {
                            statusClass =
                              "bg-slate-100 text-slate-500 dark:bg-slate-800/50 dark:text-slate-400";
                          }

                          return (
                            <button
                              key={`${date}-${timeKey}`}
                              type="button"
                              onClick={() => handleSlotPick(slot.id)}
                              disabled={isDisabled}
                              aria-label={`${formatTime(slot.start, timezoneLabel)} ${t(
                                `contact.form.status.${slot.status}`,
                              )}`}
                              data-testid={
                                slot.status === "available"
                                  ? "availability-slot-available"
                                  : undefined
                              }
                              className={`flex h-14 w-full flex-col items-center justify-center px-2 text-[11px] font-semibold transition ${
                                isDisabled && slot.status !== "available"
                                  ? "cursor-not-allowed opacity-80"
                                  : "cursor-pointer"
                              } ${statusClass}`}
                            >
                              <span>{formatTime(slot.start, timezoneLabel)}</span>
                              <span className="text-[10px] font-normal">
                                {t(`contact.form.status.${slot.status}`)}
                              </span>
                            </button>
                          );
                        })}
                      </div>
                    );
                  })}
                </div>
              ) : (
                <p className="rounded-xl border border-dashed border-slate-300 bg-slate-50 px-4 py-6 text-center text-sm text-slate-500 dark:border-slate-700 dark:bg-slate-900/40 dark:text-slate-300">
                  {t("contact.form.noAvailability")}
                </p>
              )}
            </div>

            <div className="flex flex-wrap gap-4 text-xs text-slate-500 dark:text-slate-400">
              <span className="inline-flex items-center gap-2">
                <span className="h-2.5 w-2.5 rounded-full bg-sky-500 dark:bg-sky-400" />
                {t("contact.form.legendLabels.available")}
              </span>
              <span className="inline-flex items-center gap-2">
                <span className="h-2.5 w-2.5 rounded-full bg-amber-500 dark:bg-amber-400" />
                {t("contact.form.legendLabels.reserved")}
              </span>
              <span className="inline-flex items-center gap-2">
                <span className="h-2.5 w-2.5 rounded-full bg-slate-400 dark:bg-slate-500" />
                {t("contact.form.legendLabels.blackout")}
              </span>
            </div>

            {formErrors.slotId ? (
              <span className="text-xs text-rose-500 dark:text-rose-400">
                {formErrors.slotId}
              </span>
            ) : null}
          </div>

          <button
            type="submit"
            disabled={isSubmitting || isAvailabilityLoading || isConfigLoading}
            className="inline-flex w-full items-center justify-center rounded-full bg-sky-600 px-4 py-2 text-sm font-semibold text-white transition hover:bg-sky-500 disabled:cursor-not-allowed disabled:bg-slate-300 disabled:text-slate-500 dark:bg-sky-500 dark:hover:bg-sky-400 dark:disabled:bg-slate-700"
          >
            {isSubmitting ? t("contact.form.submitting") : t("contact.form.submit")}
          </button>

          {statusMessage ? (
            <p className="rounded-xl border border-emerald-200 bg-emerald-50 p-4 text-sm text-emerald-700 dark:border-emerald-900 dark:bg-emerald-950/40 dark:text-emerald-300">
              {statusMessage}
            </p>
          ) : null}
        </form>

        <aside className="space-y-6">
          <section className="rounded-2xl border border-slate-200 bg-white/80 p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900/60">
            <h2 className="text-sm font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-400">
              {t("contact.summary.title")}
            </h2>
            <ul className="mt-3 space-y-2 text-sm text-slate-600 dark:text-slate-300">
              <li>
                {t("contact.summary.timezone", {
                  timezone: timezoneLabel || t("contact.summary.timezoneFallback"),
                })}
              </li>
              {config?.bookingWindowDays ? (
                <li>
                  {t("contact.summary.window", {
                    days: config.bookingWindowDays,
                  })}
                </li>
              ) : null}
              {config?.supportEmail ? (
                <li>
                  {t("contact.summary.supportEmail")}{" "}
                  <a
                    href={`mailto:${config.supportEmail}`}
                    className="font-medium text-sky-600 underline decoration-sky-200 underline-offset-4 dark:text-sky-400"
                  >
                    {config.supportEmail}
                  </a>
                </li>
              ) : null}
              {config?.googleCalendarId ? (
                <li>{t("contact.summary.calendarLinked")}</li>
              ) : null}
            </ul>
          </section>

          <section className="rounded-2xl border border-slate-200 bg-white/80 p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900/60">
            <h2 className="text-sm font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-400">
              {t("contact.availability.title")}
            </h2>
            {isAvailabilityLoading ? (
              <div className="mt-3 space-y-2">
                <span className="block h-4 w-40 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
                <span className="block h-4 w-32 animate-pulse rounded bg-slate-200 dark:bg-slate-700" />
              </div>
            ) : null}
            {!isAvailabilityLoading && availability?.days?.length ? (
              <ul className="mt-3 space-y-4 text-sm text-slate-600 dark:text-slate-300">
                {availability.days.slice(0, 3).map((day) => (
                  <li key={day.date}>
                    <p className="font-semibold text-slate-900 dark:text-slate-100">
                      {new Date(day.date).toLocaleDateString(i18n.language, {
                        month: "short",
                        day: "numeric",
                        year: "numeric",
                      })}
                    </p>
                    <div className="mt-2 flex flex-wrap gap-2">
                      {day.slots
                        .filter((slot) => slot.isBookable)
                        .map((slot) => (
                          <span
                            key={slot.id}
                            className="inline-flex items-center rounded-full border border-slate-300 px-3 py-1 text-xs font-medium text-slate-700 dark:border-slate-700 dark:text-slate-200"
                          >
                            {formatTime(slot.start, timezoneLabel)} –{" "}
                            {formatTime(slot.end, timezoneLabel)}
                          </span>
                        ))}
                    </div>
                  </li>
                ))}
              </ul>
            ) : null}
            {availabilityError ? (
              <p className="mt-3 text-xs text-rose-500 dark:text-rose-400">
                {t("contact.availability.error")}
              </p>
            ) : null}
          </section>

          {bookingResult ? (
            <section className="rounded-2xl border border-emerald-200 bg-emerald-50 p-6 text-sm text-emerald-700 shadow-sm dark:border-emerald-900 dark:bg-emerald-950/40 dark:text-emerald-300">
              <h2 className="text-sm font-semibold uppercase tracking-wide">
                {t("contact.bookingSummary.title")}
              </h2>
              <ul className="mt-3 space-y-2">
                <li>
                  {t("contact.bookingSummary.when", {
                    datetime: formatDateTime(
                      bookingResult.reservation.startAt,
                      timezoneLabel,
                    ),
                  })}
                </li>
                {bookingResult.reservation.lookupHash ? (
                  <li>
                    {t("contact.bookingSummary.lookup", {
                      hash: bookingResult.reservation.lookupHash,
                    })}
                  </li>
                ) : null}
                {bookingResult.calendarEventId ? (
                  <li>
                    {t("contact.bookingSummary.calendarEvent", {
                      id: bookingResult.calendarEventId,
                    })}
                  </li>
                ) : null}
                {bookingResult.supportEmail ? (
                  <li>
                    {t("contact.bookingSummary.supportEmail", {
                      email: bookingResult.supportEmail,
                    })}
                  </li>
                ) : null}
              </ul>
            </section>
          ) : null}

          {configError ? (
            <p className="rounded-xl border border-rose-200 bg-rose-50 p-4 text-sm text-rose-700 dark:border-rose-900 dark:bg-rose-950/40 dark:text-rose-300">
              {t("contact.error")}
            </p>
          ) : null}
        </aside>
      </div>
    </section>
  );
}
