package repo

import (
	"auth/pkg/connector_db"
	"auth/pkg/myJwt"
	"auth/pkg/uuid_my"
	"auth/service/structs"
	"encoding/json"
	"io"
	"strings"
)

const perfixJU = "JWT"

type jwtRepo struct {
	rClient connector_db.Redis
}

func NewJwtRepo(rClient connector_db.Redis) Jwt {
	return &jwtRepo{rClient: rClient}
}

func generateJWTID() (string, error) {
	//TODO: написать генерацию id токена
	id := uuid_my.GenerateUUID()
	return id, nil
}

func (j *jwtRepo) Generate(uuid string, keys *structs.UserKey) (*structs.JWT, error) {
	jti, err := generateJWTID()
	if err != nil {
		return nil, err
	}
	err = j.rClient.Set(strings.Join([]string{perfixJU, uuid}, "/"), jti, 0).Err()
	if err != nil {
		return nil, err
	}
	tokens, err := myJwt.GenerateJWT(uuid, jti, keys)
	if err != nil {
		return nil, err
	}
	return &structs.JWT{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken}, nil
}

func (j *jwtRepo) Validate(token string, keys *structs.UserKey) error {
	claims, err := myJwt.VerefyJWT(token, keys)
	if err != nil {
		return err
	}
	c := structs.Cl{}
	err = json.Unmarshal(claims, &c)
	if err != nil {
		return err
	}
	_, err = j.rClient.Get(strings.Join([]string{perfixJU, c.Iss}, "/")).Result()
	if err != nil {
		if err == io.EOF {
			return err
		}
		return err
	}
	return nil
}
