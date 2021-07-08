package structs

import (
	"crypto/sha256"
	"encoding/hex"
)

type Cl struct {
	Exp int64  `json:"exp"`
	Jti string `json:"jti"`
	Iat int64  `json:"iat"`
	Iss string `json:"iss"`
}

type User struct {
	Uuid     string
	Email    string
	Name     string
	Password string
}

func NewUser() *User {
	return &User{}
}

func HexPassword(in string) (out string) {
	b := sha256.Sum256([]byte(in))
	return hex.EncodeToString(b[:])
}
