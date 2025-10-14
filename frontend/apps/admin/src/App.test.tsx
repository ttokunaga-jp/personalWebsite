import { apiClient } from "@shared/lib/api-client";
import { render, screen, waitFor } from "@testing-library/react";
import { vi } from "vitest";

import App from "./App";

describe("Admin App", () => {
  beforeEach(() => {
    vi.spyOn(apiClient, "get").mockResolvedValue({
      data: { status: "operational" }
    } as { data: { status: string } });
  });

  it("renders dashboard title", () => {
    render(<App />);
    return waitFor(() => {
      expect(apiClient.get).toHaveBeenCalledWith("/health");
      expect(screen.getByRole("heading", { name: /Admin console/i })).toBeInTheDocument();
    });
  });
});
