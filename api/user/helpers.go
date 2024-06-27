package user

import (
	"bkoiki950/go-store/api/config"
	"bkoiki950/go-store/api/database"

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