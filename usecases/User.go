package usecases

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/malanak2/nextap-chat/domain"
	"github.com/malanak2/nextap-chat/ports"
)

func CreateUser(name, pass string) (domain.User, error) {
	user, err := ports.CreateUser(name, pass)
	if err != nil {
		return domain.User{}, ErrorBadParameters
	}
	return user, nil
}

func LoginUser(name, pass string) (string, error) {
	user, err := ports.UserLogin(name, pass)
	if err != nil {
		return "", err
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// Expires in about a month
		"exp":      jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
		"username": user.Username,
		"userId":   user.ID,
	})
	s, _ := t.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return s, nil
}
