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

// UserHandler encapsula las operaciones HTTP sobre la entidad Usuario
type UserHandler struct {
	queries *db.Queries
}

func NewUserHandler(q *db.Queries) *UserHandler {
	return &UserHandler{queries: q}
}

// UserRequest representa el cuerpo JSON esperado al crear un usuario
type UserRequest struct {
	Nombre          string `json:"nombre"`
	Apellido        string `json:"apellido"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	FechaNacimiento string `json:"fecha_nacimiento"` // Formato esperado: "YYYY-MM-DD"
}

// UpdateUserRequest representa el cuerpo JSON para actualizar (según queries.sql no actualiza fecha_nacimiento)
type UpdateUserRequest struct {
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type usuarioResponse struct {
	ID              int32     `json:"id"`
	Nombre          string    `json:"nombre"`
	Apellido        string    `json:"apellido"`
	Email           string    `json:"email"`
	FechaNacimiento time.Time `json:"fecha_nacimiento"`
}

func toUsuarioResponse(usuario db.Usuario) usuarioResponse {
	return usuarioResponse{
		ID:              usuario.ID,
		Nombre:          usuario.Nombre,
		Apellido:        usuario.Apellido,
		Email:           usuario.Email,
		FechaNacimiento: usuario.FechaNacimiento,
	}
}

// GET /usuarios - Listar todos los usuarios
func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	usuarios, err := h.queries.GetAllUsuarios(r.Context())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al obtener usuarios"})
		return
	}
	resultado := make([]usuarioResponse, 0, len(usuarios))
	for _, usuario := range usuarios {
		resultado = append(resultado, toUsuarioResponse(usuario))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resultado)
}

// GET /usuarios/{id} - Obtener un usuario por ID
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := idPositivo(r.PathValue("id"))
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "ID inválido"})
		return
	}

	usuario, err := h.queries.GetUserByID(r.Context(), id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Usuario no encontrado"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error interno del servidor"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(toUsuarioResponse(usuario))
}

// POST /usuarios - Crear un nuevo usuario
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Cuerpo de solicitud inválido"})
		return
	}
	if strings.TrimSpace(req.Nombre) == "" || strings.TrimSpace(req.Apellido) == "" || strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
		responderError(w, http.StatusBadRequest, "Nombre, apellido, email y contraseña son obligatorios")
		return
	}

	fecha, err := time.Parse("2006-01-02", req.FechaNacimiento) // validacion del formato
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Formato de fecha inválido, use YYYY-MM-DD"})
		return
	}

	params := db.CreateUserParams{
		Nombre:          req.Nombre,
		Apellido:        req.Apellido,
		Email:           req.Email,
		Password:        req.Password,
		FechaNacimiento: fecha,
	}

	nuevoUsuario, err := h.queries.CreateUser(r.Context(), params)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al crear usuario (verifique si el email ya existe)"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(toUsuarioResponse(nuevoUsuario))
}

// PUT /usuarios/{id} - Actualizar un usuario existente
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := idPositivo(r.PathValue("id"))
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "ID inválido"})
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Cuerpo de solicitud inválido"})
		return
	}
	if strings.TrimSpace(req.Nombre) == "" || strings.TrimSpace(req.Apellido) == "" || strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
		responderError(w, http.StatusBadRequest, "Nombre, apellido, email y contraseña son obligatorios")
		return
	}

	params := db.UpdateUserParams{
		ID:       id,
		Nombre:   req.Nombre,
		Apellido: req.Apellido,
		Email:    req.Email,
		Password: req.Password,
	}

	rows, err := h.queries.UpdateUser(r.Context(), params)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al actualizar usuario"})
		return
	}
	if rows == 0 {
		responderError(w, http.StatusNotFound, "Usuario no encontrado")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Usuario actualizado correctamente"})
}

// DELETE /usuarios/{id} - Eliminar un usuario
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := idPositivo(r.PathValue("id"))
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "ID inválido"})
		return
	}

	rows, err := h.queries.DeleteUser(r.Context(), id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al eliminar usuario"})
		return
	}
	if rows == 0 {
		responderError(w, http.StatusNotFound, "Usuario no encontrado")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
