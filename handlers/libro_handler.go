package handlers

import (
	"log"
	"net/http"

	db "servidor-go/db/sqlc"
	"servidor-go/views"
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

		renderHTML(w, r, views.IndexPage(libros))
	}
}

// handleCreateLibro maneja POST /libros - Crea el libro en la BD y devuelve el fragmento actualizado de EntityList para HTMX
func handleCreateLibro(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form, errMsg := parseLibroForm(r)
		if errMsg != "" {
			http.Error(w, errMsg, http.StatusBadRequest)
			return
		}

		params := db.CreateLibroParams{
			Titulo:           form.Titulo,
			Autor:            form.Autor,
			FechaPublicacion: form.FechaPublicacion,
			Genero:           form.Genero,
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

		renderHTML(w, r, views.EntityList(libros))
	}
}

// handleDeleteLibro maneja DELETE /libros/{id} - Elimina el libro en la BD y devuelve 200 OK vacío para HTMX
func handleDeleteLibro(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := parseLibroID(r.PathValue("id"))
		if !ok {
			http.Error(w, "ID de libro inválido", http.StatusBadRequest)
			return
		}

		rowsAffected, err := queries.DeleteLibro(r.Context(), id)
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
