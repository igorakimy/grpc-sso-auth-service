package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/igorakimy/grpc-sso-auth-service/internal/models"
	"log"
	"time"
)

func NewJWTToken(user models.User, app models.App, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{}
	claims["uid"] = user.ID
	claims["email"] = user.Email
	log.Println(duration)
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app_id"] = app.ID

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
