package tests

import (
	"database/sql"
	"fmt"
	"testing"

	db "servidor-go/db/sqlc"

	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) (*sql.DB, *db.Queries) {
	user := "user"
	pass := "pass"
	dbName := "database"
	port := "5432"

	connStr := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", user, pass, port, dbName)

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
