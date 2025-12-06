package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ras0q/go-backend-template/api"
	"github.com/ras0q/go-backend-template/core/internal/repository"
	"github.com/ras0q/go-backend-template/core/internal/services/photoapi"

	vd "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// POST /api/v1/users
func (h *Handler) CreateUser(ctx context.Context, req *api.CreateUserReq) (*api.CreateUser, error) {
	err := vd.ValidateStruct(
		req,
		vd.Field(&req.Name, vd.Required),
		vd.Field(&req.Email, vd.Required, is.Email),
	)
	if err != nil {
		return nil, &api.ErrorStatusCode{
			StatusCode: http.StatusBadRequest,
			Response: api.Error{
				Message: fmt.Sprintf("invalid request: %s", err.Error()),
			},
		}
	}

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

	photo, err := photoapi.GetPhoto()
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

	photo, err := photoapi.GetPhoto()
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
