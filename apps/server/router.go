package server

import (
	"backend/internal/handler"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(h *handler.Handler, g *echo.Group) {
	// ping API
	pingAPI := g.Group("/ping")
	{
		pingAPI.GET("", h.Ping)
	}

	// user API
	userAPI := g.Group("/users")
	{
		userAPI.GET("", h.GetUsers)
		userAPI.POST("", h.CreateUser)
		userAPI.GET("/:userID", h.GetUser)
	}
}
