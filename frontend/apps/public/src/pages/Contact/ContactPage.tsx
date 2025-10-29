import { FormEvent, useCallback, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";

import {
  useContactAvailability,
  useContactConfig,
  publicApi,
} from "../../modules/public-api";
import { formatDateTime, formatTime } from "../../utils/date";

type FormState = {
  name: string;
  email: string;
  topic: string;
  message: string;
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
  message: "",
  slotId: "",
};

export function ContactPage() {
  const { t } = useTranslation();
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
  const slots = availability?.slots ?? [];
  const topics = config?.topics ?? [];

  const [formState, setFormState] = useState<FormState>(initialFormState);
  const [formErrors, setFormErrors] = useState<FormErrors>({});
  const [statusMessage, setStatusMessage] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const timezoneLabel = useMemo(
    () => availability?.timezone ?? "",
    [availability],
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

      if (!state.message.trim() || state.message.trim().length < 20) {
        errors.message = t("contact.form.errors.messageLength");
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
      throw new Error("Recaptcha site key not configured");
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
      throw new Error("Recaptcha unavailable");
    }

    return window.grecaptcha.execute(config.recaptchaSiteKey, {
      action: "submit",
    });
  }, [config?.recaptchaSiteKey]);

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setStatusMessage(null);

    const errors = validate(formState);
    setFormErrors(errors);

    if (Object.keys(errors).length > 0) {
      return;
    }

    try {
      setIsSubmitting(true);
      const recaptchaToken = await loadRecaptchaToken();
      const response = await publicApi.createBooking({
        name: formState.name.trim(),
        email: formState.email.trim(),
        topic: formState.topic,
        message: formState.message.trim(),
        slotId: formState.slotId,
        recaptchaToken,
      });

      setStatusMessage(
        t("contact.form.success", { bookingId: response.bookingId }),
      );
      setFormState(initialFormState);
      setFormErrors({});
    } catch (error) {
      if (import.meta.env.DEV) {
        console.error(error);
      }
      setStatusMessage(t("contact.form.error"));
    } finally {
      setIsSubmitting(false);
    }
  };

  const selectSlot = (slotId: string, enabled: boolean) => {
    if (!enabled) {
      return;
    }
    setFormState((previous) => ({
      ...previous,
      slotId: previous.slotId === slotId ? "" : slotId,
    }));
  };

  return (
    <section
      id="contact"
      className="mx-auto flex w-full max-w-3xl flex-col gap-6 px-4 py-12 sm:px-8"
    >
      <header className="space-y-3">
        <p className="text-sm font-semibold uppercase tracking-[0.3em] text-sky-500 dark:text-sky-400">
          {t("contact.tagline")}
        </p>
        <h1 className="text-3xl font-bold text-slate-900 dark:text-slate-50 sm:text-4xl">
          {t("contact.title")}
        </h1>
        <p className="text-base text-slate-600 dark:text-slate-300">
          {t("contact.description")}
        </p>
      </header>
      <div className="rounded-xl border border-slate-200 bg-white/80 p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900/60">
        <h2 className="text-lg font-semibold text-slate-900 dark:text-slate-100">
          {t("contact.availability.title")}
        </h2>
        <p className="text-sm text-slate-600 dark:text-slate-300">
          {t("contact.availability.description")}
        </p>
        <div
          className="mt-4 flex flex-wrap gap-2"
          role="group"
          aria-label={t("contact.availability.groupLabel")}
        >
          {isAvailabilityLoading && (
            <>
              <span className="inline-block h-10 w-28 animate-pulse rounded-lg bg-slate-200 dark:bg-slate-700" />
              <span className="inline-block h-10 w-28 animate-pulse rounded-lg bg-slate-200 dark:bg-slate-700" />
            </>
          )}
          {!isAvailabilityLoading && slots.length === 0 ? (
            <span className="text-sm text-slate-500 dark:text-slate-400">
              {t("contact.availability.unavailable")}
            </span>
          ) : null}
          {slots.map((slot) => {
            const isSelected = formState.slotId === slot.id;
            const isEnabled = slot.isBookable;
            return (
              <button
                key={slot.id}
                type="button"
                onClick={() => selectSlot(slot.id, isEnabled)}
                disabled={!isEnabled}
                className={`rounded-lg border px-3 py-2 text-sm transition ${
                  isSelected
                    ? "border-sky-500 bg-sky-500 text-white"
                    : isEnabled
                      ? "border-slate-300 text-slate-700 hover:border-sky-400 hover:text-sky-600 dark:border-slate-700 dark:text-slate-200 dark:hover:border-sky-400 dark:hover:text-sky-300"
                      : "cursor-not-allowed border-slate-200 bg-slate-100 text-slate-400 dark:border-slate-800 dark:bg-slate-900/40 dark:text-slate-500"
                }`}
                aria-pressed={isSelected}
              >
                <span className="block font-medium">
                  {formatDateTime(slot.start, availability?.timezone)}
                </span>
                <span className="block text-xs uppercase tracking-wide text-slate-500 dark:text-slate-400">
                  {t("contact.availability.slotTo", {
                    end: formatTime(slot.end, availability?.timezone),
                  })}
                </span>
              </button>
            );
          })}
        </div>
        {formErrors.slotId ? (
          <p className="mt-2 text-xs text-rose-500 dark:text-rose-400">
            {formErrors.slotId}
          </p>
        ) : null}
        {timezoneLabel ? (
          <p className="mt-3 text-xs text-slate-500 dark:text-slate-400">
            {t("contact.availability.timezone", { timezone: timezoneLabel })}
          </p>
        ) : null}
        {availabilityError ? (
          <p
            role="alert"
            className="mt-3 text-xs text-rose-500 dark:text-rose-400"
          >
            {t("contact.availability.error")}
          </p>
        ) : null}
      </div>

      <form
        onSubmit={handleSubmit}
        noValidate
        className="space-y-6 rounded-xl border border-slate-200 bg-white/80 p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900/60"
      >
        <fieldset className="grid gap-4">
          <legend className="sr-only">{t("contact.form.legend")}</legend>
          <label className="flex flex-col gap-1 text-sm">
            <span className="font-medium text-slate-700 dark:text-slate-200">
              {t("contact.form.name")}
            </span>
            <input
              type="text"
              name="name"
              autoComplete="name"
              value={formState.name}
              onChange={handleInputChange("name")}
              className="rounded-lg border border-slate-300 px-3 py-2 text-sm text-slate-900 shadow-sm transition focus-visible:border-sky-500 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-sky-200 dark:border-slate-700 dark:bg-slate-900/40 dark:text-slate-100 dark:focus-visible:border-sky-400 dark:focus-visible:ring-sky-900/60"
              aria-invalid={Boolean(formErrors.name)}
            />
            {formErrors.name ? (
              <span className="text-xs text-rose-500 dark:text-rose-400">
                {formErrors.name}
              </span>
            ) : null}
          </label>

          <label className="flex flex-col gap-1 text-sm">
            <span className="font-medium text-slate-700 dark:text-slate-200">
              {t("contact.form.email")}
            </span>
            <input
              type="email"
              name="email"
              autoComplete="email"
              value={formState.email}
              onChange={handleInputChange("email")}
              className="rounded-lg border border-slate-300 px-3 py-2 text-sm text-slate-900 shadow-sm transition focus-visible:border-sky-500 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-sky-200 dark:border-slate-700 dark:bg-slate-900/40 dark:text-slate-100 dark:focus-visible:border-sky-400 dark:focus-visible:ring-sky-900/60"
              aria-invalid={Boolean(formErrors.email)}
            />
            {formErrors.email ? (
              <span className="text-xs text-rose-500 dark:text-rose-400">
                {formErrors.email}
              </span>
            ) : null}
          </label>

          <label className="flex flex-col gap-1 text-sm">
            <span className="font-medium text-slate-700 dark:text-slate-200">
              {t("contact.form.topic")}
            </span>
            <select
              name="topic"
              value={formState.topic}
              onChange={handleInputChange("topic")}
              disabled={topics.length === 0}
              className="rounded-lg border border-slate-300 px-3 py-2 text-sm text-slate-900 shadow-sm transition focus-visible:border-sky-500 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-sky-200 disabled:cursor-not-allowed disabled:bg-slate-100 dark:border-slate-700 dark:bg-slate-900/40 dark:text-slate-100 dark:disabled:bg-slate-900/20 dark:focus-visible:border-sky-400 dark:focus-visible:ring-sky-900/60"
              aria-invalid={Boolean(formErrors.topic)}
            >
              <option value="">{t("contact.form.topicPlaceholder")}</option>
              {topics.map((topic) => (
                <option key={topic} value={topic}>
                  {topic}
                </option>
              ))}
            </select>
            {formErrors.topic ? (
              <span className="text-xs text-rose-500 dark:text-rose-400">
                {formErrors.topic}
              </span>
            ) : null}
          </label>

          <label className="flex flex-col gap-1 text-sm">
            <span className="font-medium text-slate-700 dark:text-slate-200">
              {t("contact.form.message")}
            </span>
            <textarea
              name="message"
              rows={4}
              value={formState.message}
              onChange={handleInputChange("message")}
              className="rounded-lg border border-slate-300 px-3 py-2 text-sm text-slate-900 shadow-sm transition focus-visible:border-sky-500 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-sky-200 dark:border-slate-700 dark:bg-slate-900/40 dark:text-slate-100 dark:focus-visible:border-sky-400 dark:focus-visible:ring-sky-900/60"
              aria-invalid={Boolean(formErrors.message)}
            />
            {formErrors.message ? (
              <span className="text-xs text-rose-500 dark:text-rose-400">
                {formErrors.message}
              </span>
            ) : null}
            {config?.consentText ? (
              <p className="text-xs text-slate-500 dark:text-slate-400">
                {config.consentText}
              </p>
            ) : null}
          </label>
        </fieldset>

        <button
          type="submit"
          disabled={isSubmitting || isConfigLoading}
          className="inline-flex w-full items-center justify-center rounded-full bg-sky-600 px-6 py-3 text-sm font-semibold text-white transition hover:bg-sky-500 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 focus-visible:ring-sky-400 disabled:cursor-not-allowed disabled:bg-slate-300 dark:bg-sky-500 dark:text-slate-900 dark:hover:bg-sky-400 dark:focus-visible:ring-sky-300"
        >
          {isSubmitting
            ? t("contact.form.submitting")
            : t("contact.form.submit")}
        </button>
        {statusMessage ? (
          <p className="text-sm text-slate-600 dark:text-slate-300">
            {statusMessage}
          </p>
        ) : null}
        {configError ? (
          <p role="alert" className="text-xs text-rose-500 dark:text-rose-400">
            {t("contact.form.configError")}
          </p>
        ) : null}
      </form>
    </section>
  );
}
