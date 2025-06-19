package handlers

import (
	"TodoApp/database/dbhelper"
	"TodoApp/utils"
	"encoding/json"
	"net/http"
)

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.AuthHandler(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	user, err := dbhelper.GetUserByID(userID)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}
