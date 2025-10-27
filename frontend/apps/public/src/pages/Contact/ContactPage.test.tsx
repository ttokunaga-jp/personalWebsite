import { apiClient } from "@shared/lib/api-client";
import { screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import type { AxiosResponse } from "axios";
import { afterEach, vi } from "vitest";

import type {
  ContactAvailabilityResponse,
  ContactConfigResponse,
  CreateBookingResponse
} from "../../modules/public-api";
import { renderWithRouter } from "../../test-utils/renderWithRouter";

const createAxiosResponse = <T,>(data: T): AxiosResponse<T> =>
  ({
    data,
    status: 200,
    statusText: "OK",
    headers: {},
    config: {}
  }) as AxiosResponse<T>;

describe("ContactPage", () => {
  beforeEach(() => {
    window.grecaptcha = {
      ready: (callback: () => void) => callback(),
      execute: vi.fn().mockResolvedValue("recaptcha-token")
    };
  });

  afterEach(() => {
    delete window.grecaptcha;
  });

  it("validates form input and submits booking request", async () => {
    const user = userEvent.setup();

    const availability: ContactAvailabilityResponse = {
      timezone: "Asia/Tokyo",
      generatedAt: "2024-01-01T00:00:00Z",
      slots: [
        {
          id: "slot-1",
          start: "2025-01-05T09:00:00+09:00",
          end: "2025-01-05T09:30:00+09:00",
          isBookable: true
        }
      ]
    };

    const config: ContactConfigResponse = {
      topics: ["General"],
      recaptchaSiteKey: "test-site-key",
      minimumLeadHours: 24,
      consentText: "We will respond within 2 business days."
    };

    const bookingResponse: CreateBookingResponse = {
      bookingId: "booking-123",
      status: "pending"
    };

    const getSpy = vi.spyOn(apiClient, "get").mockImplementation((url) => {
      if (typeof url === "string" && url.includes("/v1/public/contact/availability")) {
        return Promise.resolve(createAxiosResponse(availability));
      }

      if (typeof url === "string" && url.includes("/v1/public/contact/config")) {
        return Promise.resolve(createAxiosResponse(config));
      }

      if (typeof url === "string" && url === "/health") {
        return Promise.resolve(createAxiosResponse({ status: "ok" }));
      }

      return Promise.resolve(createAxiosResponse({}));
    });

    const postSpy = vi
      .spyOn(apiClient, "post")
      .mockResolvedValue(createAxiosResponse(bookingResponse));

    await renderWithRouter({ initialEntries: ["/contact"] });

    await waitFor(() => {
      expect(getSpy).toHaveBeenCalledWith(
        "/v1/public/contact/availability",
        expect.objectContaining({ signal: expect.any(AbortSignal) })
      );
    });

    const submitButton = await screen.findByRole("button", { name: /request booking/i });
    const topicSelect = await screen.findByLabelText(/topic/i);

    await waitFor(() => {
      expect(topicSelect).not.toBeDisabled();
      expect(submitButton).not.toBeDisabled();
    });

    await user.click(submitButton);

    expect(await screen.findByText(/please provide your name/i)).toBeInTheDocument();
    expect(await screen.findByText(/an email address is required/i)).toBeInTheDocument();
    expect(await screen.findByText(/share at least 20 characters/i)).toBeInTheDocument();
    expect(await screen.findByText(/select an available time slot/i)).toBeInTheDocument();
    expect(topicSelect).toHaveAttribute("aria-invalid", "true");

    const slotButton = await screen.findByRole("button", { name: /Ends/i });
    await user.click(slotButton);

    const nameInput = await screen.findByLabelText(/your name/i);
    await user.type(nameInput, "Jane Doe");

    const emailInput = await screen.findByLabelText(/email address/i);
    await user.type(emailInput, "jane@example.com");

    await user.selectOptions(topicSelect, "General");

    const messageInput = await screen.findByLabelText(/message/i);
    await user.type(messageInput, "I'd like to discuss collaboration opportunities.");

    await user.click(submitButton);

    await waitFor(() => {
      expect(postSpy).toHaveBeenCalledWith(
        "/v1/public/contact/bookings",
        expect.objectContaining({
          name: "Jane Doe",
          email: "jane@example.com",
          topic: "General",
          slotId: "slot-1",
          recaptchaToken: "recaptcha-token"
        })
      );
    });

    expect(
      await screen.findByText(/Thank you! Your request \(ID: booking-123\)/i)
    ).toBeInTheDocument();
  });
});
