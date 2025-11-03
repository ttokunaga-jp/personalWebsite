import { screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";

import { publicApi } from "../../modules/public-api";
import { renderWithRouter } from "../../test-utils/renderWithRouter";
import {
  contactAvailabilityFixture,
  contactConfigFixture,
} from "../../test-utils/server";

declare global {
  interface Window {
    grecaptcha?: {
      ready: (callback: () => void) => void;
      execute: (
        siteKey: string,
        options: { action: string },
      ) => Promise<string>;
    };
  }
}

describe("ContactPage", () => {
  it("enforces client-side validation", async () => {
    const user = userEvent.setup();
    await renderWithRouter({ initialEntries: ["/contact"] });

    const submitButton = await screen.findByRole("button", {
      name: /request booking/i,
    });

    await user.click(submitButton);

    expect(
      await screen.findByText("Please provide your name."),
    ).toBeInTheDocument();
    expect(
      screen.getByText("An email address is required."),
    ).toBeInTheDocument();
    expect(
      screen.getByText("Select a topic to help us route your request."),
    ).toBeInTheDocument();
    expect(
      screen.getByText(
        "Share at least 20 characters so we can prepare effectively.",
      ),
    ).toBeInTheDocument();
    expect(
      screen.getByText("Select an available time slot."),
    ).toBeInTheDocument();
  });

  it("submits a booking request after collecting a Recaptcha token", async () => {
    const user = userEvent.setup();
    const firstSlot = contactAvailabilityFixture.days[0]?.slots[0];
    const createBookingMock = vi
      .spyOn(publicApi, "createBooking")
      .mockResolvedValue({
        meeting: {
          id: "bk-1",
          name: "Jane Doe",
          email: "jane.doe@example.com",
          datetime: firstSlot?.start ?? "",
          durationMinutes: firstSlot
            ? Math.round(
                (new Date(firstSlot.end).getTime() -
                  new Date(firstSlot.start).getTime()) /
                  60000,
              )
            : 30,
          meetUrl: "https://meet.example.com/mock",
          calendarEventId: "event-1",
          status: "pending",
          notes: "[Research collaboration] I would like to discuss possibilities for joint research in HRI.",
        },
        calendarEventId: "event-1",
        supportEmail: contactConfigFixture.supportEmail,
        calendarTimezone: contactConfigFixture.calendarTimezone,
      });

    window.grecaptcha = {
      ready: (callback: () => void) => {
        callback();
      },
      execute: vi.fn().mockResolvedValue("recaptcha-token-123"),
    };

    await renderWithRouter({ initialEntries: ["/contact"] });

    const nameInput = await screen.findByLabelText("Your name");
    const emailInput = await screen.findByLabelText("Email address");
    const topicSelect = await screen.findByLabelText("Topic");
    const agendaTextarea = await screen.findByLabelText("Message");

    await user.type(nameInput, "  Jane Doe  ");
    await user.type(emailInput, "jane.doe@example.com");
    await user.selectOptions(
      topicSelect,
      contactConfigFixture.topics[0]?.id ?? "",
    );
    await user.type(agendaTextarea, "I would like to discuss possibilities for joint research in HRI.");

    const slotSelect = await screen.findByLabelText("Time slot");
    await user.selectOptions(slotSelect, firstSlot?.id ?? "");

    await user.click(screen.getByRole("button", { name: /request booking/i }));

    expect(
      await screen.findByText(/Your request \(ID: bk-1\)/),
    ).toBeInTheDocument();

    expect(createBookingMock).toHaveBeenCalledWith({
      name: "Jane Doe",
      email: "jane.doe@example.com",
      topic: contactConfigFixture.topics[0]?.id ?? "",
      agenda:
        "I would like to discuss possibilities for joint research in HRI.",
      startTime: firstSlot?.start ?? "",
      durationMinutes: firstSlot
        ? Math.round(
            (new Date(firstSlot.end).getTime() -
              new Date(firstSlot.start).getTime()) /
              60000,
          )
        : 30,
      recaptchaToken: "recaptcha-token-123",
    });

    expect(window.grecaptcha?.execute).toHaveBeenCalledWith(
      contactConfigFixture.recaptchaSiteKey ?? "",
      { action: "submit" },
    );

    delete window.grecaptcha;
  });
});
