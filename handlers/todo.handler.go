package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	database "todo/db"
	"todo/internal/repository"
	"todo/models"
)

//UpdateTodo
//DeleteTodo

func GetTodo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req models.FindArgs
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Invalid JSON format: %v", err)
		SendJSONResponse(w, http.StatusBadRequest, nil, nil)
		return
	}

	if req.ID <= 0 {
		err = errors.New("invalid or missing ID")
		log.Printf("Invalid or missing ID: %v", req.ID)
		SendJSONResponse(w, http.StatusBadRequest, nil, nil)
		return
	}

	todo, err := database.Database.Queries.GetTodo(ctx, int64(req.ID))
	if err != nil {
		log.Printf("Error fetching todo with ID %d: %v", req.ID, err)
		SendJSONResponse(w, http.StatusNotFound, nil, err)
		return
	}

	sanitiseTodo := sanitiseResponse(todo)
	SendJSONResponse(w, http.StatusOK, sanitiseTodo, nil)
}

func ListTodos(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	todos, err := database.Database.Queries.ListTodos(ctx)
	if err != nil {
		log.Printf("Error fetching todos: %v", err)
		SendJSONResponse(w, http.StatusInternalServerError, nil, err)
		return
	}

	if todos == nil {
		todos = []repository.Todo{}
	}

	sanitisedTodo := sanitiseListResponse(todos)
	SendJSONResponse(w, http.StatusOK, sanitisedTodo, nil)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req repository.CreateTodoParams
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Invalid JSON format: %v", err)
		SendJSONResponse(w, http.StatusBadRequest, nil, nil)
		return
	}

	todo, err := database.Database.Queries.CreateTodo(ctx, req)
	if err != nil {
		log.Printf("Error creating todo: %v", err)
		SendJSONResponse(w, http.StatusNotFound, nil, err)
		return
	}

	sanitiseTodo := sanitiseResponse(todo)
	SendJSONResponse(w, http.StatusOK, sanitiseTodo, nil)
}

/* -------------------------------------------------------------------------------------------------- */
/*                                      Private Useful Functions                                      */
/* -------------------------------------------------------------------------------------------------- */
func sanitiseListResponse(todos []repository.Todo) []models.TodoResponse {
	responseTodos := make([]models.TodoResponse, len(todos))
	for i, todo := range todos {
		responseTodos[i] = models.TodoResponse{
			ID:          todo.ID,
			Title:       todo.Title,
			Description: todo.Description,
			Completed:   todo.Completed,
		}
	}
	return responseTodos
}

func sanitiseResponse(todo repository.Todo) models.TodoResponse {
	return models.TodoResponse{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		Completed:   todo.Completed,
	}
}
