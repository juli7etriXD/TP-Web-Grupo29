# Mi Biblioteca

Aplicación web desarrollada en Go para gestionar usuarios y controlar qué libros fueron leídos por cada uno. Permitiendo a los usuarios saber que libros leyeron, agregar un libro que leyeron y eliminar un libro que leyeron hace mucho tiempo.

## Objetivo

Permitir llevar un registro de:

- Usuarios
- Libros
- Estado de lectura (`leído / no leído`)

## Tecnologías utilizadas

- Go
- PostgreSQL
- sqlc
- Docker y Docker Compose
- HTML

## Funcionalidades principales

- Registro de usuarios
- Alta, edición y eliminación de libros
- Asociación de libros a usuarios
- Marcado de libros como leídos
- Consulta de libros leídos por usuario
- Consulta de usuarios que leyeron un libro

## Estructura del proyecto

```text
.
├── db/
│   ├── queries/
│   │   └── queries.sql          # Consultas SQL con anotaciones de 
│   ├── schema/
│   │   └── schema.sql           # Definición de tablas de PostgreSQL
│   └── sqlc/                   # Código generado por sqlc
├── static/
│   └── index.html              # Frontend estático básico
├── docker-compose.yml          # Configuración del contenedor de PostgreSQL
├── go.mod                      # Dependencias del módulo Go
├── go.sum                      # Lockfile de dependencias
├── main.go                     # Servidor HTTP principal
├── Makefile                    # Comandos útiles del proyecto
├── sqlc.yaml                   # Configuración de sqlc
├── test.sh                     # Script de validación automática
├── tests/                      # Pruebas del proyecto
└── README.md                   # Documentación del proyecto
```

## Base de datos ```(schema.sql)```

El esquema principal está definido en estas tablas:

- `libro`: Almacena la información de los libros (`id`,`titulo`,`autor`,`fecha_publicacion`,`genero`).
- `usuario`: Guarda las cuentas e información de los usuarios (`id`,`nombre`,`apellido`,`email`,`password`,`fecha_nacimiento`).
- `usuario_libro`: Relación muchos a muchos entre usuarios y libros para el seguimiento de lecturas (`usuario_id`,`libro_id`,`leido`,`fecha_lectura`).

## Requisitos previos

Antes de ejecutar el proyecto, asegúrate de tener instalados:

- [Go](https://go.dev/dl/) 1.22 o superior
- [Docker](https://www.docker.com/products/docker-desktop/) y Docker Compose
- [sqlc](https://sqlc.dev/)

## Configuración del entorno para tp2

```bash
git clone --branch tp2 --single-branch https://github.com/usuario/proyecto.git
cd TP-Web-Grupo29/
```

## Ejecución del proyecto

Antes de iniciar la aplicación:

```bash
make test
```

o

```bash
./test.sh
```

### ¿Qué hace `./test.sh`?

- borra contenedores y volúmenes anteriores
- levanta la base de datos
- genera el código SQL con `sqlc`
- compila el proyecto
- ejecuta los tests
- limpia al finalizar

### Detener los servicios

```bash
make stop
```

## Autores

- Franco Nelli
- Martin Ojeda
- Julian Rivero

