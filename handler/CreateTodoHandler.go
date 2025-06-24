package handlers

import (
	"TodoApp/database/dbhelper"
	"TodoApp/middleware"
	"TodoApp/models"
	"encoding/json"
	"net/http"
)

func CreateTodoHandler(w http.ResponseWriter, r *http.Request) {

	userID, err := middleware.AuthUserFromMiddleWare(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var todoObj struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&todoObj); err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	todo := models.Todo{
		UserID:      userID,
		Title:       todoObj.Title,
		Description: todoObj.Description,
	}

	if err := dbhelper.CreateTodo(todo); err != nil {
		http.Error(w, "Failed to create todo record", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Todo created successfully",
	})
}
