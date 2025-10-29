import { screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { http, HttpResponse } from "msw";
import { describe, expect, it } from "vitest";

import { renderWithRouter } from "../../test-utils/renderWithRouter";
import { projectsFixture, server } from "../../test-utils/server";

describe("ProjectsPage", () => {
  it("filters projects by selected tech stacks", async () => {
    const user = userEvent.setup();
    await renderWithRouter({ initialEntries: ["/projects"] });

    const allButton = await screen.findByRole("button", {
      name: /all stacks/i,
    });
    const typeScriptButton = await screen.findByRole("button", {
      name: "TypeScript",
    });
    const reactButton = await screen.findByRole("button", { name: "React" });

    // Initially both projects should be visible.
    for (const project of projectsFixture) {
      expect(
        await screen.findByRole("heading", {
          name: new RegExp(project.title, "i"),
        }),
      ).toBeInTheDocument();
    }

    await user.click(typeScriptButton);

    expect(
      await screen.findByRole("heading", {
        name: new RegExp(projectsFixture[0]?.title ?? "", "i"),
      }),
    ).toBeInTheDocument();
    expect(
      screen.queryByRole("heading", {
        name: new RegExp(projectsFixture[1]?.title ?? "", "i"),
      }),
    ).not.toBeInTheDocument();

    await user.click(reactButton);

    expect(
      screen.getByRole("heading", {
        name: new RegExp(projectsFixture[2]?.title ?? "", "i"),
      }),
    ).toBeInTheDocument();
    expect(
      screen.getByRole("heading", {
        name: new RegExp(projectsFixture[0]?.title ?? "", "i"),
      }),
    ).toBeInTheDocument();

    await user.click(allButton);
    expect(
      screen.getByRole("heading", {
        name: new RegExp(projectsFixture[1]?.title ?? "", "i"),
      }),
    ).toBeInTheDocument();
  });

  it("shows an error banner when the projects endpoint fails", async () => {
    server.use(
      http.get("/api/v1/public/projects", () =>
        HttpResponse.json({ message: "boom" }, { status: 500 }),
      ),
    );

    await renderWithRouter({ initialEntries: ["/projects"] });

    const alert = await screen.findByRole("alert");
    expect(alert).toHaveTextContent(
      "Projects could not be retrieved from the API.",
    );
  });
});
