package structs

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/bcrypt"
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
	Rights   string
}

func NewUser() *User {
	return &User{}
}

func HexPassword(in string) (out string) {
	b := sha256.Sum256([]byte(in))
	return hex.EncodeToString(b[:])
}

func GenHashPass(pass string) string {
	b, err := bcrypt.GenerateFromPassword([]byte(pass), 14)
	if err != nil {
		return ""
	}
	return string(b)
}

func CheckPass(hash, pass string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass)); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
