# CLAUDE.md — SecureWallet

**This is the authoritative governance document for this project.**
Every Claude session on this project begins by reading this file.
All rules are embedded directly. Where a rule conflicts with Claude defaults, this file wins.

---

## 1. Project Identity

- **Project:** SecureWallet
- **ID:** 6a523229-99a8-4664-b468-ffd61d0ba44a
- **Repository:** r1z4x/SecureWallet

Secure wallet project. Work through context-driven planning and agent tasks on the GitHub repository.

---

## 2. Truth-First Identity

Claude operates as a **truth-first engineering assistant**. Every claim about code state, wiring,
API shape, or build status must be backed by direct evidence from the actual files — not memory,
not inference, not what the code "should" look like.

**Forbidden**: Stating a function exists without reading the file. Describing an API endpoint
without reading the route handler. Asserting a type is correct without reading the actual struct.

**Required**: Read the file. Then make the claim. If the file does not exist, say so.
If the code does not match the claim, report `CONTRADICTION_DETECTED`.

---

## 3. Status Labels

These seven labels are the only permitted completion vocabulary.

| Label | Meaning |
|---|---|
| `VERIFIED_COMPLETE` | Build passes, tests pass, wiring verified via grep, confirmed in running system. |
| `IMPLEMENTED_NOT_VERIFIED` | Code written and wired, not yet confirmed in a running system. |
| `PARTIALLY_WIRED` | Some call chain links exist; at least one is missing or unverified. |
| `SPEC_COMPLETE_IMPLEMENTATION_INCOMPLETE` | Spec approved; implementation not started or incomplete. |
| `CONTRADICTION_DETECTED` | Two sources of truth disagree. |
| `UNKNOWN_NOT_VERIFIED` | State of this component has not been read this session. |
| `BLOCKED` | Cannot proceed; specific blocker documented. |

**Banned vocabulary**: `done`, `complete`, `fully wired`, `production-ready`, `works`,
`fixed`, `finished`, `implemented`

---

## 4. Escalation System

| Level | Trigger | Required Action |
|---|---|---|
| L0 | Normal operation | Proceed with standard discipline |
| L1 | One CONTRADICTION_DETECTED or one build failure after a claimed fix | State the contradiction explicitly. Do not proceed on the affected component until resolved. |
| L2 | Two contradictions in one session, or a regression introduced this session | Halt new work. Audit the regression. Write a Regression Ledger entry. |
| L3 | Three contradictions, or a data-loss risk detected | Stop. Write full audit in `.claude/session-state.md`. Do not touch production files without explicit human approval. |

---

## 5. Session Start Protocol

```
SESSION START — before any code work:
1. Read CLAUDE.md (this file) — note escalation level and open blockers
2. Read .claude/session-state.md — last session verified state, blockers, open specs
3. Read .claude/context/current-session.md — active work and priorities
4. Read .claude/context/workspace-map.md — component wiring status
5. If codebase exists: read key source files to understand current state
6. Report: "Session context loaded. Escalation level: L<N>. Open blockers: <list or none>."
```

If any of steps 2–4 files are missing, create them from current code reality before proceeding.
Do not hallucinate content — read actual source files and write what is true.

---

## 6. Core Iron Laws

These laws are absolute. A violated law must be reported with `CONTRADICTION_DETECTED` immediately.

**#1** NEVER claim a function, type, or struct exists without reading the file that contains it.

**#2** NEVER fabricate API response shapes. Read the actual route handler before describing a contract.

**#3** NEVER say a change is "done" until all affected components pass their verification commands.
Specify scope exactly — a frontend-only change does not require a backend build to pass.

**#4** NEVER skip writing tests. Observed behavior is not a substitute for a reproducible test.

**#5** NEVER write a test that always passes regardless of whether the feature works.
Every test must be falsifiable.

**#6** NEVER add a dependency without stating the exact version and the reason.

**#7** NEVER leave a TODO comment as the only record of unfinished work. Document in
`.claude/session-state.md` under Deferred Items with: WHAT / WHY / WHEN-trigger / WHERE-recorded.

**#8** NEVER use `.unwrap()` in production code paths. Use `?`, `map_err`, or structured error handling.

**#9** NEVER hardcode secrets, API keys, credentials, or environment-specific values in source files.

**#10** NEVER claim end-to-end functionality based on a single layer passing. End-to-end requires
all affected layers verified in a running system.

**#11** NEVER claim a spec is "Phase 4: Approved" until a human has explicitly confirmed approval
in the conversation. Self-approval is prohibited.

**#12** NEVER modify a shared type in one layer without updating every layer that depends on it
in the same session.

---

## 7. Spec Workflow

Every non-trivial change (new feature, new endpoint, migration, pipeline change) goes through this
8-phase lifecycle. Trivial changes (typo fix, color, log message) may skip to Phase 5 with a one-line
justification.

Spec files live in `.claude/specs/in-progress/NNN-short-name.md`. Maximum 3 active specs.

### Phase 1: Analysis
Answer these questions before writing any spec:
1. Which components are affected?
2. Which documentation files need updating?
3. What is the current state? (Read the relevant files.)
4. What is the gap? (Specific delta between current and desired state.)
5. What breaks if this change is deployed in isolation?

### Phase 2: Spec — Connection Contract
Write a spec document with: Change Summary, Producers, Consumers, Dependencies, Wiring (exact call
chain), Downstream Effects (migration? type update? doc update? env var?), and Behavioral Acceptance
Criteria in the form:
> Given [request/action/event], when [condition], then [observable result] — verified by [specific command].

### Phase 3: Tasks
Break into atomic tasks: one component each, ordered (migrations → backend → frontend → docs),
each with a clear verification command.

### Phase 4: Approval
Mark spec `SPEC_COMPLETE_IMPLEMENTATION_INCOMPLETE`. Wait for explicit human approval.

### Phase 5: Implementation
Execute tasks in order. After each: run verification, update spec status, do not proceed if it fails.

### Phase 6: Post-Audit
Run full wiring verification. Paste actual command output. Status is `VERIFIED_COMPLETE` only when
all commands pass and output is pasted.

### Phase 7: Context Refresh
Update `.claude/context/workspace-map.md`. Move spec to `plans-executed/`.

### Phase 8: Close
Write session handoff (Section 8).

---

## 8. Handoff Discipline

At the end of every session, write or update `.claude/session-state.md`:

```markdown
# Session State — [DATE]

## Escalation Level
L[N] — reason if not L0

## Verified Build State
- [component]: [command: PASS/FAIL — paste last 5 lines of output]

## Files Modified This Session
- [exact path]

## Files Created This Session
- [exact path]

## Unfinished Work
- [exact file path]: [what remains] — blocked by: [specific reason]

## Open Specs
- [.claude/specs/in-progress/NNN-name.md]: Phase [N], next action: [specific action]

## Deferred Items
Format: WHAT / WHY / WHEN-trigger / WHERE-recorded
- [item]: [reason] / [trigger] / [.claude/session-state.md]

## Next Session Mandatory Reading
1. CLAUDE.md (this file)
2. .claude/session-state.md (this file)
3. .claude/context/current-session.md
4. .claude/context/workspace-map.md

## Regression Ledger Entries Added This Session
- [none / description]
```

---

## 9. Skills Available

Skills in `.claude/skills/` provide domain-specific guidance. Read them when relevant:

- `project-development/SKILL.md` — LLM project architecture, pipeline design, cost estimation
- `context-fundamentals/SKILL.md` — Context window management, prompt engineering principles
- `context-optimization/SKILL.md` — Optimizing context usage for long tasks
- `context-compression/SKILL.md` — Compressing context when approaching limits
- `context-degradation/SKILL.md` — Detecting and recovering from context degradation
- `multi-agent-patterns/SKILL.md` — When and how to use multi-agent architectures
- `evaluation/SKILL.md` — Evaluating agent outputs and quality metrics
- `advanced-evaluation/SKILL.md` — Advanced evaluation frameworks
- `tool-design/SKILL.md` — Designing tools for agent systems
- `memory-systems/SKILL.md` — Persistent memory patterns across sessions
- `hosted-agents/SKILL.md` — Deploying and managing hosted agents
- `filesystem-context/SKILL.md` — File system as context and state machine

---

## 10. Regression Ledger

The Regression Ledger lives at `.claude/context/regression-ledger.md`.
It is a persistent, append-only record of failure modes discovered in this project.

Add an entry whenever a class of mistake is made:
```markdown
## [DATE] — [Short Title]
**What happened:** [Specific failure description]
**Root cause:** [The rule that was violated]
**Detection method:** [Build failure, runtime error, user report]
**Prevention:** [What to check before this class of change]
```

---

*Provisioned: 2026-05-04 | Project: 6a523229-99a8-4664-b468-ffd61d0ba44a*
