#!/bin/bash
set -euo pipefail

IMAGE_NAME=oauth2-server:dev
CONTAINER_NAME=oauth2-server-dev
KEYS_DIR=/app/keys
HOST_KEYS_DIR="$(pwd)/keys"

cd "$(dirname "$0")/.."

echo "🚀 Building Docker image: $IMAGE_NAME"
docker rmi -f "$IMAGE_NAME" > /dev/null 2>&1 || true
docker build -t "$IMAGE_NAME" .

echo "🧹 Removing existing container (if any)..."
docker rm -f "$CONTAINER_NAME" > /dev/null 2>&1 || true

echo "🐳 Running container '$CONTAINER_NAME' with volume and environment variables..."
docker run --rm --name "$CONTAINER_NAME" \
  -p 8080:8080 \
  -v "$HOST_KEYS_DIR":"$KEYS_DIR" \
  -e PATHS_PRIVATE_KEY="$KEYS_DIR/private.pem" \
  -e PATHS_PUBLIC_KEY="$KEYS_DIR/public.pem" \
  -e SERVER_PORT=8080 \
  -e JWT_TOKEN_TTL=3600 \
  -e API_TOKEN="/token" \
  -e API_INTROSPECT="/introspect" \
  -e API_JWKS=""/.well-known/jwks.json"" \
  "$IMAGE_NAME"

echo "✅ The server is now running on at http://localhost:8080"