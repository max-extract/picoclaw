---
name: maxextract-infra
description: Manage MaxExtract runtime bots via Coolify and runtime APIs with strict mutation gates.
---


**MaxExtract Infrastructure Operations**

This skill provides operational guidance for `maxextract-infra`.

**Overview**

- Purpose: Execute safe infra operations for MaxExtract runtime services with mutation gates.
- Focus: diagnostics first, action gating, rollout awareness, and blast-radius control.

**When to Use This Skill**

- User asks for deploy status or bot inventory.
- User asks for realtime runtime state or ranking.
- User requests a mutating infra action.

**Inputs**

- Required: target scope (`paper`, `live`, or `all`) and query intent (health, ranking, digest, rollout).
- Data sources: runtime APIs, script outputs, and DB ranking path when available.
- Runtime context: current bot inventory and freshness of the latest sample.

**Primary workflow**

1. Confirm scope and constraints.
2. Run minimal read-only checks first.
3. Execute focused commands for the requested outcome.
4. Return concise results with explicit confidence and risk notes.

**Quick Start**

- `cd "${PICOCLAW_ROOT:-$(pwd)}"`
- `MAXEXTRACT_USE_SSH=1 ./workspace/bin/me.sh me_bots_inventory.sh --context mycoolify --mode all --json`

**Commands**

- `git status --short`
- `MAXEXTRACT_ROOT="${MAXEXTRACT_ROOT:-$(cd .. && pwd)}"; rg -n "strategy" "$MAXEXTRACT_ROOT/strategies" "$MAXEXTRACT_ROOT/runtime" "$MAXEXTRACT_ROOT/scripts"`
- `MAXEXTRACT_ROOT="${MAXEXTRACT_ROOT:-$(cd .. && pwd)}"; rg -n "runtime" "$MAXEXTRACT_ROOT/strategies" "$MAXEXTRACT_ROOT/runtime" "$MAXEXTRACT_ROOT/scripts"`

**Examples**

```bash
cd "${PICOCLAW_ROOT:-$(pwd)}"
MAXEXTRACT_USE_SSH=1 ./workspace/bin/me.sh me_bots_digest.sh --context mycoolify --mode all --days auto
```

**Common failures and fixes**

- Missing or invalid env vars: verify script arguments and required env values.
- Partial data from one source: continue with available sources and mark missing fields as `n/a`.
- Endpoint or connectivity failure: use fallback path and label source confidence.

**Fallback behavior**

- If DB is unavailable, continue with API fallback and explain source downgrade.
- If one bot endpoint fails, keep partial output and mark missing fields `n/a`.
- If Coolify access fails, block mutating actions and return diagnostics-only result.

**Output contract**

- Always include:
- **Summary**
- **Source**
- **Window**
- **Bots**
- **Next action** when actionable
- For mutations, always include blast radius and rollback direction.

**Safety guardrails**

- Ask explicit confirmation before any mutation.
- Never fabricate healthy status when data is missing.
- Never use markdown headers or tables.

**Cross References**

- `workspace/AGENT.md`
- `workspace/USER.md`
- `workspace/memory/MEMORY.md`
