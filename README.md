# Blueprint

Claude Code plugin + parallel agent launcher for blueprint-driven development with automated iteration loops.

## Install

```bash
git clone https://github.com/JuliusBrussee/sdd-os.git ~/.blueprint
cd ~/.blueprint
./install.sh
```

This will:
1. Register the Blueprint plugin with Claude Code
2. Install the `blueprint` CLI command globally
3. Make all scripts executable

## Terminal: parallel agent launcher

```bash
blueprint --monitor                     # interactive picker → launch agents in tmux
blueprint --monitor --expanded          # one tmux window per agent with dashboards
blueprint --status                      # check progress from any terminal
blueprint --analytics                   # trends across cycles
blueprint --kill                        # stop everything, clean up worktrees
```

### Default mode (`--monitor`)

Interactive multi-select picker shows all build sites:
- **Available** — ready to launch (pre-selected)
- **In Progress** — select to resume from existing worktree
- **Done** — struck through (archived sites)

Selected sites each get:
- Their own **git worktree** (branch: `blueprint/<site-name>`)
- A **tmux pane** running Claude Code with `/bp:build`
- Auto-layout: horizontal for 2-3 agents, tiled for 4+
- Live status bar showing per-site progress

Staggered launch (5s between agents) to avoid API rate limits.

### Expanded mode (`--monitor --expanded`)

One tmux window per site with the full 3-pane layout:
- **Left (70%)** — Claude Code running `/bp:build`
- **Top-right** — live progress: tasks done, tiers, progress bar
- **Bottom-right** — live activity: iteration log, git commits

Switch between windows with `Ctrl-b <number>`.

### Analytics (`--analytics`)

Parses loop logs across all cycles and worktrees:
- Iterations to convergence per cycle
- Task outcomes (done/partial/blocked)
- Failure patterns and dead ends
- Tier distribution
- Completion velocity (tasks/iteration, success rate)

## Claude: slash commands

```
/bp:draft       →  draft blueprints (the WHAT)
/bp:architect   →  generate build site (the ORDER)
/bp:build       →  ralph loop (the BUILD)
/bp:inspect     →  gap analysis + peer review (the CHECK)
/bp:merge       →  blueprint-aware branch integration (the SHIP)
```

### 1. Draft — write blueprints

```bash
/bp:draft                       # interactive — asks what to build
/bp:draft context/refs/         # from PRDs, API docs, research
/bp:draft --from-code           # from existing codebase
```

Decomposes your project into domains. Each domain gets a blueprint with R-numbered requirements and testable acceptance criteria.

### 2. Architect — generate build site

```bash
/bp:architect                   # all blueprints
/bp:architect --filter v2       # only v2 blueprints
```

Reads blueprints, breaks requirements into tasks, maps dependencies, organizes into tiers.

### 3. Build — run the loop

```bash
/bp:build                       # implement everything
/bp:build --peer-review         # add Codex (GPT-5.4) review
/bp:build --max-iterations 30
```

Each iteration: read site → find next unblocked task → read blueprint → implement → validate → commit → loop.

### 4. Inspect — post-loop check

```bash
/bp:inspect                     # gap analysis + peer review
```

### 5. Merge — blueprint-aware branch integration

```bash
/bp:merge                       # merge all blueprint/* branches into main
```

After parallel execution, each site lives on its own `blueprint/<name>` branch. `/bp:merge` integrates them back into main:

1. Surveys all branches — commits, file overlaps, dependency order
2. Reads the **blueprints and impl tracking** for each branch
3. Merges in order: infrastructure → features → UI
4. Resolves conflicts by understanding what each blueprint intended — **keeps all features from all branches**
5. Validates after each merge (build, tests, blueprint requirements)
6. Cleans up worktrees and branches

## File structure

```
context/
├── blueprints/         # Blueprints (persist across cycles)
│   ├── blueprint-overview.md
│   └── blueprint-{domain}.md
├── sites/              # Build sites (one per plan)
│   ├── build-site-ui-v2.md
│   └── archive/        # Completed sites
├── impl/               # Progress (archived between cycles)
│   ├── impl-{domain}.md
│   ├── loop-log.md
│   └── archive/
└── refs/               # Reference materials
```

## All commands

| Command | Description |
|---------|-------------|
| **`/bp:draft`** | Draft blueprints |
| **`/bp:architect`** | Generate build site |
| **`/bp:build`** | Ralph Loop implementation |
| **`/bp:inspect`** | Gap analysis + peer review |
| **`/bp:merge`** | Blueprint-aware branch integration |
| `/bp:progress` | Check site progress |
| `/bp:gap-analysis` | Compare built vs intended |
| `/bp:revise` | Trace manual fixes to blueprints |
| `/bp:help` | Show usage |

| CLI | Description |
|-----|-------------|
| `blueprint --monitor` | Interactive picker → parallel agents in tmux |
| `blueprint --monitor --expanded` | One window per agent with dashboards |
| `blueprint --status` | Check site progress |
| `blueprint --analytics` | Trends across cycles |
| `blueprint --merge` | Shows branches ready to merge (use `/bp:merge` in Claude) |
| `blueprint --kill` | Stop all agents, clean worktrees |

## Example

See [example.md](example.md) for full sample conversations.

## License

MIT
