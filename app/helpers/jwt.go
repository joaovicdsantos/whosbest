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

	return fmt.Sprintf("%s", tokenString), nil
}

func ParseJwtToken(unparsedToken string) (map[string]interface{}, error) {

	splitedToken := strings.Split(unparsedToken, " ")
	if splitedToken[0] != prefix {
		return nil, fmt.Errorf("Invalid token")
	}

	unparsedToken = splitedToken[1]
	token, err := jwt.Parse(unparsedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
