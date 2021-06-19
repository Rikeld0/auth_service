package repo

import (
	"auth/service/structs"
	"context"
)

type Auth interface {
	Auth(ctx context.Context, email, pass string) (id, token string, err error)
}

type row interface {
	Scan(dest ...interface{}) error
}

type Users interface {
	Get(ctx context.Context, email, pass string) (*structs.User, error)
}
