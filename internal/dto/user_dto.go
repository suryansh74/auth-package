package dto

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserRegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegisterResponse struct {
	UserID      pgtype.UUID `json:"user_id"`
	Name        string      `json:"name"`
	Email       string      `json:"email"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	AccessToken string      `json:"token"`
}

type UserLoginResponse struct {
	UserID      pgtype.UUID `json:"user_id"`
	Email       string      `json:"email"`
	AccessToken string      `json:"token"`
}

type UserResponse struct {
	UserID    pgtype.UUID `json:"user_id"`
	Name      string      `json:"name"`
	Email     string      `json:"email"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}
