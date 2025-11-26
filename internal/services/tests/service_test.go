// internal/services/user_service_test.go
package services

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"

	customError "github.com/suryansh74/auth-package/internal/apperrors"
	"github.com/suryansh74/auth-package/internal/db/mock"
	"github.com/suryansh74/auth-package/internal/db/sqlc"
	"github.com/suryansh74/auth-package/internal/dto"
	"github.com/suryansh74/auth-package/internal/services"
	"github.com/suryansh74/auth-package/internal/utils"
)

func TestRegister(t *testing.T) {
	testCases := []struct {
		name          string
		request       dto.UserRegisterRequest
		buildStubs    func(mockAuth *mock.MockAuth)
		checkResponse func(t *testing.T, resp *dto.UserRegisterResponse, err error)
	}{
		{
			name: "OK",
			request: dto.UserRegisterRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
			},
			buildStubs: func(mockAuth *mock.MockAuth) {
				// User doesn't exist - return ErrNoRows
				mockAuth.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq("john@example.com")).
					Times(1).
					Return(sqlc.User{}, sql.ErrNoRows)

				// CreateUser succeeds
				hashedPassword, _ := utils.HashedPassword("password123")
				_ = sqlc.CreateUserParams{
					Name:     "John Doe",
					Email:    "john@example.com",
					Password: hashedPassword,
				}

				mockAuth.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					DoAndReturn(func(ctx interface{}, params sqlc.CreateUserParams) (sqlc.User, error) {
						return sqlc.User{
							ID:        pgtype.UUID{Valid: true},
							Name:      params.Name,
							Email:     params.Email,
							Password:  params.Password,
							CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
							UpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
						}, nil
					})
			},
			checkResponse: func(t *testing.T, resp *dto.UserRegisterResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, "John Doe", resp.Name)
				require.Equal(t, "john@example.com", resp.Email)
				require.NotZero(t, resp.CreatedAt)
				require.NotZero(t, resp.UpdatedAt)
			},
		},
		{
			name: "UserAlreadyExists",
			request: dto.UserRegisterRequest{
				Name:     "Jane Doe",
				Email:    "jane@example.com",
				Password: "password123",
			},
			buildStubs: func(mockAuth *mock.MockAuth) {
				// User exists - return a user
				existingUser := sqlc.User{
					ID:    pgtype.UUID{Valid: true},
					Name:  "Jane Doe",
					Email: "jane@example.com",
				}

				mockAuth.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq("jane@example.com")).
					Times(1).
					Return(existingUser, nil)

				// CreateUser should not be called
				mockAuth.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, resp *dto.UserRegisterResponse, err error) {
				require.Error(t, err)
				require.Nil(t, resp)
				require.Equal(t, customError.ErrUserAlreadyExist, err)
			},
		},
		{
			name: "DatabaseErrorOnCheck",
			request: dto.UserRegisterRequest{
				Name:     "Bob Smith",
				Email:    "bob@example.com",
				Password: "password123",
			},
			buildStubs: func(mockAuth *mock.MockAuth) {
				// Database connection error
				mockAuth.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq("bob@example.com")).
					Times(1).
					Return(sqlc.User{}, sql.ErrConnDone)

				// CreateUser should not be called
				mockAuth.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, resp *dto.UserRegisterResponse, err error) {
				require.Error(t, err)
				require.Nil(t, resp)
				require.Equal(t, customError.UnExpectedError, err)
			},
		},
		{
			name: "DatabaseErrorOnCreate",
			request: dto.UserRegisterRequest{
				Name:     "Alice Wonder",
				Email:    "alice@example.com",
				Password: "password123",
			},
			buildStubs: func(mockAuth *mock.MockAuth) {
				// User doesn't exist
				mockAuth.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq("alice@example.com")).
					Times(1).
					Return(sqlc.User{}, sql.ErrNoRows)

				// CreateUser fails
				mockAuth.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(sqlc.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, resp *dto.UserRegisterResponse, err error) {
				require.Error(t, err)
				require.Nil(t, resp)
				require.Equal(t, customError.UnExpectedError, err)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuth := mock.NewMockAuth(ctrl)
			tc.buildStubs(mockAuth)

			authService := services.NewAuthenticator(mockAuth)
			resp, err := authService.Register(context.Background(), tc.request)

			tc.checkResponse(t, resp, err)
		})
	}
}

func TestLogin(t *testing.T) {
	password := "password123"
	hashedPassword, _ := utils.HashedPassword(password)

	testCases := []struct {
		name          string
		request       dto.UserLoginRequest
		buildStubs    func(mockAuth *mock.MockAuth)
		checkResponse func(t *testing.T, resp *dto.UserLoginResponse, err error)
	}{
		{
			name: "OK",
			request: dto.UserLoginRequest{
				Email:    "john@example.com",
				Password: password,
			},
			buildStubs: func(mockAuth *mock.MockAuth) {
				user := sqlc.User{
					ID:        pgtype.UUID{Valid: true},
					Name:      "John Doe",
					Email:     "john@example.com",
					Password:  hashedPassword,
					CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
					UpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
				}

				mockAuth.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq("john@example.com")).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, resp *dto.UserLoginResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, "john@example.com", resp.Email)
			},
		},
		{
			name: "UserNotFound",
			request: dto.UserLoginRequest{
				Email:    "notfound@example.com",
				Password: password,
			},
			buildStubs: func(mockAuth *mock.MockAuth) {
				mockAuth.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq("notfound@example.com")).
					Times(1).
					Return(sqlc.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, resp *dto.UserLoginResponse, err error) {
				require.Error(t, err)
				require.Nil(t, resp)
				require.Equal(t, customError.ErrUserNotFound, err)
			},
		},
		{
			name: "WrongPassword",
			request: dto.UserLoginRequest{
				Email:    "john@example.com",
				Password: "wrongpassword",
			},
			buildStubs: func(mockAuth *mock.MockAuth) {
				user := sqlc.User{
					ID:        pgtype.UUID{Valid: true},
					Name:      "John Doe",
					Email:     "john@example.com",
					Password:  hashedPassword, // Correct password hash
					CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
					UpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
				}

				mockAuth.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq("john@example.com")).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, resp *dto.UserLoginResponse, err error) {
				require.Error(t, err)
				require.Nil(t, resp)
				require.Contains(t, err.Error(), "password not matched")
			},
		},
		{
			name: "DatabaseError",
			request: dto.UserLoginRequest{
				Email:    "john@example.com",
				Password: password,
			},
			buildStubs: func(mockAuth *mock.MockAuth) {
				mockAuth.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq("john@example.com")).
					Times(1).
					Return(sqlc.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, resp *dto.UserLoginResponse, err error) {
				require.Error(t, err)
				require.Nil(t, resp)
				require.Equal(t, customError.UnExpectedError, err)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuth := mock.NewMockAuth(ctrl)
			tc.buildStubs(mockAuth)

			authService := services.NewAuthenticator(mockAuth)
			resp, err := authService.Login(context.Background(), tc.request)

			tc.checkResponse(t, resp, err)
		})
	}
}
