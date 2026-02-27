# Release ops

Run release readiness checks with explicit verification and risk reporting.



## Release Ops

This skill provides operational guidance for `release-ops`.

## Overview

- Purpose: Assess release readiness with explicit verification evidence.
- Focus: changed-area checks, risk framing, and go/no-go clarity.

## When to Use This Skill

- User asks for release readiness.
- CI is red near release cut.
- Packaging/lint/test state must be summarized quickly.

## Inputs

- Required: incident/release context, impacted components, and requested audience detail.
- Evidence: command outputs, recent changes, and validation status.
- Decision context: containment urgency and acceptable risk level.

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

- If primary checks fail, switch to the closest read-only fallback command.
- If confidence remains low, report unknowns explicitly and provide the next best check.

## Output contract

- Always include summary, source, and next action when actionable.
- Never fabricate values; use `n/a` for unavailable metrics.

## Safety guardrails

- Prefer read-only diagnostics unless mutation is explicitly requested.
- Require explicit confirmation for risky or mutating operations.

## Cross References

- `workspace/AGENT.md`
- `workspace/USER.md`
- `config/skills/maxextract-infra/SKILL.md`
