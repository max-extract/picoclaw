# Maxextract strategy rollout check

Validate strategy rollout readiness from config, runtime health, and deploy routing perspective.



## MaxExtract Strategy Rollout Check

This skill provides operational guidance for `maxextract-strategy-rollout-check`.

## Overview

- Purpose: Validate readiness before and after strategy rollouts.
- Focus: deploy-route impact, gating checks, and rollback triggers.

## When to Use This Skill

- Use this skill when strategy-level operational questions require focused checks.
- Use this skill when user asks for status, drift, config, or rollout confidence.

- User asks if a strategy change is safe to ship.
- Push includes strategy, runtime, or DB migration changes.
- Post-deploy behavior needs quick go or no-go signal.

## Inputs

- Required: target scope (`paper`, `live`, or `all`) and query intent (health, ranking, digest, rollout).
- Data sources: runtime APIs, script outputs, and DB ranking path when available.
- Runtime context: current bot inventory and freshness of the latest sample.

## Primary workflow

1. Validate changed files against deploy routing rules.
2. Run strategy profile checks.
3. Run bot health and strategy state checks.
4. Return rollout verdict with explicit blockers.

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

- If full checks are unavailable, run changed-area checks and mark risk as elevated.
- If migration cannot run, block rollout and report missing prerequisite.

## Output contract

- Always include:
- **Readiness verdict**
- **Blocking checks**
- **Passed checks**
- **Post-deploy watchpoints**
- **Rollback trigger**

## Safety guardrails

- Never approve rollout when migration status is unknown.
- Never treat missing health data as healthy.

## Cross References

- `workspace/AGENT.md`
- `workspace/USER.md`
- `config/skills/maxextract-infra/SKILL.md`
