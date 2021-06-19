package myJwt

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwe"
	"github.com/lestrrat-go/jwx/jwk"
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

var (
	privatKey jwk.Key
	publicKey jwk.Key
)

func tokenBytetoString(t []byte, err error) string {
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(t)
}

func init() {
	var err error
	privatKey, err = jwk.ParseKey(ecprivkey, jwk.WithPEM(true))
	if err != nil {
		fmt.Printf("failed to create JWK: %s\n", err)
		return
	}
	_ = privatKey.Set(jwk.KeyIDKey, "mykey")
	publicKey, err = jwk.ParseKey(ecpubkey, jwk.WithPEM(true))
	if err != nil {
		fmt.Printf("failed to create JWK: %s\n", err)
		return
	}
	_ = publicKey.Set(jwk.KeyIDKey, "mykey")
}

func GenerateJWT(id string) (*JWT, error) {
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
	encrypted, err := jwe.Encrypt(buf, jwa.ECDH_ES, publicKey, jwa.A128CBC_HS256, jwa.NoCompress)
	if err != nil {
		log.Printf("failed to encrypt payload: %s", err)
		return nil, err
	}
	jt.AccessToken = tokenBytetoString(jws.Sign(encrypted, jwa.ES256, privatKey, jws.WithHeaders(hdr)))
	return jt, nil
}

func VerefyJWT(token string) (string, error) {
	verified, err := jws.Verify([]byte(token), jwa.ES256, publicKey)
	if err != nil {
		log.Printf("failed to verify message: %s", err)
		return "", err
	}
	claims, err := jwe.Decrypt(verified, jwa.ECDH_ES, privatKey)
	if err != nil {
		log.Printf("failed to decrypt: %s", err)
		return "", err
	}
	//fmt.Println(string(decrypted))
	return string(claims), nil
}
