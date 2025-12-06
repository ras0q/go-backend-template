package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/ras0q/go-backend-template/api"
	"github.com/ras0q/go-backend-template/core/internal/repository"
)

type Handler struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) NewError(ctx context.Context, err error) *api.ErrorStatusCode {
	if apiErr, ok := err.(*api.ErrorStatusCode); ok {
		return apiErr
	}

	slog.ErrorContext(ctx, "internal server error", "error", err)

	return &api.ErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: api.Error{
			Message: "internal server error",
		},
	}
}
