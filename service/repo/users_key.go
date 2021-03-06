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
	Priv []byte `json:"priv"`
	Pub  []byte `json:"pub"`
}

func (u *userKR) Get(ctx context.Context, uuid string) (*structs.UserKey, error) {
	var keys *keysStruct
	keysB, err := u.rClient.Get(strings.Join([]string{prefixKey, uuid}, "/")).Bytes()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(keysB, &keys)
	if err != nil {
		return nil, err
	}
	return gen_key.GenUserKey(uuid, keys.Priv, keys.Pub)
}

func (u *userKR) Put(ctx context.Context, uuid string) error {
	priv, pub := gen_key.GenEcdsaKey()
	keys := &keysStruct{
		Priv: priv,
		Pub:  pub,
	}
	keysB, err := json.Marshal(keys)
	if err != nil {
		return err
	}
	return u.rClient.Set(strings.Join([]string{prefixKey, uuid}, "/"), keysB, 0).Err()

}
