package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suryansh74/auth-package/internal/db/sqlc"
)

type Auth struct {
	*sqlc.Queries
	connPool *pgxpool.Pool
}

const (
	// Database connection string
	connString = "postgresql://root:secret@192.168.29.20:5432/authentication?sslmode=disable"
)

func NewAuth() *Auth {
	// Create connection pool
	connPool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}

	// Verify connection
	if err := connPool.Ping(context.Background()); err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}

	return &Auth{
		Queries:  sqlc.New(connPool),
		connPool: connPool,
	}
}

// Close closes the database connection pool
func (a *Auth) Close() {
	a.connPool.Close()
}
