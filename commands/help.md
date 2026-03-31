---
name: bp-help
description: Show Blueprint commands and usage
---

# Blueprint

## The Workflow

```
/bp:draft      →  write blueprints (the WHAT)
/bp:architect   →  generate site (the ORDER)
/bp:build       →  ralph loop (the BUILD)
/bp:inspect     →  gap analysis + peer review (the CHECK)
```

### Quick Mode

```bash
/bp:quick <describe your feature>   # runs all 4 phases end-to-end, no stops
/bp:quick add JWT auth --skip-inspect
/bp:quick refactor the API layer --peer-review
```

Streamlined draft + architect (no interactive Q&A, no user gates) followed by the full build and inspect. Best for small-to-medium features where you trust the decomposition.

## Commands

### `/bp:draft` — Write Blueprints

```bash
/bp:draft                        # interactive — asks what to build
/bp:draft context/refs/          # from reference materials (PRDs, docs)
/bp:draft --from-code            # from existing codebase (brownfield)
/bp:draft --filter v2            # only generate v2 blueprints
```

Decomposes your project into domains, writes `context/blueprints/blueprint-{domain}.md` files with R-numbered requirements and testable acceptance criteria.

### `/bp:architect` — Generate Site

```bash
/bp:architect                    # generates site from all blueprints
/bp:architect --filter v2        # only v2 blueprints
```

Reads blueprints, decomposes requirements into tasks, organizes into dependency tiers. Writes `context/sites/build-site.md`. No domain plans — just tasks and dependencies.

### `/bp:build` — Run the Loop

```bash
/bp:build                       # auto-parallel build from site
/bp:build --filter v2           # scope to v2
/bp:build --peer-review         # add Codex (GPT-5.4) review
/bp:build --max-iterations 30   # iteration limit
/bp:build --peer-review --codex-model gpt-5.4-mini
```

Auto-archives any previous cycle, then builds the site. Automatically parallelizes by grouping ready tasks into a few coherent work packets based on shared files, subsystem, and complexity. Progresses through tiers autonomously without manual intervention. If multiple build sites exist, asks which one to implement.

With `--peer-review`: alternates build and review iterations, calling Codex via MCP.

### `/bp:inspect` — Post-Loop Inspection

```bash
/bp:inspect                     # inspect everything from last loop
/bp:inspect --filter v2         # only v2
```

Runs after build completes. Does two things:
1. **Gap analysis** — compares built code against every blueprint requirement and acceptance criterion
2. **Peer review** — finds bugs, security issues, performance problems, quality gaps

Produces a verdict: APPROVE / REVISE / REJECT with prioritized findings.

### `/bp:progress` — Check Progress

```bash
/bp:progress                    # show site progress
/bp:progress --filter v2
```

Shows tasks done/ready/blocked, progress bar, current tier, and next tasks.

### `/bp:codex-review` — On-Demand Codex Review

```bash
/bp:codex-review                  # review current tier diff
/bp:codex-review --base v1.0     # review diff against a specific ref
```

Sends the current diff to Codex for adversarial review. Outputs findings in Blueprint format and appends them to `context/impl/impl-review-findings.md`. Requires Codex CLI to be installed.

### Maintenance (optional)

| Command | When |
|---------|------|
| `/bp:codex-review` | On-demand Codex adversarial review of current diff |
| `/bp:gap-analysis` | After a loop — compare built vs intended |
| `/bp:revise` | After manual code fixes — trace back to blueprints |
| `/bp:compact-specs` | When impl tracking files exceed ~500 lines |
| `/bp:archive-loop` | Manually archive a loop cycle (build does this automatically) |
| `/bp:next-session` | Generate a handoff document for next session |

### Legacy (advanced)

These still work but are superseded by the three main commands:

| Command | Replaced by |
|---------|-------------|
| `/blueprint init` | `/bp:draft` creates directories automatically |
| `/blueprint spec-from-refs` | `/bp:draft context/refs/` |
| `/blueprint spec-from-code` | `/bp:draft --from-code` |
| `/blueprint plan-from-specs` | `/bp:architect` (generates site directly, no domain plans) |
| `/blueprint implement` | `/bp:build` (one task at a time vs full loop) |
| `/blueprint spec-loop` | `/bp:build` |
| `/blueprint peer-review-loop` | `/bp:build --peer-review` |
| `/blueprint quick` | `/bp:quick` (end-to-end with no stops) |

## Skills (reference docs)

| Skill | Topic |
|-------|-------|
| `bp:methodology` | Core DABI lifecycle |
| `bp:blueprint-writing` | How to write blueprints with testable criteria |
| `bp:peer-review` | Cross-model review patterns |
| `bp:peer-review-loop` | Ralph Loop + Codex architecture |
| `bp:validation-first` | Every requirement must be auto-testable |
| `bp:convergence-monitoring` | Detecting if loop is converging or stuck |
| `bp:revision` | Tracing bugs back to blueprints |
| `bp:context-architecture` | Organizing context/ for AI agents |
| `bp:impl-tracking` | Progress tracking and dead ends |
| `blueprint:brownfield-adoption` | Adopting Blueprint on existing codebases |
| `bp:prompt-pipeline` | Designing prompt sequences |
| `blueprint:speculative-pipeline` | Staggered pipeline execution |
| `blueprint:documentation-inversion` | Agent-first documentation |
