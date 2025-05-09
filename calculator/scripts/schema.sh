#!/usr/bin/env bash
set -euo pipefail

PROJECT_ROOT="$(dirname "$(realpath "$0")")/../../"

start_temp_db() {
  echo "Starting a temp DB..."
  docker run -d --rm --name temp_db -p 5433:5432 -e POSTGRES_PASSWORD=password postgres > /dev/null 2>&1
  echo "Waiting for the DB to become ready..."
  until docker exec temp_db pg_isready -U postgres > /dev/null 2>&1; do
    sleep 0.5
  done
}

cleanup() {
  echo "Cleaning up..."
  docker stop temp_db > /dev/null 2>&1 || true
}
trap cleanup EXIT INT TERM

start_temp_db

echo "Applying the migrations..."
migrate -database "postgres://postgres:password@localhost:5433/postgres?sslmode=disable" -path internal/storage/migrations up

echo "Starting schema generation..."
export PGPASSWORD=password
npx pg-mermaid \
--host localhost \
--port 5433 \
--dbname postgres \
--username postgres \
--excluded-tables schema_migrations \
--output-path "$PROJECT_ROOT/database.md"
