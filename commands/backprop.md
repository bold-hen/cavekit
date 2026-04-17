---
name: ck-backprop
description: "Trace a bug backwards to the kit that should have caught it, propose an acceptance-criterion fix, add a regression test, and log the trace."
argument-hint: "[FAILURE_DESCRIPTION | --from-flag | --from-finding F-ID]"
allowed-tools: ["Bash(node ${CLAUDE_PLUGIN_ROOT}/scripts/cavekit-tools.cjs:*)", "Bash(git *)", "Read(*)", "Write(*)", "Edit(*)", "Grep(*)", "Glob(*)"]
---

# Cavekit Backpropagation

Manual entry point to the `backpropagation` skill. Invoke after a test
failure, review finding, or user-reported defect.

## Load the skill

Invoke the `backpropagation` skill. Follow its six steps exactly:

1. **TRACE** — map the failure to a kit R-ID.
2. **ANALYZE** — classify: missing_criterion / incomplete_criterion /
   wrong_criterion / missing_requirement.
3. **PROPOSE** — draft the spec change. **Ask the user to approve before
   writing.**
4. **GENERATE** — write a regression test that currently fails.
5. **VERIFY** — patch code until the test passes.
6. **LOG** — append an entry to `.cavekit/history/backprop-log.md`.

## Input modes

- **No args** — first, check whether `.cavekit/.auto-backprop-pending.json`
  exists. If it does, auto-consume it exactly as if `--from-flag` were
  passed, and print a one-line notice that you are honouring the pending
  flag from a failed test run. If it does not, ask the user to paste the
  failure details.
- **FAILURE_DESCRIPTION** — a free-text paragraph. Use it verbatim as the
  trace input.
- **`--from-flag`** — read `.cavekit/.auto-backprop-pending.json`. Use the
  recorded `command` and `failure_excerpt`. Delete the flag when done
  (the auto-backprop hook is idempotent, so deletion is mandatory here —
  otherwise the next stop-hook iteration will fire the directive again).
- **`--from-finding F-ID`** — read the `/ck:review-branch` or Codex review
  output and pull the finding by ID. Inline-ingest the finding's summary
  + file:line citation as the trace input.

## After the fix

Report back:

```
═══ Backprop complete ═══
Classification: {class}
Kit: {kit} → {R-ID}{optional AC-ID}
Regression test: {path}::{test name}
Fix commit: {sha}
Pattern category: {category}
Log entry: .cavekit/history/backprop-log.md #{id}
```

If three entries in one session share the same `pattern_category`, print a
warning and recommend a cross-kit amendment at the brainstorming / sketch
layer.

## Critical rules

- Never amend a kit without explicit user approval.
- The regression test must fail before the fix and pass after. Verify both.
- Commit the test separately from the fix so the log is auditable.
- Do not mark the loop's current task complete as part of backprop — that is
  the task-builder's job.
