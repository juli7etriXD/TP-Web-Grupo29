CREATE TABLE libro(
    id SERIAL PRIMARY KEY,
    titulo VARCHAR(100) NOT NULL,
    autor VARCHAR(100) NOT NULL,
    fecha_publicacion DATE NOT NULL,
    genero VARCHAR(50) NOT NULL
);