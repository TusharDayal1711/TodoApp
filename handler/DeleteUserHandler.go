package handlers

import (
	"TodoApp/database/dbhelper"
	"TodoApp/middleware"
	"encoding/json"
	"net/http"
)

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.AuthUserFromMiddleWare(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if err := dbhelper.DeleteUserByID(userID); err != nil {
		http.Error(w, "failed to delete user...", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Account Deleted....",
	})
}
