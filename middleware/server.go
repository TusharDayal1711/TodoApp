package middleware

import (
	"TodoApp/database/dbhelper"
	"context"
	"errors"
	"net/http"
)

type contextKey string

const userContextKey contextKey = "user_key"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.Header.Get("Authorization")
		if sessionID == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		userID, err := dbhelper.GetUserBySessionKey(sessionID)
		if err != nil {
			http.Error(w, "Invalid session", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), userContextKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AuthUserFromMiddleWare(r *http.Request) (int, error) {
	userID, ok := r.Context().Value(userContextKey).(int)
	if !ok {
		return 0, errors.New("could not found user id")
	}
	return userID, nil
}
