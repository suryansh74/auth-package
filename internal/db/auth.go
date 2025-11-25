package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suryansh74/auth-package/internal/db/sqlc"
)

type Auth struct {
	*sqlc.Queries
	connPool *pgxpool.Pool
}

func NewAuth(db *pgxpool.Pool) *Auth {
	return &Auth{
		Queries:  sqlc.New(db),
		connPool: db,
	}
}
