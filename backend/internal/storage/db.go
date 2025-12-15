package storage

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

type DB struct {
	*sql.DB
}

func NewDB(driver, connectionString string) (*DB, error) {
	db, err := sql.Open(driver, connectionString)
	if err != nil {
		log.Printf("STORAGE: Failed to open database connection: %v", err)
		return nil, errors.New("failed to open database connection")
	}

	if err = db.Ping(); err != nil {
		log.Printf("STORAGE: Failed to connect to database: %v", err)
		return nil, errors.New("failed to connect to database")
	}

	log.Printf("STORAGE: Database connected successfully (%s)", driver)
	return &DB{db}, nil
}
