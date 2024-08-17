package user

import (
	"github.com/Babatunde13/golangr/api/config"
	"github.com/Babatunde13/golangr/api/database"

	"github.com/golang-jwt/jwt"
)

var SECRET_KEY string = config.GetEnv("SECRET_KEY")

func getToken (user database.User) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
	})

	tokenString, err := token.SignedString([]byte(SECRET_KEY)); if err != nil {
		return ""
	}

	return tokenString
}

func verifyToken (tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	}); if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims); if !ok || !token.Valid {
		return "", err
	}

	email := claims["email"].(string)
	return email, nil
}
