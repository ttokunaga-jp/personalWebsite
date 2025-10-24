package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/takumi/personal-website/internal/errs"
	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/service"
)

// BookingHandler exposes endpoints for meeting reservations.
type BookingHandler struct {
	booking service.BookingService
}

// NewBookingHandler wires the booking service into an HTTP handler.
func NewBookingHandler(booking service.BookingService) *BookingHandler {
	return &BookingHandler{booking: booking}
}

// CreateBooking handles incoming reservation requests and orchestrates scheduling.
func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var req model.BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "invalid booking payload", err))
		return
	}

	result, err := h.booking.Book(c.Request.Context(), req)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": result,
	})
}
