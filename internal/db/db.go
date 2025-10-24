package db

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	_ "github.com/lib/pq"
)

var Queries *sqlc.Queries

// InitDB initializes the package-level database connection and the exported Queries object.
// It reads the DB_URL environment variable, opens and configures a PostgreSQL connection pool, verifies connectivity (5s timeout), assigns sqlc.New(dbConn) to Queries, and logs a fatal error if the environment variable is missing or the database cannot be reached.
func InitDB() {
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Configure pool
	dbConn.SetMaxOpenConns(25)
	dbConn.SetMaxIdleConns(5)
	dbConn.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := dbConn.PingContext(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	Queries = sqlc.New(dbConn)
	log.Println("Database initialized successfully")
}