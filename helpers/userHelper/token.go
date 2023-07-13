package userHelper

import (
	"frame/protos"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type tokenClaims struct {
	Username  string           `json:"username"`
	GroupType protos.GroupType `json:"group_type"`
	jwt.RegisteredClaims
}

func generateToken(username string, groupType protos.GroupType, key []byte) (string, error) {
	claims := tokenClaims{
		username,
		groupType,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(key)
	if err != nil {
		return "", err
	}
	return token, nil
}

func parseToken(token string, key []byte) (string, protos.GroupType, error) {
	claims := &tokenClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return "", protos.GroupTypeNil, err
	}
	return claims.Username, claims.GroupType, nil
}
