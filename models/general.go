package models

import (
	"database/sql"
	"todo/internal/repository"
)

type StatusCode int

type Response struct {
	Status int
	Error  *string
	Data   interface{}
}

type Database struct {
	DB      *sql.DB
	Queries *repository.Queries
}
