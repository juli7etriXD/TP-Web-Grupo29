package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

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

func idPositivo(value string) (int32, bool) {
	id, err := strconv.ParseInt(value, 10, 32)
	if err != nil || id <= 0 {
		return 0, false
	}
	return int32(id), true
}

func responderJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func responderError(w http.ResponseWriter, status int, mensaje string) {
	responderJSON(w, status, map[string]string{"error": mensaje})
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

	usuarios, err := h.queries.GetUsuariosPorLibro(r.Context(), libroID)
	if err != nil {
		responderError(w, http.StatusInternalServerError, "Error al obtener los usuarios del libro")
		return
	}

	// No devolver la contraseña que incluye el modelo completo de usuario.
	type usuarioResponse struct {
		ID              int32     `json:"id"`
		Nombre          string    `json:"nombre"`
		Apellido        string    `json:"apellido"`
		Email           string    `json:"email"`
		FechaNacimiento time.Time `json:"fecha_nacimiento"`
	}

	resultado := make([]usuarioResponse, 0, len(usuarios))
	for _, usuario := range usuarios {
		resultado = append(resultado, usuarioResponse{
			ID:              usuario.ID,
			Nombre:          usuario.Nombre,
			Apellido:        usuario.Apellido,
			Email:           usuario.Email,
			FechaNacimiento: usuario.FechaNacimiento,
		})
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

	err := h.queries.UpdateEstadoUsuarioLibro(r.Context(), db.UpdateEstadoUsuarioLibroParams{
		Leido:     req.Leido,
		UsuarioID: usuarioID,
		LibroID:   libroID,
	})
	if err != nil {
		responderError(w, http.StatusInternalServerError, "Error al actualizar el estado de lectura")
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

	err := h.queries.DeleteUsuarioLibro(r.Context(), db.DeleteUsuarioLibroParams{
		UsuarioID: usuarioID,
		LibroID:   libroID,
	})
	if err != nil {
		responderError(w, http.StatusInternalServerError, "Error al quitar el libro del usuario")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
