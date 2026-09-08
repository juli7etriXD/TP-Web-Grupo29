# Mi Biblioteca

Aplicación web en Go para gestionar usuarios y el seguimiento de sus libros leídos.

## Tecnologías utilizadas

- **Go** (1.22+)
- **PostgreSQL** (16)
- **sqlc** (para la generación de código Go tipado a partir de SQL)
- **HTML/CSS** básico

## Estructura del proyecto

```text
├── db/
│   ├── queries/       # Consultas SQL (CRUD)
│   ├── schema/        # Definición de tablas (PostgreSQL)
│   └── sqlc/          # Código Go generado automáticamente por sqlc
├── static/            # Archivos estáticos del frontend (HTML, etc.)
├── go.mod             # Definición del módulo y dependencias de Go
├── go.sum             # Checksums de dependencias
├── main.go            # Servidor HTTP y punto de entrada
└── sqlc.yaml          # Configuración de sqlc
```

## Ejecución con Docker (Recomendada)

## Ejecución local (sin Docker)

### Requisitos previos
- [Go](https://go.dev/dl/) (1.22 o superior)
- [PostgreSQL](https://www.postgresql.org/) en ejecución en el puerto 5432

### Pasos
1. **Crear la base de datos y tablas:**
   ```bash
   psql -U postgres -d biblioteca -f db/schema/schema.sql
   ```

2. **(Opcional) Regenerar código con sqlc:**
   Si hacés cambios en `db/schema/` o `db/queries/`:
   ```bash
   sqlc generate
   ```

3. **Iniciar el servidor:**
   ```bash
   go run .
   ```
