#!/usr/bin/env bash
set -Eeuo pipefail

# Limpieza automática de contenedores y volúmenes al salir (éxito, error o Ctrl+C)
trap 'docker compose down --volumes --remove-orphans' EXIT

# Tareas previas
echo "--> Limpiando contenedores y volúmenes previos..."
docker compose down --volumes --remove-orphans

echo "--> Ejecutando sqlc generate..."
sqlc generate

echo "--> Compilando código..."
go build ./...

echo "--> Levantando base de datos PostgreSQL..."
docker compose up -d --remove-orphans

if [ ! -f ".env" ]; then
    archivo=".env.example"
else
    archivo=".env"
fi

echo "Usando $archivo"

# Cargar las variables del archivo seleccionado en la sesión actual
if [ -f "$archivo" ]; then
    set -a
    source "$archivo"
    set +a
fi

echo "--> Esperando a que PostgreSQL esté listo..."
for attempt in $(seq 1 60); do
    if docker compose exec -T postgres sh -c "pg_isready -U \"${POSTGRES_USER:-postgres}\" -d \"${POSTGRES_DB:-tp2_db}\"" >/dev/null 2>&1; then
        echo "PostgreSQL está listo."
        break
    fi
    sleep 1
done

# Ejecutar tests
echo "--> Ejecutando tests..."
go test -v ./test/...
