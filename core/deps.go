package core

import (
	"github.com/ras0q/go-backend-template/core/internal/handler"
	"github.com/ras0q/go-backend-template/core/internal/repository"

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
