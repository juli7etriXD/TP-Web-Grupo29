package tests

import (
	"context"
	"testing"
	"time"

	db "servidor-go/db/sqlc"
)

func TestLibro_CRUD(t *testing.T) {
	dbConn, queries := setupTestDB(t)
	defer dbConn.Close()

	ctx := context.Background()

	// 1. Crear Libro
	libro, err := queries.CreateLibro(ctx, db.CreateLibroParams{
		Titulo:           "Ficciones",
		Autor:            "Jorge Luis Borges",
		FechaPublicacion: time.Date(1944, 1, 1, 0, 0, 0, 0, time.UTC),
		Genero:           "Ficción",
	})
	if err != nil {
		t.Fatalf("Error al crear libro: %v", err)
	}
	if libro.Leido {
		t.Errorf("Se esperaba estado inicial Leido=false, obtenido true")
	}

	// 2. Obtener Libro por ID
	fetchedLibro, err := queries.GetLibroByID(ctx, libro.ID)
	if err != nil {
		t.Fatalf("Error al obtener libro por ID: %v", err)
	}
	if fetchedLibro.Titulo != libro.Titulo {
		t.Errorf("Se esperaba título %s, obtenido %s", libro.Titulo, fetchedLibro.Titulo)
	}

	// 3. Actualizar Libro
	rows, err := queries.UpdateLibro(ctx, db.UpdateLibroParams{
		Titulo:           "Ficciones (Edición Revisada)",
		Autor:            "Jorge Luis Borges",
		FechaPublicacion: time.Date(1944, 1, 1, 0, 0, 0, 0, time.UTC),
		Genero:           "Cuentos",
		ID:               libro.ID,
	})
	if err != nil {
		t.Fatalf("Error al actualizar libro: %v", err)
	}
	if rows != 1 {
		t.Fatalf("Se esperaba 1 fila actualizada, obtenidas %d", rows)
	}

	updatedLibro, err := queries.GetLibroByID(ctx, libro.ID)
	if err != nil {
		t.Fatalf("Error al obtener libro actualizado: %v", err)
	}
	if updatedLibro.Titulo != "Ficciones (Edición Revisada)" || updatedLibro.Genero != "Cuentos" {
		t.Errorf("Los cambios no se guardaron correctamente: %+v", updatedLibro)
	}

	// 4. Alternar estado Leído / No leído
	toggledTrue, err := queries.ToggleLibroLeido(ctx, libro.ID)
	if err != nil {
		t.Fatalf("Error al marcar libro como leído: %v", err)
	}
	if !toggledTrue.Leido {
		t.Errorf("Se esperaba Leido=true tras primer toggle")
	}

	toggledFalse, err := queries.ToggleLibroLeido(ctx, libro.ID)
	if err != nil {
		t.Fatalf("Error al marcar libro como no leído: %v", err)
	}
	if toggledFalse.Leido {
		t.Errorf("Se esperaba Leido=false tras segundo toggle")
	}

	// 5. Eliminar Libro
	_, err = queries.DeleteLibro(ctx, libro.ID)
	if err != nil {
		t.Fatalf("Error al eliminar libro: %v", err)
	}
}
