package handler

import (
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
