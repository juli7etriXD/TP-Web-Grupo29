CREATE TABLE libro(
    id SERIAL PRIMARY KEY,
    titulo VARCHAR(100) NOT NULL,
    autor VARCHAR(100) NOT NULL,
    fecha_publicacion DATE NOT NULL,
    genero VARCHAR(50) NOT NULL
);

CREATE TABLE usuario(
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    apellido VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    fecha_nacimiento DATE NOT NULL
);

CREATE TABLE usuario_libro(
    usuario_id INT NOT NULL,
    libro_id INT NOT NULL,
    leido BOOLEAN NOT NULL DEFAULT false,
    fecha_lectura TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Opcional: saber cuándo lo leyó
    
    PRIMARY KEY (usuario_id, libro_id),
    
    FOREIGN KEY (usuario_id) REFERENCES usuario(id) 
    ON DELETE CASCADE,
    FOREIGN KEY (libro_id) REFERENCES libro(id) 
    ON DELETE CASCADE
);