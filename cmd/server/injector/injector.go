package injector

import (
	"backend/internal/handler"
	"backend/internal/repository"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type Dependency struct {
	handler *handler.Handler
}

func Inject(db *sqlx.DB) *Dependency {
	repo := repository.New(db)
	h := handler.New(repo)

	return &Dependency{
		handler: h,
	}
}

func (d *Dependency) SetupRoutes(g *echo.Group) {
	// TODO: handler.SetupRoutesを呼び出す or 直接書く？
	d.handler.SetupRoutes(g)
}
