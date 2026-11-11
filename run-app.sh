#!/usr/bin/env bash
set -Eeuo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")"

# 1. Verificar o crear .env a partir de .env.example
if [[ ! -f .env ]]; then
    if [[ ! -f .env.example ]]; then
        echo "ERROR: no existe .env.example para crear .env." >&2
        exit 1
    fi
    cp .env.example .env
    echo "--> Se creó .env desde .env.example."
fi

set -a
source .env
set +a

# 3. Tareas posteriores: limpieza garantizada de contenedores y volúmenes al salir
cleanup() {
    echo ""
    echo "--> Tareas posteriores: deteniendo contenedores y borrando volúmenes..."
    docker compose down --volumes --remove-orphans || true
}
trap cleanup EXIT

# 1. Tareas previas
echo "--> Limpiando contenedores y volúmenes previos..."
docker compose down --volumes --remove-orphans || true

echo "--> Construyendo imagen de la aplicación..."
docker compose build app

echo "--> Levantando contenedor de PostgreSQL..."
docker compose up -d postgres

echo "--> Esperando a que PostgreSQL esté listo..."
conecto=false
for attempt in $(seq 1 60); do
    if docker compose exec -T postgres pg_isready -U "${DB_USER:-postgres}" -d "${DB_NAME:-biblioteca}" >/dev/null 2>&1; then
        echo "--> PostgreSQL está listo y aceptando conexiones."
        conecto=true
        break
    fi
    sleep 1
done

if [[ "$conecto" != true ]]; then
    echo "ERROR: PostgreSQL no estuvo listo a tiempo." >&2
    exit 1
fi

echo "--> Ejecutando sqlc generate, templ generate y compilando la aplicación..."
docker compose run --rm --no-deps app sh -c "sqlc generate && templ generate && go build ./..."

# 2. Ejecutar la aplicación
echo "--> Iniciando servidor Go en http://localhost:8080 (presiona Ctrl+C para detener y limpiar)..."
docker compose up app
