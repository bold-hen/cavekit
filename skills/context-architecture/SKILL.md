---
name: context-architecture
description: |
  Progressive disclosure architecture for organizing project context so AI agents only load what they need.
  Covers the 4-level hierarchy, standard context directory structure, CLAUDE.md conventions,
  multi-repo strategy with shared context via git submodules, and 3-tier hierarchy.
  Trigger phrases: "context architecture", "progressive disclosure", "organize context for agents",
  "context directory structure", "how to structure docs for AI"
---

# Context Architecture: Progressive Disclosure for AI Agents

## Core Principle: Agents Should Only Read What They Need

AI agents have finite context windows. Loading irrelevant information wastes capacity, dilutes focus, and degrades output quality. Progressive disclosure organizes documentation as a hierarchy where agents start with high-level summaries and drill into details only when needed.

**The goal:** An agent working on authentication should load auth specs and plans, not the entire project's documentation. An agent doing gap analysis should load the spec index and implementation tracking, not every reference document.

---

## The 4-Level Hierarchy

Progressive disclosure works through four levels:

### Level 1: Index Files

Index files list all sub-documents with one-line summaries. They are the entry point for any domain.

```markdown
# Spec Overview

| Domain | File | Summary |
|--------|------|---------|
| Authentication | spec-auth.md | Registration, login, sessions, OAuth |
| Data Models | spec-data-models.md | Core entities, relationships, validation |
| API | spec-api.md | REST endpoints, formats, error handling |
```

**Rule:** An agent reads the index first. The index tells the agent which sub-documents are relevant to its task.

### Level 2: Sub-Documents

Sub-documents contain full details for one domain. They include `file:line` references back to source code or reference materials where applicable.

```markdown
# Spec: Authentication

## R1: User Registration
**Description:** New users can create accounts with email and password.
**Acceptance Criteria:**
- [ ] POST /register with valid email/password returns 201
- [ ] Duplicate email returns 409
**Source:** ref-apis.md:45-67
```

**Rule:** Sub-documents are self-contained for their domain but cross-reference related domains.

### Level 3: Agents Spider Out

Agents follow a discovery pattern:
1. Read the index
2. Identify which sub-documents are relevant to the current task
3. Read only those sub-documents
4. Follow cross-references as needed

This pattern prevents loading the entire documentation tree into the context window.

### Level 4: Context Window Efficiency

Two techniques keep documents focused and navigable:

- **Decomposition** — break large documents into domain-specific files
- **Spec compaction** — compress implementation tracking when files exceed approximately 500 lines (see `bp:impl-tracking`)

Together, these ensure that the context an agent loads is dense with relevant information rather than bloated with resolved history.

---

## Standard Context Directory Structure

Every Blueprint project uses this directory layout:

```
context/
├── CLAUDE.md               # Root context: "This is a Blueprint project. See context/ for blueprints, plans, and tracking."
├── refs/                   # Source of truth — reference materials
│   ├── prd.md              # Product requirements document
│   ├── design-doc.md       # Architecture/design documents
│   ├── api-draft.md        # API specifications
│   └── research/           # Framework research, competitive analysis
│       └── ...
├── blueprints/             # Implementation-agnostic blueprints
│   ├── CLAUDE.md           # "Blueprints define WHAT needs implementing. Never prescribe HOW."
│   ├── blueprint-overview.md    # START HERE — index of all domain blueprints
│   ├── blueprint-auth.md        # Domain: Authentication
│   ├── blueprint-data-models.md # Domain: Data Models
│   ├── blueprint-api.md         # Domain: API
│   └── blueprint-ui.md          # Domain: UI Components
├── plans/                  # Framework-specific implementation plans
│   ├── CLAUDE.md           # "Plans define HOW to implement something using {FRAMEWORK}."
│   ├── plan-build-site.md        # Tier system showing feature dependencies
│   ├── plan-known-issues.md      # Prioritized backlog (P0-P3)
│   ├── plan-auth.md              # Implementation plan: Authentication
│   ├── plan-data-models.md       # Implementation plan: Data Models
│   └── plan-api.md               # Implementation plan: API
├── impl/                   # Living implementation tracking
│   ├── CLAUDE.md           # "Impls record implementation progress. Update after every session."
│   ├── impl-all.md         # Current tracking document
│   └── archive/            # Compacted/archived tracking docs
│       └── impl-all-v1.md
└── prompts/                # DABI pipeline prompts
    ├── 001-generate-blueprints.md
    ├── 002-generate-plans.md
    └── 003-implement.md
```

### Directory Conventions

#### `refs/` — Reference Materials

- **Purpose:** Source of truth that blueprints are derived from
- **Contents:** PRDs, design documents, API drafts, language specifications, old code documentation
- **Convention:** Files are read-only after initial placement. Agents read but do not modify.
- **Naming:** Descriptive names prefixed by type (e.g., `ref-apis.md`, `ref-data-models.md`, `ref-architecture.md`)

#### `blueprints/` — Blueprints

- **Purpose:** Implementation-agnostic requirements with testable acceptance criteria
- **Contents:** One `blueprint-{domain}.md` per domain, plus `blueprint-overview.md` as the index
- **Convention:** Blueprints describe WHAT, never HOW. Every requirement has acceptance criteria.
- **Naming:** `blueprint-{domain}.md` — lowercase, hyphenated domain name
- **Updates:** Modified during Draft phase and via revision in Inspect phase

#### `plans/` — Implementation Plans

- **Purpose:** Framework-specific implementation plans derived from specs
- **Contents:** One `plan-{domain}.md` per domain, plus `plan-build-site.md` and `plan-known-issues.md`
- **Convention:** Plans reference blueprints by requirement ID (e.g., "implements blueprint-auth.md R1-R3"). Plans define implementation sequence, dependencies, and test strategies.
- **Naming:** `plan-{domain}.md` — matches corresponding blueprint domain
- **Updates:** Modified during Architect phase. May be updated by Build phase (bidirectional flow).

#### `impl/` — Implementation Tracking

- **Purpose:** Living documents recording what was built, what is pending, what failed
- **Contents:** `impl-all.md` (or domain-specific `impl-{domain}.md`), plus `archive/` for compacted history
- **Convention:** Updated after every implementation session. Dead ends are critical — they prevent agents from retrying failed approaches.
- **Naming:** `impl-{scope}.md` — scope can be "all" or a specific domain
- **Updates:** Modified during Build and Inspect phases. Compacted when exceeding 500 lines.

#### `prompts/` — Pipeline Prompts

- **Purpose:** Numbered prompts that drive each DABI phase
- **Contents:** One prompt per phase, numbered sequentially
- **Convention:** Prompts use runtime variables (`{FRAMEWORK}`, `{BUILD_COMMAND}`, `{TEST_COMMAND}`) for project-agnostic execution. Each prompt specifies its input and output directories.
- **Naming:** `NNN-{action}.md` — three-digit number, descriptive action

---

## CLAUDE.md Hierarchical Loading

CLAUDE.md files are loaded automatically by the agent when working in a directory. They load hierarchically — from the current directory AND all parent directories up to the project root.

### How Hierarchical Loading Works

```
project/
├── CLAUDE.md                    # Loaded for ALL work in the project
├── context/
│   ├── specs/
│   │   ├── CLAUDE.md            # Loaded when working in specs/
│   │   └── spec-auth.md
│   ├── plans/
│   │   ├── CLAUDE.md            # Loaded when working in plans/
│   │   └── plan-auth.md
│   └── impl/
│       ├── CLAUDE.md            # Loaded when working in impl/
│       └── impl-all.md
└── src/
    ├── CLAUDE.md                # Loaded when working in src/
    └── auth/
        ├── CLAUDE.md            # Loaded when working in src/auth/
        └── login.{ext}
```

When an agent works in `src/auth/`, it loads:
1. `project/CLAUDE.md` (project-level conventions)
2. `project/src/CLAUDE.md` (source code conventions)
3. `project/src/auth/CLAUDE.md` (auth module conventions)

### What Goes in Each CLAUDE.md

**Project root `CLAUDE.md`:**
```markdown
# Project: {Project Name}

This project uses Blueprint.

## Context Directory
- Blueprints: context/blueprints/ (WHAT to build)
- Plans: context/plans/ (HOW to build it with {FRAMEWORK})
- Impl Tracking: context/impl/ (progress, dead ends, test health)
- References: context/refs/ (source materials)

## Build & Test
- Build: {BUILD_COMMAND}
- Test: {TEST_COMMAND}
- Lint: {LINT_COMMAND}
```

**`context/blueprints/CLAUDE.md`:**
```markdown
# Blueprints

Blueprints define WHAT needs implementing. They are implementation-agnostic.

## Conventions
- Start with blueprint-overview.md for the domain index
- Each blueprint uses the R1, R2, R3... requirement numbering
- Every requirement has testable acceptance criteria
- Never prescribe HOW — that belongs in plans/
```

**`context/plans/CLAUDE.md`:**
```markdown
# Implementation Plans

Plans define HOW to implement blueprints using {FRAMEWORK}.

## Conventions
- Each plan references blueprint requirements by ID (e.g., "implements blueprint-auth.md R1")
- Plans include implementation sequence and dependencies
- plan-build-site.md shows the dependency tier system
- plan-known-issues.md tracks the prioritized backlog (P0-P3)
```

**`context/impl/CLAUDE.md`:**
```markdown
# Implementation Tracking

Impls record implementation progress. Update after every session.

## Conventions
- Record: task status, files created/modified, issues, dead ends, test health
- CRITICAL: Dead ends prevent retrying failed approaches — always document them
- Compact to under 500 lines when files get too large (archive to impl/archive/)
```

---

## Multi-Repo Strategy

For projects that need to evaluate multiple frameworks or maintain shared blueprints across implementations, Blueprint uses a **shared base context** with framework-specific repos.

### The 3-Tier Hierarchy

```
Tier 1: Shared Context (What IS)
└── shared-context/
    └── reference/              # Source materials (old code docs, PRDs, research)
        ├── ref-apis.md
        ├── ref-data-models.md
        ├── ref-ui-components.md
        └── ref-architecture.md

Tier 2: Shared Context (What MUST BE)
└── shared-context/
    └── blueprints/              # Implementation-agnostic blueprints
        ├── blueprint-overview.md
        ├── blueprint-features.md
        ├── blueprint-ui-requirements.md
        └── feature-scope.md

Tier 3: Application Context (per framework)
└── context/                    # Framework-specific
    ├── plans/                  # Plans using {FRAMEWORK}
    │   ├── plan-build-site.md
    │   └── plan-{domain}.md
    ├── impl/                   # Implementation tracking
    │   └── impl-all.md
    ├── research/               # Framework-specific research
    └── prompts/                # Prompts with {FRAMEWORK} variables
```

### How It Works

1. **Tier 1 (Reference)** and **Tier 2 (Blueprints)** live in a shared repository
2. **Tier 3 (Plans + Impl)** lives in each framework-specific repository
3. The shared repo is included as a **git submodule** in each framework repo
4. Updates to blueprints propagate to all implementations via `git submodule update`

### Directory Layout with Submodules

```
my-app-react/                   # React implementation
├── shared-context/             # Git submodule → shared specs repo
│   ├── reference/
│   └── blueprints/
├── context/                    # React-specific
│   ├── plans/
│   ├── impl/
│   └── prompts/
└── src/

my-app-vue/                     # Vue implementation
├── shared-context/             # Same git submodule → same specs
│   ├── reference/
│   └── blueprints/
├── context/                    # Vue-specific
│   ├── plans/
│   ├── impl/
│   └── prompts/
└── src/
```

### Benefits

| Benefit | How |
|---------|-----|
| **Same blueprints, different stacks** | Shared blueprints drive both implementations identically |
| **Apples-to-apples comparison** | Framework evaluation uses the same requirements |
| **Blueprint propagation** | Update blueprints once, all implementations pick up changes |
| **Independent divergence** | Each implementation's plans and tracking are independent |

### Setting Up Submodules

```bash
# Create the shared context repo
mkdir shared-context && cd shared-context
git init
mkdir -p reference specs
# ... add reference materials and specs ...
git add . && git commit -m "Initial shared context"

# In each framework repo
cd my-app-react
git submodule add <shared-context-repo-url> shared-context

cd my-app-vue
git submodule add <shared-context-repo-url> shared-context
```

> For the full multi-repo strategy reference, see `references/multi-repo-strategy.md`.

---

## Single-Repo Strategy

Most projects use a single repo. The context directory lives alongside the source code:

```
my-project/
├── CLAUDE.md               # Project root — Blueprint conventions
├── context/                # All Blueprint artifacts
│   ├── refs/
│   ├── specs/
│   ├── plans/
│   ├── impl/
│   └── prompts/
├── src/                    # Source code
├── tests/                  # Test files
└── package.json            # (or equivalent build config)
```

This is the default and recommended setup. Only use multi-repo when you need shared blueprints across implementations.

---

## Context Directory Anti-Patterns

### 1. Flat File Dump

**Wrong:** Dumping all documents into a single directory with no structure.
```
context/
├── auth-spec.md
├── auth-plan.md
├── auth-tracking.md
├── api-spec.md
├── api-plan.md
├── random-notes.md
└── old-stuff.md
```

**Right:** Use the standard directory structure. Blueprints in `blueprints/`, plans in `plans/`, tracking in `impl/`.

### 2. Missing CLAUDE.md Files

Without CLAUDE.md files, agents have no convention guidance when working in a directory. Every subdirectory of `context/` should have a CLAUDE.md.

### 3. Monolithic Documents

A single 2000-line blueprint file defeats progressive disclosure. Decompose into domains with a `blueprint-overview.md` index.

### 4. Stale Archives in Active Directories

Completed or archived content should be moved to `impl/archive/`, not left in the active tracking file consuming context window.

---

## Integration with Other Skills

### With `bp:blueprint-writing`

The context architecture defines WHERE blueprints live. The blueprint-writing skill defines HOW to write them. Blueprints go in `context/blueprints/` following the naming convention `blueprint-{domain}.md`.

### With `bp:impl-tracking`

Implementation tracking documents live in `context/impl/`. When they exceed 500 lines, compact them and archive the old version to `context/impl/archive/`.

### With `bp:validation-first`

Validation gate results are recorded in implementation tracking documents within the context structure. Phase gates reference specs by requirement ID.

### With `bp:methodology`

The context directory structure is established during the Spec phase of DABI and maintained throughout the entire lifecycle.

---

## Summary

1. **4-level hierarchy:** Index → Sub-documents → Agent discovery → Context efficiency
2. **Standard directory structure:** `refs/`, `specs/`, `plans/`, `impl/`, `prompts/` with CLAUDE.md in each
3. **CLAUDE.md is hierarchical:** Loads from current directory up to project root
4. **Multi-repo uses git submodules:** Shared specs (Tier 1-2) + framework-specific plans (Tier 3)
5. **Single-repo is the default:** Context directory alongside source code
6. **Agents spider out from indexes:** Read index first, then only relevant sub-documents
7. **Compact when large:** Archive resolved content, keep active files focused
