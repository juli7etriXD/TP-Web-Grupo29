# Guía paso a paso para completar TP3

Esta guía adapta la consigna al proyecto Biblioteca. La API ya tiene handlers para libros y usuarios, consultas CRUD generadas con sqlc y rutas para la relación usuario-libro. Falta completar el arranque HTTP, unificar la configuración de PostgreSQL y automatizar pruebas que recorran la API.

## 1. Confirmar la rama de entrega

El TP3 se entrega en una rama llamada `tp3`, no en una carpeta nueva:

```bash
git branch --show-current
```

El comando debe mostrar `tp3`. Si todavía no estás en esa rama:

```bash
git fetch origin
git switch --track origin/tp3
```

Si ya existe una rama local `tp3`, usa `git switch tp3`.

## 2. Unificar la configuración de PostgreSQL

El archivo `.env.example` define las variables de conexión. Compose convierte las variables de usuario, contraseña y base de datos en las que espera la imagen oficial de PostgreSQL, y publica el puerto configurado en el host:

```yaml
environment:
  POSTGRES_USER: ${DB_USER:-postgres}
  POSTGRES_PASSWORD: ${DB_PASSWORD:-contrasenia}
  POSTGRES_DB: ${DB_NAME:-biblioteca}
ports:
  - "${DB_PORT:-5432}:5432"
```

Usa la sintaxis `NOMBRE=valor` (sin espacios alrededor de `=`) en `.env.example`. Para trabajar localmente, crea `.env` a partir del ejemplo y no lo agregues al repositorio si contiene credenciales personales. Compose lee `.env` automáticamente; para que la aplicación y los tests también reciban esos valores, expórtalos en la terminal antes de ejecutarlos:

```bash
set -a
. ./.env
set +a
```

La aplicación y `tests/setup_test.go` deben usar `DB_USER`, `DB_PASSWORD`, `DB_HOST`, `DB_PORT` y `DB_NAME`. Al ejecutarlos desde la máquina host, `DB_HOST` normalmente será `localhost`; dentro de Compose, la aplicación tendría que usar el nombre del servicio `postgres`.

## 3. Regenerar y verificar el código sqlc

Las consultas viven en `db/queries/queries.sql` y el código generado en `db/sqlc/`. No edites manualmente los archivos generados.

```bash
sqlc generate
go build ./...
```

Las consultas actuales ofrecen operaciones CRUD para `libro` y `usuario`. El build sirve como comprobación temprana antes de probar HTTP.

## 4. Completar el arranque del servidor HTTP

En `main.go`, conserva el registro de rutas existente y completa el arranque de la aplicación:

1. Importa `database/sql`, `log`, `os`, el driver PostgreSQL `github.com/lib/pq` como importación blank, el paquete generado `servidor-go/db/sqlc` y el paquete `servidor-go/api`.
2. Abre la conexión usando el DSN construido con las variables `DB_*` y comprueba la conexión con `Ping`.
3. Crea el repositorio con `db.New(conn)`.
4. Crea un `http.ServeMux`, registra las rutas pasando las consultas y arranca `http.ListenAndServe` en `:8080` (o en una variable configurable).
5. Cierra la conexión al terminar y registra los errores con `log.Fatal`.

El `main.go` actual define `RegisterRoutes`, pero no termina de conectar las consultas con los handlers ni inicia el servidor HTTP. Por eso la aplicación todavía no queda disponible para las peticiones de curl.

Las rutas CRUD ya registradas son:

| Recurso | Operaciones |
| --- | --- |
| `/libros` | `GET`, `POST`, `GET /{id}`, `PUT /{id}`, `DELETE /{id}` |
| `/usuarios` | `GET`, `POST`, `GET /{id}`, `PUT /{id}`, `DELETE /{id}` |
| Relación usuario-libro | `POST /usuarios/{usuario_id}/libros`, `GET /usuarios/{usuario_id}/libros/leidos`, `GET /libros/{libro_id}/usuarios`, `PATCH /usuarios/{usuario_id}/libros/{libro_id}`, `DELETE /usuarios/{usuario_id}/libros/{libro_id}` |

## 5. Revisar validaciones y respuestas HTTP

Los handlers existentes ya decodifican JSON, llaman a sqlc y responden con códigos HTTP. Antes de dar la API por terminada, verifica estos casos:

- Rechazar campos obligatorios vacíos, fechas inválidas e IDs no positivos con `400 Bad Request`.
- Responder `404 Not Found` cuando un `GET`, `PUT` o `DELETE` apunta a un registro inexistente. Para distinguir un `UPDATE`/`DELETE` sin filas afectadas, usa una consulta sqlc que devuelva filas afectadas (`:execrows`) y comprueba ese valor.
- Responder `201 Created` al crear y `204 No Content` al borrar.
- Devolver listas vacías como `[]`, no como `null`.
- No devolver la contraseña en las respuestas de `GET /usuarios` ni `GET /usuarios/{id}`.
- Mantener respuestas JSON coherentes para errores y éxitos.

## 6. Levantar la base y probar manualmente

El esquema está en `db/schema/schema.sql`. PostgreSQL ejecuta ese archivo al inicializar un volumen nuevo.

```bash
docker compose down --volumes --remove-orphans
docker compose up -d
```

Espera a que `pg_isready` indique que la base está lista. Luego, con `.env` presente, exporta sus variables y ejecuta el servidor:

```bash
set -a
. ./.env
set +a
DB_HOST=localhost go run .
```

En otra terminal, prueba al menos el ciclo CRUD de libros. Ejemplo para crear:

```bash
curl -i -X POST http://localhost:8080/libros \
  -H 'Content-Type: application/json' \
  -d '{"titulo":"Ficciones","autor":"Jorge Luis Borges","fecha_publicacion":"1944-01-01","genero":"Ficción"}'
```

Usa el `id` devuelto para probar `GET /libros`, `GET /libros/{id}`, `PUT /libros/{id}` y `DELETE /libros/{id}`. Repite el ciclo para usuarios y agrega las pruebas de la relación usuario-libro si se entregan esas rutas.

## 7. Crear el archivo de pruebas HTTP

Crea `requests.sh` en la raíz del proyecto. Debe contener ejemplos reproducibles de cada endpoint que se entrega, no solo comandos de creación. Para evitar IDs fijos, captura los IDs de las respuestas `POST` con `jq` o usa Hurl con variables.

Como mínimo, prueba en orden:

1. Crear un libro y un usuario; comprobar `201` y guardar sus IDs.
2. Listar ambas colecciones y obtener cada recurso por ID; comprobar `200`.
3. Actualizar cada recurso y verificar el resultado.
4. Crear la relación usuario-libro, consultar libros leídos y usuarios por libro, actualizar el estado con `PATCH` y quitar la relación.
5. Eliminar el usuario y el libro; comprobar `204` y que las consultas posteriores devuelven `404` cuando corresponde.
6. Probar al menos un JSON inválido, un ID inválido y un recurso inexistente para comprobar los errores.

Haz que el script falle ante cualquier estado inesperado (`set -Eeuo pipefail`) y muestre qué endpoint está probando. Dale permiso de ejecución con `chmod +x requests.sh`.

## 8. Automatizar todo con `make test-api`

Amplía `test.sh` o crea `test-api.sh` para que el flujo completo sea repetible. El script debe:

1. Cargar/validar la configuración y limpiar contenedores, redes y volúmenes anteriores del entorno de test.
2. Ejecutar `sqlc generate` y compilar con `go build ./...`.
3. Levantar PostgreSQL con `docker compose up -d` y esperar de forma acotada hasta que `pg_isready` responda.
4. Arrancar la aplicación, guardar su PID y esperar a que responda antes de iniciar las pruebas HTTP.
5. Ejecutar `requests.sh` y las pruebas Go existentes con `go test ./tests/...`.
6. Tener un `trap` que, tanto ante éxito como ante error o interrupción, detenga la aplicación y ejecute `docker compose down --volumes --remove-orphans`.

La limpieza de volúmenes borra todos los datos de esa base de Compose. Usa una base exclusivamente de pruebas. El `test.sh` actual genera sqlc, compila y ejecuta pruebas directas contra la base, pero todavía no inicia el servidor ni prueba los endpoints HTTP.

Agrega en el `Makefile` un objetivo claro:

```make
.PHONY: test-api

test-api:
	./test-api.sh
```

Corrige también el objetivo `up` existente si se conserva: la forma correcta es `docker compose up -d`.

## 9. Actualizar el README

Actualiza `README.md` para que documente:

- Qué resuelve la aplicación y las entidades disponibles.
- Requisitos: Go, Docker Compose, sqlc y `jq` si `requests.sh` lo usa.
- Cómo clonar directamente la rama `tp3`.
- Cómo configurar `.env` sin publicar credenciales.
- Cómo iniciar la aplicación y la base para desarrollo.
- Tabla de endpoints, formatos JSON y códigos de respuesta relevantes.
- Cómo ejecutar `make test-api` y qué limpia al terminar.
- Cómo detener Compose manualmente.

El README actual todavía describe la configuración como TP2 y menciona `make test`; reemplaza esas instrucciones por el comando de API definitivo.

## 10. Criterio de entrega

Desde un clon limpio, en la rama `tp3`, el flujo documentado debe poder ejecutarse así:

```bash
cp .env.example .env
make test-api
```

El comando debe generar sqlc, compilar, preparar la base, esperar a PostgreSQL, iniciar la aplicación, recorrer los endpoints desde el archivo de requests, ejecutar los tests y limpiar el entorno incluso si una prueba falla. Antes de entregar, confirma `git status`, que `README.md` y `requests.sh` estén incluidos y que no se hayan agregado credenciales reales.