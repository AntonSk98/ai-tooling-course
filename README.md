# KI Entwicklungstools

Course materials on AI-powered development tools and methodologies.

## Structure

```
KI Entwicklungstools/
├── theory/          # Concepts and methodologies
└── petstore/        # Hands-on project (Go + Alpine.js)
```

## Theory

Foundational topics covered in order:

| # | Topic | Description |
|---|-------|-------------|
| 1 | [AI Agents](theory/1_agents.md) | Autonomous AI systems — planning, tool use, decision-making |
| 2 | [Model Context Protocol](theory/2_mcp.md) | Standardized protocol for connecting AI to external tools and data |
| 3 | [Spec-Driven Development](theory/3_spec-driven-development.md) | Design-first approach — OpenAPI spec as single source of truth |
| 4 | [BMAD Method](theory/4_bmad-method.md) | Structured AI-driven development across four phases |

## Practice

The [Petstore](petstore/) project applies these concepts in a full-stack application:

- **Go REST API** with clean architecture (Fiber + GORM + SQLite)
- **Alpine.js frontend** served as static files
- **Spec-driven**: models generated from OpenAPI spec via oapi-codegen
- **AI-assisted**: built and refactored using AI development tools

See [`petstore/README.md`](petstore/README.md) for setup and details.
