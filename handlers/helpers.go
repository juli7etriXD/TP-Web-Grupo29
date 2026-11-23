package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
)

type libroForm struct {
	Titulo           string
	Autor            string
	FechaPublicacion time.Time
	Genero           string
}

// renderHTML renderiza un componente templ con cabecera text/html UTF-8
func renderHTML(w http.ResponseWriter, r *http.Request, component templ.Component) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("Error al renderizar vista: %v", err)
		http.Error(w, "Error al renderizar la vista", http.StatusInternalServerError)
	}
}

// parseLibroID valida que un parámetro de ruta sea un entero positivo (int32)
func parseLibroID(value string) (int32, bool) {
	id, err := strconv.ParseInt(value, 10, 32)
	if err != nil || id <= 0 {
		return 0, false
	}
	return int32(id), true
}

// esFechaFutura comprueba si una fecha dada es posterior a la fecha actual
func esFechaFutura(fecha time.Time) bool {
	return fecha.After(time.Now())
}

// parseLibroForm parsea y valida los campos del formulario de creación/edición de libro
func parseLibroForm(r *http.Request) (libroForm, string) {
	if err := r.ParseForm(); err != nil {
		return libroForm{}, "Error al procesar el formulario"
	}

	titulo := strings.TrimSpace(r.FormValue("titulo"))
	autor := strings.TrimSpace(r.FormValue("autor"))
	fechaStr := strings.TrimSpace(r.FormValue("fecha_publicacion"))
	genero := strings.TrimSpace(r.FormValue("genero"))

	if titulo == "" || autor == "" || fechaStr == "" || genero == "" {
		return libroForm{}, "Todos los campos (título, autor, fecha y género) son obligatorios"
	}

	fecha, err := time.Parse("2006-01-02", fechaStr)
	if err != nil {
		return libroForm{}, "Formato de fecha inválido, debe ser YYYY-MM-DD"
	}

	if esFechaFutura(fecha) {
		return libroForm{}, "La fecha de publicación no puede ser futura"
	}

	return libroForm{
		Titulo:           titulo,
		Autor:            autor,
		FechaPublicacion: fecha,
		Genero:           genero,
	}, ""
}
