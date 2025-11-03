package service

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/takumi/personal-website/internal/errs"
	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

// ContactService handles contact form submissions and configuration retrieval.
type ContactService interface {
	SubmitContact(ctx context.Context, req *model.ContactRequest) (*model.ContactSubmission, error)
	GetContactSettings(ctx context.Context) (*model.ContactFormSettingsV2, error)
}

type contactService struct {
	repo     repository.ContactRepository
	settings repository.ContactFormSettingsRepository
}

func NewContactService(repo repository.ContactRepository, settings repository.ContactFormSettingsRepository) ContactService {
	return &contactService{repo: repo, settings: settings}
}

func (s *contactService) SubmitContact(ctx context.Context, req *model.ContactRequest) (*model.ContactSubmission, error) {
	if req == nil {
		return nil, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "request body is required", nil)
	}
	if strings.TrimSpace(req.Email) == "" {
		return nil, errs.New(errs.CodeInvalidInput, http.StatusBadRequest, "email is required", nil)
	}

	submission, err := s.repo.CreateSubmission(ctx, req)
	if err != nil {
		return nil, errs.New(errs.CodeInternal, http.StatusInternalServerError, "failed to queue contact request", err)
	}

	return submission, nil
}

func (s *contactService) GetContactSettings(ctx context.Context) (*model.ContactFormSettingsV2, error) {
	if s.settings == nil {
		return nil, errs.New(errs.CodeInternal, http.StatusInternalServerError, "contact settings repository not configured", nil)
	}

	settings, err := s.settings.GetContactFormSettings(ctx)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, errs.New(errs.CodeNotFound, http.StatusNotFound, "contact settings not found", err)
		}
		return nil, errs.New(errs.CodeInternal, http.StatusInternalServerError, "failed to load contact settings", err)
	}

	return settings, nil
}
