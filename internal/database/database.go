package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

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
	// Optimized: Set MaxIdleConns equal to MaxOpenConns to minimize connection churn
	// and improve performance under high concurrency.
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	// Set MaxLifetime to ensure connections are refreshed and not closed by firewalls unexpectedly.
	db.SetConnMaxLifetime(time.Hour)

	log.Println("Database connected successfully")
	return db, nil
}
