package handler

import (
	"context"
	"fmt"

	"github.com/ras0q/go-backend-template/internal/api"
	"github.com/ras0q/go-backend-template/internal/repository"
)

// POST /api/v1/users
func (h *Handler) CreateUser(ctx context.Context, req *api.CreateUserReq) (*api.CreateUser, error) {
	userID, err := h.repo.CreateUser(ctx, repository.CreateUserParams{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("create user to repository: %w", err)
	}

	res := &api.CreateUser{
		ID: userID,
	}

	return res, nil
}

// GET /api/v1/users/:userID
func (h *Handler) GetUser(ctx context.Context, params api.GetUserParams) (*api.User, error) {
	user, err := h.repo.GetUser(ctx, params.UserID)
	if err != nil {
		return nil, fmt.Errorf("get user from repository: %w", err)
	}

	photo, err := h.photo.GetPhoto(1)
	if err != nil {
		return nil, fmt.Errorf("get user icon from photoapi: %w", err)
	}

	res := &api.User{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		IconUrl: photo.ThumbnailURL,
	}

	return res, nil
}

// GET /api/v1/users
func (h *Handler) GetUsers(ctx context.Context) ([]api.User, error) {
	users, err := h.repo.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("get users from repository: %w", err)
	}

	photo, err := h.photo.GetPhoto(1)
	if err != nil {
		return nil, fmt.Errorf("get user icon from photoapi: %w", err)
	}

	res := make([]api.User, 0, len(users))
	for _, user := range users {
		res = append(res, api.User{
			ID:      user.ID,
			Name:    user.Name,
			Email:   user.Email,
			IconUrl: photo.ThumbnailURL,
		})
	}

	return res, nil
}
