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
var Connection *sql.DB

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
	Connection = dbConn
	log.Println("Database initialized successfully")
}
