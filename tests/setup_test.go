package tests

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	db "servidor-go/db/sqlc"

	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) (*sql.DB, *db.Queries) {
	getEnv := func(key, defaultValue string) string {
		if value := os.Getenv(key); value != "" {
			return value
		}
		return defaultValue
	}

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "contrasenia"),
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_NAME", "biblioteca"),
	)

	dbConn, err := sql.Open("postgres", connStr)

	if err != nil {
		t.Fatalf("Error al conectar con la base de datos:%v", err)
	}

	if err := dbConn.Ping(); err != nil {
		dbConn.Close()
		t.Fatalf("La base de datos no está respondiendo:%v", err)
	}
	return dbConn, db.New(dbConn)
}
