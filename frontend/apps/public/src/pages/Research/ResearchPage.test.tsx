import { screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { http, HttpResponse } from "msw";
import { describe, expect, it } from "vitest";

import { renderWithRouter } from "../../test-utils/renderWithRouter";
import {
  cloneFixture,
  researchEntriesFixture,
  server,
} from "../../test-utils/server";

describe("ResearchPage", () => {
  it("filters research entries by tags", async () => {
    const user = userEvent.setup();
    const additionalEntry = {
      ...cloneFixture(researchEntriesFixture[0]),
      id: "research-ml",
      title: "Self-Supervised Models for Robotics",
      tags: ["machine learning"],
      contentMarkdown:
        "### Summary\n\nInvestigated multitask learning objectives.",
      contentHtml:
        "<h3>Summary</h3><p>Investigated multitask learning objectives.</p>",
    };

    server.use(
      http.get("/api/v1/public/research", () =>
        HttpResponse.json({
          data: [researchEntriesFixture[0], additionalEntry],
        }),
      ),
    );

    await renderWithRouter({ initialEntries: ["/research"] });

    const allButton = await screen.findByRole("button", {
      name: /all topics/i,
    });
    const mlButton = await screen.findByRole("button", {
      name: /machine learning/i,
    });

    expect(
      await screen.findByRole("heading", {
        name: researchEntriesFixture[0]?.title ?? "",
      }),
    ).toBeInTheDocument();
    expect(
      screen.getByRole("heading", { name: /Self-Supervised Models/i }),
    ).toBeInTheDocument();

    await user.click(mlButton);

    expect(
      screen.queryByRole("heading", {
        name: researchEntriesFixture[0]?.title ?? "",
      }),
    ).not.toBeInTheDocument();
    expect(
      screen.getByRole("heading", { name: /Self-Supervised Models/i }),
    ).toBeInTheDocument();

    await user.click(allButton);

    expect(
      await screen.findByRole("heading", {
        name: researchEntriesFixture[0]?.title ?? "",
      }),
    ).toBeInTheDocument();
  });

  it("renders API error feedback", async () => {
    server.use(
      http.get("/api/v1/public/research", () =>
        HttpResponse.json({ message: "error" }, { status: 500 }),
      ),
    );

    await renderWithRouter({ initialEntries: ["/research"] });

    const alert = await screen.findByRole("alert");
    expect(alert).toHaveTextContent(
      "Research entries could not be loaded. Please retry shortly.",
    );
  });
});
