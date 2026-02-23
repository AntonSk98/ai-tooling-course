# Petstore

A full-stack pet store application with a Go REST API backend and an Alpine.js browser frontend.

## Project Structure

```
petstore/
├── backend/                    # Go REST API (Fiber + GORM + SQLite)
│   ├── cmd/api/                # Entry point
│   ├── internal/               # Handlers, services, repositories, models
│   ├── pkg/                    # Configuration
│   ├── tests/                  # Integration tests
│   ├── docs/                   # Auto-generated Swagger docs
│   ├── Makefile                # Build/test/generate tasks
│   ├── .cursorrules            # AI assistant instructions
│   ├── IMPLEMENTATION.md       # Detailed implementation guide
│   └── petstore_swagger.yml    # OpenAPI 3.0 spec (source of truth)
├── ui/                         # Frontend (Alpine.js + Tailwind CSS)
│   ├── index.html              # Pet management
│   ├── categories.html         # Category management
│   └── agent.md                # Frontend architecture docs
└── README.md                   # This file
```

## Quick Start

### Prerequisites

- Go 1.22+

### Run

```bash
cd backend
go mod download
go run cmd/api/main.go
```

Open in browser:
- **UI**: http://localhost:3000/index.html
- **Swagger**: http://localhost:3000/swagger/index.html

### Test

```bash
cd backend
go test ./tests/integration/... -v
```

## Backend

Go REST API built with [Fiber](https://gofiber.io/), [GORM](https://gorm.io/), and SQLite. Follows clean 3-layer architecture (Handler → Service → Repository).

### API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/health` | Health check |
| `POST` | `/pet` | Add a new pet |
| `PUT` | `/pet` | Update an existing pet |
| `GET` | `/pet/{petId}` | Get pet by ID |
| `GET` | `/pet/findByStatus?status=` | Find pets by status |
| `POST` | `/pet/{petId}?name=&status=` | Update pet with form data |
| `DELETE` | `/pet/{petId}` | Delete a pet |
| `POST` | `/category` | Add a new category |
| `PUT` | `/category` | Update an existing category |
| `GET` | `/category/listAll` | Get all categories |
| `DELETE` | `/category/{categoryId}` | Delete a category |

### Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `3000` | Server port |
| `DB_PATH` | `./petstore.db` | SQLite database path |

### Makefile (from `backend/`)

```bash
make run               # Run the application
make test              # Run integration tests
make test-coverage     # Tests with coverage report
make models            # Regenerate models from OpenAPI spec
make swagger           # Regenerate Swagger docs
make generate          # Regenerate both
make build             # Build binary
make verify            # Build + test
make install-tools     # Install oapi-codegen, swag, air
make help              # Show all commands
```

For detailed backend documentation, see [`backend/IMPLEMENTATION.md`](backend/IMPLEMENTATION.md).

## Frontend

Browser-based UI served as static files by the Go backend.

- **Pet Management** (`index.html`): List, filter by status, create, edit, delete pets
- **Category Management** (`categories.html`): List, create, edit, delete categories
- **Stack**: [Alpine.js](https://alpinejs.dev/) 3.x (reactivity) + [Tailwind CSS](https://tailwindcss.com/) 3.x (styling via CDN)
- **API calls**: All via `fetch()` — no HTMX, no build step, no bundler

For frontend architecture details, see [`ui/agent.md`](ui/agent.md).

## Technology Stack

| Layer | Technology |
|-------|-----------|
| API Framework | Fiber v2.52.11 |
| ORM | GORM v1.25.7 |
| Database | SQLite (pure Go, no CGO) |
| Code Generation | oapi-codegen v2.5.1 |
| API Docs | Swagger UI (swaggo/swag) |
| Frontend | Alpine.js 3.x + Tailwind CSS 3.x |

## License

Apache 2.0
