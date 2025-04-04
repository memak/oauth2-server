// README.md
# Simple OAuth2 Server

## Features
- Client Credentials Grant (RFC 6749)
- JWT Access Tokens (RFC 7519)
- RS256 Signing
- JWKS endpoint (RFC 7517)
- Token Introspection (RFC 7662)

## Important notes
Usually private/public keys would not be saved in the repo. This is only done here as an example. The keys are mounted to the Docker image during runtime and as a mounted secret in Kubernetes cluster as it would be done in production

## Usage
Either start a local Go server manually, or start a local Docker server, or start a local Kubernetes cluster including the Docker image with the following commands. After starting the server you can use the provided `curl` commands below

### Start the Server locally
```sh
go run ./cmd/server
```

### Start the Server with local Docker
Builds the docker image and starts it with mounting the keys directory
```sh
./scripts/dockerstart.sh
```

### Start the Server with local Kubernetes and Helm
Creates a secret for the keys and starts the kubernetes cluster. If you use `kind` (Kubernets in Docker) or `minikube` (own VM), you have to run `kind load docker-image oauth2-server:dev` or `minikube image load oauth2-server:dev` respectively. If you use Docker Desktop you do not have to do anything additionally. Check if you have the `helm` command, otherwise please install first.
```sh
./scripts/helmstart.sh
```
## Use the APIs
After starting the server you can use the following APIs

### Get Token
```sh
curl -u client_id:secret -X POST http://localhost:8080/token -d 'grant_type=client_credentials'
```

### Get JWKS
```sh
curl http://localhost:8080/jwks
```

### Introspect Token
```sh
curl -X POST http://localhost:8080/introspect -d 'token=ey...'
```

### Token expiration
If you want to test token expiration, you can change the token ttl value in .env for the local server, and in dockerstart.sh for the Docker image or in values.yaml for Kubernetes cluster

---
