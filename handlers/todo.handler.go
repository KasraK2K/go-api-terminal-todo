package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	database "todo/db"
	"todo/internal/repository"
	"todo/models"
)

func GetTodo(w http.ResponseWriter, r *http.Request) {
	var req models.FindArgs
	if err := decodeJSONBody(w, r, &req); err != nil {
		return
	}

	if req.ID <= 0 {
		log.Printf("Invalid ID: %v", req.ID)
		SendJSONResponse(w, http.StatusBadRequest, nil, errors.New("invalid or missing ID"))
		return
	}

	executeQuery(w, func() (models.TodoResponse, error) {
		todo, err := database.Database.Queries.GetTodo(r.Context(), int64(req.ID))
		return sanitiseTodo(todo), err
	})
}

func ListTodos(w http.ResponseWriter, r *http.Request) {
	executeQuery(w, func() ([]models.TodoResponse, error) {
		todos, err := database.Database.Queries.ListTodos(r.Context())
		if err != nil {
			return nil, err
		}
		return sanitiseTodos(todos), nil
	})
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req repository.CreateTodoParams
	if err := decodeJSONBody(w, r, &req); err != nil {
		return
	}

	executeQuery(w, func() (models.TodoResponse, error) {
		todo, err := database.Database.Queries.CreateTodo(r.Context(), req)
		return sanitiseTodo(todo), err
	})
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	var req repository.UpdateTodoParams
	if err := decodeJSONBody(w, r, &req); err != nil {
		return
	}

	existingTodo, err := database.Database.Queries.GetTodo(r.Context(), req.ID)
	if err != nil {
		log.Printf("Error fetching existing todo: %v", err)
		SendJSONResponse(w, http.StatusNotFound, nil, err)
		return
	}

	if req.Title == "" {
		req.Title = existingTodo.Title
	}
	if req.Description == "" {
		req.Description = existingTodo.Description
	}

	executeQuery(w, func() (models.TodoResponse, error) {
		todo, err := database.Database.Queries.UpdateTodo(r.Context(), req)
		return sanitiseTodo(todo), err
	})
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	var req models.FindArgs
	if err := decodeJSONBody(w, r, &req); err != nil {
		return
	}

	if req.ID <= 0 {
		log.Printf("Invalid ID: %v", req.ID)
		SendJSONResponse(w, http.StatusBadRequest, nil, errors.New("invalid or missing ID"))
		return
	}

	executeQuery(w, func() (string, error) {
		err := database.Database.Queries.DeleteTodo(r.Context(), int64(req.ID))
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Todo with id %d successfully deleted", req.ID), nil
	})
}

/* -------------------------------------------------------------------------------------------------- */
/*                                      Private Useful Functions                                      */
/* -------------------------------------------------------------------------------------------------- */
func sanitiseTodos(todos []repository.Todo) []models.TodoResponse {
	response := make([]models.TodoResponse, len(todos))
	for i, todo := range todos {
		response[i] = sanitiseTodo(todo)
	}
	return response
}

func sanitiseTodo(todo repository.Todo) models.TodoResponse {
	return models.TodoResponse{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		Completed:   todo.Completed,
	}
}
