package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
	"strconv"
	"time"

	"github.com/lib/pq"
)

// idPositivo valida que un parámetro de ruta sea un entero positivo (int32)
func idPositivo(value string) (int32, bool) {
	id, err := strconv.ParseInt(value, 10, 32)
	if err != nil || id <= 0 {
		return 0, false
	}
	return int32(id), true
}

// responderJSON escribe los encabezados y serializa la respuesta a JSON
func responderJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// responderError escribe una respuesta JSON estándar con el mensaje de error
func responderError(w http.ResponseWriter, status int, mensaje string) {
	responderJSON(w, status, map[string]string{"error": mensaje})
}

// esErrorDuplicado verifica si el error de PostgreSQL corresponde a violación de clave única (código 23505)
func esErrorDuplicado(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505"
	}
	return false
}

// esErrorClaveForanea verifica si el error de PostgreSQL corresponde a violación de clave foránea (código 23503)
func esErrorClaveForanea(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23503"
	}
	return false
}

// esEmailValido valida que el email tenga un formato sintáctico correcto
func esEmailValido(email string) bool {
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}
	return addr.Address == email
}

// esFechaFutura comprueba si una fecha dada es posterior a la fecha actual
func esFechaFutura(fecha time.Time) bool {
	return fecha.After(time.Now())
}
