package services

import (
	"context"
	"database/sql"
	"errors"

	customError "github.com/suryansh74/auth-package/internal/appError"
	"github.com/suryansh74/auth-package/internal/db"
	"github.com/suryansh74/auth-package/internal/db/sqlc"
	"github.com/suryansh74/auth-package/internal/dto"
	"github.com/suryansh74/auth-package/internal/utils"
)

type AuthService interface {
	// TODO: add token in return along with response
	Register(dto.UserRegisterRequest) (*dto.UserRegisterResponse, error)
	Login(dto.UserLoginRequest) (*dto.UserLoginResponse, error)
}

type Authenticator struct {
	auth db.Auth
}

func NewAuthenticator(auth db.Auth) AuthService {
	return &Authenticator{
		auth: auth,
	}
}

func (a *Authenticator) Register(req dto.UserRegisterRequest) (*dto.UserRegisterResponse, error) {
	// check if user is already existed
	exists, _, err := a.userExists(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, customError.ErrUserAlreadyExist
	}

	// create hash password
	hashedPassword, err := utils.HashedPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// insert into table
	arg := sqlc.CreateUserParams{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}
	user, err := a.auth.CreateUser(context.Background(), arg)
	if err != nil {
		return nil, customError.UnExpectedError
	}

	userResponse := dto.UserRegisterResponse{
		Name:      user.Name,
		Email:     req.Email,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
	}
	return &userResponse, nil
}

func (a *Authenticator) Login(req dto.UserLoginRequest) (*dto.UserLoginResponse, error) {
	// check if user is existed or not
	exists, user, err := a.userExists(req.Email)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, customError.ErrUserNotFound
	}

	err = utils.CheckPassword(req.Password, user.Password)
	if err != nil {
		return nil, errors.New("password not matched")
	}
	userResponse := dto.UserLoginResponse{
		Email: req.Email,
	}
	return &userResponse, nil
}

func (a *Authenticator) userExists(email string) (bool, *sqlc.User, error) {
	user, err := a.auth.GetUserByEmail(context.Background(), email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil, nil
		}
		return false, nil, customError.UnExpectedError
	}
	return true, &user, nil
}
