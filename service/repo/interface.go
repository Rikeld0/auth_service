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
	// SaveUserIDAndIP запоминаем связку id пользователя и ip адрес устройства ????
	SaveUserIDAndIP(ctx context.Context, uuid, ip string) error
	// FindUserId ищем пользователя по ip адрессу
	FindUserId(ctx context.Context, ip string) (string, error)
}

type UsersKey interface {
	// Get получаем ecdsa ключи и генерируем пользовательские
	Get(ctx context.Context, uuid string) (*structs.UserKey, error)
	// Put генерируем и кладем в базу ключи ecdsa
	Put(ctx context.Context, uuid string) error
}

type Jwt interface {
	// Generate генерация токен
	Generate(uuid string, keys *structs.UserKey) (*structs.JWT, error)
	// Validate проверка токена
	Validate(token string, keys *structs.UserKey) error
}
