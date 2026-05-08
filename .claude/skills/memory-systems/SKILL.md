---
name: memory-systems
description: This skill should be used when the user asks to "remember past work", "find what we did before", "search previous sessions", "persist knowledge across runs", or mentions cross-session memory, claude-mem, or looking up prior decisions/bugs/features.
---

# Memory Systems (claude-mem)

Agents running in this system have access to `claude-mem` — a persistent cross-session observation store mounted at `/claude-mem-data`. Use it to avoid re-solving problems already solved, to carry context across task runs, and to understand the project's decision history.

## When to Activate

Activate this skill when:
- A task references something that might have been worked on before
- Debugging a bug that may have appeared in prior runs
- Making an architectural decision that prior sessions may have already evaluated
- Needing to understand what changed recently in the project

## Architecture

`claude-mem` runs as an MCP subprocess (spawned by `caveman-shrink`) inside the `claude-runner` container. Data persists in the `/claude-mem-data` volume. Available as MCP tools — not HTTP.

```
claude CLI (in claude-runner container)
    └── MCP: claude-mem mcp  (subprocess via caveman-shrink)
            └── /claude-mem-data volume  (persistent SQLite)
```

## Memory Model

claude-mem stores **observations** — structured records of past work:

| Type | When recorded |
|------|---------------|
| `🔵 discovery` | New understanding, architecture findings |
| `🟣 feature` | Feature implemented |
| `🔴 bugfix` | Bug found and fixed |
| `✅ change` | Configuration or structural change |
| `⚖️ decision` | Architectural or design decision |
| `🚨 security_alert` | Security issue found |

Each observation has: ID, timestamp, title, narrative, facts, concepts, files referenced.

## 3-Layer Retrieval Workflow

**Never fetch full observations without filtering first — 10x token savings.**

### Step 1: Search — Get Index

```
search(query="authentication", project="AIDevelopmentWorkflow", limit=20)
```

Returns a lightweight table (~50-100 tokens/result):
```
| ID    | Time  | T  | Title                        | Read |
|-------|-------|----|------------------------------|------|
| #4190 | 1:10a | 🔴 | Fixed SSL cert bypass        | ~75  |
| #4180 | 1:07a | 🔵 | Claude CLI Docker auth       | ~120 |
```

Parameters:
- `query` — search term
- `project` — always set to `"AIDevelopmentWorkflow"` for this system
- `limit` — default 20, max 100
- `obs_type` — filter: `bugfix`, `feature`, `decision`, `discovery`, `change`
- `dateStart` / `dateEnd` — YYYY-MM-DD
- `orderBy` — `"date_desc"` (default), `"relevance"`

### Step 2: Timeline — Get Surrounding Context

```
timeline(anchor=4190, depth_before=3, depth_after=3, project="AIDevelopmentWorkflow")
```

Returns observations before and after the anchor to understand what led to and followed from a finding. Use when a single result doesn't give enough context.

### Step 3: Fetch — Full Details for Selected IDs Only

```
get_observations(ids=[4190, 4180])
```

Always use `get_observations` for 2+ items — single request vs N requests.

Returns full observation: narrative, facts, concepts, file list (~500-1000 tokens each).

## Common Patterns

**"Was this bug seen before?"**
```
search(query="<error message or symptom>", obs_type="bugfix", project="AIDevelopmentWorkflow")
```

**"What decisions were made about X?"**
```
search(query="<component/topic>", obs_type="decision", project="AIDevelopmentWorkflow")
```

**"What changed recently?"**
```
search(type="observations", dateStart="2026-05-01", project="AIDevelopmentWorkflow", limit=30)
```

**"Why does this code look like this?"**
```
search(query="<filename or function>", project="AIDevelopmentWorkflow")
# then timeline around the relevant result
```

## When NOT to Use

- Current session context is enough — don't search memory for information already in the conversation
- Simple factual lookups about the codebase — use `codebase-memory` MCP or Read tool instead
- Writing new observations — that's handled by the claude-mem worker automatically; agents don't need to write manually

## Memory Layers

| Layer | Implementation | Scope |
|-------|---------------|-------|
| Working | Active context window | Current turn only |
| Short-term | Message history + tool outputs | Current task run |
| Long-term | claude-mem `/claude-mem-data` | Cross-session, persisted |
| Codebase | `codebase-memory` MCP (`/cbm-cache`) | Static code knowledge graph |

Use `codebase-memory` MCP for "where is X defined in the code." Use `claude-mem` for "what happened when we worked on X."

## Gotchas

1. **Always set `project` filter**: Without it, results include observations from all projects. Always pass `project="AIDevelopmentWorkflow"`.
2. **Search before fetching**: Fetching observations directly by guessing IDs wastes tokens. Always search first to get the index.
3. **Observations reflect past state**: Memory records what was true when recorded. Code may have changed since. Verify current state against the actual files.
4. **`get_observations` for batches**: Fetching IDs one by one costs N round-trips. Always batch into a single `get_observations([id1, id2, ...])` call.
5. **timeline is for context, not just facts**: Use `timeline` when you need to understand WHY a decision was made, not just WHAT was decided.

## Integration

- filesystem-context - Use filesystem for within-session scratch; use claude-mem for cross-session persistence
- context-optimization - Memory retrieval results are tool outputs — apply observation masking after processing
- multi-agent-patterns - Sub-agents can each query claude-mem independently; no shared state conflicts since it's read-only during a task run

---

## Skill Metadata

**Created**: 2026-05-02
**Author**: AIDevelopmentWorkflow
**Version**: 1.0.0
