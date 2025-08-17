package jwt

import (
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var SecretKey = []byte("secret")

func GenerateToken(username string) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return tokenString, nil
}
