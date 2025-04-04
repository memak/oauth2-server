#!/bin/bash

set -euo pipefail  # Stop the script on any error

RELEASE_NAME=oauth2-server
CHART_PATH=./oauth2-helm
NAMESPACE=default
PORT_LOCAL=8080
PORT_TARGET=80

echo "🔄 Cleaning up previous deployment (if any)..."
helm uninstall $RELEASE_NAME || true
kubectl delete deployment $RELEASE_NAME -n $NAMESPACE --ignore-not-found
kubectl delete service $RELEASE_NAME -n $NAMESPACE --ignore-not-found

echo "🚀 Installing Helm release..."
helm install $RELEASE_NAME $CHART_PATH

echo "⏳ Waiting for pod to be ready..."
if kubectl wait --for=condition=ready pod -l app.kubernetes.io/instance=$RELEASE_NAME -n $NAMESPACE --timeout=20s; then
  echo "✅ Pod is ready!"
else
  echo "⚠️  Pod not ready in time. You can check with: kubectl get pods"
  exit 1
fi

echo "🌐 Port forwarding service $RELEASE_NAME on http://localhost:$PORT_LOCAL"
kubectl port-forward service/$RELEASE_NAME $PORT_LOCAL:$PORT_TARGET