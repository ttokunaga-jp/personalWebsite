import { screen, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";

import { publicApi } from "../../modules/public-api";
import { renderWithRouter } from "../../test-utils/renderWithRouter";
import { contactAvailabilityFixture, contactConfigFixture } from "../../test-utils/server";

declare global {
  interface Window {
    grecaptcha?: {
      ready: (callback: () => void) => void;
      execute: (siteKey: string, options: { action: string }) => Promise<string>;
    };
  }
}

async function getFirstAvailableSlotButton() {
  const slotGroup = await screen.findByRole("group", { name: /available time slots/i });
  const buttons = within(slotGroup).getAllByRole("button");
  const enabledButton = buttons.find((button) => !button.hasAttribute("disabled"));
  if (!enabledButton) {
    throw new Error("No bookable slot button found");
  }
  return enabledButton;
}

describe("ContactPage", () => {
  it("enforces client-side validation", async () => {
    const user = userEvent.setup();
    await renderWithRouter({ initialEntries: ["/contact"] });

    const submitButton = await screen.findByRole("button", { name: /request booking/i });

    await user.click(submitButton);

    expect(await screen.findByText("Please provide your name.")).toBeInTheDocument();
    expect(screen.getByText("An email address is required.")).toBeInTheDocument();
    expect(screen.getByText("Select a topic to help us route your request.")).toBeInTheDocument();
    expect(
      screen.getByText("Share at least 20 characters so we can prepare effectively.")
    ).toBeInTheDocument();
    expect(screen.getByText("Select an available time slot.")).toBeInTheDocument();
  });

  it("submits a booking request after collecting a Recaptcha token", async () => {
    const user = userEvent.setup();
    const createBookingMock = vi
      .spyOn(publicApi, "createBooking")
      .mockResolvedValue({ bookingId: "bk-slot-1", status: "pending" });

    window.grecaptcha = {
      ready: (callback: () => void) => {
        callback();
      },
      execute: vi.fn().mockResolvedValue("recaptcha-token-123")
    };

    await renderWithRouter({ initialEntries: ["/contact"] });

    const nameInput = await screen.findByLabelText("Your name");
    const emailInput = await screen.findByLabelText("Email address");
    const topicSelect = await screen.findByLabelText("Topic");
    const textboxes = await screen.findAllByRole("textbox");
    const messageTextarea = textboxes.find((element) => element.getAttribute("name") === "message");
    if (!messageTextarea) {
      throw new Error("Message textarea not found");
    }

    await user.type(nameInput, "  Jane Doe  ");
    await user.type(emailInput, "jane.doe@example.com");
    await user.selectOptions(topicSelect, contactConfigFixture.topics[0] ?? "");
    await user.type(
      messageTextarea,
      "I would like to discuss possibilities for joint research in HRI."
    );

    const slotButton = await getFirstAvailableSlotButton();
    await user.click(slotButton);

    await user.click(screen.getByRole("button", { name: /request booking/i }));

    expect(await screen.findByText(/Your request \(ID: bk-slot-1\)/)).toBeInTheDocument();

    expect(createBookingMock).toHaveBeenCalledWith({
      name: "Jane Doe",
      email: "jane.doe@example.com",
      topic: contactConfigFixture.topics[0] ?? "",
      message: "I would like to discuss possibilities for joint research in HRI.",
      slotId: contactAvailabilityFixture.slots[0]?.id ?? "",
      recaptchaToken: "recaptcha-token-123"
    });

    expect(window.grecaptcha?.execute).toHaveBeenCalledWith(
      contactConfigFixture.recaptchaSiteKey,
      { action: "submit" }
    );

    delete window.grecaptcha;
  });
});
