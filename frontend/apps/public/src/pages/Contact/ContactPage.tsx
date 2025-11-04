import { FormEvent, useCallback, useId, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";

import {
  publicApi,
  useContactAvailability,
  useContactConfig,
} from "../../modules/public-api";
import { formatDateTime, formatTime } from "../../utils/date";

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

  const availableSlots = useMemo(() => {
    if (!availability?.days) {
      return [];
    }
    return availability.days.flatMap((day) =>
      day.slots
        .filter((slot) => slot.isBookable)
        .map((slot) => ({
          id: slot.id || slot.start,
          day: day.date,
          start: slot.start,
          end: slot.end,
        })),
    );
  }, [availability]);

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
      }

      return errors;
    },
    [t],
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
      const selectedSlot = availableSlots.find(
        (slot) => slot.id === formState.slotId,
      );
      if (!selectedSlot) {
        setFormErrors({ slotId: t("contact.form.errors.slotRequired") });
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
          id: response.meeting.id,
          email: response.supportEmail ?? response.meeting.email,
        }),
      );
      setFormState(initialFormState);
    } catch {
      setStatusMessage(t("contact.form.error"));
    } finally {
      setIsSubmitting(false);
    }
  };

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

          <label className="flex flex-col gap-2">
            <span className="text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-400">
              {t("contact.form.slot")}
            </span>
            <select
              name="slotId"
              value={formState.slotId}
              onChange={handleInputChange("slotId")}
              required
              className="rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 focus:border-sky-400 focus:outline-none focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-900 dark:text-slate-100 dark:focus:border-sky-500 dark:focus:ring-sky-900/30"
            >
              <option value="">{t("contact.form.slotPlaceholder")}</option>
              {availableSlots.map((slot) => (
                <option key={slot.id} value={slot.id}>
                  {formatDateTime(slot.start, timezoneLabel)} (
                  {formatTime(slot.start, timezoneLabel)} -{" "}
                  {formatTime(slot.end, timezoneLabel)})
                </option>
              ))}
            </select>
            {formErrors.slotId ? (
              <span className="text-xs text-rose-500 dark:text-rose-400">
                {formErrors.slotId}
              </span>
            ) : null}
          </label>

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
                      bookingResult.meeting.datetime,
                      timezoneLabel,
                    ),
                  })}
                </li>
                {bookingResult.calendarEventId ? (
                  <li>
                    {t("contact.bookingSummary.calendarEvent", {
                      id: bookingResult.calendarEventId,
                    })}
                  </li>
                ) : null}
                {bookingResult.meeting.meetUrl ? (
                  <li>
                    <a
                      href={bookingResult.meeting.meetUrl}
                      target="_blank"
                      rel="noreferrer"
                      className="font-medium text-sky-600 underline decoration-sky-200 underline-offset-4 dark:text-sky-400"
                    >
                      {t("contact.bookingSummary.joinLink")}
                    </a>
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
