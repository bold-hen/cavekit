---
name: ck-review-branch
description: "Review the current branch end-to-end — kit compliance (Pass 1) plus code quality (Pass 2). Uses Karpathy guardrails as the Pass-1 filter. Optionally dispatches Codex for a second-opinion review."
argument-hint: "[--base REF] [--codex] [--strict]"
allowed-tools: ["Bash(git *)", "Bash(node ${CLAUDE_PLUGIN_ROOT}/scripts/cavekit-tools.cjs:*)", "Bash(${CLAUDE_PLUGIN_ROOT}/scripts/bp-config.sh:*)", "Read(*)", "Grep(*)", "Glob(*)", "Agent(ck:inspector,ck:verifier,ck:researcher)"]
---

# Cavekit Review Branch

Two-pass review of the current branch against kits and code-quality standards.
Use before merging to main, or as a tier gate during `/ck:make`.

## Setup

1. Compute the base ref:
   ```bash
   BASE="${ARGUMENTS_BASE:-$(git merge-base HEAD $(git symbolic-ref refs/remotes/origin/HEAD | sed 's@^refs/remotes/origin/@@'))}"
   git diff "$BASE"...HEAD --stat
   ```

2. Resolve execution model:
   ```bash
   EXECUTION_MODEL=$("${CLAUDE_PLUGIN_ROOT}/scripts/bp-config.sh" model execution)
   ```

## Pass 1 — Kit Compliance (Karpathy Pass 1)

Invoke the `karpathy-guardrails` skill first. For each file changed in the
diff, check:

- Does every diff line trace to a kit acceptance criterion? List the kit /
  R-ID / AC-ID for each hunk.
- Any silent assumptions? Any "while I'm here" edits? Any files touched
  outside the task's stated scope?
- Are the acceptance criteria verifiable, and did the author verify them?

Classify findings by severity: **CRITICAL** (blocks merge) / **IMPORTANT**
(should fix) / **MINOR** (optional).

If Pass 1 has any CRITICAL finding, **do not run Pass 2**. Report Pass 1 only
and stop.

## Pass 2 — Code Quality

Only runs if Pass 1 is clean of CRITICAL findings. Check:

- Naming, structure, and API shape.
- Test quality: do the new tests assert on the new behaviour, or just on
  implementation details?
- Error handling at system boundaries; no bare `try/except`.
- Security: input validation, secrets, authn/authz, injection.
- Observability: logs/metrics for new code paths.

Classify findings by the same severity scale.

## Optional Codex second pass

If `--codex` is passed and `codex` is on $PATH, dispatch an independent
review via the existing `codex-review.sh` script:

```bash
if [[ -x "${CLAUDE_PLUGIN_ROOT}/scripts/codex-review.sh" ]]; then
  "${CLAUDE_PLUGIN_ROOT}/scripts/codex-review.sh" --base "$BASE"
fi
```

Diff the two findings sets. Anything unique to Codex surfaces as an
`IMPORTANT` item in the final report. Agreement across both reviewers raises
a finding's confidence tier.

## Report shape

```
═══ Branch Review ═══
Base: <sha>  →  HEAD: <sha>
Files:  {N}   Tests: {N}   Diff lines: +{a} -{b}

Pass 1 — Kit Compliance
  CRITICAL: {n}
  IMPORTANT: {n}
  MINOR: {n}
  {findings, each referencing <kit>:<R-ID>[.<AC-ID>]}

Pass 2 — Code Quality     {skipped if Pass 1 blocked}
  CRITICAL: {n}
  IMPORTANT: {n}
  MINOR: {n}

Verdict: {PROCEED | BLOCKED}
```

With `--strict`, any IMPORTANT finding also blocks. Default gates only on
CRITICAL.

## Fix cycle (tier gate use)

When `/ck:review-branch` runs as a tier gate during `/ck:make` and the
verdict is `BLOCKED`:

1. Turn each blocking finding into a fix task description:
   ```
   FIX-{n}: {finding summary}
     Cite: {file:line}
     Kit: {kit}:R{id}{.AC-id?}
     Severity: CRITICAL|IMPORTANT
   ```
2. Hand the fix tasks back to the loop. The stop-hook will route the next
   wave with these as additional `ck:task-builder` prompts (use the same
   model tier the original task used).
3. Track cycle count across iterations (`review_fix_cycle` field in
   `.cavekit/state.md`). After each fix wave, re-run `/ck:review-branch`:
   - If clean → advance the tier.
   - If still blocked and cycle < 2 → another fix wave.
   - If still blocked and cycle == 2 → emit
     `ADVANCE_WITH_FINDINGS` with the list of unresolved items, log to
     `.cavekit/history/backprop-log.md` as candidate kit amendments, and
     let the tier advance. This matches the existing
     `bp_review_fix_cycle` guard from `scripts/codex-gate.sh`.

## Critical rules

- Read actual diff hunks before citing findings. Never cite a file you didn't
  read.
- Every finding must cite `file:line` and a kit R-ID.
- Do not attempt to fix findings here — that is a follow-up task. This
  command reports only.
