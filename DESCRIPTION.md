# About This Repository

This is a complete production-ready authorization system built with OpenFGA. If you're looking to implement fine-grained access control in your application, this repository gives you everything you need to get started properly.

## What's Inside

Four complete authorization models covering different real-world scenarios:

- RBAC: Traditional role-based access for teams and organizations
- ABAC: Attribute-based control when roles aren't enough
- API: Protecting API endpoints with scopes and tokens
- SaaS: Multi-tenant patterns with complete tenant isolation

Each model comes with example data, documentation, and ready-to-run code.

## Why Use This

Authorization is hard to get right. Most projects either oversimplify (everyone's an admin) or overcomplicate (tangled permission logic everywhere). OpenFGA lets you define clear, auditable permissions that scale as your product grows.

This repository saves you months of design and implementation work by showing you how real systems handle authorization.

## Quick Start

1. Start the local environment with Docker: `docker-compose -f deployment/docker/docker-compose.yml up`
2. Pick a model from the models/ folder that matches your use case
3. Load the model and example relationships
4. Check the examples/ folder to see how to integrate with your app

Takes about 10 minutes to go from zero to working authorization.

## Security Focused

Every pattern here follows security best practices:
- Least privilege by default
- Complete audit trails
- Tenant isolation you can trust
- No hardcoded permissions
- Production-ready deployment configs

## Real World Examples

The models handle things you'll actually encounter:
- Users in multiple teams
- Nested permissions (admins can do what members can do, plus more)
- Time-based or context-based access
- API tokens with specific scopes
- Complete multi-tenant customer isolation

## Not Just Theory

This isn't academic. The patterns come from real production systems handling thousands of users and complex permission requirements.

## Perfect For

- Building the authorization layer for a new product
- Refactoring tangled permission logic in existing code
- Migrating from simple role-based to more flexible systems
- Learning how real authorization systems work

## Getting Help

Each model folder has detailed documentation. The docs/ folder has architecture explanations, security guidelines, and integration patterns. The examples/ folder shows Go, Python, and Node.js implementations.

Start with GETTING_STARTED.md if you're new to OpenFGA.
