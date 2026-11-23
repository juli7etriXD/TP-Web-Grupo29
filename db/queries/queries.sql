--CRUD

--Create

-- name: CreateLibro :one
INSERT INTO libro (titulo, autor, fecha_publicacion, genero) 
VALUES ($1, $2, $3, $4) RETURNING *;

-- Read

-- name: GetLibroByID :one
SELECT * FROM libro WHERE id = $1;

-- name: GetAllLibros :many
SELECT * FROM libro ORDER BY id;

-- Update

-- name: UpdateLibro :execrows
UPDATE libro
SET titulo = $1, autor = $2, fecha_publicacion = $3, genero = $4
WHERE id = $5;

-- name: ToggleLibroLeido :one
UPDATE libro
SET leido = NOT leido
WHERE id = $1
RETURNING *;

-- Delete

-- name: DeleteLibro :execrows
DELETE FROM libro
WHERE id = $1;
