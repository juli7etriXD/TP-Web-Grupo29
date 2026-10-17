package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	db "servidor-go/db/sqlc"
)

// BookHandler encapsula las operaciones HTTP sobre la entidad Libro
type BookHandler struct {
	queries *db.Queries
}

func NewBookHandler(q *db.Queries) *BookHandler {
	return &BookHandler{queries: q}
}

// BookRequest representa el cuerpo JSON esperado al crear o actualizar un libro
type BookRequest struct {
	Titulo           string `json:"titulo"`
	Autor            string `json:"autor"`
	FechaPublicacion string `json:"fecha_publicacion"` // Formato esperado: "YYYY-MM-DD"
	Genero           string `json:"genero"`
}

// GET /libros - Listar todos los libros
func (h *BookHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	libros, err := h.queries.GetAllLibros(r.Context())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al obtener libros"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(libros)
}

// GET /libros/{id} - Obtener un libro por ID
func (h *BookHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "ID inválido"})
		return
	}

	libro, err := h.queries.GetLibroByID(r.Context(), int32(id))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Libro no encontrado"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error interno del servidor"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(libro)
}

// POST /libros - Crear un libro nuevo
func (h *BookHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req BookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Cuerpo de solicitud inválido"})
		return
	}

	fecha, err := time.Parse("2006-01-02", req.FechaPublicacion)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Formato de fecha inválido, use YYYY-MM-DD"})
		return
	}

	params := db.CreateLibroParams{
		Titulo:           req.Titulo,
		Autor:            req.Autor,
		FechaPublicacion: fecha,
		Genero:           req.Genero,
	}

	nuevoLibro, err := h.queries.CreateLibro(r.Context(), params)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al crear libro"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nuevoLibro)
}

// PUT /libros/{id} - Actualizar un libro existente
func (h *BookHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "ID inválido"})
		return
	}

	var req BookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Cuerpo de solicitud inválido"})
		return
	}

	fecha, err := time.Parse("2006-01-02", req.FechaPublicacion)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Formato de fecha inválido, use YYYY-MM-DD"})
		return
	}

	params := db.UpdateLibroParams{
		ID:               int32(id),
		Titulo:           req.Titulo,
		Autor:            req.Autor,
		FechaPublicacion: fecha,
		Genero:           req.Genero,
	}

	err = h.queries.UpdateLibro(r.Context(), params)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al actualizar libro"})
		return
	}

	// Retornamos el libro actualizado o mensaje de éxito
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Libro actualizado correctamente"})
}

// DELETE /libros/{id} - Eliminar un libro
func (h *BookHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "ID inválido"})
		return
	}

	err = h.queries.DeleteLibro(r.Context(), int32(id))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al eliminar libro"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
