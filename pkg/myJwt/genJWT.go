package myJwt

import (
	"auth/service/structs"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwe"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/lestrrat-go/jwx/jwt"
	"log"
	"time"
)

//go:embed key/ecprivkey.pem
var ecprivkey []byte

//go:embed key/ecpubkey.pem
var ecpubkey []byte

type JWT struct {
	AccessToken  string
	RefreshToken string
}

func tokenBytetoString(t []byte, err error) string {
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(t)
}

func GenerateJWT(id string, keys *structs.UserKey) (*JWT, error) {
	jt := &JWT{}
	tokenAccess := jwt.New()
	_ = tokenAccess.Set(jwt.IssuerKey, id)
	_ = tokenAccess.Set(jwt.IssuedAtKey, time.Now().Unix())
	_ = tokenAccess.Set(jwt.ExpirationKey, time.Now().Add(time.Minute*20).Unix())
	buf, err := json.Marshal(tokenAccess)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	hdr := jws.NewHeaders()
	if err = hdr.Set(jws.TypeKey, `JWT`); err != nil {
		fmt.Println(err)
		return nil, err
	}
	encrypted, err := jwe.Encrypt(buf, jwa.ECDH_ES, keys.PublicKey, jwa.A128CBC_HS256, jwa.NoCompress)
	if err != nil {
		log.Printf("failed to encrypt payload: %s", err)
		return nil, err
	}
	jt.AccessToken = tokenBytetoString(jws.Sign(encrypted, jwa.ES256, keys.PrivateKey, jws.WithHeaders(hdr)))
	return jt, nil
}

func VerefyJWT(token string, keys *structs.UserKey) (string, error) {
	verified, err := jws.Verify([]byte(token), jwa.ES256, keys.PublicKey)
	if err != nil {
		log.Printf("failed to verify message: %s", err)
		return "", err
	}
	claims, err := jwe.Decrypt(verified, jwa.ECDH_ES, keys.PrivateKey)
	if err != nil {
		log.Printf("failed to decrypt: %s", err)
		return "", err
	}
	//fmt.Println(string(decrypted))
	return string(claims), nil
}
