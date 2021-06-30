package repo

import (
	_interface "auth/pkg/connector_db/interface"
	"auth/pkg/gen_key"
	"auth/service/structs"
	"context"
)

const (
	userKeyQueryGet = `SELECT * FROM main.userkey `
)

type userKR struct {
	conn _interface.DB
}

func NewUserKR(conn _interface.DB) UsersKey {
	return &userKR{conn: conn}
}

type keys struct {
	uuid string
	priv string
	pub  string
}

func (u *userKR) parse(rows row) (*structs.UserKey, error) {
	out := &keys{}
	err := rows.Scan(
		&out.uuid,
		&out.priv,
		&out.pub,
	)
	if err != nil {
		return nil, err
	}
	return gen_key.GenUserKey(out.uuid, []byte(out.priv), []byte(out.pub))
}

func (u *userKR) Get(ctx context.Context, uuid string) (*structs.UserKey, error) {
	return u.parse(u.conn.QueryRow(ctx, userKeyQueryGet+`WHERE uuid_my=$1`, uuid))
}

func (u *userKR) Put(ctx context.Context, uuid string, priv, pub string) error {
	return u.conn.Exec(ctx, `INSERT INTO main.userkey (uuid_my, priv, pub) VALUES ($1, $2, $3)`, uuid, priv, pub)

}
