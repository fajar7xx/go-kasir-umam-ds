package database

import (
	"database/sql"
	"fmt"
	"log"

	// _ "github.com/lib/pq"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func InitDB(connectionString string) (*sql.DB, error) {
	// open database
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		fmt.Printf("sql.open: %v", err)
		return nil, err
	}

	// connection test
	err = db.Ping()
	if err != nil {
		fmt.Printf("db.Ping: %v", err)
		fmt.Println("error ping")
		return nil, err
	}

	// set connection pool settings (optional but recommended)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("Database connected successfully")
	return db, nil
}
