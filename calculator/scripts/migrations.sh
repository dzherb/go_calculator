#!/usr/bin/env bash
set -e

GO_PROJECT_ROOT="$(dirname "$(realpath "$0")")/../"
MIGRATIONS_DIR="$GO_PROJECT_ROOT/internal/storage/migrations"

read -rp "Enter migration name (use lowercase and underscores): " NAME

if [[ -z "$NAME" ]]; then
  echo "Migration name cannot be empty"
  exit 1
fi

# Replace spaces with underscores (optional)
NAME="${NAME// /_}"

echo
migrate create -ext sql -dir "$MIGRATIONS_DIR" -seq "$NAME"

echo
echo "Created new migration: $NAME"