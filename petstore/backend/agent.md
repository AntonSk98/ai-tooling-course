# Petstore Backend - Go REST API

## Overview

REST API backend for the Petstore application:
- **Go 1.22+** with **Fiber v2.52.11** web framework
- **GORM v1.25.7** ORM with **SQLite** (pure Go, no CGO)
- **oapi-codegen v2.5.1** for model generation from OpenAPI spec
- **Swagger UI** via swaggo/swag for interactive API docs

## Architecture

Clean 3-layer architecture with strict separation of concerns:

```
HTTP Request
    ↓
Handler (internal/api/handlers/)  → Parse HTTP, return JSON. No business logic.
    ↓
Service (internal/domain/)        → Validation, business rules. HTTP-independent.
    ↓
Repository (internal/repository/) → GORM CRUD. Model ↔ Entity conversion.
    ↓
Database (SQLite with GORM)
```

## File Structure

```
backend/
├── cmd/api/main.go              → Entry point, wires everything together
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   │   ├── pet.go           → 6 pet endpoints with Swagger annotations
│   │   │   └── category.go     → 4 category endpoints with Swagger annotations
│   │   └── router.go           → Route registration + middleware (CORS, logger, recover)
│   ├── models/generated.go     → Auto-generated from OpenAPI spec (DO NOT EDIT)
│   ├── models/requests.go      → Additional request types
│   ├── domain/
│   │   ├── pet.go              → Pet business logic + validation
│   │   └── category.go         → Category business logic
│   ├── repository/
│   │   ├── pet.go              → Pet CRUD (GORM)
│   │   └── category.go         → Category CRUD (GORM)
│   └── database/
│       ├── connection.go       → DB connection, migration, close
│       └── entities.go         → GORM entities (PetEntity, CategoryEntity)
├── pkg/config/config.go        → Config loader (PORT, DB_PATH from env)
├── tests/integration/          → Integration tests (in-memory SQLite)
│   ├── pet_test.go
│   ├── category_test.go
│   └── helpers.go
├── docs/                       → Auto-generated Swagger docs
├── petstore_swagger.yml        → OpenAPI 3.0 spec (source of truth)
├── oapi-codegen.yaml           → Code generation config
├── Makefile                    → Build/test/generate tasks
├── go.mod                      → Go dependencies
└── agent.md                    → This documentation
```

## API Endpoints

### Pet (6 endpoints)

| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| POST | `/pet` | `AddPet` | Create new pet |
| PUT | `/pet` | `UpdatePet` | Update existing pet |
| GET | `/pet/findByStatus?status=` | `FindPetsByStatus` | Find pets by status |
| GET | `/pet/:petId` | `GetPetByID` | Get pet by ID |
| POST | `/pet/:petId` | `UpdatePetWithForm` | Update pet with form data |
| DELETE | `/pet/:petId` | `DeletePet` | Delete pet |

### Category (4 endpoints)

| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| POST | `/category` | `AddCategory` | Create new category |
| PUT | `/category` | `UpdateCategory` | Update existing category |
| GET | `/category/listAll` | `GetAllCategories` | Get all categories |
| DELETE | `/category/:categoryId` | `DeleteCategory` | Delete category |

### Other

- `GET /health` — Health check
- `GET /swagger/*` — Swagger UI

**Note**: DELETE /category exists in code but NOT in `petstore_swagger.yml` (added outside design-first workflow).

## Key Patterns

### Model vs Entity

- **API Models** (`models/generated.go`): `Id` (lowercase d), `*PetStatus` enum, pointer fields — match API contract
- **DB Entities** (`database/entities.go`): `ID` (uppercase), string status, manual fields (no `gorm.Model` embed), timestamps + soft delete
- **Conversion** happens in repository layer: `toModel()` / entity creation

### Error Responses

All errors return JSON: `{"error": "message"}`
- 400: Bad request / invalid input
- 404: Not found
- 422: Validation failed
- 500: Internal error

### Validation Rules

| Entity | Rules |
|--------|-------|
| Pet | Name required, status ∈ {available, pending, sold}, ID required for updates |
| Category | Name required, name must be unique, ID required for updates |

### Design-First Approach

`petstore_swagger.yml` is the single source of truth:
1. Edit OpenAPI spec
2. Regenerate models: `oapi-codegen -config oapi-codegen.yaml petstore_swagger.yml`
3. Update repository/domain/handlers as needed
4. Regenerate Swagger docs: `swag init -g cmd/api/main.go --output docs --parseDependency --parseInternal`

## Database

- **SQLite** with pure Go driver (`github.com/glebarez/sqlite`) — no CGO required
- **Schema auto-migrates** on startup
- **Soft delete**: records marked with `deleted_at`, never physically removed
- **Tables**: `pets` (with `category_id` FK) and `categories`

## Testing

```bash
cd backend
go test ./tests/integration/... -v
```

- Each test gets fresh in-memory SQLite via `setupTestApp(t)`
- Covers: happy path, validation errors, 404s, edge cases

## Running

```bash
cd backend
go mod download
go run cmd/api/main.go
```

Server starts on `http://localhost:3000`

## Makefile Commands

```bash
make run               # Run the application
make test              # Run integration tests
make models            # Regenerate models from OpenAPI spec
make swagger           # Regenerate Swagger docs
make generate          # Regenerate both
make build             # Build binary
make verify            # Build + test
make install-tools     # Install oapi-codegen, swag, air
```

## Dependencies

```
github.com/gofiber/fiber/v2 v2.52.11
github.com/gofiber/swagger v1.1.1
gorm.io/gorm v1.25.7
github.com/glebarez/sqlite v1.11.0  # Pure Go, no CGO
```

## Guidelines

- Always refer to `petstore_swagger.yml` for field names and types
- Generated models use `Id` not `ID` — this is oapi-codegen convention
- Use `string(*pet.Status)` when saving PetStatus enum to DB
- Handle errors explicitly, never ignore them
- Add Swagger annotations to all new handlers
- Write integration tests for new endpoints
