# Mi Biblioteca (TP6 - Interfaces Reactivas con HTMX)

Aplicación web desarrollada en **Go** utilizando **Server-Side Rendering (SSR)** con **`templ`**, **HTMX**, **`Pico.css`**, **`sqlc`** y **PostgreSQL**. Evoluciona la aplicación SSR del TP5 para brindar una experiencia interactiva tipo **SPA (Single Page Application)** —creación, edición inline, cambio de estado de lectura y borrado de libros en tiempo real sin recargas completas de página— sin escribir archivos JavaScript propios (`.js`).

## Funcionalidades

- **Visualización del catálogo**: listado dinámico de libros almacenados en PostgreSQL ordenados por ID.
- **Creación sin recarga de página**: envío asíncrono del formulario (`POST /libros`) que actualiza únicamente el contenedor de la lista con el fragmento HTML renderizado por el servidor y limpia el formulario tras un guardado exitoso (incluye validación de campos obligatorios y control de fechas futuras).
- **Edición inline sin recarga de página**: botón **Editar** por fila (`GET /libros/{id}/edit`) que convierte la fila en inputs editables para **Guardar** (`PUT /libros/{id}`) o **Cancelar** (`GET /libros/{id}`) actualizando únicamente ese `<tr>`.
- **Estado de lectura (`Leído` / `No leído`)**: botón interactivo por fila (`PATCH /libros/{id}/leido`) que alterna el estado `leido` en PostgreSQL y actualiza la fila al instante.
- **Borrado sin recarga de página**: eliminación individual desde cada fila (`DELETE /libros/{id}`) con confirmación previa (`hx-confirm`) y remoción inmediata del elemento del DOM al recibir `200 OK` vacío.

## Entidades

- **Libro**: `id`, `titulo`, `autor`, `fecha_publicacion` (`YYYY-MM-DD`), `genero` y `leido` (`boolean`, por defecto `false`).

## Tecnologías

- **Go 1.22** (`net/http` con enrutamiento por método y parámetros de ruta)
- **Templ** (`github.com/a-h/templ`) para componentes HTML tipados en Go
- **HTMX 2** (vía CDN) para interactividad declarativa basada en hipermedia
- **Pico.css** (vía CDN) para estilos semánticos
- **PostgreSQL 15**
- **sqlc** para generación de código SQL con tipado seguro en Go
- **Docker** y **Docker Compose**

## Requisitos

- **Docker** y **Docker Compose**.
- **No es necesario instalar Go, `sqlc` ni `templ` en la máquina host**: todo se genera, compila y ejecuta dentro de los contenedores Docker.
- Conexión a Internet para descargar imágenes/dependencias la primera vez y cargar **HTMX** y **Pico.css** desde CDN.

## Obtener el proyecto

```bash
git clone --branch tp6 --single-branch https://github.com/juli7etriXD/TP-Web-Grupo29.git
cd TP-Web-Grupo29
```

## Cómo ejecutar la aplicación

Desde la raíz del proyecto, ejecuta cualquiera de los siguientes comandos equivalentes:

```bash
make run-app
```

*(o directamente `./run-app.sh`)*

### ¿Qué automatiza `make run-app`?

El script [`run-app.sh`](run-app.sh) gestiona automáticamente todo el ciclo de vida de la aplicación:

1. **Tareas previas**:
   - Crea `.env` a partir de `.env.example` si no existe.
   - Detiene y elimina contenedores y volúmenes previos (`docker compose down --volumes --remove-orphans`).
   - Construye la imagen de la aplicación (`app`) con Go, `sqlc` y `templ`.
   - Levanta PostgreSQL (`postgres`) y espera activamente mediante `pg_isready` hasta que la base de datos acepte conexiones.
   - Ejecuta `sqlc generate`, `templ generate` y compila el proyecto (`go build ./...`).
2. **Ejecución**:
   - Inicia el servidor HTTP escuchando en **[http://localhost:8080](http://localhost:8080)**.
3. **Tareas posteriores**:
   - Al presionar `Ctrl+C` (o al finalizar el proceso), un `trap` automático detiene los contenedores y elimina contenedores y volúmenes (`docker compose down --volumes --remove-orphans`).

### Otros comandos disponibles en el `Makefile`

- `make up`: levanta la base de datos y la aplicación en segundo plano (`docker compose up -d --build`).
- `make stop`: detiene los contenedores y elimina los volúmenes.
- `make restart`: reinicia y reconstruye el entorno completo.
- `make test`: levanta PostgreSQL y la aplicación en Docker, compila y ejecuta la suite de pruebas Go (`go test -v ./tests/...`), limpiando todo al finalizar.
- `make generate` / `make build` / `make run`: útiles si se cuenta con Go, `sqlc` y `templ` instalados localmente.

---

## Arquitectura y Flujo Reactivo con HTMX

En esta entrega final se elimina el patrón clásico *Post/Redirect/Get* con recarga total y se adopta el intercambio de fragmentos HTML (**Hypermedia-Driven Application**), con el servidor modularizado en `main.go` (arranque y conexión) y el paquete `handlers/` (rutas, handlers y helpers de validación/render):

### Rutas HTTP (`handlers/routes.go`)

| Método y Ruta | Respuesta | Descripción |
| --- | --- | --- |
| `GET /` | `200 OK` (`text/html`) | Consulta todos los libros con `queries.GetAllLibros` y renderiza `views.IndexPage(libros)`. |
| `POST /libros` | `200 OK` (`text/html`) | Valida el formulario (incluyendo fechas futuras), inserta con `CreateLibro`, consulta `GetAllLibros` y devuelve **únicamente el fragmento HTML** de `views.EntityList(libros)`. |
| `GET /libros/{id}` | `200 OK` (`text/html`) | Obtiene el libro por ID (`GetLibroByID`) y devuelve la fila en modo lectura `views.EntityRow(libro)` (utilizado al cancelar una edición). |
| `GET /libros/{id}/edit` | `200 OK` (`text/html`) | Obtiene el libro por ID (`GetLibroByID`) y devuelve la fila editable `views.EntityRowEdit(libro)`. |
| `PUT /libros/{id}` | `200 OK` (`text/html`) | Valida los campos editados, actualiza con `UpdateLibro` y devuelve la fila actualizada `views.EntityRow(libro)`. |
| `PATCH /libros/{id}/leido` | `200 OK` (`text/html`) | Alterna el estado `leido` de forma atómica con `ToggleLibroLeido` y devuelve `views.EntityRow(libro)`. |
| `DELETE /libros/{id}` | `200 OK` (cuerpo vacío) | Elimina el registro mediante `DeleteLibro` y devuelve `200 OK` vacío (o `400`/`404` si el ID es inválido o inexistente). |
| `GET /libros` | `303 See Other` | Redirige a `/` por comodidad de navegación. |
| `GET /static/` | Archivos estáticos | Sirve recursos estáticos auxiliares (sin archivos `.js`). |

### Componentes Templ (`views/`)

- **`views/layout.templ`**: estructura HTML5 base que incluye **Pico.css** y el script de **HTMX** (`https://unpkg.com/htmx.org@2.0.4`) en el `<head>`.
- **`views/index.templ`**: página principal (`IndexPage`) que compone `@Layout`, `@EntityForm()` y `@EntityList(libros)`.
- **`views/entity_form.templ`**: formulario reactivo con `hx-post="/libros"`, `hx-target="#lista-libros"`, `hx-swap="outerHTML"` y `hx-on::after-request` para limpiar el formulario tras una creación exitosa.
- **`views/entity_list.templ`**:
  - `EntityList`: renderiza el catálogo dentro de `<article id="lista-libros">`.
  - `EntityRow`: fila de lectura con botón interactivo **Leído / No leído** (`hx-patch="/libros/{id}/leido"`), botón **Editar** (`hx-get="/libros/{id}/edit"`) y botón **Borrar** (`hx-delete="/libros/{id}"`, `hx-confirm`), todos operando sobre `hx-target="closest tr"` y `hx-swap="outerHTML"`.
  - `EntityRowEdit`: fila con inputs inline precargados y botones **Guardar** (`hx-put="/libros/{id}"`, `hx-include="closest tr"`) y **Cancelar** (`hx-get="/libros/{id}"`).

---

## Verificación Manual

1. Ejecuta `make run-app` y abre [http://localhost:8080](http://localhost:8080) en el navegador.
2. Abre **DevTools → Network** (pestaña Red).
3. **Crear**: completa el formulario y haz clic en **Agregar Libro** → se dispara `POST /libros` (`200 OK`) devolviendo `<article id="lista-libros">...`, el libro aparece como **No leído** y el formulario se limpia sin recargar la página.
4. **Leído / No leído**: haz clic en **No leído** → se dispara `PATCH /libros/{id}/leido` (`200 OK`) y la fila cambia a **Leído** (y viceversa).
5. **Editar inline**: haz clic en **Editar** (`GET /libros/{id}/edit`), modifica campos y presiona **Guardar** (`PUT /libros/{id}`), o prueba **Cancelar** (`GET /libros/{id}`).
6. **Borrar**: haz clic en **Borrar** y acepta la confirmación → se dispara `DELETE /libros/{id}` (`200 OK` vacío) y la fila desaparece del DOM sin recargar.
7. Recarga manualmente (`F5`) para confirmar que todos los cambios persisten en PostgreSQL.

### Solución de problemas

- **Error de resolución DNS al construir `app`**: si Docker muestra `Could not resolve host` o `lookup ... i/o timeout` al descargar dependencias de Go, revisa tu conexión a Internet, VPN/proxy o la configuración DNS del daemon de Docker y vuelve a ejecutar `make run-app`.
- **Puerto `8080` o `5432` ocupado**: asegúrate de no tener otra instancia de PostgreSQL o del servidor ocupando esos puertos antes de ejecutar `make run-app`.

---

## Estructura del proyecto

```text
.
├── db/
│   ├── queries/             # Consultas SQL para sqlc (CreateLibro, GetAllLibros, UpdateLibro, ToggleLibroLeido, DeleteLibro)
│   ├── schema/              # Esquema inicial PostgreSQL (tabla libro con columna leido)
│   └── sqlc/                # Código Go generado por sqlc
├── handlers/                # Paquete de rutas, handlers SSR/HTMX y helpers de validación/render
│   ├── helpers.go
│   ├── libro_handler.go
│   └── routes.go
├── static/                  # Recursos estáticos CSS (sin archivos .js)
├── tests/                   # Tests de integración Go (libro_test.go)
├── views/                   # Componentes .templ y archivos *_templ.go generados
│   ├── entity_form.templ
│   ├── entity_list.templ
│   ├── index.templ
│   └── layout.templ
├── .dockerignore
├── .env.example
├── .gitignore
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── main.go                  # Punto de entrada: configuración, conexión a PostgreSQL y arranque HTTP
├── Makefile                 # Atajos para ejecutar, compilar y testear
├── README.md                # Documentación principal del proyecto
├── run-app.sh               # Script automatizado de ciclo de vida completo
├── sqlc.yaml                # Configuración de sqlc
└── test.sh                  # Script de pruebas
```

## Autores

- Franco Nelli
- Martin Ojeda
- Julian Rivero
