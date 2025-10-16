package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/takumi/personal-website/internal/errs"
	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/service"
)

type ContactHandler struct {
	submission   service.ContactService
	availability service.AvailabilityService
}

func NewContactHandler(submission service.ContactService, availability service.AvailabilityService) *ContactHandler {
	return &ContactHandler{submission: submission, availability: availability}
}

func (h *ContactHandler) SubmitContact(c *gin.Context) {
	var req model.ContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "invalid contact payload", err))
		return
	}

	submission, err := h.submission.SubmitContact(c.Request.Context(), &req)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"data": submission,
	})
}

func (h *ContactHandler) GetAvailability(c *gin.Context) {
	var opts service.AvailabilityOptions

	if start := c.Query("startDate"); start != "" {
		parsed, err := time.Parse("2006-01-02", start)
		if err != nil {
			respondError(c, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "startDate must be in YYYY-MM-DD format", err))
			return
		}
		opts.StartDate = parsed
	}

	if days := c.Query("days"); days != "" {
		value, err := strconv.Atoi(days)
		if err != nil || value <= 0 {
			respondError(c, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "days must be a positive integer", err))
			return
		}
		opts.Days = value
	}

	availability, err := h.availability.GetAvailability(c.Request.Context(), opts)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": availability,
	})
}
