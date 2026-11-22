# Mi Biblioteca (TP6 - Interfaces Reactivas con HTMX)

Aplicación web desarrollada en **Go** utilizando **Server-Side Rendering (SSR)** con **`templ`**, **HTMX**, **`Pico.css`**, **`sqlc`** y **PostgreSQL**. Evoluciona la aplicación SSR del TP5 para brindar una experiencia interactiva tipo **SPA (Single Page Application)** —creación y borrado de libros en tiempo real sin recargas completas de página— sin escribir archivos JavaScript propios (`.js`).

## Funcionalidades

- **Visualización del catálogo**: listado dinámico de libros almacenados en PostgreSQL.
- **Creación sin recarga de página**: envío asíncrono del formulario (`POST /libros`) que actualiza únicamente el contenedor de la lista con el fragmento HTML renderizado por el servidor y limpia el formulario tras un guardado exitoso.
- **Borrado sin recarga de página**: eliminación individual desde cada fila (`DELETE /libros/{id}`) con confirmación previa y remoción inmediata del elemento del DOM al recibir `200 OK` vacío.

## Entidades

- **Libro**: `id`, `titulo`, `autor`, `fecha_publicacion` (`YYYY-MM-DD`) y `genero`.

## Tecnologías

- **Go 1.22** (`net/http` con enrutamiento por método y parámetros de ruta `DELETE /libros/{id}`)
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

En esta entrega final se elimina el patrón clásico *Post/Redirect/Get* con recarga total y se adopta el intercambio de fragmentos HTML (**Hypermedia-Driven Application**):

### Rutas HTTP (`main.go`)

| Método y Ruta | Respuesta | Descripción |
| --- | --- | --- |
| `GET /` | `200 OK` (`text/html`) | Consulta todos los libros con `queries.GetAllLibros` y renderiza la página completa `views.IndexPage(libros)`. |
| `POST /libros` | `200 OK` (`text/html`) | Valida los campos del formulario, inserta el libro con `queries.CreateLibro`, vuelve a consultar `GetAllLibros` y devuelve **únicamente el fragmento HTML** de `views.EntityList(libros)` (sin redirección). |
| `DELETE /libros/{id}` | `200 OK` (cuerpo vacío) | Extrae el `id` con `r.PathValue("id")`, valida que sea un entero positivo, elimina el registro mediante `queries.DeleteLibro` y devuelve `200 OK` vacío (o `400`/`404` si el ID es inválido o inexistente). |
| `GET /libros` | `303 See Other` | Redirige a `/` por comodidad de navegación. |
| `GET /static/` | Archivos estáticos | Sirve recursos estáticos auxiliares (solo CSS, sin archivos `.js`). |

### Componentes Templ (`views/`)

- **`views/layout.templ`**: estructura HTML5 base que incluye **Pico.css** y el script de **HTMX** (`https://unpkg.com/htmx.org@2.0.4`) en el `<head>`.
- **`views/index.templ`**: página principal (`IndexPage`) que compone `@Layout`, `@EntityForm()` y `@EntityList(libros)`.
- **`views/entity_form.templ`**: formulario reactivo configurado con atributos HTMX:
  - `hx-post="/libros"`: envía los datos del formulario de forma asíncrona.
  - `hx-target="#lista-libros"`: apunta al contenedor `<article id="lista-libros">` de la lista.
  - `hx-swap="outerHTML"`: reemplaza el contenedor completo por el nuevo HTML devuelto por el servidor (manejando limpiamente la transición de lista vacía a tabla).
  - `hx-on::after-request="if(event.detail.successful) this.reset()"`: limpia los campos del formulario tras una creación exitosa.
- **`views/entity_list.templ`**: renderiza el catálogo dentro de `<article id="lista-libros">` e incluye en cada fila un botón **Borrar** configurado con:
  - `hx-delete={ fmt.Sprintf("/libros/%d", libro.ID) }`: dispara `DELETE /libros/{id}`.
  - `hx-target="closest tr"`: selecciona la fila `<tr>` correspondiente al libro.
  - `hx-swap="outerHTML"`: reemplaza la fila por la respuesta vacía del servidor, eliminándola del DOM.
  - `hx-confirm`: solicita confirmación al usuario antes de eliminar el registro.

---

## Verificación Manual

1. Ejecuta `make run-app` y abre [http://localhost:8080](http://localhost:8080) en el navegador.
2. Abre **DevTools → Network** (pestaña Red).
3. Completa el formulario y haz clic en **Agregar Libro**:
   - Se dispara `POST /libros` con respuesta `200 OK` conteniendo únicamente el fragmento `<article id="lista-libros">...`.
   - El libro aparece al instante en la tabla y el formulario se limpia **sin recargar la página**.
4. Haz clic en **Borrar** en cualquier fila y acepta la confirmación:
   - Se dispara `DELETE /libros/{id}` con respuesta `200 OK` vacía.
   - La fila desaparece inmediatamente del DOM **sin recargar la página**.
5. Recarga manualmente la página (`F5`) para confirmar que los cambios persisten en PostgreSQL.

### Solución de problemas

- **Error de resolución DNS al construir `app`**: si Docker muestra `Could not resolve host` o `lookup ... i/o timeout` al descargar dependencias de Go, revisa tu conexión a Internet, VPN/proxy o la configuración DNS del daemon de Docker y vuelve a ejecutar `make run-app`.
- **Puerto `8080` o `5432` ocupado**: asegúrate de no tener otra instancia de PostgreSQL o del servidor ocupando esos puertos antes de ejecutar `make run-app`.

---

## Estructura del proyecto

```text
.
├── api/                     # Handlers previos de entregas anteriores
├── db/
│   ├── queries/             # Consultas SQL para sqlc (CreateLibro, GetAllLibros, DeleteLibro, etc.)
│   ├── schema/              # Esquema inicial PostgreSQL (schema.sql)
│   └── sqlc/                # Código Go generado por sqlc
├── static/                  # Recursos estáticos CSS (sin archivos .js)
├── tests/                   # Tests de integración Go
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
├── main.go                  # Servidor HTTP y handlers SSR/HTMX
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
