package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
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
		responderError(w, http.StatusInternalServerError, "Error al obtener libros")
		return
	}
	if libros == nil {
		libros = []db.Libro{}
	}

	responderJSON(w, http.StatusOK, libros)
}

// GET /libros/{id} - Obtener un libro por ID
func (h *BookHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := idPositivo(r.PathValue("id"))
	if !ok {
		responderError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	libro, err := h.queries.GetLibroByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			responderError(w, http.StatusNotFound, "Libro no encontrado")
			return
		}
		responderError(w, http.StatusInternalServerError, "Error interno del servidor")
		return
	}

	responderJSON(w, http.StatusOK, libro)
}

// POST /libros - Crear un libro nuevo
func (h *BookHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req BookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responderError(w, http.StatusBadRequest, "Cuerpo de solicitud inválido")
		return
	}

	req.Titulo = strings.TrimSpace(req.Titulo)
	req.Autor = strings.TrimSpace(req.Autor)
	req.Genero = strings.TrimSpace(req.Genero)

	if req.Titulo == "" || req.Autor == "" || req.Genero == "" {
		responderError(w, http.StatusBadRequest, "Título, autor y género son obligatorios")
		return
	}

	fecha, err := time.Parse("2006-01-02", req.FechaPublicacion)
	if err != nil {
		responderError(w, http.StatusBadRequest, "Formato de fecha inválido, use YYYY-MM-DD")
		return
	}

	// Regla de negocio (Punto 3)
	if esFechaFutura(fecha) {
		responderError(w, http.StatusBadRequest, "La fecha de publicación no puede ser futura")
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
		responderError(w, http.StatusInternalServerError, "Error al crear libro")
		return
	}

	responderJSON(w, http.StatusCreated, nuevoLibro)
}

// PUT /libros/{id} - Actualizar un libro existente
func (h *BookHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := idPositivo(r.PathValue("id"))
	if !ok {
		responderError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var req BookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responderError(w, http.StatusBadRequest, "Cuerpo de solicitud inválido")
		return
	}

	req.Titulo = strings.TrimSpace(req.Titulo)
	req.Autor = strings.TrimSpace(req.Autor)
	req.Genero = strings.TrimSpace(req.Genero)

	if req.Titulo == "" || req.Autor == "" || req.Genero == "" {
		responderError(w, http.StatusBadRequest, "Título, autor y género son obligatorios")
		return
	}

	fecha, err := time.Parse("2006-01-02", req.FechaPublicacion)
	if err != nil {
		responderError(w, http.StatusBadRequest, "Formato de fecha inválido, use YYYY-MM-DD")
		return
	}

	if esFechaFutura(fecha) {
		responderError(w, http.StatusBadRequest, "La fecha de publicación no puede ser futura")
		return
	}

	params := db.UpdateLibroParams{
		ID:               id,
		Titulo:           req.Titulo,
		Autor:            req.Autor,
		FechaPublicacion: fecha,
		Genero:           req.Genero,
	}

	rows, err := h.queries.UpdateLibro(r.Context(), params)
	if err != nil {
		responderError(w, http.StatusInternalServerError, "Error al actualizar libro")
		return
	}
	if rows == 0 {
		responderError(w, http.StatusNotFound, "Libro no encontrado")
		return
	}

	responderJSON(w, http.StatusOK, map[string]string{"mensaje": "Libro actualizado correctamente"})
}

// DELETE /libros/{id} - Eliminar un libro
func (h *BookHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := idPositivo(r.PathValue("id"))
	if !ok {
		responderError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	rows, err := h.queries.DeleteLibro(r.Context(), id)
	if err != nil {
		responderError(w, http.StatusInternalServerError, "Error al eliminar libro")
		return
	}
	if rows == 0 {
		responderError(w, http.StatusNotFound, "Libro no encontrado")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
