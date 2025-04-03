// README.md
# Simple OAuth2 Server

## Features
- Client Credentials Grant (RFC 6749)
- JWT Access Tokens (RFC 7519)
- RS256 Signing
- JWKS endpoint (RFC 7517)
- Token Introspection (RFC 7662)

## Usage

### Start the Server
```sh
go run ./cmd/server
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
