cd "$(dirname "$0")/.."

echo "🚀 Creating Kubernetes secret for OAuth2 keys ..."

kubectl create secret generic oauth2-keys \
  --from-file=private.pem=keys/private.pem \
  --from-file=public.pem=keys/public.pem