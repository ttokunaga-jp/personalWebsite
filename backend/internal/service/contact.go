package service

import (
	"context"
	"net/http"
	"strings"

	"github.com/takumi/personal-website/internal/errs"
	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

// ContactService handles contact form submissions.
type ContactService interface {
	SubmitContact(ctx context.Context, req *model.ContactRequest) (*model.ContactSubmission, error)
}

type contactService struct {
	repo repository.ContactRepository
}

func NewContactService(repo repository.ContactRepository) ContactService {
	return &contactService{repo: repo}
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
