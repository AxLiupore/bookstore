package utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtSecret = []byte("yijiansanlian")

type Claims struct {
	Id        uint   `json:"id"`
	Username  string `json:"username"`
	Authority int    `json:"authority"`
	jwt.StandardClaims
}

// GenerateToken 签发 token
func GenerateToken(id uint, username string, authority int) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(24 * time.Hour)
	claims := Claims{
		Id:        id,
		Username:  username,
		Authority: authority,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "FanOne-Mall",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}
