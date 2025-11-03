import { describe, expect, it } from "vitest";

import { resources } from "./i18n";

describe("i18n resources", () => {
  it("exposes newly added home translations", () => {
    expect(resources.en.translation.home.tech.title).toBeDefined();
    expect(resources.ja.translation.home.tech.title).toBeDefined();
    expect(resources.en.translation.home.work.title).toBeDefined();
    expect(resources.ja.translation.home.work.title).toBeDefined();
  });

  it("includes extended profile sections", () => {
    expect(resources.en.translation.profile.sections.social.title).toBeDefined();
    expect(resources.ja.translation.profile.sections.social.title).toBeDefined();
  });

  it("supports research filters and kinds", () => {
    expect(resources.en.translation.research.filters.research).toBe("Research");
    expect(resources.ja.translation.research.filters.blog).toBe("ブログ");
    expect(resources.en.translation.research.kind.blog).toBe("Blog");
  });

  it("exposes booking summary translations", () => {
    expect(resources.en.translation.contact.summary.title).toBe("Booking summary");
    expect(resources.ja.translation.contact.summary.title).toBe("予約に関する情報");
    expect(resources.en.translation.contact.bookingSummary.title).toBe("Reservation details");
  });
});
