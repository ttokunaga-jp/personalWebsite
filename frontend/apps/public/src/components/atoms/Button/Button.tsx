import type { ButtonHTMLAttributes } from "react";

import { classNames } from "../../../lib/classNames";

type Variant = "primary" | "ghost";
type Size = "md" | "sm";

export type ButtonProps = {
  variant?: Variant;
  size?: Size;
  block?: boolean;
} & ButtonHTMLAttributes<HTMLButtonElement>;

const baseStyles =
  "inline-flex items-center justify-center gap-2 rounded-full font-medium transition focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 focus-visible:ring-sky-500 dark:focus-visible:ring-sky-400";

const variantStyles: Record<Variant, string> = {
  primary:
    "bg-sky-600 text-white hover:bg-sky-500 active:bg-sky-700 dark:bg-sky-500 dark:text-slate-900 dark:hover:bg-sky-400",
  ghost:
    "border border-slate-300 text-slate-700 hover:border-sky-400 hover:text-sky-600 dark:border-slate-700 dark:text-slate-200 dark:hover:border-sky-400 dark:hover:text-sky-300"
};

const sizeStyles: Record<Size, string> = {
  md: "px-5 py-2 text-sm",
  sm: "px-4 py-1.5 text-xs"
};

export function Button({
  variant = "primary",
  size = "md",
  block = false,
  className,
  ...props
}: ButtonProps) {
  return (
    <button
      className={classNames(
        baseStyles,
        variantStyles[variant],
        sizeStyles[size],
        block && "w-full",
        className
      )}
      {...props}
    />
  );
}
