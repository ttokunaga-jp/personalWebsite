import type { ButtonHTMLAttributes } from "react";

import { classNames } from "../../../lib/classNames";

export type IconButtonProps = ButtonHTMLAttributes<HTMLButtonElement>;

export function IconButton({ className, ...props }: IconButtonProps) {
  return (
    <button
      type="button"
      className={classNames(
        "inline-flex h-10 w-10 items-center justify-center rounded-full border border-slate-300 text-slate-600 transition hover:border-sky-400 hover:text-sky-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-sky-500 focus-visible:ring-offset-2 dark:border-slate-700 dark:text-slate-200 dark:hover:border-sky-400 dark:hover:text-sky-300 dark:focus-visible:ring-sky-400",
        className,
      )}
      {...props}
    />
  );
}
