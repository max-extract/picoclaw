# Maxextract strategy backtest sync

Compare runtime strategy behavior with research backtests to detect silent drift.



## MaxExtract Strategy Backtest Sync

This skill provides operational guidance for `maxextract-strategy-backtest-sync`.

## Overview

- Purpose: Compare runtime behavior against research/backtest assumptions to catch drift.
- Focus: metric deltas, window alignment, and confidence-qualified conclusions.

## When to Use This Skill

- Use this skill when strategy-level operational questions require focused checks.
- Use this skill when user asks for status, drift, config, or rollout confidence.

- Live or paper behavior diverges from expected edge.
- Strategy update was made without explicit backtest sync.
- User asks if research and runtime are still aligned.

## Inputs

- Required: target scope (`paper`, `live`, or `all`) and query intent (health, ranking, digest, rollout).
- Data sources: runtime APIs, script outputs, and DB ranking path when available.
- Runtime context: current bot inventory and freshness of the latest sample.

## Primary workflow

1. Identify strategy and market under review.
2. Run relevant backtest script.
3. Compare core metrics and signal assumptions with runtime metrics.
4. Flag likely drift points and priority fixes.

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

- If backtest dependencies fail, compare config and signal inputs only.
- If runtime history window is too short, mark confidence as low.

## Output contract

- Always include:
- **Backtest window**
- **Runtime window**
- **Metric deltas**
- **Likely drift cause**
- **Confidence**

## Safety guardrails

- Do not declare strategy invalid from one short window.
- Separate data-quality issues from logic drift.

## Cross References

- `workspace/AGENT.md`
- `workspace/USER.md`
- `config/skills/maxextract-infra/SKILL.md`
