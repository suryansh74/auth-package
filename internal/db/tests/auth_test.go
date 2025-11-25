package tests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/suryansh74/auth-package/internal/db/sqlc"
	"github.com/suryansh74/auth-package/internal/utils"
)

func createRandomUser(t *testing.T) sqlc.User {
	arg := sqlc.CreateUserParams{
		Name:     utils.RandomString(6),
		Email:    utils.RandomEmail(),
		Password: utils.RandomPassword(8),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.UpdatedAt)

	// checking wheater args and inserted values are same or not
	require.Equal(t, arg.Name, user.Name)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Password, user.Password)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUserByEmail(t *testing.T) {
	user := createRandomUser(t)
	returnedUser, err := testQueries.GetUserByEmail(context.Background(), user.Email)

	require.NoError(t, err)
	require.NotEmpty(t, returnedUser)

	// checking wheater args and inserted values are same or not
	require.Equal(t, user.ID, returnedUser.ID)
	require.Equal(t, user.Email, returnedUser.Email)
	require.Equal(t, user.Name, returnedUser.Name)
	require.WithinDuration(t, user.CreatedAt.Time, returnedUser.CreatedAt.Time, time.Second)
	require.WithinDuration(t, user.UpdatedAt.Time, returnedUser.UpdatedAt.Time, time.Second)
}

func TestGetUserByID(t *testing.T) {
	user := createRandomUser(t)
	returnedUser, err := testQueries.GetUser(context.Background(), user.ID)

	require.NoError(t, err)
	require.NotEmpty(t, returnedUser)

	// checking wheater args and inserted values are same or not
	require.Equal(t, user.ID, returnedUser.ID)
	require.Equal(t, user.Email, returnedUser.Email)
	require.Equal(t, user.Name, returnedUser.Name)
	require.WithinDuration(t, user.CreatedAt.Time, returnedUser.CreatedAt.Time, time.Second)
	require.WithinDuration(t, user.UpdatedAt.Time, returnedUser.UpdatedAt.Time, time.Second)
}
