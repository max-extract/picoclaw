# Tmux

Perform minimal tmux diagnostics for stuck or long-running MaxExtract operations.



## tmux Diagnostics

This skill provides operational guidance for `tmux`.

## Overview

- Purpose: Provide minimal terminal-session diagnostics when direct execution is insufficient.
- Focus: session visibility, stuck-command inspection, and low-noise status checks.

## When to Use This Skill

- Long-running command appears stuck.
- User asks whether a background process is still active.
- Live output sampling is needed for troubleshooting.

## Inputs

- Required: user intent and desired outcome.
- Context: relevant environment/mode and current repository state.
- Evidence: latest command or test outputs used for conclusions.

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

- If tmux is unavailable, run equivalent direct command checks and report that tmux path was skipped.
- If pane output is empty, verify process state with non-tmux commands.

## Output contract

- Always include:
- **Session status**
- **Observed command state**
- **Confidence**
- **Next check**

## Safety guardrails

- Do not use tmux as default for routine monitoring.
- Do not mutate sessions unless user explicitly requests it.
- Never use markdown headers or tables.

## Cross References

- `workspace/AGENT.md`
- `config/skills/trading-dashboard/SKILL.md`
- `${MAXEXTRACT_ROOT}/scripts/`
