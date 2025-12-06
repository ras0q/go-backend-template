package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/ras0q/go-backend-template/internal/api"
	"github.com/ras0q/go-backend-template/internal/repository"
	"github.com/ras0q/go-backend-template/internal/service/photo"
)

type Handler struct {
	photo *photo.Service
	repo  *repository.Repository
}

func New(
	photo *photo.Service,
	repo *repository.Repository,
) *Handler {
	return &Handler{
		photo,
		repo,
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
