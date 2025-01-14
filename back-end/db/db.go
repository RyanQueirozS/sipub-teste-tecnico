// This file handles the db and context with a singleton approach, perhaps not
// the best with clean-architecture, but it gives nice results fast. Ideally
// the db and ctx should be injected through a dependency of each repository,
// and/or usecase but that would introduce more boilerplate and even more
// complexity to an already somewhat complex architecture.

package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db  *sql.DB
	ctx context.Context
)

func InitializeDB(dsn string) error {
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("Database connection established")
	return nil
}

func GetDB() *sql.DB {
	return db
}

func CloseDB() {
	if err := db.Close(); err != nil {
		log.Printf("Failed to close database: %v", err)
	}
}

// https://pkg.go.dev/context#Context
func InitializeContext() {
	ctx = context.Background()
}

func GetContext() context.Context {
	return ctx
}
