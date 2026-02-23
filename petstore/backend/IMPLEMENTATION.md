# Petstore API - Implementation Guide

This document explains the implementation of the Petstore REST API built with Go, Fiber, and SQLite.

## Architecture Overview

The project follows **Clean Architecture** principles with clear separation of concerns:

```
HTTP Request
    ↓
Handler (API Layer) - Handles HTTP requests/responses
    ↓
Service (Domain Layer) - Business logic and validation
    ↓
Repository (Data Layer) - Database operations
    ↓
Database (SQLite with GORM)
```

## Technology Stack

- **Go 1.22+** - Programming language
- **Fiber v2.52.11** - Fast Express-inspired web framework
- **GORM v1.25.7** - Go ORM for database operations
- **SQLite** (glebarez/sqlite) - Pure Go database driver (no CGO)
- **oapi-codegen v2.5.1** - Generate Go models from OpenAPI spec
- **swaggo/swag v1.16.6** - Generate Swagger documentation from code annotations
- **gofiber/swagger v1.1.1** - Serve Swagger UI in Fiber
- **Alpine.js 3.x** - Frontend reactivity and state management (CSR)
- **Tailwind CSS 3.x** - Utility-first CSS framework (CDN)

## Project Structure Explained

```
petstore/
├── backend/                         # Go backend
│   ├── cmd/api/main.go              # Application entry point - wires everything together
│   ├── docs/                        # Auto-generated Swagger documentation
│   │   ├── docs.go                  # Swagger docs Go file
│   │   ├── swagger.json             # Swagger JSON spec
│   │   └── swagger.yaml             # Swagger YAML spec
│   ├── internal/
│   │   ├── models/generated.go      # Auto-generated from OpenAPI spec
│   │   ├── database/
│   │   │   ├── connection.go        # DB setup and auto-migration
│   │   │   └── entities.go          # GORM entity models (database schema)
│   │   ├── repository/              # Data access layer
│   │   │   ├── pet.go               # Pet CRUD operations
│   │   │   └── category.go          # Category CRUD operations
│   │   ├── domain/                  # Business logic layer
│   │   │   ├── pet.go               # Pet service with validation
│   │   │   └── category.go          # Category service
│   │   └── api/
│   │       ├── handlers/            # HTTP handlers (controllers)
│   │       │   ├── pet.go           # 6 pet endpoints with Swagger annotations
│   │       │   └── category.go      # 4 category endpoints with Swagger annotations
│   │       └── router.go            # Route registration + Swagger UI route
│   ├── pkg/
│   │   └── config/
│   │       └── config.go            # Configuration (Port, DBPath from env vars)
│   ├── tests/
│   │   └── integration/             # Integration tests for all endpoints
│   ├── petstore_swagger.yml         # OpenAPI 3.0 specification (source of truth)
│   ├── oapi-codegen.yaml            # Code generation configuration
│   ├── Makefile                     # Build/test/generate tasks
│   ├── go.mod                       # Go dependencies
│   ├── .cursorrules                 # AI assistant instructions
│   └── IMPLEMENTATION.md            # This file — detailed implementation guide
├── ui/                              # Frontend (Client-Side Rendering)
│   ├── index.html                   # Pet management (Alpine.js + fetch)
│   ├── categories.html              # Category management (Alpine.js + fetch)
│   └── agent.md                     # Frontend architecture documentation
└── README.md
```

## Key Design Decisions

### 1. Design-First Approach

**OpenAPI spec is the single source of truth.**

- API contracts defined in `petstore_swagger.yml`
- Models generated automatically with `oapi-codegen`
- Ensures API documentation always matches implementation

**Regenerating models when spec changes:**
```bash
cd backend
oapi-codegen -config oapi-codegen.yaml petstore_swagger.yml
```

### 2. Three-Layer Architecture

#### Layer 1: HTTP Handlers (`internal/api/handlers/`)
- Parses HTTP requests
- Validates basic input format
- Calls domain services
- Formats HTTP responses
- **Never contains business logic**

Example:
```go
func (h *PetHandler) AddPet(c *fiber.Ctx) error {
    var pet models.Pet
    if err := c.BodyParser(&pet); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
    }
    
    createdPet, err := h.service.CreatePet(&pet)
    if err != nil {
        return c.Status(422).JSON(fiber.Map{"error": err.Error()})
    }
    
    return c.Status(200).JSON(createdPet)
}
```

#### Layer 2: Domain Services (`internal/domain/`)
- Business logic and validation
- Validates business rules (e.g., status must be valid enum)
- Orchestrates repository operations
- **Independent of HTTP layer**

Example:
```go
func (s *PetService) CreatePet(pet *models.Pet) (*models.Pet, error) {
    if pet.Name == "" {
        return nil, errors.New("pet name is required")
    }
    
    if pet.Status != nil {
        if err := s.validateStatus(string(*pet.Status)); err != nil {
            return nil, err
        }
    }
    
    return s.repo.Create(pet)
}
```

#### Layer 3: Repository (`internal/repository/`)
- Database operations (CRUD)
- Converts between API models and database entities
- Handles GORM queries
- **No business logic**

Example:
```go
func (r *PetRepository) Create(pet *models.Pet) (*models.Pet, error) {
    entity := &database.PetEntity{
        Name:   pet.Name,
        Status: string(*pet.Status),  // Convert enum to string
    }
    
    if err := r.db.Create(entity).Error; err != nil {
        return nil, err
    }
    
    return r.toModel(entity), nil  // Convert entity back to API model
}
```

### 3. Model vs Entity Pattern

**API Models** (`internal/models/generated.go`):
- Generated from OpenAPI spec
- Used in HTTP layer and domain layer
- Match API contract exactly
- Use `Id` field naming (oapi-codegen convention)
- Enums like `PetStatus` are typed

**Database Entities** (`internal/database/entities.go`):
- GORM models with database tags and manual field definitions
- Use `ID` field naming (GORM convention)
- Enums stored as strings in database
- Include timestamps (CreatedAt, UpdatedAt) and soft delete (DeletedAt) as explicit fields

**Example conversion:**
```go
// API Model (from generated.go)
type Pet struct {
    Id       *int64      `json:"id,omitempty"`
    Name     string      `json:"name"`
    Status   *PetStatus  `json:"status,omitempty"`  // Enum type
    Category *Category   `json:"category,omitempty"`
}

// Database Entity (entities.go)
type PetEntity struct {
    ID         uint           `gorm:"primaryKey;autoIncrement"`
    Name       string         `gorm:"not null;size:100"`
    Status     string         `gorm:"size:20;default:'available'"`
    CategoryID *uint          `gorm:"index"`
    Category   *CategoryEntity `gorm:"foreignKey:CategoryID"`
    CreatedAt  time.Time
    UpdatedAt  time.Time
    DeletedAt  gorm.DeletedAt `gorm:"index"`
}

// Conversion in repository
func (r *PetRepository) toModel(entity *PetEntity) *models.Pet {
    id := int64(entity.ID)
    status := models.PetStatus(entity.Status)  // String -> Enum
    
    return &models.Pet{
        Id:     &id,    // uint -> int64
        Name:   entity.Name,
        Status: &status,
    }
}
```

### 4. Pure Go SQLite Driver

Using `github.com/glebarez/sqlite` instead of `mattn/go-sqlite3` because:
- ✅ No CGO required (works on Windows without GCC)
- ✅ Cross-platform compatibility
- ✅ Easier to build and deploy
- ✅ Same GORM interface

## API Endpoints

### Pet Endpoints (6)

| Method | Path | Description | Handler |
|--------|------|-------------|---------|
| POST | `/pet` | Create new pet | `AddPet` |
| PUT | `/pet` | Update existing pet | `UpdatePet` |
| GET | `/pet/:id` | Get pet by ID | `GetPetByID` |
| GET | `/pet/findByStatus` | Find pets by status | `FindPetsByStatus` |
| POST | `/pet/:id` | Update pet with form | `UpdatePetWithForm` |
| DELETE | `/pet/:id` | Delete pet | `DeletePet` |

### Category Endpoints (4)

| Method | Path | Description | Handler |
|--------|------|-------------|----------|
| POST | `/category` | Create new category | `AddCategory` |
| PUT | `/category` | Update existing category | `UpdateCategory` |
| GET | `/category/listAll` | Get all categories | `GetAllCategories` |
| DELETE | `/category/:categoryId` | Delete category | `DeleteCategory` |

**Note**: DELETE /category was added outside the design-first workflow and is not in `petstore_swagger.yml`.

## Database Schema

### Pets Table
```sql
CREATE TABLE pets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    status TEXT DEFAULT 'available',
    category_id INTEGER,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    FOREIGN KEY (category_id) REFERENCES categories(id)
);
```

### Categories Table
```sql
CREATE TABLE categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME
);
```

**Soft Delete**: Records are never actually deleted, just marked with `deleted_at` timestamp. Entities define fields manually (ID, CreatedAt, UpdatedAt, DeletedAt) rather than embedding `gorm.Model`.

## Validation Rules

### Pet Validation
- ✅ Name is required (non-empty string)
- ✅ Status must be: "available", "pending", or "sold"
- ✅ ID required for updates

### Category Validation
- ✅ Name is required (non-empty string)
- ✅ Name must be unique
- ✅ ID required for updates

## Error Handling

HTTP status codes used:

- **200 OK** - Successful operation
- **400 Bad Request** - Invalid JSON or malformed request
- **404 Not Found** - Resource doesn't exist
- **422 Unprocessable Entity** - Validation failed
- **500 Internal Server Error** - Database or server error

Error response format:
```json
{
  "error": "Description of what went wrong"
}
```

## Testing Strategy

### Integration Tests (`tests/integration/`)

Every endpoint has integration tests covering:
1. ✅ Happy path (successful operation)
2. ✅ Validation errors (missing/invalid fields)
3. ✅ Not found errors (non-existent IDs)
4. ✅ Edge cases

**Test isolation**: Each test gets a fresh in-memory database via `setupTestApp(t)`

**Running tests:**
```bash
go test ./tests/integration/... -v
```

## Running the Application

### Start the server:
```bash
cd backend
go run cmd/api/main.go
```

Server starts on `http://localhost:3000`

### Environment:
- Default database: `petstore.db` (created automatically)
- Database auto-migrates schema on startup
- CORS enabled for all origins (development mode)

### Running tests:
```bash
cd backend
go test ./tests/integration/... -v
```

## Development Workflow

### Using the Makefile

The project includes a comprehensive Makefile in `backend/` for common tasks:

```bash
# Show all available commands
make help

# Development
make run               # Run the application
make dev               # Run with live reload (requires air)
make build             # Build binary

# Code generation
make models            # Regenerate models from OpenAPI spec
make swagger           # Regenerate Swagger documentation
make generate          # Regenerate both models and Swagger docs

# Testing
make test              # Run integration tests
make test-coverage     # Run tests with coverage

# Maintenance
make clean             # Clean artifacts and database
make deps              # Update dependencies
make verify            # Build and test everything
make install-tools     # Install oapi-codegen, swag, air
```

**Windows users:** If `make` is not installed, use the direct commands shown below.

### 1. Modifying the API

If you need to change the API (add fields, new endpoints, etc.):

```bash
cd backend

# 1. Edit the OpenAPI spec
nano petstore_swagger.yml

# 2. Regenerate models
make models
# Or: oapi-codegen -config oapi-codegen.yaml petstore_swagger.yml

# 3. Regenerate Swagger docs (if handler annotations changed)
make swagger
# Or: swag init -g cmd/api/main.go --output docs --parseDependency --parseInternal

# 4. Update repositories if model fields changed
# (Fix any compilation errors in repository conversions)

# 5. Update domain services if validation changed

# 6. Update handlers if needed

# 7. Verify everything works
make verify
# Or: go build ./... && go test ./tests/integration/... -v
```

### 2. Adding a New Endpoint

Example: Add `GET /pet/findByName?name=Balu`

```bash
cd backend

# 1. Add to petstore_swagger.yml
# 2. Regenerate models (may generate new param types)
make models
# 3. Add repository method: FindByName(name string) ([]*models.Pet, error)
# 4. Add domain method: FindPetsByName(name string) ([]*models.Pet, error)
# 5. Add handler with Swagger annotations: FindPetsByName(c *fiber.Ctx) error
# 6. Register route in router.go: app.Get("/pet/findByName", handler.FindPetsByName)
# 7. Regenerate Swagger docs
make swagger
# 8. Write tests in pet_test.go
make test
```

## Swagger UI Documentation

The API includes **interactive Swagger UI documentation** for easy exploration and testing.

### Accessing Swagger UI

Once the server is running (`go run cmd/api/main.go`), open your browser to:

**http://localhost:3000/swagger/index.html**

### Features

- **Interactive API Documentation**: View all endpoints with descriptions
- **Try It Out**: Execute API requests directly from the browser
- **Request/Response Schemas**: See exact data structures with examples
- **Parameter Descriptions**: Understand what each field means
- **Response Codes**: Know what HTTP status codes to expect

### Swagger Architecture

The Swagger documentation is generated from code annotations:

```
Go Source Code with Annotations
         ↓
   swag CLI tool
         ↓
Generated docs/ directory
    - docs.go
    - swagger.json
    - swagger.yaml
         ↓
Imported in main.go (_ "petstore/docs")
         ↓
Served by gofiber/swagger middleware
         ↓
Accessible at /swagger/index.html
```

### Swagger Annotations

All handlers include Swagger annotations for automatic documentation generation.

**Example from `pet.go`:**

```go
// @Summary Add a new pet to the store
// @Description Add a new pet with the input payload
// @Tags pets
// @Accept json
// @Produce json
// @Param pet body models.Pet true "Pet object that needs to be added"
// @Success 200 {object} models.Pet
// @Failure 400 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Router /pet [post]
func (h *PetHandler) AddPet(c *fiber.Ctx) error {
    // Implementation
}
```

**Annotation breakdown:**
- `@Summary`: Short description shown in endpoint list
- `@Description`: Detailed explanation of what the endpoint does
- `@Tags`: Groups related endpoints (e.g., "pets", "categories")
- `@Accept`/`@Produce`: Content types (application/json)
- `@Param`: Request parameters
  - Format: `name location type required "description"`
  - Locations: `body`, `path`, `query`, `header`
- `@Success`: Successful response with HTTP code and type
- `@Failure`: Error responses with HTTP codes
- `@Router`: Route path and HTTP method

**API-level annotations in `main.go`:**

```go
// @title Petstore API
// @version 1.0
// @description REST API for a petstore application built with Go Fiber
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@petstore.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /
// @schemes http
```

### Regenerating Swagger Documentation

Whenever you modify handler annotations or add new endpoints:

```bash
cd backend
swag init -g cmd/api/main.go --output docs --parseDependency --parseInternal
```

**What happens:**
1. Scans `cmd/api/main.go` for general API info
2. Parses all files for handler annotations
3. Discovers models from handler signatures
4. Generates `docs/docs.go`, `docs/swagger.json`, `docs/swagger.yaml`

**When to regenerate:**
- Added new endpoints
- Modified handler annotations
- Changed request/response models
- Updated API metadata (title, version, contact, etc.)

**Recommended workflow:**
```bash
# 1. Add or modify handler with Swagger annotations
# 2. Regenerate docs
swag init -g cmd/api/main.go --output docs --parseDependency --parseInternal
# 3. Verify
go build ./...
# 4. Test in browser
go run cmd/api/main.go
# Open http://localhost:3000/swagger/index.html
```

### Swagger UI Route Configuration

In `internal/api/router.go`:

```go
import "github.com/gofiber/swagger"

func SetupRoutes(app *fiber.App, petHandler *handlers.PetHandler, categoryHandler *handlers.CategoryHandler) {
    // ... other routes ...
    
    // Swagger UI
    app.Get("/swagger/*", swagger.HandlerDefault)
}
```

The wildcard `/*` allows Swagger UI to serve all its static assets (CSS, JS, etc.).

### Common Issues

**Problem**: "404 Not Found" at `/swagger/index.html`  
**Solution**: Ensure `docs` package is imported in `main.go`:
```go
_ "petstore/docs"
```

**Problem**: Swagger shows old endpoint definitions  
**Solution**: Regenerate docs with `swag init`

**Problem**: Model not showing in Swagger UI  
**Solution**: Ensure model is referenced in handler annotation:
```go
// @Success 200 {object} models.Pet  ← Must reference the model
```

**Problem**: `swag` command not found  
**Solution**: Install swag CLI:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

## Common Issues and Solutions

### Issue: "undefined: models.ID"
**Cause**: Generated models use `Id` (lowercase 'd'), not `ID`  
**Solution**: Use `pet.Id`, `category.Id` in your code

### Issue: "cannot use *pet.Status (type models.PetStatus) as string"
**Cause**: Generated enums are typed, database expects string  
**Solution**: Convert: `string(*pet.Status)` when saving to DB

### Issue: "go-sqlite3 requires cgo"
**Cause**: Using wrong SQLite driver  
**Solution**: Use `github.com/glebarez/sqlite` (already configured)

### Issue: Tests failing with unique constraint violations
**Cause**: Tests sharing database state  
**Solution**: Each test should call `setupTestApp(t)` to get fresh DB

## Code Generation Details

### oapi-codegen Configuration (`oapi-codegen.yaml`)

```yaml
package: models
output: internal/models/generated.go
generate:
  models: true              # Generate type definitions
  embedded-spec: false      # Don't embed full spec
output-options:
  skip-prune: false         # Prune unreferenced generated types
```

### What Gets Generated

From this OpenAPI definition:
```yaml
Pet:
  type: object
  required: [name]
  properties:
    id:
      type: integer
      format: int64
    name:
      type: string
    status:
      type: string
      enum: [available, pending, sold]
```

Generates:
```go
type Pet struct {
    Id     *int64     `json:"id,omitempty"`
    Name   string     `json:"name"`
    Status *PetStatus `json:"status,omitempty"`
}

type PetStatus string

const (
    PetStatusAvailable PetStatus = "available"
    PetStatusPending   PetStatus = "pending"
    PetStatusSold      PetStatus = "sold"
)
```

## Performance Considerations

- **SQLite** is sufficient for development and small deployments
- **GORM** uses connection pooling automatically
- **Soft delete** means queries need `WHERE deleted_at IS NULL` (GORM handles this)
- **Indexes** on foreign keys and frequently queried fields

## Security Notes

Current implementation is for **development/learning**:
- ⚠️ No authentication or authorization
- ⚠️ CORS allows all origins
- ⚠️ No rate limiting
- ⚠️ No input sanitization beyond validation

For production, add:
- JWT or session-based authentication
- Role-based access control
- CORS whitelist
- Request rate limiting
- SQL injection protection (GORM provides this)
- XSS protection

## Next Steps

To extend this project:
1. Add authentication (JWT tokens)
2. Add photo upload for pets (OpenAPI spec has `photoUrls` field)
3. Add pagination for list endpoints
4. Add filtering and sorting
5. Add unit tests for domain and repository layers
6. ✅ **API documentation UI (Swagger UI)** - Already implemented!
7. Deploy with Docker
8. Add logging and monitoring
9. Migrate to PostgreSQL for production

## Frontend

The project includes a client-side rendered (CSR) frontend served as static files.

### Architecture
- **Rendering**: All rendering happens in the browser (NOT server-side rendering)
- **Static serving**: Go Fiber serves files from `../ui` via `app.Static("/", "../ui")`
- **API calls**: Alpine.js components use `fetch()` to call REST endpoints
- **Styling**: Tailwind CSS via CDN with custom gradients and animations

### Pages
- `index.html` - Pet management: list, filter by status, create, edit, delete pets
- `categories.html` - Category management: list, create, edit, delete categories

### Error Handling
- Backend returns `{"error": "message"}`
- UI checks both `error.error` and `error.message` fields
- Loading/error/empty states shown conditionally via Alpine.js directives

## Resources

- [Go Fiber Documentation](https://docs.gofiber.io/)
- [GORM Documentation](https://gorm.io/docs/)
- [oapi-codegen GitHub](https://github.com/oapi-codegen/oapi-codegen)
- [OpenAPI 3.0 Specification](https://swagger.io/specification/)
- [Swaggo Documentation](https://github.com/swaggo/swag)
- [Swagger UI](https://swagger.io/tools/swagger-ui/)
