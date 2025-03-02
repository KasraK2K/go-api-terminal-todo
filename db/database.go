package database

import (
	"context"
	"database/sql"
	_ "embed"
	"log"
	_ "modernc.org/sqlite"
	"todo/internal/repository"
	"todo/models"
)

//go:embed schemas/*.sql
var ddl string
var Database *models.Database

func newDatabase(dbFile string) (*models.Database, error) {
	ctx := context.Background()

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}

	// Create tables if they don’t exist
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return nil, err
	}

	return &models.Database{
		DB:      db,
		Queries: repository.New(db),
	}, nil
}

func InitDatabase() {
	var err error
	Database, err = newDatabase("./db/database.db")
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
}
