package rsa

import (
	"crypto/rsa"
	"io/ioutil"

	jwt "github.com/dgrijalva/jwt-go"
)

//RSA -
type RSA struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

//NewRSA -
func NewRSA(prikey, pubkey string) (*RSA, error) {
	vbytes, err := ioutil.ReadFile(prikey)
	if err != nil {
		return nil, err
	}
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(vbytes)
	if err != nil {
		return nil, err
	}
	bbytes, err := ioutil.ReadFile(pubkey)
	if err != nil {
		return nil, err
	}
	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(bbytes)
	if err != nil {
		return nil, err
	}

	return &RSA{
		privateKey: signKey,
		publicKey:  verifyKey,
	}, nil
}

//GetPrivateKey -
func (s *RSA) GetPrivateKey() *rsa.PrivateKey {
	return s.privateKey
}

//GetPublicKey -
func (s *RSA) GetPublicKey() *rsa.PublicKey {
	return s.publicKey
}
