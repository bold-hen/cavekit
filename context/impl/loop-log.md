### Iteration 1 — 2026-03-20
- **Task:** T-004 — Merge main into existing worktree branch on build start
- **Tier:** 0
- **Status:** DONE
- **Files:** scripts/setup-build.sh
- **Validation:** Build P, Tests P (syntax check), Acceptance: merge on existing worktree P, first-time creation unchanged P
- **Next:** T-005, T-006, T-007, T-008

### Iteration 2 — 2026-03-20
- **Task:** T-005 (merge conflicts), T-006 (merge logging), T-007 (env forwarding), T-008 (env verification)
- **Tier:** 0
- **Status:** DONE
- **Files:** scripts/setup-build.sh
- **Validation:** Build P, Tests P (syntax), Acceptance: T-005 conflict handling P (abort+options+exit), T-006 merge logging P (up-to-date/merged output), T-007 symlinks on creation P, T-008 symlink verify on resume P
- **Note:** T-005/T-006 were already satisfied by T-004's implementation. T-007/T-008 use symlinks for auto-sync and are verified on every start.
- **Next:** T-001, T-002, T-003

### Iteration 3 — 2026-03-20
- **Task:** T-001 (changelog format), T-002 (draft generates changelog), T-003 (changelog append logic)
- **Tier:** 0
- **Status:** DONE
- **Files:** commands/draft.md, commands/revise.md
- **Validation:** Build P, Acceptance: template has ## Changelog P, revise appends entries P, since-last-revision from changelog P, append-only P
- **Next:** T-009 — Scan commits since last revision
