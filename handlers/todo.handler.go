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

	todo, err := database.Database.Queries.GetTodo(r.Context(), int64(req.ID))
	handleResponse[repository.Todo](w, todo, err)
}

func ListTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := database.Database.Queries.ListTodos(r.Context())
	handleResponse[[]repository.Todo](w, todos, err)

}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req repository.CreateTodoParams
	if err := decodeJSONBody(w, r, &req); err != nil {
		return
	}

	todo, err := database.Database.Queries.CreateTodo(r.Context(), req)
	handleResponse[repository.Todo](w, todo, err)
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

	todo, err := database.Database.Queries.UpdateTodo(r.Context(), req)
	handleResponse[repository.Todo](w, todo, err)
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

	_, err := database.Database.Queries.GetTodo(r.Context(), int64(req.ID))
	if err != nil {
		SendJSONResponse(w, http.StatusNotFound, nil, err)
		return
	}

	err = database.Database.Queries.DeleteTodo(r.Context(), int64(req.ID))
	message := fmt.Sprintf("Todo with id %d successfully deleted", req.ID)
	handleResponse[string](w, message, err)
}
