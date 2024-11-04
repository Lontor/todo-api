package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Lontor/todo-api/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type contextKey string

func JWTAuth(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				next.ServeHTTP(w, r)
				return
			}

			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
			claims := jwt.MapClaims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(secretKey), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			userID, ok := claims["userID"].(string)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			role, ok := claims["role"].(string)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "userID", uuid.MustParse(userID))
			ctx = context.WithValue(ctx, "role", model.UserType(role))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
