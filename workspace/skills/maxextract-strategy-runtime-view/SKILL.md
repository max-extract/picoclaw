---
name: maxextract-strategy-runtime-view
description: Build a strategy-focused runtime view using bot state, health, and recent activity.
---


**MaxExtract Strategy Runtime View**

This skill provides operational guidance for `maxextract-strategy-runtime-view`.

**Overview**

- Purpose: Produce a strategy-centric runtime snapshot across markets and modes.
- Focus: health/status alignment, stale signals, and per-market behavior visibility.

**When to Use This Skill**

- Use this skill when strategy-level operational questions require focused checks.
- Use this skill when user asks for status, drift, config, or rollout confidence.

- User asks how a specific strategy is performing now.
- Bot seems alive but behavior looks off.
- Operations review needs strategy-first status.

**Inputs**

- Required: target scope (`paper`, `live`, or `all`) and query intent (health, ranking, digest, rollout).
- Data sources: runtime APIs, script outputs, and DB ranking path when available.
- Runtime context: current bot inventory and freshness of the latest sample.

**Primary workflow**

1. Pull bot inventory and API state.
2. Pull bot-level report for targeted strategy-market pairs.
3. Summarize status by strategy first, then by market.
4. Highlight stale or contradictory signals.

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

- If API endpoints fail, use DB-based scripts and mark source as `api_fallback`.
- If bot report fails for one market, continue with remaining markets.

**Output contract**

- Always include:
- **Strategy summary**
- **Per-market status**
- **Data source**
- **Freshness timestamp**
- **Next action** only when actionable.

**Safety guardrails**

- Do not call mutation endpoints.
- Prefer DB for historical ranking when API and DB disagree.

**Cross References**

- `workspace/AGENT.md`
- `workspace/USER.md`
- `config/skills/maxextract-infra/SKILL.md`
