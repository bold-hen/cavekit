---
name: methodology
description: |
  Core Blueprint methodology — the master skill that teaches the DABI lifecycle
  and routes to all sub-skills. Covers the Specify Before Building principle, the scientific method analogy,
  the four-phase DABI lifecycle, decision matrix for when to use Blueprint, and build pipeline analogy.
  Trigger phrases: "use Blueprint", "blueprint methodology", "start Blueprint project", "blueprint methodology",
  "how should I structure this project for AI agents"
---

# Blueprint Methodology

## Core Principle: Specify Before Building

**Always define what you want before telling agents how to build it. Go through a blueprint stage — never jump straight from raw requirements to implementation.**

Blueprint is a methodology for building software with AI coding agents that **puts blueprints at the center of the development process — code is derived from them, not the other way around**. Whether starting from scratch or modernizing an existing system, the principle is the same:

- **Greenfield projects:** reference material → blueprints → code
- **Rewrites:** old code → blueprints → new code

In both cases, the blueprints become a living contract that agents consume to continuously build, validate, and refine the application.

### Why Blueprints Are the First-Class Citizen

| Property | Benefit |
|----------|---------|
| **Structured** | Organized as a navigable tree, enabling agents to load only what they need |
| **Human-legible** | Engineers can audit requirements at a higher level than code |
| **Stack-independent** | Decoupled from any single framework or language |
| **Independently evolvable** | Blueprints can be refined without touching implementation |
| **Verifiable** | Every requirement includes acceptance criteria agents can check |

> **Key Insight:** Well-written blueprints with strong validation make your application reproducible — any agent can rebuild it from the blueprints alone. Think of it as continuous regeneration.

---

## The Scientific Method Analogy

LLMs are inherently non-deterministic — like running an experiment, each individual call may yield different results. But through the right methodology — clear hypotheses, controlled conditions, and repeated trials — we extract reliable, reproducible outcomes from a stochastic process.

**Blueprint applies the scientific method to software construction — hypothesize, test, observe, refine.**

| Layer | Analogy | What It Does |
|-------|---------|-------------|
| **LLM calls** | Individual experiments | Each run may produce different results; no single output is authoritative |
| **Blueprints** | Hypotheses | Define what you expect to observe — the predicted behavior |
| **Validation gates** | Controlled conditions | Ensure reproducibility by constraining what counts as a valid outcome |
| **Convergence loops** | Repeated trials | Build statistical confidence through successive passes |
| **Implementation tracking** | Lab notebook | Record what was tried, what worked, and what failed |
| **Revision** | Revising the hypothesis | When results contradict expectations, update the theory upstream |

The outcome: a disciplined, repeatable engineering process layered on top of probabilistic generation.

---

## The 5 DABI Phases

DABI stands for **Draft, Architect, Build, Inspect**. Each phase has dedicated prompts that drive it.

| Phase | Input | Output | AI Role | Human Role |
|-------|-------|--------|---------|------------|
| **Draft** | Source materials, domain knowledge, existing systems | Implementation-agnostic blueprints | Extract requirements, structure knowledge | Verify blueprints capture intent accurately |
| **Architect** | Blueprints + framework research | Framework-specific implementation plans | Design architecture, break down work, order steps | Approve architectural choices |
| **Build** | Plans + blueprints | Working code + tests + tracking docs | Write code, run tests, check against blueprints | Watch for drift and blockers |
| **Inspect** | Failed validations, gaps, manual fixes | Updated blueprints/plans via revision | Identify root causes, propagate fixes upstream | Evaluate outcomes, set priorities |
| **Monitor** | Running application, git history | Issues, anomalies, progress reports | Scan for regressions, surface metrics | Interpret reports, guide next steps |

### Phase Transitions

Each phase has **gate conditions** that must be met before moving to the next:

1. **Draft → Architect:** All domains have blueprints with testable acceptance criteria. Human has reviewed for completeness.
2. **Architect → Build:** Plans reference blueprints, define implementation sequence, and include test strategies. Architecture decisions validated.
3. **Build → Inspect:** Code builds, tests pass at current coverage level, implementation tracking is up to date.
4. **Inspect → Monitor:** Convergence detected (changes decreasing iteration-over-iteration). Remaining changes are trivial.
5. **Monitor → Draft (cycle):** Gap found or new requirement identified. Revise blueprints and restart the cycle.

The **Inspect** phase is where the human serves as **reviewer and decision-maker**, not hands-on coder. You monitor the process, request changes as needed, and make systemic improvements to blueprints and prompts.

> For the full DABI phase reference, see `references/dabi-phases.md`.

---

## Decision Matrix: When to Use Blueprint

### Full Blueprint

Use when the project has significant scope, evolving requirements, or needs autonomous agent execution.

| Indicator | Threshold |
|-----------|-----------|
| Codebase size | 50+ source files |
| Requirements | Evolving, multi-domain |
| Agent coordination | Multi-agent or multi-prompt pipelines |
| Environment | Production, security-sensitive, brownfield |
| Team structure | Multi-team or cross-team |
| Execution mode | Long-running autonomous work (overnight, unattended) |

**What you get:** Full DABI lifecycle, context directory with blueprints/plans/impl tracking, prompt pipeline, convergence loops, revision, validation gates.

### Lightweight Blueprint

Use when scope is moderate — too complex for ad-hoc but not worth a full pipeline.

| Indicator | Threshold |
|-----------|-----------|
| Codebase size | 5-50 files |
| Requirements | Mostly clear, focused |
| Agent coordination | Single agent, possibly with sub-agents |
| Execution mode | Interactive with occasional iteration loops |

**What you do:**
1. Write a focused `context/blueprints/blueprint-task.md` capturing requirements
2. Add a `context/plans/plan-task.md` sequencing the implementation
3. Skip full DABI — just run an iteration loop against the plan

This is the "Blueprint floor" — most of the benefit without the overhead of a full multi-phase pipeline.

### Skip Blueprint

Use when the task is trivially small.

| Indicator | Threshold |
|-----------|-----------|
| Codebase size | Less than 5 files |
| Task type | One-off tools, simple bug fixes, exploratory prototypes |
| Implementation | Fits comfortably in one agent session without needing external references |

**Heuristic:** If the whole task fits in one context window with room to spare, full Blueprint adds more overhead than value.

### Growth Path

Start with lightweight Blueprint even if the project is small. If the scope expands, you already have the structure in place to scale up. It is much harder to retrofit blueprints onto a large codebase than to grow a blueprint directory from the beginning.

---

## The CI Pipeline Analogy

Blueprint mirrors a **build pipeline** — each stage transforms input into validated output, with feedback loops that propagate corrections upstream:

```
Traditional CI/CD:
  Code → Build → Test → Deploy

Blueprint AI Pipeline:
  Blueprint Change
    → Generate Plans (iteration loop)
    → Generate Implementation (iteration loop)
    → Validate (Tests + Review)
    → Human Audit (Monitor & Steer)
    → [Gap Found]
    → Revise
    → Blueprint Change (cycle repeats)
```

Every stage can run as an iteration loop — the same prompt executed repeatedly until output stabilizes. The iteration loop is what transforms nondeterministic LLM output into predictable, validated software.

### The Iteration Loop

The iteration loop is the fundamental execution unit in Blueprint. Execute the same prompt against the same codebase multiple times until the delta between runs approaches zero.

**Mechanics:**
1. Execute a prompt against the current codebase
2. The agent inspects git history and tracking documents to understand what has already been done
3. The agent applies changes and commits its progress
4. Return to step 1

**Convergence signal:** A shrinking volume of modifications across successive passes — the diff gets smaller each time until only cosmetic changes remain. You are looking for diminishing returns, not absolute zero.

**When the loop isn't stabilizing, the problem is upstream — fix the inputs (specs, validation, coordination), not the iteration count.**

If the diff is not shrinking between runs:
- Blueprints are ambiguous (agents interpret them differently each time)
- Validation criteria are too loose (the agent has no way to confirm it got things right)
- Multiple agents are overwriting each other's work (ownership boundaries are unclear)

---

## Cross-References to Sub-Skills

Blueprint is composed of techniques that work together. This methodology skill is the index — each sub-skill below is self-contained but cross-references others.

### Foundation Skills

| Skill | Purpose | When to Use |
|-------|---------|-------------|
| `bp:blueprint-writing` | Write implementation-agnostic blueprints with testable acceptance criteria | Draft phase — always the first step |
| `bp:context-architecture` | Organize context for progressive disclosure | Project setup and ongoing maintenance |
| `bp:impl-tracking` | Track implementation progress, dead ends, test health | Build and Inspect phases |
| `bp:validation-first` | Design validation gates agents can execute | All phases — validation is continuous |

### Pipeline Skills

| Skill | Purpose | When to Use |
|-------|---------|-------------|
| `bp:prompt-pipeline` | Design numbered prompt pipelines for DABI | Setting up automation |
| `bp:revision` | Trace bugs back to blueprints and fix at the source | Inspect phase — after finding gaps |
| `blueprint:brownfield-adoption` | Adopt Blueprint on existing codebases | Starting Blueprint on legacy projects |

### Advanced Skills

| Skill | Purpose | When to Use |
|-------|---------|-------------|
| `bp:peer-review` | Use a second agent to challenge the first | Quality gates, architecture review |
| `blueprint:speculative-pipeline` | Stagger pipeline stages for parallelism | Optimizing long pipelines |
| `bp:convergence-monitoring` | Detect convergence vs ceiling | Monitoring iteration loops |
| `blueprint:documentation-inversion` | Turn documentation into agent-consumable skills | Library/module documentation |

### Integration with Existing Skills

Blueprint works **with** existing skills, not as a replacement:

| Existing Skill | Blueprint Integration |
|----------------|-----------------|
| `superpowers:brainstorming` | Use during blueprint generation to explore requirements |
| `superpowers:writing-plans` | Use during plan generation for structured planning |
| `superpowers:test-driven-development` | TDD-within-Blueprint: blueprint acceptance criteria become failing tests |
| `superpowers:verification-before-completion` | Use for gate validation in every phase |
| `superpowers:executing-plans` | Use during implementation phase |
| `superpowers:dispatching-parallel-agents` | Use for agent team coordination |

---

## Quick Start

### For a New Project (Greenfield)

1. **Set up context directory:**
   ```
   context/
   ├── refs/           # Source materials (PRDs, language specs, research)
   ├── blueprints/     # Implementation-agnostic blueprints
   ├── plans/          # Framework-specific implementation plans
   ├── impl/           # Living implementation tracking
   └── prompts/        # DABI pipeline prompts
   ```

2. **Write blueprints** from your reference materials (see `bp:blueprint-writing`)
3. **Generate plans** from blueprints (see `bp:prompt-pipeline`)
4. **Implement** with validation gates (see `bp:validation-first`)
5. **Track progress** in implementation documents (see `bp:impl-tracking`)
6. **Iterate** — when gaps are found, revise blueprints (see `bp:revision`)

### For an Existing Project (Brownfield)

1. **Set up context directory** (same structure as above)
2. **Designate existing codebase as reference material**
3. **Generate blueprints from code** (see `blueprint:brownfield-adoption`)
4. **Validate blueprints match behavior** — run tests against generated blueprints
5. **Proceed with normal DABI** — future changes flow through blueprints first

---

## Summary

Blueprint is not a tool — it is a methodology. The core loop is simple:

1. **Describe what you want** (blueprints with testable criteria)
2. **Let agents build it** (plans → implementation → validation)
3. **Fix the blueprints, not the code** (revision)
4. **Repeat until converged** (iteration loops)

Agents become more capable the more precisely you constrain them — clear blueprints, automated validation, and structured iteration loops let them operate with increasing autonomy. None of this eliminates the need for software engineers. Your judgment on architecture, your ability to write precise blueprints, and your instinct for what "done" looks like are the inputs that make the whole system function. Blueprint is a force multiplier: one engineer's clarity of thought, scaled across an entire implementation pipeline.
