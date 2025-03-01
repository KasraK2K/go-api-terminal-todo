package handlers

import (
	"encoding/json"
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

	// Encode the response and handle possible errors
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
