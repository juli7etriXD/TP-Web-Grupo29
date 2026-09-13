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

	// 2. Obtener Libro por ID
	fetchedLibro, err := queries.GetLibroByID(ctx, libro.ID)
	if err != nil {
		t.Fatalf("Error al obtener libro por ID: %v", err)
	}
	if fetchedLibro.Titulo != libro.Titulo {
		t.Errorf("Se esperaba título %s, obtenido %s", libro.Titulo, fetchedLibro.Titulo)
	}

	// 3. Eliminar Libro
	err = queries.DeleteLibro(ctx, libro.ID)
	if err != nil {
		t.Fatalf("Error al eliminar libro: %v", err)
	}
}
