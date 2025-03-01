package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"todo/models"
)

func SendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}, err error) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	var errorMessage *string
	if err != nil {
		msg := err.Error()
		errorMessage = &msg
	}

	response := models.Response{
		Status: statusCode,
		Error:  errorMessage,
		Data:   data,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func decodeJSONBody[T any](w http.ResponseWriter, r *http.Request, reference *T) error {
	err := json.NewDecoder(r.Body).Decode(reference)
	if err != nil {
		log.Printf("Invalid JSON format: %v", err)
		SendJSONResponse(w, http.StatusBadRequest, nil, errors.New("invalid JSON format"))
		return err
	}
	return nil
}

func executeQuery[T any](w http.ResponseWriter, queryFunc func() (T, error)) {
	result, err := queryFunc()
	if err != nil {
		log.Printf("Database error: %v", err)
		SendJSONResponse(w, http.StatusNotFound, nil, err)
		return
	}

	SendJSONResponse(w, http.StatusOK, result, nil)
}
