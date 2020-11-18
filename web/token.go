package web

import (
	"time"

	"ditto.co.jp/awsman/config"
	"ditto.co.jp/awsman/utils"
	"github.com/dgrijalva/jwt-go"
)

func (s *Service) generateToken(name string, isadmin bool) (*Token, error) {
	conf, _ := config.Load()

	tm := time.Duration(conf.Account.Expires)
	rftm := time.Duration(tm + 6)
	//create token
	token := jwt.New(jwt.SigningMethodRS512)

	//set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = 1
	claims["name"] = name
	claims["admin"] = isadmin
	claims["exp"] = utils.NowJST().Add(time.Hour * tm).Unix()

	//generate encoded token
	t, err := token.SignedString(s._rsa.GetPrivateKey())
	if err != nil {
		return nil, err
	}

	expires := int64(time.Hour * rftm / time.Millisecond)
	//refresh token
	refreshtoken := jwt.New(jwt.SigningMethodRS512)
	rclaims := refreshtoken.Claims.(jwt.MapClaims)
	rclaims["sub"] = 1
	rclaims["exp"] = utils.NowJST().Add(time.Hour * rftm).Unix()

	//generate encoded token
	rt, err := token.SignedString(s._rsa.GetPrivateKey())
	if err != nil {
		return nil, err
	}

	return &Token{
		AccessToken:  t,
		RefreshToken: rt,
		TokenType:    "bearer",
		Expires:      expires,
	}, nil

}
