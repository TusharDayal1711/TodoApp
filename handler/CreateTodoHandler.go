package handlers

import (
	"TodoApp/database/dbhelper"
	"TodoApp/models"
	"TodoApp/utils"
	"encoding/json"
	"net/http"
)

func CreateTodoHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.AuthHandler(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	var todoObj struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&todoObj); err != nil {
		http.Error(w, "invalid JSON input", http.StatusBadRequest)
		return
	}
	todo := models.Todo{
		UserID:      userID,
		Title:       todoObj.Title,
		Description: todoObj.Description,
	}
	if err := dbhelper.CreateTodo(todo); err != nil {
		http.Error(w, "failed to create todo record", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Todo created successfully",
	})
}
