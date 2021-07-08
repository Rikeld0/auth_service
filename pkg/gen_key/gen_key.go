package gen_key

import (
	"auth/service/structs"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/lestrrat-go/jwx/jwk"
)

func encode(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) ([]byte, []byte) {
	x509Encoded, _ := x509.MarshalECPrivateKey(privateKey)
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: x509Encoded})

	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(publicKey)
	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})

	return pemEncoded, pemEncodedPub
}

func GenEcdsaKey() ([]byte, []byte) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil
	}
	pubKey := &privKey.PublicKey
	return encode(privKey, pubKey)
}

func GenUserKey(id string, privE, pubE []byte) (*structs.UserKey, error) {
	privatKey, err := jwk.ParseKey(privE, jwk.WithPEM(true))
	if err != nil {
		fmt.Printf("failed to create JWK: %s\n", err)
		return nil, err
	}
	_ = privatKey.Set(jwk.KeyIDKey, id)
	publicKey, err := jwk.ParseKey(pubE, jwk.WithPEM(true))
	if err != nil {
		fmt.Printf("failed to create JWK: %s\n", err)
		return nil, err
	}
	_ = publicKey.Set(jwk.KeyIDKey, id)
	return &structs.UserKey{PrivateKey: privatKey, PublicKey: publicKey}, nil
}
