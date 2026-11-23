package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	db "servidor-go/db/sqlc"
	handlers "servidor-go/handlers"

	_ "github.com/lib/pq"
)

// Función auxiliar para leer variables con valor por defecto
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "contrasenia")
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_NAME", "biblioteca")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbName)

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("error en sql.Open: %w", err)
	}
	defer conn.Close()

	if err := conn.Ping(); err != nil {
		return fmt.Errorf("error al conectar con la BD: %w", err)
	}

	queries := db.New(conn)
	mux := http.NewServeMux()
	handlers.RegisterRoutes(mux, queries)

	log.Println("Conexión a la base de datos exitosa; servidor escuchando en :8080 (Modo SSR + HTMX con Templ)")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		return fmt.Errorf("error al iniciar el servidor HTTP: %w", err)
	}
	return nil
}
