# OpenFGA Authorization Model

Production-ready authorization models and implementation guide for RBAC, ABAC, API, and SaaS applications. Everything you need to build fine-grained access control into your system.

## What is This

This repository contains complete, working authorization patterns built on OpenFGA. If you're building an application and need to control who can do what, this gives you tested patterns and code to start from.

OpenFGA is built on Google's Zanzibar architecture. It lets you define permissions as relationships between users, actions, and resources. Simple to understand, scales to millions of relationships, and auditable.

## What You Get

Four complete authorization models covering real scenarios:

RBAC - Role-Based Access Control
- Users get assigned roles (admin, member, viewer)
- Roles define what you can do
- Simple and straightforward
- Good for teams and organizations
- Works when your permission structure is predictable

ABAC - Attribute-Based Access Control
- Decisions based on attributes (department, clearance level, resource type)
- Flexible and powerful
- Handles complex policies
- Scales with your organization
- Harder to understand and maintain

API - API Authorization
- Protect endpoints with API keys and scopes
- OAuth2-like token model
- Rate limiting per key
- Service-to-service authentication
- Production patterns for API security

SaaS - Multi-Tenant Authorization
- Complete isolation between customers
- Nested organizations, workspaces, projects
- Invite system with role assignment
- Audit logging everything
- Designed for product companies

Each model includes:
- The authorization model definition
- Example relationships showing how it works
- Detailed documentation
- Integration examples in Go, Python, Node.js

## Getting Started

Start local in 5 minutes:

```bash
cd deployment/docker
docker-compose up -d
```

This runs:
- OpenFGA server on port 8080
- Postgres database
- Playground UI on port 3000

Then load a model:

```bash
cd models/rbac
openfga model write --api-url http://localhost:8080 model.fga
```

Add example relationships:

```bash
openfga relationship write --api-url http://localhost:8080 relations.txt
```

Check if someone has permission:

```bash
openfga check --api-url http://localhost:8080 \
  --user "user:alice" \
  --relation "admin" \
  --object "organization:acme"
```

Returns: true or false

That's it. You have a working authorization system.

## How It Works

Authorization is about relationships. Alice is an admin of ACME. Bob is a member of the engineering team. The engineering team has access to the API project.

You define these relationships in OpenFGA. Then you ask: Can Bob access the API project? OpenFGA figures it out by following the relationships.

Three parts:
1. Define your model (what types of things exist, what relationships are valid)
2. Add relationships (alice is admin of acme)
3. Check (can bob view this document?)

## Directory Layout

models/ - Authorization model definitions
- rbac/ - Role-based access control example
- abac/ - Attribute-based access control example
- api/ - API key and scope authorization
- saas/ - Multi-tenant SaaS patterns

Each model includes:
- model.fga - The permission model
- relations.txt - Example relationships
- README.md - How the model works

examples/ - How to integrate with your code
- go/ - Go client example
- python/ - Python client example
- nodejs/ - Node.js client example

deployment/ - How to run in production
- docker/ - Local development setup
- kubernetes/ - Production Kubernetes deployment

docs/ - Deep dives
- GETTING_STARTED.md - Step-by-step tutorial
- ARCHITECTURE.md - How OpenFGA works inside
- BEST_PRACTICES.md - Security and design patterns
- DESCRIPTION.md - Why use this repository

## Picking a Model

Start with RBAC if:
- You have a traditional org structure
- Users have defined roles
- Permissions are mostly static
- Your team is comfortable with role-based thinking

Use ABAC if:
- Roles aren't flexible enough
- You need dynamic conditions
- Users have multiple attributes
- Policies change frequently

Use API if:
- You're building an API
- You need token-based auth
- Scopes and permissions matter
- Service-to-service access is important

Use SaaS if:
- You're building a multi-customer product
- Complete customer isolation is required
- You have complex org hierarchies
- You manage multiple workspaces

Most apps use RBAC or API. Start there.

## Integration Patterns

Typical flow in your application:

```go
// User makes request
user := extractUser(request)
resource := getResource(request)

// Check with OpenFGA
allowed, err := fga.Check(ctx, user, "edit", resource)

// Allow or deny
if !allowed {
  return forbidden()
}

performOperation(resource)
```

For APIs:

```go
// Extract token from header
token := request.Header.Get("Authorization")

// Validate and extract scopes
scopes := validateToken(token)

// Check scope includes required permission
if !hasScope(scopes, "users:write") {
  return unauthorized()
}

// Process request
handleRequest(request)
```

See examples/ for complete implementations.

## Security

This is built with security in mind.

Key practices included:
- Least privilege by default
- Complete audit trail
- Tenant isolation that actually works
- No permission logic in code
- Separate credentials per environment

Read docs/BEST_PRACTICES.md and docs/SECURITY.md for detailed guidance.

The patterns here are used in production systems. Not theoretical stuff.

## Production Deployment

Docker Compose for local development. Kubernetes configs in deployment/kubernetes/.

For production:
- Use a managed Postgres instance
- Run OpenFGA behind a load balancer
- Set up monitoring and alerts
- Use TLS for all communication
- Implement audit logging
- Regular backups

See deployment/kubernetes/README.md for setup details.

## Testing Your Model

Each model includes test cases showing what should and shouldn't be allowed.

Test your changes:
```bash
cd tests
openfga test test_rbac.yaml
openfga test test_api.yaml
```

Write your own tests in the same format.

## Performance

Authorization checks should be fast.

Typical numbers:
- Check latency: 2-5ms with local database
- Throughput: 10,000+ checks per second
- With caching: 100,000+ checks per second

Database performance matters most. Add indexes on frequently queried relations.

## Troubleshooting

Authorization returning false when it should be true?

1. Check the relationship exists in database
2. Verify the model definition is correct
3. Confirm user and object identifiers match exactly
4. Use the expand operation to see the permission tree

Slow checks?

1. Monitor database query time
2. Add indexes
3. Implement caching (5-15 minute TTL)
4. Use batch operations

See docs/GETTING_STARTED.md for detailed troubleshooting.

## Migrating from Existing Authorization

From LDAP or Active Directory?
- Map groups to roles
- Import relationships
- Run both systems in parallel
- Cutover when confident

From simple role checks in code?
- Define permission model
- Load existing permissions
- Update code to call OpenFGA
- Gradually migrate endpoints

From OAuth2 scopes?
- Map scopes to OpenFGA permissions
- Use API authorization model
- Update token validation

See migrations/ for examples.

## Learning More

Start here:
- GETTING_STARTED.md - Step-by-step tutorial
- models/rbac/README.md - Understanding role-based access
- docs/ARCHITECTURE.md - How OpenFGA works

Then dive deeper:
- docs/BEST_PRACTICES.md - Design patterns
- docs/SECURITY.md - Security hardening
- examples/ - Real code

OpenFGA docs: https://openfga.dev/docs
Zanzibar paper: https://research.google/pubs/zanzibar/

## Using This Repository

Each model is independent. You don't need all four. Pick one that matches your use case and go from there.

Models are starting points. Adapt them to your needs. The patterns work but your authorization might have unique requirements.

Test thoroughly before going to production. Authorization is security-critical. Run through edge cases and permission combinations.

## Questions

Each model folder has a README with details.
Examples show how to code against it.
Docs explain the thinking behind the patterns.

Start with the Quick Start section above. Get something running locally first. Then read the model docs for your use case.

## Contributing

These patterns came from real systems. If you have improvements, changes, or new patterns, contributions welcome.

Just keep it production-focused. Things that work in theory but not in practice don't belong here.

## License

Use as you like. No restrictions.
