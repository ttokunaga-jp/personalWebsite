import { ChangeEvent } from "react";

import type { LocalizedField } from "../types";

type LocalizedTextFieldProps = {
  id: string;
  label: string;
  value: LocalizedField;
  onChange: (value: LocalizedField) => void;
  required?: boolean;
  multiline?: boolean;
  helperText?: string;
};

export function LocalizedTextField({
  id,
  label,
  value,
  onChange,
  required,
  multiline,
  helperText,
}: LocalizedTextFieldProps) {
  const handleChange =
    (key: keyof LocalizedField) => (event: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
      onChange({
        ...value,
        [key]: event.target.value,
      });
    };

  const InputComponent = multiline ? "textarea" : "input";
  const baseProps = multiline
    ? {
        rows: 3,
      }
    : {};

  return (
    <fieldset className="flex flex-col gap-2 rounded-lg border border-slate-200 bg-white/60 p-4 shadow-sm transition dark:border-slate-700 dark:bg-slate-900/60">
      <legend className="text-sm font-semibold text-slate-700 dark:text-slate-200">
        {label}
        {required && <span className="ml-1 text-rose-500">*</span>}
      </legend>
      <div className="grid gap-3 md:grid-cols-2">
        <label className="flex flex-col gap-1 text-sm text-slate-600 dark:text-slate-300">
          <span className="font-medium">日本語</span>
          <InputComponent
            {...baseProps}
            id={`${id}-ja`}
            value={value.ja}
            onChange={handleChange("ja")}
            required={required}
            className="w-full rounded-md border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 shadow-sm transition focus:border-sky-500 focus:outline-none focus:ring-2 focus:ring-sky-500 dark:border-slate-700 dark:bg-slate-900 dark:text-slate-100"
          />
        </label>
        <label className="flex flex-col gap-1 text-sm text-slate-600 dark:text-slate-300">
          <span className="font-medium">English</span>
          <InputComponent
            {...baseProps}
            id={`${id}-en`}
            value={value.en}
            onChange={handleChange("en")}
            required={required}
            className="w-full rounded-md border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 shadow-sm transition focus:border-sky-500 focus:outline-none focus:ring-2 focus:ring-sky-500 dark:border-slate-700 dark:bg-slate-900 dark:text-slate-100"
          />
        </label>
      </div>
      {helperText && (
        <p className="text-xs text-slate-500 dark:text-slate-400">{helperText}</p>
      )}
    </fieldset>
  );
}
