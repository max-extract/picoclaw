# Trading dashboard

Query MaxExtract runtime bots for health, state, and performance with safe fallbacks.



## Trading Dashboard Operations

This skill provides operational guidance for `trading-dashboard`.

## Overview

- Purpose: Monitor runtime bot health, state, and performance from a strategy/operator perspective.
- Focus: inventory, realtime API state, digest generation, and ranking context.

## When to Use This Skill

- User asks for health across bots.
- User asks for strategy or market performance snapshot.
- User needs unified operator digest.

## Inputs

- Required: target scope (`paper`, `live`, or `all`) and query intent (health, ranking, digest, rollout).
- Data sources: runtime APIs, script outputs, and DB ranking path when available.
- Runtime context: current bot inventory and freshness of the latest sample.

## Primary workflow

1. Confirm scope and constraints.
2. Run minimal read-only checks first.
3. Execute focused commands for the requested outcome.
4. Return concise results with explicit confidence and risk notes.

## Quick Start

- `cd "${PICOCLAW_ROOT:-$(pwd)}"`
- `MAXEXTRACT_USE_SSH=1 ./workspace/bin/me.sh me_bots_inventory.sh --context mycoolify --mode all --json`

## Commands

- `git status --short`
- `MAXEXTRACT_ROOT="${MAXEXTRACT_ROOT:-$(cd .. && pwd)}"; rg -n "strategy" "$MAXEXTRACT_ROOT/strategies" "$MAXEXTRACT_ROOT/runtime" "$MAXEXTRACT_ROOT/scripts"`
- `MAXEXTRACT_ROOT="${MAXEXTRACT_ROOT:-$(cd .. && pwd)}"; rg -n "runtime" "$MAXEXTRACT_ROOT/strategies" "$MAXEXTRACT_ROOT/runtime" "$MAXEXTRACT_ROOT/scripts"`

## Examples

```bash
cd "${PICOCLAW_ROOT:-$(pwd)}"
MAXEXTRACT_USE_SSH=1 ./workspace/bin/me.sh me_bots_digest.sh --context mycoolify --mode all --days auto
```

## Common failures and fixes

- Missing or invalid env vars: verify script arguments and required env values.
- Partial data from one source: continue with available sources and mark missing fields as `n/a`.
- Endpoint or connectivity failure: use fallback path and label source confidence.

## Fallback behavior

- If API endpoints fail, use configured public runtime URLs and mark source.
- If DB ranking fails, continue with API fallback and state limitation clearly.

## Output contract

- Always include:
- **Summary**
- **Source** (`db` or `api_fallback`)
- **Window** (`auto` or explicit days)
- **Bots**
- **Next action** when actionable
- Use `n/a` for missing metrics.

## Safety guardrails

- Never execute mutation endpoints unless explicitly requested.
- Never hardcode bot names or endpoints in user-facing output.
- Prefer DB as authority for historical ranking conflicts.

## Cross References

- `config/skills/maxextract-infra/SKILL.md`
- `workspace/USER.md`
- `../AGENTS.md`
