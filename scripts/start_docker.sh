#!/bin/sh
set -e

echo "Starting Gogogo system..."

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$PROJECT_ROOT"

echo "Preparing volumes..."
mkdir -p scripts/gogogo-api/config
mkdir -p scripts/gogogo-api/storage
mkdir -p scripts/postgres

echo "Starting docker compose..."
docker compose -f docker/docker-compose.yml up -d --build

echo "Docker system started successfully"
