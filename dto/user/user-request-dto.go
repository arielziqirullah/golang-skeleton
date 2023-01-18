package user

import "time"

type UserUpdateRequestDTO struct {
	ID        uint64 `json:"id" form:"id"`
	Name      string `json:"name" form:"name" validate:"required"`
	Email     string `json:"email" form:"email" validate:"email"`
	Password  string `json:"password,omitempty" form:"password,omitempty"`
	UpdatedAt time.Time
}
