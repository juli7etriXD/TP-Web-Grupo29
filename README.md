# Mi Biblioteca (TP5 - Renderizado del Lado del Servidor con Templ)

Aplicación web desarrollada en **Go** utilizando **Server-Side Rendering (SSR)** con **`templ`**, **`Pico.css`**, **`sqlc`** y **PostgreSQL**. Permite visualizar el catálogo de libros y registrar nuevos ejemplares mediante formularios HTML tradicionales aplicando el patrón **Post/Redirect/Get (PRG)**, sin JavaScript en el cliente.

## Entidades

- **Libro**: título, autor, fecha de publicación (`YYYY-MM-DD`) y género.

## Tecnologías

- **Go** (net/http)
- **Templ** (`github.com/a-h/templ`) para generación de vistas HTML con tipado estricto
- **Pico.css** (vía CDN) para estilos semánticos
- **PostgreSQL 15**
- **sqlc** para generación de código SQL en Go
- **Docker** y **Docker Compose**

## Requisitos

- Docker y Docker Compose.
- **No es necesario instalar Go, `sqlc` ni `templ` en la máquina host**: todo se genera, compila y ejecuta dentro de los contenedores Docker.
- Conexión a Internet para descargar imágenes/dependencias la primera vez y cargar Pico.css desde CDN.

## Obtener el proyecto

```bash
git clone --branch tp5 --single-branch https://github.com/juli7etriXD/TP-Web-Grupo29.git
cd TP-Web-Grupo29
```

## Cómo ejecutar la aplicación

Desde la raíz del proyecto, ejecuta cualquiera de los siguientes comandos equivalentes:

```bash
make run-app
```
*(o directamente `./run-app.sh`)*

### ¿Qué automatiza `make run-app`?
El script [`run-app.sh`](run-app.sh) ejecuta de forma automática todo el ciclo de vida:
1. **Tareas previas**:
   - Crea `.env` a partir de `.env.example` si no existe.
   - Limpia contenedores y volúmenes previos (`docker compose down --volumes --remove-orphans`).
   - Construye la imagen de la aplicación (`app`) con Go, `sqlc` y `templ`.
   - Levanta el motor de base de datos PostgreSQL (`postgres`) y espera activamente con `pg_isready` hasta que esté listo para aceptar conexiones.
   - Ejecuta `sqlc generate`, `templ generate` y compila el proyecto (`go build ./...`).
2. **Ejecución**:
   - Compila e inicia el servidor web escuchando en **[http://localhost:8080](http://localhost:8080)**.
3. **Tareas posteriores**:
   - Al presionar `Ctrl+C` (o finalizar la ejecución), un `trap` automático detiene los contenedores y elimina contenedores y volúmenes (`docker compose down --volumes --remove-orphans`).

---

### Otros comandos disponibles en el `Makefile`

- `make up`: levanta la base de datos y la aplicación en segundo plano (`docker compose up -d --build`).
- `make stop`: detiene los contenedores y elimina los volúmenes.
- `make test`: ejecuta la suite de pruebas Go dentro del contenedor.
- `make generate` / `make build` / `make run`: útiles si tienes Go, `sqlc` y `templ` instalados localmente.

## Funcionamiento y Rutas SSR

En esta versión la aplicación abandona la arquitectura API JSON + JavaScript en el cliente para renderizar HTML directamente desde el servidor Go:

| Método y Ruta | Descripción |
| --- | --- |
| `GET /` | Consulta todos los libros en PostgreSQL mediante `sqlc` (`GetAllLibros`) y renderiza la vista `views.IndexPage` con `templ`. |
| `POST /libros` | Procesa los datos del formulario (`application/x-www-form-urlencoded`), valida y guarda el libro con `CreateLibro` y redirige a `/` mediante `303 See Other` (**patrón Post/Redirect/Get**). |
| `GET /libros` | Redirige a `/` para comodidad de navegación. |
| `GET /static/` | Sirve archivos estáticos auxiliares si se requieren. |

### Componentes Templ (`views/`)
- `views/layout.templ`: estructura base HTML5 (`<head>`, `<body>`, barra de navegación y enlace a Pico.css).
- `views/entity_form.templ`: formulario clásico para crear un libro (`action="/libros" method="POST"`).
- `views/entity_list.templ`: componente que recibe `[]db.Libro` y renderiza la tabla del catálogo o un estado vacío.
- `views/index.templ`: página principal (`IndexPage`) que compone el layout, el formulario y la lista.

### Si falla la construcción de `app`
Si la construcción de la imagen falla con `Could not resolve host` o `lookup ... i/o timeout` al acceder a `github.com` o `proxy.golang.org`, Docker tuvo un fallo de resolución DNS. Verifica tu conexión a Internet, VPN/proxy y la configuración DNS del daemon de Docker o Docker Desktop, y vuelve a intentar `make run-app`.

## Estructura del proyecto

```text
.
├── api/                     # Handlers previos de TP3
├── db/
│   ├── queries/             # Consultas SQL para sqlc
│   ├── schema/              # Esquema inicial PostgreSQL
│   └── sqlc/                # Código Go generado por sqlc
├── static/                  # Recursos estáticos CSS (sin archivos .js)
├── tests/                   # Tests de integración
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
├── main.go
├── Makefile
├── README.md
├── run-app.sh
├── sqlc.yaml
└── test.sh
```

## Autores

- Franco Nelli
- Martin Ojeda
- Julian Rivero
