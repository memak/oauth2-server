// README.md
# Simple OAuth2 Server

## Features
- Client Credentials Grant (RFC 6749)
- JWT Access Tokens (RFC 7519)
- RS256 Signing
- JWKS endpoint (RFC 7517)
- Token Introspection (RFC 7662)

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

### Start the Server with local Kubernetes
Creates a secret for the keys and starts the kubernetes cluster. If you use `kind` (Kubernets in Docker) or `minikube` (own VM), you have to run `kind load docker-image oauth2-server:dev` or `minikube image load oauth2-server:latest` respectively. If you use Docker Desktop you  
```sh
./scripts/k8sbuild.sh
```

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

---
