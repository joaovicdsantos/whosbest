package helpers

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joaovicdsantos/whosbest-api/app/models"
)

var (
	jwtSecret = os.Getenv("JWT_SECRET")
	prefix    = "Bearer"
)

func CreateJwtToken(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseJwtToken(unparsedToken string) (map[string]interface{}, error) {

	splitedToken := strings.Split(unparsedToken, " ")
	if splitedToken[0] != prefix || len(splitedToken) != 2 {
		return nil, fmt.Errorf("invalid token")
	}

	unparsedToken = splitedToken[1]
	token, err := jwt.Parse(unparsedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
