import { apiClient, type ApiClientPromise } from "@shared/lib/api-client";
import { screen, waitFor, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { act } from "react";
import { vi } from "vitest";

import { canonicalProfileEn } from "./modules/profile-content";
import { renderWithRouter } from "./test-utils/renderWithRouter";

let user = userEvent.setup();

type ApiResponse<T> = Awaited<ApiClientPromise<T>>;

describe("App", () => {
  beforeEach(() => {

    vi.spyOn(apiClient, "get").mockImplementation((url) => {
      if (typeof url === "string" && url === "/health") {
        return Promise.resolve({
          data: { status: "healthy" },
        } as ApiResponse<{ status: string }>);
      }

      if (typeof url === "string" && url.includes("/v1/public/profile")) {
        return Promise.resolve({
          data: { data: null },
        } as ApiResponse<{ data: null }>);
      }

      return Promise.resolve({
        data: {},
      } as ApiResponse<Record<string, unknown>>);
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

    expect(
      await screen.findByText(canonicalProfileEn.headline ?? ""),
    ).toBeInTheDocument();
    expect(await screen.findByText("healthy")).toBeInTheDocument();
  });

  it("toggles between light and dark themes", async () => {
    await act(async () => {
      await renderWithRouter();
    });

    const [toggle] = await screen.findAllByRole("button", {
      name: /switch to dark theme/i,
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

    const [japaneseButton] = await screen.findAllByRole("button", {
      name: "日本語",
    });

    await act(async () => {
      await user.click(japaneseButton);
    });

    const homeLinks = await screen.findAllByRole("link", { name: "ホーム" });
    expect(homeLinks.length).toBeGreaterThan(0);
  });

  it("navigates to the profile page via navigation links", async () => {
    await act(async () => {
      await renderWithRouter();
    });

    const mobileMenuToggle = await screen.findByRole("button", {
      name: /toggle navigation menu/i,
    });
    await act(async () => {
      await user.click(mobileMenuToggle);
    });

    const mobileNav = await screen.findByRole("navigation", {
      name: /mobile navigation/i,
    });
    const profileLink = await within(mobileNav).findByRole("link", {
      name: /profile/i,
    });
    await act(async () => {
      await user.click(profileLink);
    });

    expect(
      await screen.findByRole("heading", {
        level: 1,
        name: canonicalProfileEn.displayName,
      }),
    ).toBeInTheDocument();
  });
});
