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
	UserId   int
	Username string
	jwt.RegisteredClaims
}

type ChangeUsername struct {
	Username string
}

type EditMessage struct {
	Content string
}

type User struct {
	ID       int
	Username string
}

type Message struct {
	ID      int
	Content string
}

type MessageWithAuthor struct {
	Message
	Author User
}
