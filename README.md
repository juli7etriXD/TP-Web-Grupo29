# Mi Biblioteca

API web desarrollada en Go para administrar usuarios, libros y la relación de lectura entre ambos. Permite registrar usuarios y libros, asociarlos, indicar si un libro fue leído y consultar libros leídos por usuario o usuarios asociados a un libro.

## Entidades

- Libro: título, autor, fecha de publicación y género.
- Usuario: nombre, apellido, email, contraseña y fecha de nacimiento. La API no expone la contraseña en las respuestas.
- Usuario-libro: relación entre ambas entidades, con estado de lectura (`leido`) y fecha de lectura.

## Tecnologías

- Go
- PostgreSQL
- sqlc
- Docker Compose
- Hurl
- HTML

## Requisitos

- Docker y Docker Compose.
- No hace falta instalar Go, `sqlc`, Hurl ni `curl` en la máquina host: se usan desde el contenedor `app`.
- Se necesita conexión a Internet para descargar las imágenes y construir `app` la primera vez.

## Obtener el proyecto

```bash
git clone --branch tp3 --single-branch https://github.com/juli7etriXD/TP-Web-Grupo29.git
cd TP-Web-Grupo29
cp .env.example .env
```

`.env.example` es la plantilla versionada con valores de desarrollo. Cópiala como `.env` para crear tu configuración local; modifica ese archivo solo si necesitas otros valores (`DB_USER`, `DB_PASSWORD`, `DB_PORT` o `DB_NAME`). No publiques credenciales personales ni agregues `.env` al repositorio. En Compose, la aplicación usa el nombre de servicio `postgres` como host de la base de datos.

## Ejecutar para desarrollo

El proyecto se levanta con dos contenedores:

- `postgres`: base de datos PostgreSQL.
- `app`: servicio Go que incluye `go`, `sqlc` y `hurl` para compilar y correr la API.

Desde la raíz del proyecto:

```bash
docker compose up -d --build
```

Compose lee `.env` automáticamente. El comando construye y levanta los dos servicios: PostgreSQL (`postgres`) y la aplicación (`app`). Al iniciar, `app` ejecuta `sqlc generate` y luego `go run .`. La API queda disponible en `http://localhost:8080`.

Para detener los contenedores manteniendo el volumen de datos:

```bash
docker compose down
```

Para ver los registros de la aplicación:

```bash
docker compose logs -f app
```

`make stop` también detiene los contenedores, pero además **borra el volumen de PostgreSQL y todos sus datos**.

## API HTTP

Las fechas se envían como `YYYY-MM-DD`. Las solicitudes con cuerpo deben incluir `Content-Type: application/json`. Los IDs de las rutas son enteros positivos. Las respuestas de error son JSON con la propiedad `error`.

### Validaciones y reglas de negocio:
- Las fechas (`fecha_nacimiento`, `fecha_publicacion`) no pueden ser futuras y deben respetar el formato `YYYY-MM-DD`. La fecha de nacimiento no puede ser anterior al año 1900.
- El `email` debe poseer un formato sintáctico válido.
- La `password` debe contener al menos 6 caracteres.
- Conflictos de unicidad (como registrar un email ya existente o asociar un libro ya vinculado) devuelven `409 Conflict`.
- Consultas o asociaciones sobre entidades no encontradas devuelven `404 Not Found`.

| Método y ruta | Descripción | JSON de entrada | Respuestas principales |
| --- | --- | --- | --- |
| `GET /libros` | Listar libros | — | `200` |
| `GET /libros/{id}` | Obtener un libro | — | `200`, `400`, `404` |
| `POST /libros` | Crear libro | `{"titulo":"Ficciones","autor":"Jorge Luis Borges","fecha_publicacion":"1944-01-01","genero":"Ficción"}` | `201`, `400` |
| `PUT /libros/{id}` | Actualizar libro | mismo formato que al crear | `200`, `400`, `404` |
| `DELETE /libros/{id}` | Eliminar libro | — | `204`, `400`, `404` |
| `GET /usuarios` | Listar usuarios | — | `200` |
| `GET /usuarios/{id}` | Obtener un usuario | — | `200`, `400`, `404` |
| `POST /usuarios` | Crear usuario | `{"nombre":"Ana","apellido":"Pérez","email":"ana@example.com","password":"clave123","fecha_nacimiento":"1998-04-15"}` | `201`, `400`, `409` |
| `PUT /usuarios/{id}` | Actualizar usuario | `{"nombre":"Ana","apellido":"Pérez","email":"ana@example.com","password":"clave456"}` | `200`, `400`, `404`, `409` |
| `DELETE /usuarios/{id}` | Eliminar usuario | — | `204`, `400`, `404` |
| `POST /usuarios/{usuario_id}/libros` | Asociar libro a usuario | `{"libro_id":1,"leido":false}` | `201`, `400`, `404`, `409` |
| `GET /usuarios/{usuario_id}/libros/leidos` | Listar libros leídos por usuario | — | `200`, `400`, `404` |
| `GET /libros/{libro_id}/usuarios` | Listar usuarios asociados a un libro | — | `200`, `400`, `404` |
| `PATCH /usuarios/{usuario_id}/libros/{libro_id}` | Cambiar estado de lectura | `{"leido":true}` | `200`, `400`, `404` |
| `DELETE /usuarios/{usuario_id}/libros/{libro_id}` | Quitar asociación | — | `204`, `400`, `404` |

Los errores inesperados del servidor o la base de datos pueden responder `500`.

## Pruebas

Ejecuta el flujo integrado desde la raíz:

```bash
make test-api
```

Este es el comando de prueba para entregar. Levanta PostgreSQL y la app con Docker Compose; dentro del contenedor `app` genera el código con `sqlc`, compila el proyecto, ejecuta las peticiones de `request.hurl` con Hurl y corre `go test -v ./tests/...`. No requiere Go, `sqlc`, Hurl ni `curl` instalados en la máquina local.

Al terminar, tanto si las pruebas pasan como si fallan, el objetivo detiene los contenedores y ejecuta `docker compose down --volumes --remove-orphans`. Esto borra el volumen de PostgreSQL y sus datos; no lo uses si necesitas conservarlos.

Para ejecutar solo las pruebas HTTP manualmente, primero levanta ambos servicios con `docker compose up -d --build`. Luego ejecuta:

```bash
docker compose exec -T app sh -lc 'hurl --test --variable "email=ana.$(date +%s%N)@example.com" request.hurl'
```

Cuando termines, usa `docker compose down` para detenerlos conservando los datos.

### Si falla la construcción de `app`

Si la construcción falla con `Could not resolve host` o `lookup ... i/o timeout` al acceder a `github.com` o `proxy.golang.org`, Docker no pudo resolver el DNS para descargar Hurl o los módulos Go. Las pruebas todavía no se ejecutaron. Comprueba la conexión a Internet, VPN/proxy y DNS de Docker; después vuelve a ejecutar `make test-api`. Si el host tiene Internet pero el error persiste dentro del build, reinicia Docker y revisa la configuración DNS del daemon o de Docker Desktop. El log puede mostrar que el servidor DNS de Docker agotó el tiempo de espera. `--no-cache` por sí solo no corrige un fallo de DNS.

## Estructura del proyecto

```text
.
├── api/
│   ├── helpers.go
│   ├── libro_handler.go
│   ├── usuario_handler.go
│   └── usuario_libro_handler.go
├── db/
│   ├── queries/
│   ├── schema/
│   └── sqlc/
├── static/
│   └── index.html
├── tests/
│   ├── libro_test.go
│   ├── setup_test.go
│   ├── usuario_libro_test.go
│   └── usuario_test.go
├── .dockerignore
├── .env.example
├── .gitignore
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── main.go
├── Makefile
├── README.md
├── request.hurl
├── sqlc.yaml
└── test.sh
```

## Autores

- Franco Nelli
- Martin Ojeda
- Julian Rivero
