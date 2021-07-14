package usecase

import (
	"auth/service/structs"
	"context"
)

type Auth interface {
	// Auth аутенфикация пользовтеля
	Auth(ctx context.Context) (context.Context, error)
}

type User interface {
	Auth
	// SignIn найти пользователя
	SignIn(ctx context.Context, user *structs.User) (*structs.JWT, error)
	// SignUp добавить пользователя
	SignUp(ctx context.Context, user *structs.User) error
	// SignOut выход из системы
	SignOut(ctx context.Context) error
	GetMsg(ctx context.Context) (string, error)
}
