package injector

import (
	"backend/internal/handler"
	"backend/internal/repository"

	"github.com/jmoiron/sqlx"
)

type Dependency struct {
	Handler *handler.Handler
}

func Inject(db *sqlx.DB) *Dependency {
	repo := repository.New(db)
	h := handler.New(repo)

	return &Dependency{
		Handler: h,
	}
}
