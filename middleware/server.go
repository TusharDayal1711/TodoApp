package middleware

import (
	"TodoApp/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type contextKey string

const userContextKey contextKey = "user_key"

//	func AuthMiddleware(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			sessionID := r.Header.Get("Authorization")
//			if sessionID == "" {
//				http.Error(w, "Authorization header required", http.StatusUnauthorized)
//				return
//			}
//
//			userID, err := dbhelper.GetUserBySessionKey(sessionID)
//			if err != nil {
//				http.Error(w, "Invalid session", http.StatusUnauthorized)
//				return
//			}
//			ctx := context.WithValue(r.Context(), userContextKey, userID)
//			next.ServeHTTP(w, r.WithContext(ctx))
//		})
//	}
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get("Authorization")
		if accessToken == "" {
			http.Error(w, "missing access token", http.StatusUnauthorized)
			return
		}

		userID, err := utils.ParseJWT(accessToken)
		if err != nil {
			if strings.Contains(err.Error(), "invalid or expired token") {
				refreshToken := r.Header.Get("refresh_token")
				if refreshToken == "" {
					http.Error(w, "access token expired, and refresh token missing", http.StatusUnauthorized)
					return
				}

				userID, err = utils.ParseRefreshToken(refreshToken)
				if err != nil {
					http.Error(w, "invalid or expired refresh token", http.StatusUnauthorized)
					return
				}

				valid, err := utils.IsRefreshTokenValid(userID, refreshToken)
				if err != nil || !valid {
					http.Error(w, "refresh token invalid", http.StatusUnauthorized)
					return
				}
				fmt.Println("Refresh token valid...creating new access token...")
				newAccessToken, err := utils.GenerateJWT(userID)
				if err != nil {
					http.Error(w, "failed to generate new access token", http.StatusInternalServerError)
					return
				}
				fmt.Println("new access token generated ::", newAccessToken)
				w.Header().Set("new access token", newAccessToken)
			} else {
				http.Error(w, "unauthorized: "+err.Error(), http.StatusUnauthorized)
				return
			}
		}

		ctx := context.WithValue(r.Context(), userContextKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AuthUserFromMiddleWare(r *http.Request) (int, error) {
	userID, ok := r.Context().Value(userContextKey).(int)
	if !ok {
		return 0, errors.New("user ID not found in context")
	}
	return userID, nil
}

func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.Header.Get("Refresh-Token")
	if refreshToken == "" {
		http.Error(w, "refresh token required", http.StatusUnauthorized)
		return
	}

	userID, err := utils.ParseRefreshToken(refreshToken)
	if err != nil {
		http.Error(w, "invalid or expired refresh token", http.StatusUnauthorized)
		return
	}
	newAccessToken, err := utils.GenerateJWT(userID)
	if err != nil {
		http.Error(w, "could not generate new token", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": newAccessToken,
	})
}
