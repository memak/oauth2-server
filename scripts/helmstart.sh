#!/bin/bash

set -euo pipefail  # Stop the script on any error

RELEASE_NAME=oauth2-server
IMAGE_NAME=${RELEASE_NAME}:dev
CHART_PATH=./oauth2-helm
NAMESPACE=default
PORT_LOCAL=8080
PORT_TARGET=80

echo "🔄 Cleaning up previous deployment (if any)..."
helm uninstall $RELEASE_NAME || true
kubectl delete deployment $RELEASE_NAME -n $NAMESPACE --ignore-not-found
kubectl delete service $RELEASE_NAME -n $NAMESPACE --ignore-not-found

echo "🚀 Building Docker image: $IMAGE_NAME"
docker rmi -f "$IMAGE_NAME" > /dev/null 2>&1 || true
docker build -t "$IMAGE_NAME" .

if ! kubectl get secret oauth2-keys -n $NAMESPACE > /dev/null 2>&1; then
  echo "🔑 Creating secret oauth2-keys..."
  kubectl create secret generic oauth2-keys \
    --from-file=private.pem=keys/private.pem \
    --from-file=public.pem=keys/public.pem
else
  echo "🔑 Secret oauth2-keys already exists. Skipping creation."
fi

echo "🚀 Installing Helm release..."
helm install $RELEASE_NAME $CHART_PATH

echo "⏳ Waiting up to 10s for pod to be ready..."
if kubectl wait --for=condition=ready pod -l app.kubernetes.io/instance=$RELEASE_NAME -n $NAMESPACE --timeout=10s; then
  echo "✅ Pod is ready!"
else
  echo "⚠️ Pod not ready in time."
fi

echo "🌐 Port forwarding service $RELEASE_NAME on http://localhost:$PORT_LOCAL"
kubectl port-forward service/$RELEASE_NAME $PORT_LOCAL:$PORT_TARGET