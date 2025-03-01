package routes

import (
	"net/http"

	"todo/handlers"
)

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	apiMuxPrefix := http.NewServeMux()
	apiMuxPrefix.HandleFunc("GET /todos", handlers.ListTodos)
	apiMuxPrefix.HandleFunc("POST /todos", handlers.GetTodo)
	apiMuxPrefix.HandleFunc("POST /todos/new", handlers.CreateTodo)
	apiMuxPrefix.HandleFunc("PATCH /todos", handlers.UpdateTodo)
	apiMuxPrefix.HandleFunc("DELETE /todos", handlers.DeleteTodo)

	mux.Handle("/api/", http.StripPrefix("/api", apiMuxPrefix))
	return mux
}
