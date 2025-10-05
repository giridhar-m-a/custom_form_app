package db

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// Global variables
var Queries *sqlc.Queries
var Ctx context.Context

// InitDB initializes the database connection and sqlc queries
func InitDB() {
	ctx := context.Background()

	// Replace with your DB URL or read from env
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Optional: ping DB to ensure connection is working
	if err := dbConn.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	Queries = sqlc.New(dbConn)
	Ctx = ctx

	log.Println("Database initialized successfully")
}
