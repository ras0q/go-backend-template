package core

import (
	"github.com/ras0q/go-backend-template/core/internal/handler"
	"github.com/ras0q/go-backend-template/core/internal/repository"
	photo_service "github.com/ras0q/go-backend-template/core/internal/service/photo"

	"github.com/jmoiron/sqlx"
)

type Deps struct {
	Handler *handler.Handler
}

func InjectDeps(db *sqlx.DB) *Deps {
	photo := photo_service.NewPhotoService()
	repo := repository.New(db)
	h := handler.New(photo, repo)

	return &Deps{
		Handler: h,
	}
}
