package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/malanak2/nextap-chat/domain"
)

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		// Parse token with claims into JwtClaims type
		token, err := jwt.ParseWithClaims(tokenString, &domain.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token ("+err.Error()+")", http.StatusUnauthorized)
			return
		}
		// Sends UserId to the handler func
		ctx := context.WithValue(r.Context(), "userId", token.Claims.(*domain.JwtClaims).UserId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
