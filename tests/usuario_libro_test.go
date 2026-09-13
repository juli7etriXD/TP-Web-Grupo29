package tests

import (
	"context"
	"testing"
	"time"

	db "servidor-go/db/sqlc"
)

func TestUsuarioLibro_Relacion(t *testing.T) {
	dbConn, queries := setupTestDB(t)
	defer dbConn.Close()

	ctx := context.Background()

	// 1. Crear datos previos (Usuario y Libro) requeridos para la relación
	user, err := queries.CreateUser(ctx, db.CreateUserParams{
		Nombre:          "Carlos",
		Apellido:        "Gómez",
		Email:           "carlos_relacion@example.com",
		Password:        "clave456",
		FechaNacimiento: time.Date(1995, 8, 20, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("Error al crear usuario para relación: %v", err)
	}
	defer queries.DeleteUser(ctx, user.ID)

	libro, err := queries.CreateLibro(ctx, db.CreateLibroParams{
		Titulo:           "El Aleph",
		Autor:            "Jorge Luis Borges",
		FechaPublicacion: time.Date(1949, 1, 1, 0, 0, 0, 0, time.UTC),
		Genero:           "Ficción",
	})
	if err != nil {
		t.Fatalf("Error al crear libro para relación: %v", err)
	}
	defer queries.DeleteLibro(ctx, libro.ID)

	// 2. Asociar Usuario y Libro
	_, err = queries.CreateUsuarioLibro(ctx, db.CreateUsuarioLibroParams{
		UsuarioID: user.ID,
		LibroID:   libro.ID,
		Leido:     false,
	})
	if err != nil {
		t.Fatalf("Error al asociar usuario y libro: %v", err)
	}

	// 3. Actualizar Estado a Leído
	err = queries.UpdateEstadoUsuarioLibro(ctx, db.UpdateEstadoUsuarioLibroParams{
		Leido:     true,
		UsuarioID: user.ID,
		LibroID:   libro.ID,
	})
	if err != nil {
		t.Fatalf("Error al actualizar estado de lectura: %v", err)
	}

	// 4. Consultar Libros Leídos del Usuario
	librosLeidos, err := queries.GetLibrosLeidosPorUsuario(ctx, user.ID)
	if err != nil {
		t.Fatalf("Error al obtener libros leídos: %v", err)
	}
	if len(librosLeidos) != 1 {
		t.Errorf("Se esperaba 1 libro leído, se obtuvieron %d", len(librosLeidos))
	}
}
