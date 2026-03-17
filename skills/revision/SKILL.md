---
name: revision
description: >
  The technique of tracing bugs and manual fixes back to blueprints and prompts, then fixing at the source
  so the iteration loop can reproduce the fix autonomously. Covers the 6-step revision process,
  commit classification, blueprint-level root cause analysis, and regression test generation.
  Trigger phrases: "revise", "revision", "trace bug to blueprint", "fix the blueprint not the code",
  "why did this bug happen", "update blueprints from bug"
---

# Revision: Tracing Bugs Back to Blueprints

In Blueprint, revision means tracing a production defect upstream through the blueprint chain until you find the gap that allowed it. In practice, when the built software has bugs or gaps, you trace the issue back to the blueprints and prompts and fix at the source -- not just in code.

**Key insight:** When a fix lives only in code with no corresponding blueprint update, the next iteration loop may reintroduce the same defect. The goal is that blueprints plus the iteration loop can reproduce any fix autonomously.

---

## 1. Why Revision Matters

Without revision, every bug fix is a one-off patch. The next time the iteration loop runs, it may reintroduce the bug because nothing in the blueprints or plans prevents it.

With revision:
- Bug fixes become **blueprint improvements** that persist across all future iterations
- The iteration loop becomes **self-correcting** -- it learns from every manual intervention
- Blueprints become **progressively more complete** over time
- The gap between "what blueprints describe" and "what works" **shrinks monotonically**

```
Without revision:
  Bug found -> Fix code -> Bug may return next iteration

With revision:
  Bug found -> Fix code -> Update blueprint -> Re-run iteration loop -> Fix emerges from blueprints alone
```

---

## 2. The 6-Step Revision Process

This is the complete process for tracing a bug back to its blueprint-level root cause and closing the loop.

### Step 1: Identify and Fix the Defect

Locate the bug -- whether through manual testing, automated failures, user reports, or monitoring alerts -- and resolve it through normal debugging. This produces a working code change, but the job is far from done: until the underlying blueprint gap is closed, this fix is fragile.

```bash
# The fix produces commits that we will analyze
git log --oneline -5
# a1b2c3d Fix: connection pool exhaustion under concurrent load
# e4f5g6h Fix: missing rate limit headers in API responses
```

### Step 2: Analyze What the Blueprint Missed

This is the pivotal step. Ask: **"Where in the blueprint chain did this requirement slip through?"**

Break the analysis into four dimensions:
- **WHAT** changed (files, functions, observable behavior)
- **WHY** it was wrong (which assumption proved false)
- **The RULE** (the invariant that should have been stated)
- **The LAYER** (which blueprint, plan, or prompt should have contained this)

Example analysis:

```markdown
## Revision Analysis: Database Connection Pooling

**WHAT changed:** Added pool size limits and idle timeout in `src/db/pool.ts`
**WHY:** The data layer blueprint assumed unlimited connections; under load the database
         rejected new connections once the server-side limit was reached
**RULE:** "The database module MUST configure a bounded connection pool with
          idle timeout and max-connection limits matching the deployment target"
**LAYER:** blueprint-data.md (no mention of pool configuration), plan-data.md (no task for pool tuning)
**Blueprint implications:** Add requirement R5 to blueprint-data.md covering connection pool settings
```

### Step 3: Update the Blueprint

Add the missing requirement or constraint to the appropriate blueprint file. Focus on acceptance criteria that are concrete enough for the iteration loop to act on:

```markdown
# In context/blueprints/blueprint-data.md, add:

### R5: Database Connection Pool Configuration
**Description:** The database module must use a bounded connection pool
with configurable limits to prevent resource exhaustion under load.
**Acceptance Criteria:**
- [ ] Maximum pool size is configurable and defaults to a sensible value
- [ ] Idle connections are reaped after a configurable timeout
- [ ] Pool exhaustion returns a clear error rather than hanging indefinitely
- [ ] Connection health checks run before returning a connection from the pool
**Dependencies:** R1 (database client setup), R2 (environment configuration)
```

### Step 4: Propagate Changes to Plans and Tracking

Trace the blueprint update through every downstream context file:

1. **Identify affected plan files:** Which plans govern the changed source paths?
2. **Update plans:** Add or close tasks reflecting the new requirement.
3. **Update impl tracking:** Record the revision event and its root cause.
4. **Annotate:** Mark updated sections with revision metadata so future reviews can trace lineage.

```markdown
# In context/plans/plan-data.md, add:

### T-DATA-005: Configure bounded connection pool
- **Status:** DONE (revised from manual fix a1b2c3d)
- **Blueprint:** R5 in blueprint-data.md
- **Files:** src/db/pool.ts
- **Acceptance criteria:**
  - [ ] Max pool size enforced
  - [ ] Idle timeout configured
  - [ ] Exhaustion handled gracefully
```

### Step 5: Apply Systemic Prompt Improvements (If Pattern Detected)

When the defect represents a recurring class of problem rather than a one-off, elevate the fix to the prompt level so it applies across all domains:

**Signs you are looking at a pattern:**
- The same category of bug has surfaced in more than one module
- The gap is structural (e.g., no specs anywhere address resource limits)
- A missing validation gate allowed the issue through

**Example systemic fix:**

```markdown
# In prompt 003, add to the validation section:

## Resource Management Validation
For every external resource integration, verify:
- [ ] Connection or handle limits are bounded and configurable
- [ ] Idle resources are cleaned up on a timeout
- [ ] Exhaustion scenarios return actionable errors
- [ ] Resource lifecycle is covered by tests under load
```

### Step 6: Verify and Lock In

Run the iteration loop against the updated blueprints to prove the fix emerges from blueprints alone, then generate regression tests to prevent future recurrence:

```bash
# Proof step: remove the manual fix and re-run from specs
git stash  # temporarily remove the manual fix
iteration-loop context/prompts/003-generate-impl-from-plans.md -n 5 -t 1h
# Verify the fix appears in the generated implementation

# If it does NOT, the blueprint update is insufficient -- return to Step 3
```

Once verified, create regression tests:

```bash
# Generate tests targeting the updated blueprint
{TEST_COMMAND} --blueprint context/blueprints/blueprint-data.md

# Or manually create a regression test
# tests/db/connection-pool-limits.test.ts
```

The regression tests should:
- Map directly to the acceptance criteria from Step 3
- Fail if the fix is reverted
- Run as part of the standard test suite going forward

---

## 3. Revision Analysis (Automated)

The revision analysis automates Steps 2-4 by examining recent git history.

### 3.1 Classify Commits

Analyze recent commits and classify each as:

| Classification | Meaning | Action |
|---------------|---------|--------|
| **Manual fix** | Human or interactive agent fixed a bug | Trace back to blueprint -- this is a revision target |
| **Iteration loop** | Automated iteration loop made the change | No action -- this is the system working as intended |
| **Infrastructure** | Build config, CI, tooling changes | No action -- not blueprint-related |

**How to classify:**
- Commits from iteration loop sessions have predictable patterns (automated commit messages, batch changes)
- Manual fixes are typically single-issue, focused commits with descriptive messages
- Infrastructure changes touch config files, build scripts, CI pipelines

### 3.2 Analyze Each Manual Fix

For each commit classified as a manual fix, determine:

```markdown
## Commit: abc1234 "Fix: auth token not refreshing on 401"

### WHAT changed
- File: src/auth/client.ts
- Function: handleApiResponse()
- Behavior: Added 401 detection and token refresh logic

### WHY it was wrong
- The auth module did not handle 401 responses
- Tokens would expire and never refresh, causing cascading auth failures

### RULE (invariant that should have been specified)
- "Authentication tokens must be refreshed automatically on 401 responses"

### LAYER (which context file should have caught this)
- blueprint-auth.md: Missing requirement for error-based token refresh
- plan-auth.md: No task for 401 handling

### Blueprint Implications
- Add R7 to blueprint-auth.md: Token Refresh on Authentication Failure
- Add T-AUTH-007 to plan-auth.md: Implement token refresh on 401
```

### 3.3 Discover Affected Plan Files

Dynamically discover which plan files govern the changed source paths:

```
Changed file: src/auth/client.ts
  -> Matches pattern: src/auth/*
  -> Governed by: plan-auth.md
  -> Blueprint: blueprint-auth.md

Changed file: src/data/api.ts
  -> Matches pattern: src/data/*
  -> Governed by: plan-data.md
  -> Blueprint: blueprint-data.md
```

Use file ownership tables (from prompts) or directory conventions to map source files to plan/blueprint files.

### 3.4 Update Context Files

For each revision target, update:

1. **Blueprint file:** Add missing requirement with acceptance criteria
2. **Plan file:** Add task referencing the new requirement
3. **Impl tracking:** Record the revision event

```markdown
# In context/impl/impl-auth.md, add:

## Revision Log
| Date | Commit | Issue | Blueprint Update | Plan Update |
|------|--------|-------|-------------|-------------|
| 2026-03-14 | abc1234 | 401 not handled | R7 added to blueprint-auth.md | T-AUTH-007 added |
```

### 3.5 Run Tests

After updating context files, run the test suite to verify nothing broke:

```bash
{BUILD_COMMAND}
{TEST_COMMAND}
```

### 3.6 Generate Regression Tests

For each revision target, generate a regression test that:
- Tests the specific acceptance criteria from the new blueprint requirement
- Would fail if the fix were reverted
- Is included in the standard test suite going forward

---

## 4. Patterns and Anti-Patterns

### Signs the process is working

| Pattern | What You Observe |
|---------|--------------------|
| **Declining manual intervention** | Each iteration cycle requires fewer hand-applied fixes because blueprints capture more of the ground truth |
| **Broader blueprint coverage per fix** | A single revision event adds constraints that block an entire family of related defects, not just one |
| **Cross-domain prevention** | Prompt-level adjustments made after a bug in one module prevent analogous bugs from appearing in other modules |
| **Autonomous reproducibility** | After a blueprint update, the iteration loop independently produces the same correction that a human applied manually |

### Warning signs and remedies

| Anti-Pattern | Symptom | Remedy |
|-------------|---------|-----|
| **Code-only patches** | The same category of defect resurfaces across iterations | Follow the full 6-step process; never stop after the code fix in Step 1 |
| **Overly specific blueprint additions** | Each revision prevents only the exact bug encountered, while slight variations slip through | Formulate the RULE as a general invariant, not a narrow patch |
| **Skipping verification** | Blueprints are updated but nobody confirms the iteration loop can reproduce the fix independently | Always execute Step 6; a blueprint that does not drive correct generation is incomplete |
| **Brittle over-specification** | Blueprints dictate implementation minutiae, causing breakage on minor refactors | Constrain the WHAT and WHY; leave the HOW to the implementation |
| **Accumulated revision debt** | A backlog of manual fixes sits un-traced, growing with each sprint | Set a cadence (e.g., end of each iteration) to clear the backlog; debt compounds quickly |

---

## 5. When NOT to Revise

Not every code fix needs revision:

- **One-off environment issues** (wrong config, missing dependency) -- these are infrastructure, not blueprint gaps
- **Typos and formatting** -- trivial fixes that do not reflect missing requirements
- **Exploratory changes** during prototyping -- blueprints are still being formed
- **Performance optimizations** that do not change behavior -- unless performance is a blueprint requirement

**Rule of thumb:** If the iteration loop could plausibly reintroduce the bug, revise. If not, skip it.

---

## 6. Revision and Convergence

Revision directly improves convergence:

```
Iteration 1: 350 lines changed, 8 manual fixes needed
  -> Revise all 8 fixes into blueprints
Iteration 2: 140 lines changed, 3 manual fixes needed
  -> Revise 3 fixes
Iteration 3: 30 lines changed, 1 manual fix needed
  -> Revise 1 fix
Iteration 4: 10 lines changed, 0 manual fixes needed
  -> Convergence achieved
```

Every revision cycle tightens the blueprints, so the iteration loop settles into a stable solution in fewer passes. If convergence is not improving, the most likely cause is that manual fixes are being applied without tracing them back to blueprints.

**Stalled convergence paired with ongoing manual fixes is a clear sign of revision debt.** The blueprints have not absorbed the lessons from past corrections, so the loop keeps regenerating flawed output that demands human repair.

---

## 7. Integration with Other Blueprint Skills

- **Convergence monitoring:** Use `bp:convergence-monitoring` to detect when manual fixes are decreasing (good) or increasing (revision debt).
- **Prompt pipeline:** Revision may trigger changes to prompts (Step 6), which affects the `bp:prompt-pipeline` design.
- **Validation-first design:** Stronger validation gates catch issues earlier, reducing the need for revision.
- **Gap analysis:** Systematic gap analysis (`/bp:gap-analysis`) identifies revision targets proactively, rather than waiting for bugs.

---

## Cross-References

- **Convergence patterns:** See `references/convergence-patterns.md` for how revision drives convergence.
- **Prompt pipeline:** See `bp:prompt-pipeline` skill for how prompt 006 (rewrite pattern) implements automated revision.
- **Impl tracking:** See `bp:impl-tracking` skill for the revision log format in implementation tracking documents.
- **Validation gates:** See `bp:validation-first` skill for validation layers that catch issues before they require revision.
