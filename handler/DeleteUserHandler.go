package handlers

import (
	"TodoApp/database/dbhelper"
	"TodoApp/utils"
	"encoding/json"
	"net/http"
)

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.AuthHandler(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
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
