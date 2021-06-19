package usecase

import (
	"auth/service/structs"
	"context"
)

type User interface {
	SignIn(ctx context.Context, user *structs.User) (string, error)
}
