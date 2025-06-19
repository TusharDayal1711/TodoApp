package handlers

import (
	"TodoApp/database/dbhelper"
	"TodoApp/utils"

	"encoding/json"
	"net/http"
)

func DeleteTodoRecord(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.AuthHandler(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var todoObj struct {
		ItemId int `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&todoObj); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err = dbhelper.DeleteTodoByID(userID, todoObj.ItemId)
	if err != nil {
		http.Error(w, "Could not delete todo", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Todo deleted successfully",
	})
}
