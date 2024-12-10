package jwt

import (
	"CheckUser/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// NewToken creates new JWT token for given user and app.
func NewToken(userName string, password string, app models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["userName"] = userName
	claims["password"] = password
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app_id"] = app.ID

	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
