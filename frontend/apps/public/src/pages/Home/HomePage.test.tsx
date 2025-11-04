import { screen, waitFor } from "@testing-library/react";
import { http, HttpResponse } from "msw";
import { describe, expect, it } from "vitest";

import { renderWithRouter } from "../../test-utils/renderWithRouter";
import { profileFixture, server } from "../../test-utils/server";

describe("HomePage", () => {
  it("renders profile information, featured links, and health status from the API", async () => {
    await renderWithRouter({ initialEntries: ["/"] });

    expect(
      await screen.findByRole("heading", {
        level: 1,
        name: profileFixture.displayName,
      }),
    ).toBeInTheDocument();

    if (profileFixture.headline) {
      expect(
        await screen.findByText(profileFixture.headline, { exact: false }),
      ).toBeInTheDocument();
    }

    expect(await screen.findByText("ok")).toBeInTheDocument();

    const primarySocial = profileFixture.socialLinks[0];
    if (primarySocial) {
      const socialLinks = await screen.findAllByRole("link", {
        name: primarySocial.label,
      });
      expect(
        socialLinks.some((link) => link.getAttribute("href") === primarySocial.url),
      ).toBe(true);
    }
  });

  it("surfaces API errors for profile and health endpoints", async () => {
    server.use(
      http.get("/api/v1/public/profile", () =>
        HttpResponse.json({ message: "failed" }, { status: 500 }),
      ),
      http.get("/api/health", () =>
        HttpResponse.json({ status: "down" }, { status: 503 }),
      ),
    );

    await renderWithRouter({ initialEntries: ["/"] });

    await waitFor(async () => {
      expect(await screen.findByRole("alert")).toHaveTextContent(
        "We were unable to load the latest profile details.",
      );
    });

    await waitFor(() => {
      expect(screen.getByText("unreachable")).toBeInTheDocument();
    });
  });
});
