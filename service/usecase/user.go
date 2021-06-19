package usecase

import (
	"auth/pkg/myJwt"
	"auth/service/repo"
	"auth/service/structs"
	"context"
	"fmt"
)

type userService struct {
	repo repo.Users
}

func NewUserService(repo repo.Users) User {
	return &userService{
		repo: repo,
	}
}

func (u *userService) SignIn(ctx context.Context, user *structs.User) (string, error) {
	usr, err := u.repo.Get(ctx, user.Email, user.Password)
	if err != nil {
		return "", err
	}
	jwtToken, err := myJwt.GenerateJWT(usr.Uuid)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return jwtToken.AccessToken, nil
}
