package handlers

import (
	"net/http"

	db "servidor-go/db/sqlc"
)

// RegisterRoutes inicializa y registra las rutas SSR y HTMX
func RegisterRoutes(mux *http.ServeMux, queries *db.Queries) {
	// Ruta principal GET /
	mux.HandleFunc("GET /{$}", handleIndex(queries))

	// Manejo del formulario POST /libros (devuelve fragmento HTML para HTMX)
	mux.HandleFunc("POST /libros", handleCreateLibro(queries))

	// Edición inline de libros (devuelven la fila <tr> correspondiente para HTMX)
	mux.HandleFunc("GET /libros/{id}", handleGetLibroRow(queries))
	mux.HandleFunc("GET /libros/{id}/edit", handleEditLibroRow(queries))
	mux.HandleFunc("PUT /libros/{id}", handleUpdateLibro(queries))

	// Alternancia del estado leído / no leído
	mux.HandleFunc("PATCH /libros/{id}/leido", handleToggleLibroLeido(queries))

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
