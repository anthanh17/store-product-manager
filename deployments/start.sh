#!/bin/bash
set -e

echo "Waiting for database to be ready..."
until pg_isready -h postgres -p 5432 -U root; do
  sleep 1
done

echo "Running database migrations..."
migrate -path /app/migration -database "postgresql://root:secret@postgres:5432/bfast?sslmode=disable" up

echo "Starting the application..."
exec /app/main
