package routers

import (
	handlers "TodoApp/handler"
	"TodoApp/middleware"
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

	//router.HandleFunc("/refresh-token", middleware.RefreshTokenHandler).Methods("GET")

	protected := router.PathPrefix("/").Subrouter()
	protected.Use(middleware.JWTAuthMiddleware)

	//user routes
	protected.HandleFunc("/fetch-user-info", handlers.GetUserHandler).Methods("GET")
	protected.HandleFunc("/logout", handlers.LogoutHandler).Methods("POST")
	protected.HandleFunc("/delete-profile", handlers.DeleteUser).Methods("DELETE")

	//routes todos
	protected.HandleFunc("/create-todos", handlers.CreateTodoHandler).Methods("POST")
	protected.HandleFunc("/mark-todo-complete", handlers.MarkTodoCompletedHandler).Methods("POST")
	protected.HandleFunc("/fetch-todos-info", handlers.GetTodoByStatus).Methods("GET")
	protected.HandleFunc("/todos/update", handlers.UpdateTodoHandler).Methods("PUT")
	protected.HandleFunc("/todos/delete", handlers.DeleteTodoRecord).Methods("DELETE")

	return router
}
