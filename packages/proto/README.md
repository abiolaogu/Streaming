# Protocol Buffers Package

gRPC service definitions and message types for StreamVerse.

## Services

- `ContentService` - Content management operations
- `AuthService` - Authentication operations

## Usage

### Generate Go Code

```bash
make generate
```

### Generate TypeScript Code

```bash
make generate-ts
```

## Files

- `content.proto` - Content service definitions
- `auth.proto` - Auth service definitions

## Dependencies

- protoc compiler
- protoc-gen-go
- protoc-gen-go-grpc
- protoc-gen-ts (for TypeScript)

