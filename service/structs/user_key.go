package structs

import "github.com/lestrrat-go/jwx/jwk"

type UserKey struct {
	PrivateKey jwk.Key
	PublicKey  jwk.Key
}

func NewUserKey() *UserKey {
	return &UserKey{}
}
