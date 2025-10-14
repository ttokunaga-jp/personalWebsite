import { apiClient } from "@shared/lib/api-client";
import { render, screen, waitFor } from "@testing-library/react";
import { vi } from "vitest";

import App from "./App";

describe("App", () => {
  beforeEach(() => {
    vi.spyOn(apiClient, "get").mockResolvedValue({
      data: { status: "healthy" }
    } as { data: { status: string } });
  });

  it("renders welcome text", () => {
    render(<App />);
    return waitFor(() => {
      expect(apiClient.get).toHaveBeenCalledWith("/health");
      expect(screen.getByText(/personal website/i)).toBeInTheDocument();
    });
  });
});
