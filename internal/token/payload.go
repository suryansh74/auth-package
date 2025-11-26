// Package token response for token functionality
package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token is expired")
)

// Payload contain payload data of token
type Payload struct {
	ID        uuid.UUID   `json:"id"`
	UserID    pgtype.UUID `json:"user_id"`
	Email     string      `json:"email"`
	IssuedAt  time.Time   `json:"issued_at"`
	ExpiredAt time.Time   `json:"expired_at"`
}

func NewPayload(userID pgtype.UUID, email string, duration time.Duration) (*Payload, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &Payload{
		ID:        id,
		UserID:    userID,
		Email:     email,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
