---
name: peer-review-loop
description: >
  Peer Review Ralph Loop — combines Blueprint blueprints with a Ralph Loop and true cross-model peer review
  review using Codex (OpenAI) as an MCP server. Claude builds from specs; Codex reviews peer reviewly.
  Covers setup, MCP configuration, iteration patterns, convergence detection, and completion criteria.
  Triggers: "peer review loop", "ralph loop with codex", "blueprint ralph", "peer review build loop",
  "cross-model loop", "codex peer reviewer", "blueprint to ralph loop"
---

# Peer Review Loop — Blueprint + Ralph Loop + Codex Peer reviewer

Run a Blueprint blueprint through a Ralph Loop where Claude builds and Codex peer reviewly reviews.
This is the most rigorous automated quality process available: every few iterations, a completely
different model (different training data, different biases, different blind spots) challenges
your implementation.

---

## Why This Works

| Factor | Single-Model Loop | Peer Review Loop |
|--------|-------------------|------------------|
| Blind spots | Same model, same blind spots every iteration | Two models catch different classes of issues |
| Blueprint drift | Builder may silently deviate from blueprint | Peer reviewer checks blueprint compliance explicitly |
| Quality floor | Converges to "good enough for one model" | Converges to "survives cross-examination" |
| Dead ends | May retry failed approaches | Peer reviewer flags repeated patterns |

---

## Architecture

```
┌─────────────────────────────────────────────────────┐
│                   Ralph Loop                         │
│  (Stop hook feeds same prompt each iteration)        │
│                                                      │
│  ┌──────────┐    ┌──────────────┐    ┌────────────┐ │
│  │  Claude   │───▶│ Build from   │───▶│  Commit    │ │
│  │  (Build)  │    │ blueprint +  │    │  changes   │ │
│  └──────────┘    └──────────────┘    └──────┬─────┘ │
│       ▲                                      │       │
│       │                                      ▼       │
│  ┌──────────┐    ┌──────────────┐    ┌────────────┐ │
│  │  Fix      │◀──│ Parse        │◀──│  Codex MCP │ │
│  │  findings │    │ findings     │    │  (Review)  │ │
│  └──────────┘    └──────────────┘    └────────────┘ │
│                                                      │
│  Completion: all blueprint requirements met +         │
│              no CRITICAL/HIGH findings               │
└─────────────────────────────────────────────────────┘
```

---

## Quick Start

```bash
# Basic: implement a blueprint with peer review
/bp:peer-review-loop context/blueprints/blueprint-auth.md

# With options
/bp:peer-review-loop context/blueprints/blueprint-api.md --max-iterations 20 --codex-model gpt-5.4-mini

# Review-only mode (review existing code, don't build new)
/bp:peer-review-loop context/blueprints/blueprint-api.md --review-only

# Review every iteration instead of every 2nd
/bp:peer-review-loop context/blueprints/blueprint-auth.md --review-interval 1
```

---

## What the Command Does

1. **Validates** the blueprint file exists and Codex CLI is installed
2. **Configures** Codex as an MCP server in `.mcp.json` (if not already configured)
3. **Builds** a Ralph Loop prompt that embeds:
   - The blueprint path and related plan/impl files
   - Instructions to alternate between build and review iterations
   - The peer review prompt template for Codex
   - Completion criteria tied to blueprint acceptance criteria
4. **Starts** the Ralph Loop via the stop hook mechanism

---

## Codex MCP Server

The command configures Codex as an MCP server automatically:

```json
{
  "mcpServers": {
    "codex-peer reviewer": {
      "command": "codex",
      "args": ["mcp-server", "-c", "model=\"gpt-5.4\""]
    }
  }
}
```

Claude calls this MCP server on review iterations to get peer review feedback. The MCP server
exposes Codex as a tool that accepts prompts and returns responses — Claude sends the blueprint +
code diff, Codex returns findings.

### Changing the Codex Model

Use `--codex-model` to specify which OpenAI model Codex should use:

```bash
/bp:peer-review-loop blueprint.md --codex-model gpt-5.4-mini    # faster, cheaper
/bp:peer-review-loop blueprint.md --codex-model gpt-5.4          # default, most capable
```

---

## Iteration Pattern

```
Iteration 1: BUILD  — Read blueprint, implement first requirement
Iteration 2: REVIEW — Call Codex MCP, get findings, fix CRITICAL/HIGH
Iteration 3: BUILD  — Continue implementing, address remaining findings
Iteration 4: REVIEW — Call Codex MCP again, new findings on new code
...
Iteration N: BUILD  — All requirements met, all findings fixed
             → outputs <promise>SPEC COMPLETE</promise>
```

The review interval is configurable. Default is every 2nd iteration.
Use `--review-interval 1` for maximum rigor (review every iteration).

---

## Peer Review Findings File

Review findings are tracked in `context/peer-review-findings.md`:

```markdown
# Peer Review Findings

## Latest Review: Iteration 4 — 2026-03-14T10:30:00Z
### Reviewer: Codex (gpt-5.4)

| # | Severity | File | Issue | Status |
|---|----------|------|-------|--------|
| 1 | CRITICAL | src/auth.ts:L42 | Missing input validation on token | FIXED |
| 2 | HIGH | src/auth.ts:L67 | Race condition in session refresh | FIXED |
| 3 | MEDIUM | src/auth.ts:L15 | Unused import | NEW |
| 4 | LOW | src/auth.ts:L3 | Comment typo | WONTFIX |

## History
### Iteration 2
| # | Severity | File | Issue | Status |
|---|----------|------|-------|--------|
| 1 | CRITICAL | src/auth.ts:L20 | SQL injection in login query | FIXED |
```

---

## Completion Criteria

The loop exits when the completion promise is output. The prompt instructs Claude
to ONLY output it when ALL of these are true:

- All blueprint requirements (R-numbers) have been implemented
- All acceptance criteria pass
- No CRITICAL or HIGH peer review findings remain unfixed
- Build passes
- Tests pass
- At least one review iteration completed with no new CRITICAL/HIGH findings

---

## Modes

### Build + Review (default)
Alternates between implementing blueprint requirements and calling Codex for review.
Use for greenfield implementation from a blueprint.

### Review Only (`--review-only`)
Skips building. Each iteration calls Codex to review existing code against the blueprint,
then fixes issues found. Use when code already exists and you want peer review QA.

---

## Prerequisites

1. **Codex CLI installed**: `npm install -g @openai/codex`
2. **OpenAI API key configured**: Codex needs authentication (via `codex login` or env var)
3. **Blueprint context directory**: Blueprint file must exist at the given path
4. **Ralph Loop plugin**: The ralph-loop plugin must be installed (provides the stop hook)

---

## Convergence Signals

The peer review loop has converged when:
- Codex's findings drop to zero or only LOW/MEDIUM severity
- Code diffs between iterations are minimal
- All blueprint requirements confirmed as met by both Claude and Codex

If the loop hits max iterations without converging:
- Check `context/peer-review-findings.md` for persistent issues
- Consider whether the blueprint needs clarification
- Run `/bp:revise` to trace issues back to blueprints

---

## Cross-References

- **peer-review** — The underlying peer review patterns and prompt templates
- **convergence-monitoring** — How to detect convergence vs ceiling
- **validation-first** — Validation gates that run on every build iteration
- **impl-tracking** — How implementation progress is tracked across iterations
- **Ralph Loop** — The underlying Ralph Loop mechanism
