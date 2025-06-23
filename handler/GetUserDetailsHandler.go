package handlers

import (
	"TodoApp/database/dbhelper"
	"TodoApp/middleware"
	"encoding/json"
	"net/http"
)

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	//userID, err := utils.AuthHandler(r)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusUnauthorized)
	//	return
	//}

	userID, err := middleware.AuthUserFromMiddleWare(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := dbhelper.GetUserByID(userID)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}
