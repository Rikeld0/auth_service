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
	// SignIn найти уч запись
	SignIn(ctx context.Context, user *structs.User) (*structs.JWT, error)
	// SignUp добавить уч запись
	SignUp(ctx context.Context, user *structs.User) error
	GetMsg(ctx context.Context) (string, error)
}
