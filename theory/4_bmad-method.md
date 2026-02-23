# BMAD Method

## What is the BMAD Method?

The **BMAD Method** (Breakthrough Method of Agile AI-Driven Development) is an open-source framework that structures how AI agents build software — from initial idea all the way to working code. Instead of just asking an AI "build me an app" and hoping for the best, BMAD gives you **specialized AI agents**, **guided workflows**, and **structured documents** that progressively build context across four phases.

Think of it like this:
- **Without BMAD**: You throw a vague request at an AI and get inconsistent results
- **With BMAD**: You follow a structured process where each step produces a document that feeds into the next step — so the AI always knows exactly what to build and why

## The Problem It Solves

### Without BMAD

**Scenario**: You ask an AI to build a SaaS product.

```
You: "Build me a project management tool like Trello"

AI: "Sure! Here's a React app with..."
    → Makes assumptions about tech stack you didn't want
    → Picks database you'd never use in production
    → Skips authentication entirely
    → Creates 3 files, then loses context
    → Next chat session: AI forgot everything
    → You end up with inconsistent, half-baked code
```

**Result**: Wasted time, inconsistent code, no plan, no architecture, constant re-explaining.

### With BMAD

```
Phase 1 (Analysis):
  You brainstorm with the Analyst agent
  → Output: product-brief.md (vision, goals, users)

Phase 2 (Planning):
  PM agent creates requirements
  → Output: PRD.md (features, user stories, constraints)

Phase 3 (Solutioning):
  Architect agent designs the system
  → Output: architecture.md (tech stack, patterns, decisions)
  PM agent breaks work into stories
  → Output: epics/ folder with detailed stories

Phase 4 (Implementation):
  Dev agent implements story by story
  → Each story file has FULL context (PRD + architecture + specific task)
  → Code is consistent because every agent reads the same plan
```

**Result**: Structured, consistent, well-architected software — built by AI agents that always have the right context.

## How It Works — Simple Explanation

### The Film Production Analogy

Building software with BMAD is like producing a film:

**Without a process** (chaotic indie film):
- Director says "let's just start filming!"
- No script, no storyboard, no budget
- Actors don't know their lines
- Every scene is improvised
- Result: unwatchable mess

**With BMAD** (professional production):
1. **Writer** (Analyst) creates the concept and pitch
2. **Screenwriter** (PM) writes the detailed script with scenes
3. **Director** (Architect) plans shot lists, locations, crew
4. **Actors & Crew** (Dev Agent) execute scene by scene, following the script
5. **Editor** (Code Review) reviews each scene for quality

Each role is a **specialized AI agent** that's an expert in its domain.

## The Four Phases

### Phase 1: Analysis (Optional)

**Purpose**: Explore the idea before committing to building it.

**Agent**: Analyst 🔍

**Workflows**:
- **Brainstorming**: Guided ideation with an AI coach
- **Research**: Market and technical validation
- **Product Brief**: Capture the strategic vision

**Output**: `product-brief.md` — What are we building and why?

```
Example output:
┌──────────────────────────────────────┐
│ Product Brief: PetStore App          │
│                                      │
│ Vision: Simple pet adoption platform │
│ Users: Pet owners, shelters          │
│ Key Features: Browse, adopt, manage  │
│ Success Metric: 100 adoptions/month  │
└──────────────────────────────────────┘
```

### Phase 2: Planning (Required)

**Purpose**: Define exactly what to build.

**Agent**: PM (Product Manager) 📋

**Workflows**:
- **Create PRD**: Product Requirements Document with features, user stories, constraints
- **UX Design**: UI/UX specifications (if the project has a frontend)

**Output**: `PRD.md` — Detailed requirements

```
Example output:
┌──────────────────────────────────────┐
│ PRD: PetStore App                    │
│                                      │
│ FR-1: Users can list pets            │
│ FR-2: Users can filter by status     │
│ FR-3: Users can add new pets         │
│ NFR-1: Response time < 200ms        │
│ NFR-2: SQLite for simplicity         │
│ Constraints: No auth for v1          │
└──────────────────────────────────────┘
```

### Phase 3: Solutioning

**Purpose**: Design how to build it, then break it into stories.

**Agent**: Architect 🏗️ → PM 📋

**Workflows**:
1. **Create Architecture**: Tech stack, patterns, API design, data models
2. **Create Epics & Stories**: Break PRD into implementable units of work
3. **Implementation Readiness Check**: Validate everything is coherent

**Output**: `architecture.md` + `epics/` folder with story files

```
Example output:
┌──────────────────────────────────────┐
│ Architecture                         │
│                                      │
│ Stack: Go + Fiber + SQLite + GORM    │
│ Pattern: Clean 3-Layer Architecture  │
│ API: REST (OpenAPI 3.0 spec)         │
│ Frontend: Alpine.js + Tailwind CSS   │
│                                      │
│ Epics:                               │
│   Epic 1: Pet CRUD (5 stories)       │
│   Epic 2: Category Management (3)    │
│   Epic 3: Frontend UI (4 stories)    │
└──────────────────────────────────────┘
```

### Phase 4: Implementation

**Purpose**: Build it — one story at a time.

**Agents**: Scrum Master 📊 → Developer 💻 → Code Review 🔍

**Workflows**:
1. **Sprint Planning**: Initialize tracking (`sprint-status.yaml`)
2. **Create Story**: Prepare the next story file with full context
3. **Dev Story**: AI implements the story (code + tests)
4. **Code Review**: Validate implementation quality
5. **Retrospective**: Review after epic completion

**The Build Cycle** (repeat for each story):

```
┌─────────────────────────────────────────────────┐
│                                                 │
│   Scrum Master          Developer               │
│   creates story  ──→    implements code  ──→    │
│   (story-001.md)        (writes + tests)        │
│                              │                  │
│                              ▼                  │
│                         Code Review             │
│                         (quality check)         │
│                              │                  │
│                              ▼                  │
│                         Next Story              │
│                                                 │
└─────────────────────────────────────────────────┘
```

## The Secret: Context Engineering

The real power of BMAD is how it manages **context** for AI agents.

### The Problem with AI and Context

```
Without structured context:

  Chat 1: "Build a REST API" → AI makes assumptions
  Chat 2: "Add authentication" → AI forgot the tech stack from Chat 1
  Chat 3: "Fix the database" → AI doesn't know the schema
  
  Result: Every chat starts from scratch. AI has amnesia.
```

### How BMAD Solves It

```
Each phase produces a document. Each document feeds into the next phase.

  product-brief.md    ──→  PM agent reads it to create PRD
       PRD.md         ──→  Architect reads it to design architecture  
    architecture.md   ──→  PM reads both to create stories
     story-001.md     ──→  Dev agent reads ALL of the above to implement

The Dev agent implementing story-001 has:
  ✅ The product vision (from brief)
  ✅ The exact requirements (from PRD)
  ✅ The technical decisions (from architecture)
  ✅ The specific task details (from story file)
  
  → Agent has COMPLETE context. No guessing.
```

This is called **progressive context building** — each phase enriches the context for the next.

## The 12+ Specialized Agents

BMAD uses specialized AI agent personas, each an expert in one domain:

| Agent | Role | What They Do |
|-------|------|-------------|
| 🔍 Analyst | Research & Discovery | Brainstorming, market research, product briefs |
| 📋 PM | Product Manager | PRD, requirements, epics and stories |
| 🎨 UX Designer | User Experience | UI/UX specifications and design |
| 🏗️ Architect | System Design | Architecture, tech stack, patterns |
| 💻 Developer | Implementation | Code, tests, story implementation |
| 📊 Scrum Master | Project Management | Sprint planning, tracking, retrospectives |
| 🔍 Code Reviewer | Quality | Validate code quality and standards |
| 🧭 BMad-Help | Guide | Intelligent assistant that knows what to do next |

Each agent is loaded via a slash command (e.g., `/bmad-agent-bmm-pm`) and has deep expertise in its domain.

## Scale-Adaptive Intelligence

BMAD adapts to project size — you don't use the full process for a bug fix:

| Track | When to Use | What You Do |
|-------|-------------|-------------|
| **Quick Flow** | Bug fixes, small features (1-15 stories) | Write a quick tech-spec → implement |
| **BMad Method** | Products, platforms (10-50+ stories) | Full 4-phase process |
| **Enterprise** | Complex systems (30+ stories) | Full process + security + DevOps |

```
Bug fix?
  → /bmad-bmm-quick-spec → /bmad-bmm-quick-dev → Done

New product?
  → Analysis → Planning → Architecture → Stories → Implementation
```

## Real-World Example

**Scenario**: Build the Petstore API (the project in this course).

### Without BMAD
```
You: "Build me a pet store REST API in Go"
AI:  Creates a single main.go with everything mixed together
     Uses different naming than the spec
     Forgets to handle errors
     No tests
     You spend 3 hours fixing AI mistakes
```

### With BMAD
```
Phase 1 (Skip — we already know what we want)

Phase 2 — PM Agent:
  Input: "Pet store API for managing pets and categories"
  Output: PRD.md with:
    - 6 pet endpoints, 4 category endpoints
    - Validation rules (name required, status enum)
    - Error handling requirements
    - Test coverage requirements

Phase 3 — Architect Agent:
  Input: PRD.md
  Output: architecture.md with:
    - Go + Fiber + GORM + SQLite
    - Clean 3-layer architecture
    - OpenAPI spec as source of truth
    - Entity vs Model pattern
    
  PM Agent (again):
  Input: PRD.md + architecture.md
  Output: 3 epics, 12 stories:
    Epic 1: Pet CRUD (6 stories)
    Epic 2: Category Management (4 stories)  
    Epic 3: API Documentation (2 stories)

Phase 4 — Dev Agent:
  Story 1: "Implement POST /pet endpoint"
    → Has PRD (knows validation rules)
    → Has architecture (knows 3-layer pattern)
    → Has story (knows exact acceptance criteria)
    → Generates: handler + service + repository + test
    → Code review: ✅ Approved
    
  Story 2: "Implement GET /pet/:id endpoint"
    → Same rich context
    → Consistent with Story 1's patterns
    → ...continues for all 12 stories
```

**Result**: Consistent, well-tested, properly architected code — with AI doing the heavy lifting while you guide the process.

## Connection to Other Concepts

### BMAD + Agents
BMAD IS an agent framework. Each persona (PM, Architect, Developer) is a specialized agent with specific system prompts, tools, and expertise. The framework orchestrates when each agent is invoked.

### BMAD + MCP
MCP servers can provide context to BMAD agents — e.g., a Jira MCP server lets the PM agent read existing tickets, or a filesystem MCP server lets the Dev agent read and write code.

### BMAD + Spec-Driven Development
BMAD and spec-driven development complement each other:
- BMAD's Architect agent can produce an OpenAPI spec during Phase 3
- The spec becomes part of the architecture context
- The Dev agent uses the spec to generate consistent implementations
- See `spec-driven-development.md` for details

## Key Takeaways

```
1. BMAD = Structured process for AI-driven development
2. Four phases: Analysis → Planning → Solutioning → Implementation
3. Specialized agents: PM, Architect, Developer, Scrum Master, etc.
4. Context engineering: Each phase produces documents that feed the next
5. Scale-adaptive: Quick Flow for small tasks, full process for big projects
6. The AI doesn't guess — it has complete context at every step
```

## Resources

- **GitHub**: https://github.com/bmad-code-org/BMAD-METHOD (37k+ ⭐)
- **Docs**: https://docs.bmad-method.org
- **Install**: `npx bmad-method install`
- **Discord**: https://discord.gg/gk8jAdXWmj
- **YouTube**: https://www.youtube.com/@BMadCode
