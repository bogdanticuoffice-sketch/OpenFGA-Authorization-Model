# Getting Started with OpenFGA

Complete tutorial for implementing OpenFGA authorization in your application.

## Prerequisites

- Docker and Docker Compose
- Git
- curl or Postman for API testing
- Basic understanding of authorization concepts

## Step 1: Start OpenFGA

```bash
cd deployment/docker
cp .env.example .env
# Edit .env and set a real POSTGRES_PASSWORD
docker compose up -d
```

Services running:
- OpenFGA HTTP API: http://localhost:8080
- OpenFGA gRPC: localhost:8081
- Playground (dev only): http://localhost:3000

The `migrate` service runs automatically before `openfga` starts and applies
the required database schema.

## Step 2: Create an Authorization Store

```bash
curl -s -X POST http://localhost:8080/stores \
  -H 'Content-Type: application/json' \
  -d '{"name": "my-store"}' | jq .
```

Response:
```json
{
  "id": "01HVND...",
  "name": "my-store",
  "created_at": "2026-02-21T10:00:00Z",
  "updated_at": "2026-02-21T10:00:00Z"
}
```

Save the store ID — you will need it for all subsequent requests.

## Step 3: Write an Authorization Model

The OpenFGA REST API accepts a JSON body with `schema_version` and
`type_definitions`. The `.fga` files in `models/` are DSL representations;
to write them via the API you must either:

- Use the [OpenFGA CLI](https://github.com/openfga/cli): `fga model write --store-id <id> --file models/rbac/model.fga`
- Or convert the DSL to JSON using the [openfga/language](https://github.com/openfga/language) transformer

Example using the CLI:

```bash
fga model write \
  --api-url http://localhost:8080 \
  --store-id 01HVND... \
  --file models/rbac/model.fga
```

Response:
```json
{
  "authorization_model_id": "01HVNE..."
}
```

## Step 4: Write Relationships (Tuples)

```bash
curl -s -X POST http://localhost:8080/stores/01HVND.../write \
  -H 'Content-Type: application/json' \
  -d '{
    "writes": {
      "tuple_keys": [
        {
          "user": "user:alice",
          "relation": "admin",
          "object": "organization:acme"
        },
        {
          "user": "user:bob",
          "relation": "member",
          "object": "organization:acme"
        }
      ]
    }
  }' | jq .
```

## Step 5: Check Authorization

```bash
curl -s -X POST http://localhost:8080/stores/01HVND.../check \
  -H 'Content-Type: application/json' \
  -d '{
    "tuple_key": {
      "user": "user:alice",
      "relation": "admin",
      "object": "organization:acme"
    }
  }' | jq .
```

Response:
```json
{
  "allowed": true
}
```

## Step 6: List Accessible Objects

```bash
curl -s -X POST http://localhost:8080/stores/01HVND.../list-objects \
  -H 'Content-Type: application/json' \
  -d '{
    "user": "user:alice",
    "relation": "admin",
    "type": "organization"
  }' | jq .
```

Response:
```json
{
  "objects": ["organization:acme"]
}
```

## Using the Playground

OpenFGA includes a visual playground for development. It is enabled by default
in the Docker Compose setup and **disabled in the Kubernetes deployment**.

1. Open http://localhost:3000
2. Create an authorization model
3. Add relationships
4. Run authorization checks
5. Visualize the permission graph

## Code Example — Go

See `examples/go/` for a working Go client. Run with:

```bash
cd examples/go
go mod tidy
go run main.go
```

The example demonstrates:
- Creating a store
- Writing an authorization model via TypeDefinitions
- Writing relationship tuples
- Checking access
- Listing accessible objects

## API Reference

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/stores` | Create a store |
| GET | `/stores/{id}` | Get store details |
| POST | `/stores/{id}/authorization-models` | Write a model |
| GET | `/stores/{id}/authorization-models/{model-id}` | Get a model |
| POST | `/stores/{id}/write` | Write or delete tuples |
| POST | `/stores/{id}/check` | Check authorization |
| POST | `/stores/{id}/expand` | Expand a relation |
| POST | `/stores/{id}/list-objects` | List objects a user can access |
| POST | `/stores/{id}/list-users` | List users with access to an object |

## Troubleshooting

### Authorization always returns false

1. Confirm the store ID in your request matches the created store
2. Verify the authorization model was written successfully
3. Check that the tuple (user, relation, object) exists — query `/read` to inspect stored tuples
4. Ensure user and object identifiers match exactly (case-sensitive)

### Model write fails with validation error

1. Check the DSL syntax against the [OpenFGA language spec](https://openfga.dev/docs/configuration-language)
2. Confirm all referenced types are defined in the model
3. Ensure `from` expressions reference relations that exist on the linked type

### Performance

1. Monitor database query times via the OpenFGA metrics endpoint (`/metrics`)
2. Use `ListObjects` sparingly on large datasets — it performs a graph traversal
3. Implement application-level caching for repeated checks on stable data
4. Use batch writes to reduce round-trips when writing many tuples at once

## Next Steps

1. Choose the appropriate model for your use case from `models/`
2. Integrate the Go SDK (or Python/Node.js SDK) into your application
3. Add audit logging by reading from `/stores/{id}/read-changes`
4. Test authorization rules thoroughly before production deployment

See `docs/ARCHITECTURE.md` for advanced topics.
