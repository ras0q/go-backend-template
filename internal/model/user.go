package model

import "time"

type (
	User struct {
		ID string `json:"id" gorm:"size:32;primary_key"`

		Email string `json:"email" gorm:"size:191;unique" validate:"omitempty,email"`

		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"-"`
	}
)
