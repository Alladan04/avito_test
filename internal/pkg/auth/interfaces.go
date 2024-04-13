package auth

import (
	"context"
	"errors"
	"time"

	"github.com/Alladan04/avito_test/internal/models"
)

var (
	ErrCreatingUser     = errors.New("this username is already taken")
	ErrIncorrectPayload = errors.New("incorrect data format")
	ErrUserNotFound     = errors.New("user not found")
	ErrWrongPassword    = errors.New("wrong password")
	ErrWrongUserData    = errors.New("wrong username or password")
)

type AuthRepo interface {
	GetUserByUsername(context.Context, string) (models.User, error)
	AddUser(context.Context, models.User) error
}
type AuthUsecase interface {
	SignIn(context.Context, models.UserForm) (models.User, string, time.Time, error)
	SignUp(context.Context, models.UserForm) (models.User, string, time.Time, error)
}
