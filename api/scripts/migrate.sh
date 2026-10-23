#!/usr/bin/env sh
set -e

# Run migrations against PG_URL or DATABASE_URL
DB_URL="${DATABASE_URL:-$PG_URL}"

if [ -z "$DB_URL" ]; then
  echo "Error: Neither DATABASE_URL nor PG_URL is set."
  echo "Usage: DATABASE_URL=\"postgres://...\" ./migrate.sh"
  exit 1
fi

echo "Checking for goose..."
if ! command -v goose >/dev/null 2>&1; then
  echo "Installing goose..."
  go install github.com/pressly/goose/v3/cmd/goose@latest
  export PATH="$PATH:$(go env GOPATH)/bin"
fi

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
SCHEMA_DIR="$SCRIPT_DIR/../internal/sql/schema"

echo "Applying migrations from $SCHEMA_DIR..."
goose -dir "$SCHEMA_DIR" postgres "$DB_URL" up

echo "Migrations completed successfully."
