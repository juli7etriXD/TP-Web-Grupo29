--CRUD

--Create

-- name: CreateUser :one
INSERT INTO usuario (nombre, apellido, email, password, fecha_nacimiento)
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: CreateLibro :one
INSERT INTO libro (titulo, autor, fecha_publicacion, genero) 
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: CreateUsuarioLibro :one
INSERT INTO usuario_libro (usuario_id, libro_id, leido) 
VALUES ($1, $2, $3) RETURNING *;

-- Read

-- name: GetUserByID :one
SELECT * FROM usuario WHERE id = $1;

-- name: GetLibroByID :one
SELECT * FROM libro WHERE id = $1;

-- name: GetAllLibros :many
SELECT * FROM libro;

-- name: GetAllUsuarios :many
SELECT * FROM usuario;

-- name: GetLibrosLeidosPorUsuario :many
SELECT l.* 
FROM libro l JOIN usuario_libro ul ON l.id = ul.libro_id
WHERE ul.usuario_id = $1 AND ul.leido = true;

-- name: GetUsuariosPorLibro :many
SELECT u.*
FROM usuario u JOIN usuario_libro ul ON u.id = ul.usuario_id
WHERE ul.libro_id = $1 AND ul.leido = true;

-- Update

-- name: UpdateUser :execrows
UPDATE usuario
SET nombre = $1, apellido = $2, email = $3, password = $4
WHERE id = $5;

-- name: UpdateLibro :execrows
UPDATE libro
SET titulo = $1, autor = $2, fecha_publicacion = $3, genero = $4
WHERE id = $5;

-- name: UpdateEstadoUsuarioLibro :execrows
UPDATE usuario_libro
SET leido = $1, fecha_lectura = CURRENT_TIMESTAMP
WHERE usuario_id = $2 AND libro_id = $3;

-- Delete

-- name: DeleteUser :execrows
DELETE FROM usuario
WHERE id = $1;

-- name: DeleteLibro :execrows
DELETE FROM libro
WHERE id = $1;

-- name: DeleteUsuarioLibro :execrows
DELETE FROM usuario_libro
WHERE usuario_id = $1 AND libro_id = $2;

