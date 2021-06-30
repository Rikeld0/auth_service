package repo

import (
	"auth/service/structs"
	"context"
)

type row interface {
	Scan(dest ...interface{}) error
}

type Users interface {
	// Get получение пользователя по email и pass
	Get(ctx context.Context, email, pass string) (*structs.User, error)
	// GetUserId получение пользователя по id
	GetUserId(ctx context.Context, uuid string) (*structs.User, error)
	// InsertUser добавить пользователя
	InsertUser(ctx context.Context, user *structs.User) (string, error)
}

type UsersKey interface {
	// Get получить ключи для пользователя
	Get(ctx context.Context, uuid string) (*structs.UserKey, error)
	// Put положить ключи в базу
	Put(ctx context.Context, uuid string, priv, pub string) error
}
