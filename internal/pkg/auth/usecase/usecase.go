package usecase

import (
	"context"
	"time"

	"github.com/Alladan04/avito_test/internal/models"
	"github.com/Alladan04/avito_test/internal/pkg/auth"
	"github.com/Alladan04/avito_test/internal/pkg/utils"
)

const (
	JWTLifeTime      = 24 * time.Hour
	defaultImagePath = "default.jpg"
)

type AuthUsecase struct {
	repo auth.AuthRepo
}

func NewAuthUsecase(repo auth.AuthRepo) *AuthUsecase {
	return &AuthUsecase{
		repo: repo,
	}
}

func (uc *AuthUsecase) SignUp(ctx context.Context, data models.UserForm) (models.User, string, time.Time, error) {

	currentTime := time.Now().UTC()
	expTime := currentTime.Add(JWTLifeTime)

	newUser := models.User{
		Username:   data.Username,
		Password:   utils.GetHash(data.Password),
		CreateTime: currentTime,
		IsAdmin:    data.IsAdmin,
	}

	err := uc.repo.AddUser(ctx, newUser)
	if err != nil {

		return models.User{}, "", currentTime, auth.ErrCreatingUser
	}

	token, err := utils.GenToken(newUser, JWTLifeTime)
	if err != nil {
		return models.User{}, "", currentTime, err
	}

	return newUser, token, expTime, nil
}

func (uc *AuthUsecase) SignIn(ctx context.Context, data models.UserForm) (models.User, string, time.Time, error) {

	currentTime := time.Now().UTC()
	expTime := currentTime.Add(JWTLifeTime)

	user, err := uc.repo.GetUserByUsername(ctx, data.Username)
	if err != nil {
		return models.User{}, "", currentTime, auth.ErrUserNotFound
	}
	if user.Password != utils.GetHash(data.Password) {

		return models.User{}, "", currentTime, auth.ErrWrongUserData
	}

	token, err := utils.GenToken(user, JWTLifeTime)
	if err != nil {
		return models.User{}, "", currentTime, err
	}

	return user, token, expTime, nil
}
