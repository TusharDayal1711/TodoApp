package handlers

import (
	"TodoApp/database/dbhelper"
	"TodoApp/middleware"
	"encoding/json"
	"net/http"
)

func GetTodoByStatus(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.AuthUserFromMiddleWare(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	status := r.URL.Query().Get("status")
	if status == "" {
		http.Error(w, "query parameter status is required", http.StatusBadRequest)
		return
	}
	todos, err := dbhelper.GetTodosStatus(userID, status)
	if err != nil {
		http.Error(w, "could not fetch todo data", http.StatusInternalServerError)
		return
	}
	res, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		http.Error(w, "encoding failed", http.StatusInternalServerError)
		return
	}
	w.Write(res)
}

//func GetCompletedTask(w http.ResponseWriter, r *http.Request) {
//	sessionID := r.Header.Get("Authorization")
//	if sessionID == "" {
//		http.Error(w, "Authorization header required", http.StatusUnauthorized)
//		return
//	}
//
//	userID, err := dbhelper.GetUserBySessionKey(sessionID)
//	if err != nil {
//		http.Error(w, "Invalid session...", http.StatusUnauthorized)
//		return
//	}
//
//	todos, err := dbhelper.GetCompletedTodosByUserID(userID)
//	if err != nil {
//		http.Error(w, "Could not fetch completed todos", http.StatusInternalServerError)
//		return
//	}
//
//	json.NewEncoder(w).Encode(todos)
//}
//
//func GetInCompletedTask(w http.ResponseWriter, r *http.Request) {
//	sessionID := r.Header.Get("Authorization")
//	if sessionID == "" {
//		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
//		return
//	}
//
//	userID, err := dbhelper.GetUserBySessionKey(sessionID)
//	if err != nil {
//		http.Error(w, "Invalid session", http.StatusUnauthorized)
//		return
//	}
//
//	todos, err := dbhelper.GetInCompletedTodosByUserID(userID)
//	if err != nil {
//		http.Error(w, "Could not fetch completed todos", http.StatusInternalServerError)
//		return
//	}
//	json.NewEncoder(w).Encode(todos)
//}

//func GetUserTodosHandler(w http.ResponseWriter, r *http.Request) {
//	sessionId := r.Header.Get("Authorization")
//	if sessionId == "" {
//		http.Error(w, "Authorization header is required", http.StatusUnauthorized)
//		return
//	}
//	userID, err := dbhelper.GetUserBySessionKey(sessionId)
//	if err != nil {
//		http.Error(w, "Invalid session", http.StatusUnauthorized)
//		return
//	}
//	todos, err := dbhelper.GetTodosByID(userID)
//	if err != nil {
//		http.Error(w, "Failed to fetch todos", http.StatusInternalServerError)
//		return
//	}
//	w.Header().Set("Content-Type", "application/json")
//	json.NewEncoder(w).Encode(todos)
//}
