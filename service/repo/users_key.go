package repo

import (
	"auth/pkg/connector_db"
	"auth/pkg/gen_key"
	"auth/service/structs"
	"context"
	"encoding/json"
	"strings"
)

const prefixKey = "UserKey"

type userKR struct {
	rClient connector_db.Redis
}

func NewUserKR(rClient connector_db.Redis) UsersKey {
	return &userKR{
		rClient: rClient,
	}
}

type keysStruct struct {
	priv []byte
	pub  []byte
}

func (u *userKR) Get(ctx context.Context, uuid string) (*structs.UserKey, error) {
	var keys *keysStruct
	keysB, err := u.rClient.Get(strings.Join([]string{prefixKey, uuid}, "/")).Bytes()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(keysB, &keys)
	return gen_key.GenUserKey(uuid, keys.priv, keys.pub)
}

func (u *userKR) Put(ctx context.Context, uuid string) error {
	priv, pub := gen_key.GenEcdsaKey()
	keys := &keysStruct{
		priv: priv,
		pub:  pub,
	}
	keysB, err := json.Marshal(keys)
	if err != nil {
		return err
	}
	return u.rClient.Set(strings.Join([]string{prefixKey, uuid}, "/"), keysB, 0).Err()

}
