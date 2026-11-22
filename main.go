package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	db "servidor-go/db/sqlc"
	"servidor-go/views"

	_ "github.com/lib/pq"
)

// handleIndex maneja GET / - Obtiene los libros con sqlc y renderiza la vista principal con Templ
func handleIndex(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		libros, err := queries.GetAllLibros(r.Context())
		if err != nil {
			log.Printf("Error al consultar libros: %v", err)
			http.Error(w, "Error al obtener libros de la base de datos", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := views.IndexPage(libros).Render(r.Context(), w); err != nil {
			log.Printf("Error al renderizar página: %v", err)
			http.Error(w, "Error al renderizar la vista", http.StatusInternalServerError)
		}
	}
}

// handleCreateLibro maneja POST /libros - Crea el libro en la BD y devuelve el fragmento actualizado de EntityList para HTMX
func handleCreateLibro(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
			return
		}

		titulo := strings.TrimSpace(r.FormValue("titulo"))
		autor := strings.TrimSpace(r.FormValue("autor"))
		fechaStr := strings.TrimSpace(r.FormValue("fecha_publicacion"))
		genero := strings.TrimSpace(r.FormValue("genero"))

		if titulo == "" || autor == "" || fechaStr == "" || genero == "" {
			http.Error(w, "Todos los campos (título, autor, fecha y género) son obligatorios", http.StatusBadRequest)
			return
		}

		fecha, err := time.Parse("2006-01-02", fechaStr)
		if err != nil {
			http.Error(w, "Formato de fecha inválido, debe ser YYYY-MM-DD", http.StatusBadRequest)
			return
		}

		params := db.CreateLibroParams{
			Titulo:           titulo,
			Autor:            autor,
			FechaPublicacion: fecha,
			Genero:           genero,
		}

		if _, err := queries.CreateLibro(r.Context(), params); err != nil {
			log.Printf("Error al crear libro en la BD: %v", err)
			http.Error(w, "Error al guardar el libro", http.StatusInternalServerError)
			return
		}

		libros, err := queries.GetAllLibros(r.Context())
		if err != nil {
			log.Printf("Error al consultar libros tras creación: %v", err)
			http.Error(w, "Error al obtener la lista actualizada de libros", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := views.EntityList(libros).Render(r.Context(), w); err != nil {
			log.Printf("Error al renderizar lista de libros: %v", err)
			http.Error(w, "Error al renderizar la lista actualizada", http.StatusInternalServerError)
		}
	}
}

// handleDeleteLibro maneja DELETE /libros/{id} - Elimina el libro en la BD y devuelve 200 OK vacío para HTMX
func handleDeleteLibro(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 32)
		if err != nil || id <= 0 {
			http.Error(w, "ID de libro inválido", http.StatusBadRequest)
			return
		}

		rowsAffected, err := queries.DeleteLibro(r.Context(), int32(id))
		if err != nil {
			log.Printf("Error al eliminar libro %d en la BD: %v", id, err)
			http.Error(w, "Error al eliminar el libro", http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			http.Error(w, "Libro no encontrado", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// RegisterRoutes inicializa y registra las rutas SSR y HTMX
func RegisterRoutes(mux *http.ServeMux, queries *db.Queries) {
	// Ruta principal GET /
	mux.HandleFunc("GET /{$}", handleIndex(queries))

	// Manejo del formulario POST /libros (devuelve fragmento HTML para HTMX)
	mux.HandleFunc("POST /libros", handleCreateLibro(queries))

	// Borrado de libro DELETE /libros/{id} (devuelve 200 OK vacío para HTMX)
	mux.HandleFunc("DELETE /libros/{id}", handleDeleteLibro(queries))

	// Redirección amigable si se accede a GET /libros
	mux.HandleFunc("GET /libros", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	// Servir archivos estáticos (por ejemplo styles.css o imágenes)
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))
}

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
	RegisterRoutes(mux, queries)

	log.Println("Conexión a la base de datos exitosa; servidor escuchando en :8080 (Modo SSR + HTMX con Templ)")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		return fmt.Errorf("error al iniciar el servidor HTTP: %w", err)
	}
	return nil
}
