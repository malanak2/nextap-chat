package domain

import "github.com/golang-jwt/jwt/v5"

type CreateUser struct {
	Username string
	Password string
}

type UserLogin struct {
	Username string
	Password string
}

type SendMessage struct {
	Channel int
	Content string
}

type JwtClaims struct {
	UserId   int    `json:"userId"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type ChangeUsername struct {
	Username string
}
