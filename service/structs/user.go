package structs

import (
	"crypto/sha256"
	"encoding/hex"
)

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
