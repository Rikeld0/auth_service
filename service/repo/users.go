package repo

import (
	"auth/service/structs"
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"github.com/jackc/pgx"
)

const (
	token     = "user"
	userQuery = `SELECT uuid, email, name FROM main.users `
)

func getUserCTX(ctx context.Context) (*structs.User, error) {
	buf, ok := ctx.Value(token).([]byte)
	if !ok {
		return nil, errors.New("no user")
	}
	out := &structs.User{}
	if err := gob.NewDecoder(bytes.NewReader(buf)).Decode(out); err != nil {
		return nil, err
	}
	return out, nil
}

type userR struct {
	conn *pgx.Conn
}

func NewUserDB(conn *pgx.Conn) Users {
	u := &userR{
		conn: conn,
	}
	return u
}

func (u *userR) parse(rows row) (out *structs.User, err error) {
	out = structs.NewUser()
	err = rows.Scan(
		&out.Uuid,
		&out.Email,
		&out.Name,
	)
	return
}

func (u *userR) Get(ctx context.Context, email, pass string) (*structs.User, error) {
	return u.parse(u.conn.QueryRow(userQuery+`WHERE email=$1 AND password=$2`, email, structs.HexPassword(pass)))
}
