import { render, screen } from "@testing-library/react";
import App from "./App";
import "./modules/i18n";

describe("Admin App", () => {
  it("renders dashboard title", () => {
    render(<App />);
    expect(screen.getByText(/Admin console/i)).toBeInTheDocument();
  });
});
