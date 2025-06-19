package utils

import (
	"TodoApp/database/dbhelper"
	"errors"
	"net/http"
)

// authentication
func AuthHandler(r *http.Request) (int, error) {
	sessionID := r.Header.Get("Authorization")
	if sessionID == "" {
		return 0, errors.New("authorization header required")
	}
	userID, err := dbhelper.GetUserBySessionKey(sessionID)
	if err != nil {
		return 0, errors.New("invalid session")
	}
	return userID, nil
}
