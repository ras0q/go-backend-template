package handler

import (
	"context"
	"strings"

	"github.com/ras0q/go-backend-template/internal/api"
)

// GET /api/v1/ping
func (h *Handler) Ping(_ context.Context) (api.PingOK, error) {
	r := strings.NewReader("pong")

	res := api.PingOK{
		Data: r,
	}

	return res, nil
}
