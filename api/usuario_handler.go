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
		responderError(w, http.StatusInternalServerError, "Error al obtener usuarios")
		return
	}
	resultado := make([]usuarioResponse, 0, len(usuarios))
	for _, usuario := range usuarios {
		resultado = append(resultado, toUsuarioResponse(usuario))
	}

	responderJSON(w, http.StatusOK, resultado)
}

// GET /usuarios/{id} - Obtener un usuario por ID
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := idPositivo(r.PathValue("id"))
	if !ok {
		responderError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	usuario, err := h.queries.GetUserByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			responderError(w, http.StatusNotFound, "Usuario no encontrado")
			return
		}
		responderError(w, http.StatusInternalServerError, "Error interno del servidor")
		return
	}

	responderJSON(w, http.StatusOK, toUsuarioResponse(usuario))
}

// POST /usuarios - Crear un nuevo usuario
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responderError(w, http.StatusBadRequest, "Cuerpo de solicitud inválido")
		return
	}

	req.Nombre = strings.TrimSpace(req.Nombre)
	req.Apellido = strings.TrimSpace(req.Apellido)
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	if req.Nombre == "" || req.Apellido == "" || req.Email == "" || req.Password == "" {
		responderError(w, http.StatusBadRequest, "Nombre, apellido, email y contraseña son obligatorios")
		return
	}

	// Reglas de negocio (Punto 3)
	if !esEmailValido(req.Email) {
		responderError(w, http.StatusBadRequest, "Formato de email inválido")
		return
	}

	if len(req.Password) < 6 {
		responderError(w, http.StatusBadRequest, "La contraseña debe tener al menos 6 caracteres")
		return
	}

	fecha, err := time.Parse("2006-01-02", req.FechaNacimiento)
	if err != nil {
		responderError(w, http.StatusBadRequest, "Formato de fecha inválido, use YYYY-MM-DD")
		return
	}
	if esFechaFutura(fecha) {
		responderError(w, http.StatusBadRequest, "La fecha de nacimiento no puede ser futura")
		return
	}
	if fecha.Year() < 1900 {
		responderError(w, http.StatusBadRequest, "La fecha de nacimiento es inválida (año anterior a 1900)")
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
		// Código HTTP 409 Conflict si el email ya existe (Punto 2)
		if esErrorDuplicado(err) {
			responderError(w, http.StatusConflict, "El email ya se encuentra registrado")
			return
		}
		responderError(w, http.StatusInternalServerError, "Error al crear usuario")
		return
	}

	responderJSON(w, http.StatusCreated, toUsuarioResponse(nuevoUsuario))
}

// PUT /usuarios/{id} - Actualizar un usuario existente
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := idPositivo(r.PathValue("id"))
	if !ok {
		responderError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responderError(w, http.StatusBadRequest, "Cuerpo de solicitud inválido")
		return
	}

	req.Nombre = strings.TrimSpace(req.Nombre)
	req.Apellido = strings.TrimSpace(req.Apellido)
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	if req.Nombre == "" || req.Apellido == "" || req.Email == "" || req.Password == "" {
		responderError(w, http.StatusBadRequest, "Nombre, apellido, email y contraseña son obligatorios")
		return
	}

	if !esEmailValido(req.Email) {
		responderError(w, http.StatusBadRequest, "Formato de email inválido")
		return
	}

	if len(req.Password) < 6 {
		responderError(w, http.StatusBadRequest, "La contraseña debe tener al menos 6 caracteres")
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
		if esErrorDuplicado(err) {
			responderError(w, http.StatusConflict, "El email ya se encuentra registrado por otro usuario")
			return
		}
		responderError(w, http.StatusInternalServerError, "Error al actualizar usuario")
		return
	}
	if rows == 0 {
		responderError(w, http.StatusNotFound, "Usuario no encontrado")
		return
	}

	responderJSON(w, http.StatusOK, map[string]string{"mensaje": "Usuario actualizado correctamente"})
}

// DELETE /usuarios/{id} - Eliminar un usuario
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := idPositivo(r.PathValue("id"))
	if !ok {
		responderError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	rows, err := h.queries.DeleteUser(r.Context(), id)
	if err != nil {
		responderError(w, http.StatusInternalServerError, "Error al eliminar usuario")
		return
	}
	if rows == 0 {
		responderError(w, http.StatusNotFound, "Usuario no encontrado")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
