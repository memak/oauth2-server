#!/bin/bash
set -e
cd "$(dirname "$0")/.."

echo "🚀 Building and starting Docker image for local testing ..."

docker rmi -f oauth2-server:dev || true

docker build -t oauth2-server:dev .

docker run -p 8080:8080 -v $(pwd)/keys:/app/keys -e PRIVATE_KEY_PATH=/app/keys/private.pem \
  -e PUBLIC_KEY_PATH=/app/keys/public.pem oauth2-server:dev
