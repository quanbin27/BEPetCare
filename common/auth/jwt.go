package auth

import (
	"github.com/golang-jwt/jwt"
	"strconv"
	"time"
)

func CreateJWT(secret []byte, userID int32, seconds int64) (string, error) {
	expiration := time.Second * time.Duration(seconds)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   strconv.Itoa(int(userID)),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
