package handlers

import (
	"TodoApp/database/dbhelper"
	"TodoApp/middleware"
	"TodoApp/models"
	"TodoApp/utils"
	"encoding/json"
	"net/http"
)

func UpdateTodoHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.AuthUserFromMiddleWare(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var todoObj struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		IsCompleted bool   `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&todoObj); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	todo := models.Todo{
		ID:          todoObj.ID,
		UserID:      userID,
		Title:       todoObj.Title,
		Description: todoObj.Description,
		Status:      todoObj.IsCompleted,
	}
	if err := dbhelper.UpdateTodoDetails(todo); err != nil {
		if err.Error() == "invalid todo id" {
			http.Error(w, "Invalid todo ID", http.StatusNotFound)
			return
		}
		http.Error(w, "could not update todo", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Todo updated successfully",
	})
}

// Marking toggle task
func MarkTodoCompletedHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.AuthHandler(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var input struct {
		TodoID int `json:"todo_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if err := dbhelper.MarkCompleted(input.TodoID, userID); err != nil {
		http.Error(w, "failed to mark todo as completed", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Todo marked as completed",
	})
}
