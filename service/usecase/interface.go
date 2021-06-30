package usecase

import (
	"auth/pkg/myJwt"
	"auth/service/structs"
	"context"
)

type Auth interface {
	// Auth аутенфикация пользовтеля
	Auth(ctx context.Context) (context.Context, error)
}

type User interface {
	Auth
	SignIn(ctx context.Context, user *structs.User) (*myJwt.JWT, error)
	SignUp(ctx context.Context, user *structs.User) error
	GetMsg(ctx context.Context) (string, error)
}
