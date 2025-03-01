package models

import (
	"database/sql"
	"todo/internal/repository"
)

type StatusCode int

type Response struct {
	Status int         `json:"status"`
	Error  *string     `json:"error"`
	Data   interface{} `json:"data"`
}

type Database struct {
	DB      *sql.DB
	Queries *repository.Queries
}

type FindArgs struct {
	ID int `json:"id"`
}
