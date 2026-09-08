CREATE TABLE libro(
    id INT PRIMARY KEY SERIAL,
    titulo VARCHAR(100) NOT NULL,
    autor VARCHAR(100) NOT NULL,
    fecha_publicacion DATE NOT NULL,
    genero VARCHAR(50) NOT NULL,
)

CREATE TABLE usuario(
    id INT PRIMARY KEY SERIAL,
    nombre VARCHAR(100) NOT NULL,
    apellido VARCHAR(100) NOT NULL,
    correo_electronico VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(20) NOT NULL,
    fecha_nacimiento DATE NOT NULL
)

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