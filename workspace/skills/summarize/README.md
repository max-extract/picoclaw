# Summarize

Produce concise incident and deployment summaries with clear decisions and ownership.



## Incident And Deploy Summaries

This skill provides operational guidance for `summarize`.

## Overview

- Purpose: Convert operational context into concise handoff-ready summaries.
- Focus: status clarity, ownership/ETA, and single-decision outputs.

## When to Use This Skill

- User asks for incident recap.
- User asks for deployment outcome summary.
- User needs a handoff note for another operator.

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

- If root cause is unknown, state unknown explicitly and give containment plan.
- If data sources conflict, prefer DB for ranking and flag confidence.

## Output contract

- Always include:
- **What happened**
- **Current status**
- **Root cause confidence**
- **Next action**
- **Owner and ETA**
- Keep output short and operational.

## Safety guardrails

- Do not hide uncertainty.
- Do not provide multiple conflicting next decisions.
- Never use markdown headers or tables.

## Cross References

- `workspace/AGENT.md`
- `config/skills/trading-dashboard/SKILL.md`
- `config/skills/maxextract-infra/SKILL.md`
