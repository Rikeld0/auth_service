package repo

import (
	"auth/pkg/connector_db"
	"auth/pkg/uuid_my"
	"auth/service/structs"
	"context"
	"errors"
	"io"
	"strconv"
	"strings"
)

const (
	userQuery = `SELECT uuid_my, email, name FROM main.users `
	prefixUIP = "UserIP"
)

type userR struct {
	conn    connector_db.Postgre
	rCliebt connector_db.Redis
}

func NewUserDB(conn connector_db.Postgre, rCliebt connector_db.Redis) Users {
	u := &userR{
		conn:    conn,
		rCliebt: rCliebt,
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

func (u *userR) SaveUserIDAndIP(ctx context.Context, uuid, ip string) error {
	return u.rCliebt.Set(strings.Join([]string{prefixUIP, ip}, "/"), uuid, 0).Err()
}

func (u *userR) FindUserId(ctx context.Context, ip string) (string, error) {
	uuid, err := u.rCliebt.Get(strings.Join([]string{prefixUIP, ip}, "/")).Result()
	if err != nil {
		if err == io.EOF {
			return "", errors.New("client not found")
		}
		return "", err
	}
	return uuid, nil
}
