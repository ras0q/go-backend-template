package core

import (
	"backend/core/internal/handler"
	"backend/frontend"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(h *handler.Handler, e *echo.Echo) {
	e.StaticFS("/", frontend.Static)

	v1API := e.Group("/api/v1")

	// ping API
	pingAPI := v1API.Group("/ping")
	{
		pingAPI.GET("", h.Ping)
	}

	// user API
	userAPI := v1API.Group("/users")
	{
		userAPI.GET("", h.GetUsers)
		userAPI.POST("", h.CreateUser)
		userAPI.GET("/:userID", h.GetUser)
	}
}
