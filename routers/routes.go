package routers

import (
	handlers "TodoApp/handler"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func GetRoutes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "connection established...")
	}).Methods("GET")

	router.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	router.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	router.HandleFunc("/fetch-user-info", handlers.GetUserHandler).Methods("GET")
	router.HandleFunc("/logout", handlers.LogoutHandler).Methods("POST")
	router.HandleFunc("/delete-profile", handlers.DeleteUser).Methods("DELETE")

	router.HandleFunc("/create-todos", handlers.CreateTodoHandler).Methods("POST")
	router.HandleFunc("/mark-todo-complete", handlers.MarkTodoCompletedHandler).Methods("POST")
	router.HandleFunc("/fetch-todos-info", handlers.GetTodoByStatus).Methods("GET")
	router.HandleFunc("/todos/update", handlers.UpdateTodoHandler).Methods("PUT")
	router.HandleFunc("/todos/delete", handlers.DeleteTodoRecord).Methods("DELETE")

	return router
}
