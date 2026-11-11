package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
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

// handleCreateLibro maneja POST /libros - Parsea el formulario tradicional, guarda en BD y redirige (PRG)
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

		// Patrón Post/Redirect/Get (PRG): Redirige a GET / para evitar reenvíos al refrescar
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// RegisterRoutes inicializa y registra las rutas SSR
func RegisterRoutes(mux *http.ServeMux, queries *db.Queries) {
	// Ruta principal GET /
	mux.HandleFunc("GET /{$}", handleIndex(queries))

	// Manejo del formulario POST /libros
	mux.HandleFunc("POST /libros", handleCreateLibro(queries))

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

	log.Println("Conexión a la base de datos exitosa; servidor escuchando en :8080 (Modo SSR con Templ)")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		return fmt.Errorf("error al iniciar el servidor HTTP: %w", err)
	}
	return nil
}
