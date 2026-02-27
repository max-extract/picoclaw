---
name: database
description: Access MaxExtract historical ranking data with DB-first logic and explicit API fallback.
---


**Database Ranking Operations**

This skill provides operational guidance for `database`.

**Overview**

- Purpose: Provide DB-first historical ranking and PnL interpretation with explicit fallbacks.
- Focus: ROI authority, data-source arbitration, and window-aware summaries.

**When to Use This Skill**

- User asks for ROI or PnL ranking.
- User asks for historical trade summaries.
- Data source confidence must be explicit.

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

- If DB query errors, return API fallback output and include DB error reason.
- If history window is unclear, use `--days auto`.

**Output contract**

- Always include:
- **Summary**
- **Source** (`db` or `api_fallback`)
- **Window**
- **Top bots**
- **Bottom bots** when requested
- Use `n/a` for unavailable metrics.

**Safety guardrails**

- Do not mix DB and API ranking without declaring authority.
- Do not infer ROI from missing trade history.
- Never use markdown headers or tables.

**Cross References**

- `config/skills/maxextract-infra/SKILL.md`
- `config/skills/trading-dashboard/SKILL.md`
- `workspace/USER.md`
