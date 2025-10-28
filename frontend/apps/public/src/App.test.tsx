import { apiClient } from "@shared/lib/api-client";
import { screen, waitFor, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import type { AxiosResponse } from "axios";
import { act } from "react";
import { vi } from "vitest";

import { renderWithRouter } from "./test-utils/renderWithRouter";

let user = userEvent.setup();

describe("App", () => {
  beforeEach(() => {
    const profileResponse = {
      name: "Takumi Asano",
      headline: "Research Engineer",
      summary: "Building reliable software around human-centered research.",
      affiliations: [
        {
          id: "aff-1",
          organization: "Kyoto University",
          department: "Graduate School of Informatics",
          role: "Research Fellow",
          startDate: "2023-04-01",
          endDate: null,
          isCurrent: true
        }
      ],
      workHistory: [],
      skillGroups: [],
      communities: [],
      socialLinks: []
    };

    vi.spyOn(apiClient, "get").mockImplementation((url) => {
      if (typeof url === "string" && url === "/health") {
        return Promise.resolve({
          data: { status: "healthy" }
        } as AxiosResponse<{ status: string }>);
      }

      if (typeof url === "string" && url.includes("/v1/public/profile")) {
        return Promise.resolve({
          data: { data: profileResponse }
        } as AxiosResponse<{ data: typeof profileResponse }>);
      }

      return Promise.resolve({
        data: {}
      } as AxiosResponse<Record<string, unknown>>);
    });

    user = userEvent.setup();
  });

  it("renders home page scaffold and fetches API health", async () => {
    await act(async () => {
      await renderWithRouter();
    });

    await waitFor(() => {
      expect(apiClient.get).toHaveBeenCalledWith("/health");
    });

    expect(await screen.findByText(/Research Engineer/i)).toBeInTheDocument();
    expect(await screen.findByText("healthy")).toBeInTheDocument();
  });

  it("toggles between light and dark themes", async () => {
    await act(async () => {
      await renderWithRouter();
    });

    const [toggle] = await screen.findAllByRole("button", {
      name: /switch to dark theme/i
    });

    await act(async () => {
      await user.click(toggle);
    });

    await waitFor(() => {
      expect(document.documentElement.classList.contains("dark")).toBe(true);
      expect(toggle).toHaveAttribute("title", "Switch to light theme");
    });
  });

  it("switches language to Japanese", async () => {
    await act(async () => {
      await renderWithRouter();
    });

    const [japaneseButton] = await screen.findAllByRole("button", { name: "日本語" });

    await act(async () => {
      await user.click(japaneseButton);
    });

    expect(await screen.findByText("人を軸にしたイノベーション")).toBeInTheDocument();
  });

  it("navigates to the profile page via navigation links", async () => {
    await act(async () => {
      await renderWithRouter();
    });

    const mobileMenuToggle = await screen.findByRole("button", {
      name: /toggle navigation menu/i
    });
    await act(async () => {
      await user.click(mobileMenuToggle);
    });

    const mobileNav = await screen.findByRole("navigation", { name: /mobile navigation/i });
    const profileLink = await within(mobileNav).findByRole("link", { name: /profile/i });
    await act(async () => {
      await user.click(profileLink);
    });

    expect(
      await screen.findByText(/Professional profile/)
    ).toBeInTheDocument();
  });
});
