import { apiClient } from "@shared/lib/api-client";
import { screen, waitFor, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { act } from "react";
import { vi } from "vitest";

import { renderWithRouter } from "./test-utils/renderWithRouter";

let user = userEvent.setup();

describe("App", () => {
  beforeEach(() => {
    vi.spyOn(apiClient, "get").mockResolvedValue({
      data: { status: "healthy" }
    } as { data: { status: string } });

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
      await screen.findByText(/Crafting research-driven products/i)
    ).toBeInTheDocument();
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

    expect(
      await screen.findByText("研究を軸にしたプロダクトと体験を創出します。")
    ).toBeInTheDocument();
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
