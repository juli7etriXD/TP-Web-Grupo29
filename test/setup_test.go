package tests

import (
	"database/sql"
	"testing"

	// Ajusta la ruta al nombre de tu módulo en go.mod + la ubicación de sqlc
	db "servidor-go/db/sqlc"

	_ "github.com/lib/pq"
)

// setupTestDB centraliza la conexión a PostgreSQL para reutilizarla en todos los tests
func setupTestDB(t *testing.T) (*sql.DB, *db.Queries) {
	connStr := "postgres://postgres:postgres@localhost:5432/tp2_db?sslmode=disable"
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Error al conectar con la base de datos: %v", err)
	}

	if err := dbConn.Ping(); err != nil {
		dbConn.Close()
		t.Fatalf("La base de datos no está respondiendo: %v", err)
	}

	return dbConn, db.New(dbConn)
}

/*package tests


import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"
	db "module servidor-go/tp2"
)


func TestUsuarioCRUD(t *testing.T) {
	queries, _ := setupTestDB(t)
	ctx := context.Background()

	var createdUser db.Usuario

	// 1. CREATE
	t.Run("CreateUser", func(t *testing.T) {
		arg := db.CreateUserParams{
			Nombre:          "Martin",
			Apellido:        "Fowler",
			Email:           "martin@refactoring.com",
			Password:        "secret123",
			FechaNacimiento: time.Date(1963, 12, 18, 0, 0, 0, 0, time.UTC),
		}

		user, err := queries.CreateUser(ctx, arg)
		if err != nil {
			t.Fatalf("Error al crear usuario: %v", err)
		}

		if user.ID == 0 || user.Email != arg.Email {
			t.Errorf("Usuario creado con datos incorrectos: %+v", user)
		}

		createdUser = user
	})



-----------------------------------------
package tests


import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"
	db "module servidor-go/tp2"
)


func TestDominioBiblioteca_CRUD(t *testing.T) {
	connStr := "postgres://postgres:postgres@localhost:5432/tp2_db?sslmode=disable"
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Error al conectar con la base de datos: %v", err)
	}
	defer dbConn.Close()

	if err := dbConn.Ping(); err != nil {
		t.Fatalf("La base de datos no está respondiendo: %v", err)
	}

	queries := db.New(dbConn)
	ctx := context.Background()

	// 1. Test Crear Usuario
	user, err := queries.CreateUser(ctx, db.CreateUserParams{
		Nombre:          "Ana",
		Apellido:        "Pérez",
		Email:           "ana@example.com",
		Password:        "clave123",
		FechaNacimiento: time.Date(1998, 4, 15, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("Error al crear usuario: %v", err)
	}

	// 2. Test Crear Libro
	libro, err := queries.CreateLibro(ctx, db.CreateLibroParams{
		Titulo:           "Ficciones",
		Autor:            "Jorge Luis Borges",
		FechaPublicacion: time.Date(1944, 1, 1, 0, 0, 0, 0, time.UTC),
		Genero:           "Ficción",
	})
	if err != nil {
		t.Fatalf("Error al crear libro: %v", err)
	}

	// 3. Test Asociar Usuario y Libro
	_, err = queries.CreateUsuarioLibro(ctx, db.CreateUsuarioLibroParams{
		UsuarioID: user.ID,
		LibroID:   libro.ID,
		Leido:     false,
	})
	if err != nil {
		t.Fatalf("Error al crear relación usuario_libro: %v", err)
	}

	// 4. Test Actualizar Estado a Leído
	err = queries.UpdateEstadoUsuarioLibro(ctx, db.UpdateEstadoUsuarioLibroParams{
		Leido:     true,
		UsuarioID: user.ID,
		LibroID:   libro.ID,
	})
	if err != nil {
		t.Fatalf("Error al actualizar estado de lectura: %v", err)
	}

	// 5. Test Consultar Libros Leídos del Usuario
	librosLeidos, err := queries.GetLibrosLeidosPorUsuario(ctx, user.ID)
	if err != nil {
		t.Fatalf("Error al obtener libros leídos: %v", err)
	}
	if len(librosLeidos) != 1 {
		t.Errorf("Se esperaba 1 libro leído, se obtuvieron %d", len(librosLeidos))
	}

	// 6. Test Eliminar Libro (el CASCADE eliminará la relación)
	err = queries.DeleteLibro(ctx, libro.ID)
	if err != nil {
		t.Fatalf("Error al eliminar libro: %v", err)
	}
}
*/
