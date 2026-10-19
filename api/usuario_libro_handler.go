package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	db "servidor-go/db/sqlc"
)

type UsuarioLibroHandler struct {
	queries *db.Queries
}

func NewUsuarioLibroHandler(q *db.Queries) *UsuarioLibroHandler {
	return &UsuarioLibroHandler{queries: q}
}

type crearRelacionRequest struct {
	LibroID int32 `json:"libro_id"`
	Leido   bool  `json:"leido"`
}

type actualizarEstadoRequest struct {
	Leido bool `json:"leido"`
}

// POST /usuarios/{usuario_id}/libros
// JSON: {"libro_id": 1, "leido": false}
func (h *UsuarioLibroHandler) Create(w http.ResponseWriter, r *http.Request) {
	usuarioID, ok := idPositivo(r.PathValue("usuario_id"))
	if !ok {
		responderError(w, http.StatusBadRequest, "ID de usuario inválido")
		return
	}

	var req crearRelacionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.LibroID <= 0 {
		responderError(w, http.StatusBadRequest, "Cuerpo de solicitud o ID de libro inválido")
		return
	}

	relacion, err := h.queries.CreateUsuarioLibro(r.Context(), db.CreateUsuarioLibroParams{
		UsuarioID: usuarioID,
		LibroID:   req.LibroID,
		Leido:     req.Leido,
	})
	if err != nil {
		if esErrorDuplicado(err) {
			responderError(w, http.StatusConflict, "El libro ya está asociado al usuario")
			return
		}
		if esErrorClaveForanea(err) {
			responderError(w, http.StatusNotFound, "El usuario o el libro especificado no existe")
			return
		}
		responderError(w, http.StatusInternalServerError, "Error al asociar el libro al usuario")
		return
	}

	responderJSON(w, http.StatusCreated, relacion)
}

// GET /usuarios/{usuario_id}/libros/leidos
func (h *UsuarioLibroHandler) GetLibrosLeidosPorUsuario(w http.ResponseWriter, r *http.Request) {
	usuarioID, ok := idPositivo(r.PathValue("usuario_id"))
	if !ok {
		responderError(w, http.StatusBadRequest, "ID de usuario inválido")
		return
	}

	// Verificar si el usuario existe antes de consultar sus libros leídos
	if _, err := h.queries.GetUserByID(r.Context(), usuarioID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			responderError(w, http.StatusNotFound, "Usuario no encontrado")
			return
		}
		responderError(w, http.StatusInternalServerError, "Error al verificar existencia del usuario")
		return
	}

	libros, err := h.queries.GetLibrosLeidosPorUsuario(r.Context(), usuarioID)
	if err != nil {
		responderError(w, http.StatusInternalServerError, "Error al obtener los libros leídos")
		return
	}
	if libros == nil {
		libros = []db.Libro{}
	}

	responderJSON(w, http.StatusOK, libros)
}

// GET /libros/{libro_id}/usuarios
func (h *UsuarioLibroHandler) GetUsuariosPorLibro(w http.ResponseWriter, r *http.Request) {
	libroID, ok := idPositivo(r.PathValue("libro_id"))
	if !ok {
		responderError(w, http.StatusBadRequest, "ID de libro inválido")
		return
	}

	// Verificar si el libro existe antes de consultar sus usuarios asociados
	if _, err := h.queries.GetLibroByID(r.Context(), libroID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			responderError(w, http.StatusNotFound, "Libro no encontrado")
			return
		}
		responderError(w, http.StatusInternalServerError, "Error al verificar existencia del libro")
		return
	}

	usuarios, err := h.queries.GetUsuariosPorLibro(r.Context(), libroID)
	if err != nil {
		responderError(w, http.StatusInternalServerError, "Error al obtener los usuarios del libro")
		return
	}

	resultado := make([]usuarioResponse, 0, len(usuarios))
	for _, usuario := range usuarios {
		resultado = append(resultado, toUsuarioResponse(usuario))
	}

	responderJSON(w, http.StatusOK, resultado)
}

// PATCH /usuarios/{usuario_id}/libros/{libro_id}
// JSON: {"leido": true}
func (h *UsuarioLibroHandler) UpdateEstado(w http.ResponseWriter, r *http.Request) {
	usuarioID, ok := idPositivo(r.PathValue("usuario_id"))
	if !ok {
		responderError(w, http.StatusBadRequest, "ID de usuario inválido")
		return
	}
	libroID, ok := idPositivo(r.PathValue("libro_id"))
	if !ok {
		responderError(w, http.StatusBadRequest, "ID de libro inválido")
		return
	}

	var req actualizarEstadoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responderError(w, http.StatusBadRequest, "Cuerpo de solicitud inválido")
		return
	}

	rows, err := h.queries.UpdateEstadoUsuarioLibro(r.Context(), db.UpdateEstadoUsuarioLibroParams{
		Leido:     req.Leido,
		UsuarioID: usuarioID,
		LibroID:   libroID,
	})
	if err != nil {
		responderError(w, http.StatusInternalServerError, "Error al actualizar el estado de lectura")
		return
	}
	if rows == 0 {
		responderError(w, http.StatusNotFound, "Relación usuario-libro no encontrada")
		return
	}

	responderJSON(w, http.StatusOK, map[string]string{
		"mensaje": "Estado de lectura actualizado correctamente",
	})
}

// DELETE /usuarios/{usuario_id}/libros/{libro_id}
func (h *UsuarioLibroHandler) Delete(w http.ResponseWriter, r *http.Request) {
	usuarioID, ok := idPositivo(r.PathValue("usuario_id"))
	if !ok {
		responderError(w, http.StatusBadRequest, "ID de usuario inválido")
		return
	}
	libroID, ok := idPositivo(r.PathValue("libro_id"))
	if !ok {
		responderError(w, http.StatusBadRequest, "ID de libro inválido")
		return
	}

	rows, err := h.queries.DeleteUsuarioLibro(r.Context(), db.DeleteUsuarioLibroParams{
		UsuarioID: usuarioID,
		LibroID:   libroID,
	})
	if err != nil {
		responderError(w, http.StatusInternalServerError, "Error al quitar el libro del usuario")
		return
	}
	if rows == 0 {
		responderError(w, http.StatusNotFound, "Relación usuario-libro no encontrada")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
