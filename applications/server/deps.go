package server

import (
	"backend/internal/handler"
	"backend/internal/repository"

	"github.com/jmoiron/sqlx"
)

type Deps struct {
	Handler *handler.Handler
}

func InjectDeps(db *sqlx.DB) *Deps {
	repo := repository.New(db)
	h := handler.New(repo)

	return &Deps{
		Handler: h,
	}
}
