import { render, screen } from "@testing-library/react";
import App from "./App";
import "./modules/i18n";

describe("App", () => {
  it("renders welcome text", () => {
    render(<App />);
    expect(screen.getByText(/personal website/i)).toBeInTheDocument();
  });
});
