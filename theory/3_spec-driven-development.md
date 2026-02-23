# Spec-Driven Development

## What is Spec-Driven Development?

Spec-Driven Development (also called **Design-First** or **API-First** development) is an approach where you **write the specification before writing any code**. The specification (usually an OpenAPI/Swagger YAML file) becomes the single source of truth, and code is generated or validated against it.

Think of it like building a house:
- **Code-First**: Start building walls, figure out the floor plan as you go
- **Spec-First**: Draw the blueprint first, then build exactly to spec

## The Problem It Solves

### Without Spec-Driven Development

**Scenario**: Two teams — frontend and backend — building a pet store app.

```
Week 1:
  Backend dev: "I'll create POST /pets with a 'petName' field"
  Frontend dev: "I'll send POST /pet with a 'name' field"
  → They discover the mismatch in Week 3 during integration
  → 2 days wasted fixing naming, paths, and response formats

Week 2:
  Backend dev: "I changed the response from 200 to 201 for creation"
  Frontend dev: (doesn't know, code breaks)
  → Bug discovered by QA in Week 4

Week 3:
  New developer joins: "Where's the API documentation?"
  Team: "Read the code... and the Slack messages from January"
```

**Result**: Miscommunication, integration bugs, outdated docs, frustration.

### With Spec-Driven Development

```
Day 1:
  Team writes OpenAPI spec together:
  
  POST /pet:
    body: { name: string (required), status: enum[available,pending,sold] }
    response 200: { id: int64, name: string, status: string }
    response 400: { error: string }

Day 2:
  Backend: Generates models from spec, implements handlers
  Frontend: Generates API client from spec, builds UI
  → Both use the SAME contract — zero mismatches

Day 5:
  New developer joins: "Here's the spec — it IS the docs"
  → Productive in 30 minutes
```

**Result**: Aligned teams, no integration surprises, living documentation.

## How It Works

### The Spec is the Source of Truth

```
                    petstore_swagger.yml
                    (OpenAPI 3.0 Spec)
                          │
            ┌─────────────┼─────────────┐
            ▼             ▼             ▼
      Code Generation   Swagger UI   Validation
      (oapi-codegen)    (auto docs)  (contract tests)
            │             │             │
            ▼             ▼             ▼
      Go Structs     Interactive    "Does my API
      Pet, Category  API Browser    match the spec?"
```

### Step-by-Step Workflow

**1. Write the Specification**

```yaml
# petstore_swagger.yml
openapi: "3.0.0"
info:
  title: Petstore API
  version: "1.0"

paths:
  /pet:
    post:
      summary: Add a new pet
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/Pet'
      responses:
        '200':
          description: Pet created successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Pet'

components:
  schemas:
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

**2. Generate Code from the Spec**

```bash
# Generate Go models
oapi-codegen -config oapi-codegen.yaml petstore_swagger.yml
```

This automatically creates:
```go
// generated.go — DO NOT EDIT
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

**3. Implement Business Logic Using Generated Models**

```go
// You write this part — the spec gives you the types
func (s *PetService) CreatePet(pet *models.Pet) (*models.Pet, error) {
    if pet.Name == "" {
        return nil, errors.New("pet name is required")
    }
    return s.repo.Create(pet)
}
```

**4. Generate Documentation Automatically**

```bash
# Swagger UI is auto-generated from code annotations
swag init -g cmd/api/main.go --output docs
```

Open `http://localhost:3000/swagger/index.html` → interactive docs that always match your code.

## Real-World Analogy

### The Architect's Blueprint

**Building a house without blueprints** (Code-First):
- Electrician runs wires wherever seems right
- Plumber puts pipes where there's space
- They meet in the same wall → rip out and redo
- Inspector arrives: "Where are the plans?" → "Uh..."

**Building with blueprints** (Spec-First):
- Everyone works from the same drawing
- Electrician knows where plumber's pipes go before starting
- Inspector can verify against the plan
- New contractor can join and immediately understand the building

The **OpenAPI spec is your blueprint**.

## What Can Be Generated from a Spec?

| What | Tool | Language |
|------|------|----------|
| Server models/types | oapi-codegen | Go |
| Server stubs | openapi-generator | Go, Java, Python, etc. |
| Client SDKs | openapi-generator | TypeScript, Python, Java, etc. |
| API documentation | Swagger UI, ReDoc | Interactive HTML |
| Mock servers | Prism | Any (runs standalone) |
| Contract tests | Schemathesis | Python |
| Postman collections | openapi-to-postman | Postman |

## Spec-Driven Development with AI Agents

This is where spec-driven development becomes **extremely powerful with AI tools**:

### Why AI Loves Specs

```
Without spec:
  You: "Build me a REST API for a pet store"
  AI: "What fields should Pet have? What endpoints? What status codes? 
       What validation? JSON or XML? What about error handling?"
  → 10 rounds of back-and-forth to clarify requirements

With spec:
  You: "Implement this API" [attaches petstore_swagger.yml]
  AI: Reads the spec → Knows EXACTLY:
      - Every endpoint, method, and path
      - Every field name, type, and constraint
      - Every response code and format
      - Every validation rule
  → Generates correct implementation in one shot
```

### The AI Workflow

```
1. Human writes/reviews the OpenAPI spec
   (This is the creative, decision-making part)

2. AI agent receives the spec as context
   (Via .cursorrules, MCP, or direct upload)

3. AI generates:
   - Repository layer (CRUD operations)
   - Service layer (business logic + validation)
   - Handler layer (HTTP endpoints)
   - Integration tests
   - Swagger annotations
   
4. Human reviews and adjusts
   (Focus on edge cases and business rules)
```

### Example: How Our Petstore Uses This

```
petstore_swagger.yml          ← Human-written contract
        │
        ├─→ oapi-codegen      ← Generates models/generated.go
        │
        ├─→ AI agent          ← Reads spec to understand exact API contract
        │   ├─→ Implements handlers matching spec paths
        │   ├─→ Implements validation matching spec constraints  
        │   └─→ Writes tests covering spec scenarios
        │
        └─→ swag init         ← Generates Swagger UI docs
```

The spec acts as a **shared language** between human and AI — unambiguous, machine-readable, complete.

## Benefits

### 1. **Contract Clarity**
Everyone (humans AND AI) agrees on the API shape before coding starts.

### 2. **Parallel Development**
Frontend and backend can work simultaneously using the same spec. Frontend can even use a mock server generated from the spec.

### 3. **Living Documentation**
The spec IS the documentation. It's never outdated because the code is generated from it.

### 4. **Code Generation**
Eliminate boilerplate. Types, validators, and client SDKs come free.

### 5. **AI-Friendly**
AI assistants can read the spec and generate accurate implementations without guessing.

### 6. **Change Management**
Want to add a field? Change the spec first → regenerate → compiler tells you what broke.

## Challenges

### 1. **Upfront Investment**
Writing a good spec takes time. Teams used to "just coding" may resist.

### 2. **Spec Drift**
If developers bypass the spec (adding endpoints directly in code), the spec becomes stale. Discipline is needed.

### 3. **Learning Curve**
OpenAPI syntax has a learning curve, especially for complex schemas (oneOf, allOf, discriminators).

### 4. **Not Everything Fits**
Some APIs are highly dynamic or event-driven — REST specs don't cover WebSockets, GraphQL, or streaming well.

## When to Use Spec-Driven Development

✅ **Good For:**
- Multi-team projects (frontend + backend + mobile)
- Public APIs that external developers consume
- AI-assisted development (give AI the spec as context)
- Microservices with many inter-service contracts
- Projects requiring API documentation

❌ **Less Useful For:**
- Quick prototypes or hackathons
- Solo projects where you're the only consumer
- Non-REST APIs (gRPC, GraphQL have their own spec systems)
- Rapidly changing experimental APIs

## Tools Ecosystem

### Writing Specs
- **Swagger Editor**: Web-based YAML editor with live preview
- **Stoplight Studio**: Visual API designer (GUI for OpenAPI)
- **VS Code + OpenAPI extension**: Syntax highlighting and validation

### Generating Code
- **oapi-codegen**: Go models and server interfaces
- **openapi-generator**: 50+ languages, server stubs and client SDKs
- **swagger-codegen**: Original code generator (now superseded by openapi-generator)

### Validating
- **Spectral**: Lint your OpenAPI spec for best practices
- **Prism**: Mock server + request validation from spec
- **Schemathesis**: Auto-generate API tests from spec

## Connection to Other Concepts

### Spec-Driven + Agents
An AI agent can use the OpenAPI spec as its "instruction manual" — knowing exactly what endpoints exist, what data to send, and what to expect back.

### Spec-Driven + MCP
An MCP server could expose an OpenAPI spec as a resource, allowing any AI tool to discover and understand an API automatically.

### Spec-Driven + BMAD Method
The BMAD method (see `bmad-method.md`) takes spec-driven thinking further — using structured documents (stories, PRDs, architecture docs) as specs that AI agents consume at each development phase.

## Summary

```
Traditional:  Idea → Code → Hope it works → Write docs (maybe)
Spec-Driven:  Idea → Spec → Generate code → Docs come free
With AI:      Idea → Spec → AI implements from spec → Human reviews
```

The spec is the contract. Everything flows from it. The more precise your spec, the better your AI-generated code.
