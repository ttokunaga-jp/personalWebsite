const monthYearFormatter = new Intl.DateTimeFormat(undefined, {
  year: "numeric",
  month: "short"
});

const fullDateFormatter = new Intl.DateTimeFormat(undefined, {
  year: "numeric",
  month: "short",
  day: "numeric"
});

export function formatDateRange(
  startIso: string,
  endIso?: string | null,
  presentLabel: string = "Present"
): string {
  if (!startIso) {
    return "";
  }

  let startLabel = "";
  let endLabel = "";

  try {
    startLabel = monthYearFormatter.format(new Date(startIso));
  } catch {
    startLabel = startIso;
  }

  if (!endIso) {
    endLabel = presentLabel;
  } else {
    try {
      endLabel = monthYearFormatter.format(new Date(endIso));
    } catch {
      endLabel = endIso;
    }
  }

  return `${startLabel} – ${endLabel}`;
}

export function formatDate(isoDate: string): string {
  try {
    return fullDateFormatter.format(new Date(isoDate));
  } catch {
    return isoDate;
  }
}

export function formatDateTime(isoDate: string, timeZone?: string): string {
  try {
    return new Intl.DateTimeFormat(undefined, {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
      timeZone
    }).format(new Date(isoDate));
  } catch {
    return isoDate;
  }
}

export function formatTime(isoDate: string, timeZone?: string): string {
  try {
    return new Intl.DateTimeFormat(undefined, {
      hour: "2-digit",
      minute: "2-digit",
      timeZone
    }).format(new Date(isoDate));
  } catch {
    return isoDate;
  }
}
