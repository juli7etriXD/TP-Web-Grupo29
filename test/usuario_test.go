package tests

import (
	"context"
	"testing"
	"time"

	db "servidor-go/db/sqlc"
)

func TestUsuario_CRUD(t *testing.T) {
	dbConn, queries := setupTestDB(t)
	defer dbConn.Close()

	ctx := context.Background()

	// 1. Crear Usuario
	user, err := queries.CreateUser(ctx, db.CreateUserParams{
		Nombre:          "Ana",
		Apellido:        "Pérez",
		Email:           "ana_test@example.com",
		Password:        "clave123",
		FechaNacimiento: time.Date(1998, 4, 15, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("Error al crear usuario: %v", err)
	}

	// 2. Obtener Usuario por ID
	fetchedUser, err := queries.GetUserByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("Error al obtener usuario por ID: %v", err)
	}
	if fetchedUser.Email != user.Email {
		t.Errorf("Se esperaba email %s, obtenido %s", user.Email, fetchedUser.Email)
	}

	// 3. Eliminar Usuario
	err = queries.DeleteUser(ctx, user.ID)
	if err != nil {
		t.Fatalf("Error al eliminar usuario: %v", err)
	}
}
