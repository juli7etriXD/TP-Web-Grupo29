package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"servidor-go/api"
	db "servidor-go/db/sqlc"
)

// RegisterRoutes inicializa todos los handlers y registra los endpoints en el mux
func RegisterRoutes(mux *http.ServeMux, queries *db.Queries) {
	bookHandler := api.NewBookHandler(queries)
	userHandler := api.NewUserHandler(queries)
	usuarioLibroHandler := api.NewUsuarioLibroHandler(queries)

	mux.HandleFunc("GET /libros", bookHandler.GetAll)
	mux.HandleFunc("GET /libros/{id}", bookHandler.GetByID)
	mux.HandleFunc("POST /libros", bookHandler.Create)
	mux.HandleFunc("PUT /libros/{id}", bookHandler.Update)
	mux.HandleFunc("DELETE /libros/{id}", bookHandler.Delete)

	mux.HandleFunc("GET /usuarios", userHandler.GetAll)
	mux.HandleFunc("GET /usuarios/{id}", userHandler.GetByID)
	mux.HandleFunc("POST /usuarios", userHandler.Create)
	mux.HandleFunc("PUT /usuarios/{id}", userHandler.Update)
	mux.HandleFunc("DELETE /usuarios/{id}", userHandler.Delete)

	mux.HandleFunc("POST /usuarios/{usuario_id}/libros", usuarioLibroHandler.Create)
	mux.HandleFunc("GET /usuarios/{usuario_id}/libros/leidos", usuarioLibroHandler.GetLibrosLeidosPorUsuario)
	mux.HandleFunc("GET /libros/{libro_id}/usuarios", usuarioLibroHandler.GetUsuariosPorLibro)
	mux.HandleFunc("PATCH /usuarios/{usuario_id}/libros/{libro_id}", usuarioLibroHandler.UpdateEstado)
	mux.HandleFunc("DELETE /usuarios/{usuario_id}/libros/{libro_id}", usuarioLibroHandler.Delete)

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/", fileServer)
}

// Función auxiliar para leer variables con valor por defecto
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	// 1. Leer variables individuales
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "contrasenia")
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_NAME", "biblioteca")

	// 2. Construir la cadena de conexión (DSN)
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbName)

	// 3. Conectar a PostgreSQL
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Error en sql.Open: %v", err)
	}
	defer conn.Close()

	if err := conn.Ping(); err != nil {
		log.Fatalf("Error al conectar con la BD: %v", err)
	}

	fmt.Println("¡Conexión a la base de datos exitosa!")
}
