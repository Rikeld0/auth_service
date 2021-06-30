package repo

import (
	_interface "auth/pkg/connector_db/interface"
	"auth/pkg/uuid_my"
	"auth/service/structs"
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"strconv"
	"strings"
)

const (
	token     = "user"
	userQuery = `SELECT uuid_my, email, name FROM main.users `
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
	conn _interface.DB
}

func NewUserDB(conn _interface.DB) Users {
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
	return u.parse(u.conn.QueryRow(ctx, userQuery+`WHERE email=$1 AND password=$2`, email, structs.HexPassword(pass)))
}

func (u *userR) GetUserId(ctx context.Context, uuid string) (*structs.User, error) {
	return u.parse(u.conn.QueryRow(ctx, userQuery+`WHERE uuid_my=$1`, uuid))
}

func (u *userR) InsertUser(ctx context.Context, user *structs.User) (string, error) {
	var (
		sqlArr []string
		arg    []interface{}
	)
	sqlArr = append(sqlArr, strconv.Quote(`uuid`))
	arg = append(arg, uuid_my.GenerateNameUUID(user.Name))
	sqlArr = append(sqlArr, strconv.Quote(`email`))
	arg = append(arg, user.Email)
	sqlArr = append(sqlArr, strconv.Quote(`name`))
	arg = append(arg, user.Name)
	if user.Password == "" {
		return "", errors.New("empty password")
	}
	sqlArr = append(sqlArr, strconv.Quote(`password`))
	arg = append(arg, structs.HexPassword(user.Password))
	v := make([]string, 0, len(sqlArr))
	for i := range sqlArr {
		v = append(v, "$"+strconv.Itoa(i+1))
	}
	var uuid string
	err := u.conn.QueryRow(ctx, `INSERT INTO main.users (`+strings.Join(sqlArr, ",")+`) VALUES (`+strings.Join(v, ",")+`) RETURNING uuid`, arg...).Scan(&uuid)
	if err != nil {
		return "", err
	}
	return uuid, nil
}
