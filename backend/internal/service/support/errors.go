package support

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/takumi/personal-website/internal/errs"
	"github.com/takumi/personal-website/internal/repository"
)

// MapRepositoryError converts repository-level errors into AppError values with consistent semantics.
func MapRepositoryError(err error, resource string) *errs.AppError {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, repository.ErrNotFound):
		return errs.New(errs.CodeNotFound, http.StatusNotFound, fmt.Sprintf("%s not found", resource), err)
	case errors.Is(err, repository.ErrInvalidInput):
		return errs.New(errs.CodeInvalidInput, http.StatusBadRequest, fmt.Sprintf("invalid %s input", resource), err)
	case errors.Is(err, repository.ErrConflict):
		return errs.New(errs.CodeConflict, http.StatusConflict, fmt.Sprintf("%s conflict", resource), err)
	case errors.Is(err, repository.ErrDuplicate):
		return errs.New(errs.CodeConflict, http.StatusConflict, fmt.Sprintf("%s already exists", resource), err)
	default:
		return errs.New(errs.CodeInternal, http.StatusInternalServerError, fmt.Sprintf("failed to process %s", resource), err)
	}
}
