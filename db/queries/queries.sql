--CRUD

--Create

-- Create User
INSERT INTO usuario (nombre, apellido, email, password) 
VALUES ($1, $2, $3, $4) RETURNING *;

-- Create Libro
INSERT INTO libro (titulo, autor, fecha_publicacion, genero) 
VALUES ($1, $2, $3, $4) RETURNING *;

-- Create Usuario_Libro
INSERT INTO usuario_libro (usuario_id, libro_id, leido) 
VALUES ($1, $2, $3) RETURNING *;

-- Read

-- Get User by ID
SELECT * FROM usuario WHERE id = $1;

-- Get Libro by ID
SELECT * FROM libro WHERE id = $1;

-- Get all Libros
SELECT * FROM libro;

-- Get all Usuarios
SELECT * FROM usuario;

-- Get all Libros read by a User
SELECT l.* 
FROM libro l JOIN usuario_libro ul ON l.id = ul.libro_id
WHERE ul.usuario_id = $1 AND ul.leido = true;

-- Get all User reading a specific Libro
SELECT u.*
FROM usuario u JOIN usuario_libro ul ON u.id = ul.usuario_id
WHERE ul.libro_id = $1 AND ul.leido = true;

-- Update

-- Update User_Libro Status read
UPDATE usuario_libro
SET leido = $1, fecha_lectura = CURRENT_TIMESTAMP
WHERE usuario_id = $2 AND libro_id = $3 RETURNING *;

-- Delete

-- Delete User
DELETE FROM usuario
WHERE id = $1 RETURNING *;

-- Delete Libro
DELETE FROM libro
WHERE id = $1 RETURNING *;

-- Delete User_Libro
DELETE FROM usuario_libro
WHERE usuario_id = $1 AND libro_id = $2 RETURNING *;

