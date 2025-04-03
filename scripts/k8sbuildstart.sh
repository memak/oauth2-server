cd "$(dirname "$0")/.."

echo "🚀 Starting Kubernetes cluster with  oauth2-server ..."

kubectl create secret generic oauth2-keys \
  --from-file=private.pem=keys/private.pem \
  --from-file=public.pem=keys/public.pem

kubectl apply -f manifests/deployment.yaml
kubectl port-forward service/oauth2-server 8080:80