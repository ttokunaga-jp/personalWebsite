import { screen, within } from "@testing-library/react";
import { http, HttpResponse } from "msw";
import { describe, expect, it } from "vitest";

import { renderWithRouter } from "../../test-utils/renderWithRouter";
import { profileFixture, server } from "../../test-utils/server";

describe("ProfilePage", () => {
  it("orders affiliations and renders skill groups", async () => {
    await renderWithRouter({ initialEntries: ["/profile"] });

    const affiliationHeading = await screen.findByRole("heading", {
      name: /Affiliations/i,
    });
    const affiliationCard = affiliationHeading.closest("article");
    if (!affiliationCard) {
      throw new Error("Affiliation card not found");
    }

    const items = await within(affiliationCard).findAllByRole("listitem");
    expect(items[0]).toHaveTextContent(
      profileFixture.affiliations[0]?.name ?? "",
    );
    expect(items[1]).toHaveTextContent(
      profileFixture.affiliations[1]?.name ?? "",
    );

    const skillsHeading = await screen.findByRole("heading", {
      name: /Skills and capabilities/i,
    });
    const skillsCard = skillsHeading.closest("article");
    if (!skillsCard) {
      throw new Error("Skills card not found");
    }

    expect(
      await within(skillsCard).findByText(
        profileFixture.techSections[0]?.members[0]?.tech.displayName ?? "",
        { exact: false },
      ),
    ).toBeInTheDocument();
  });

  it("displays an error banner when profile API fails", async () => {
    server.use(
      http.get("/api/v1/public/profile", () =>
        HttpResponse.json({ message: "denied" }, { status: 500 }),
      ),
    );

    await renderWithRouter({ initialEntries: ["/profile"] });

    const alert = await screen.findByRole("alert");
    expect(alert).toHaveTextContent(
      "We could not refresh the profile data. Please reload the page.",
    );
  });
});
