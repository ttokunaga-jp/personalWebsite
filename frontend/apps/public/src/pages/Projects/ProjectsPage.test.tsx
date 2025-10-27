import { apiClient } from "@shared/lib/api-client";
import { screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import type { AxiosResponse } from "axios";
import { vi } from "vitest";

import type { Project } from "../../modules/public-api";
import { renderWithRouter } from "../../test-utils/renderWithRouter";

const createAxiosResponse = <T,>(data: T): AxiosResponse<T> =>
  ({
    data,
    status: 200,
    statusText: "OK",
    headers: {},
    config: {}
  }) as AxiosResponse<T>;

describe("ProjectsPage", () => {
  it("filters projects by tech stack selection", async () => {
    const user = userEvent.setup();
    const projects: Project[] = [
      {
        id: "proj-1",
        title: "AI Research Dashboard",
        subtitle: "Observability for experiments",
        description: "A dashboard aligning experiment metadata with AI lab metrics.",
        techStack: ["React", "TypeScript", "Go"],
        category: "Research",
        period: {
          start: "2023-01-01",
          end: null
        },
        links: [
          { label: "Repository", url: "https://example.com/repo", type: "repo" },
          { label: "Demo", url: "https://example.com/demo", type: "demo" }
        ]
      },
      {
        id: "proj-2",
        title: "Cloud IaC Platform",
        description: "Composable Terraform modules for data-intensive workloads.",
        techStack: ["Terraform", "Go"],
        category: "Platform",
        period: {
          start: "2022-05-01",
          end: "2023-06-01"
        },
        links: [{ label: "Docs", url: "https://example.com/docs", type: "article" }]
      }
    ];

    const getSpy = vi.spyOn(apiClient, "get").mockImplementation((url) => {
      if (typeof url === "string" && url.includes("/v1/public/projects")) {
        return Promise.resolve(createAxiosResponse(projects));
      }

      if (typeof url === "string" && url === "/health") {
        return Promise.resolve(createAxiosResponse({ status: "ok" }));
      }

      return Promise.resolve(createAxiosResponse({}));
    });

    await renderWithRouter({ initialEntries: ["/projects"] });

    await waitFor(() => {
      expect(getSpy).toHaveBeenCalledWith(
        "/v1/public/projects",
        expect.objectContaining({ signal: expect.any(AbortSignal) })
      );
    });

    expect(
      await screen.findByRole("heading", { name: /AI Research Dashboard/i })
    ).toBeInTheDocument();
    expect(
      await screen.findByRole("heading", { name: /Cloud IaC Platform/i })
    ).toBeInTheDocument();

    const terraformFilter = await screen.findByRole("button", { name: "Terraform" });
    await user.click(terraformFilter);

    await waitFor(() => {
      expect(
        screen.queryByRole("heading", { name: /AI Research Dashboard/i })
      ).not.toBeInTheDocument();
    });

    expect(
      screen.getByRole("heading", { name: /Cloud IaC Platform/i })
    ).toBeInTheDocument();

    const allStacksButton = await screen.findByRole("button", { name: /All stacks/i });
    await user.click(allStacksButton);

    expect(
      await screen.findByRole("heading", { name: /AI Research Dashboard/i })
    ).toBeInTheDocument();
  });
});
