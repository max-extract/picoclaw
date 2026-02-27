---
name: github
description: Support MaxExtract release, incident, and change-log workflows with concise verified summaries.
---


**GitHub Workflow**

This skill provides operational guidance for `github`.

**Overview**

- Purpose: Summarize code and deployment deltas into operator-ready change narratives.
- Focus: impact area mapping, verification evidence, and residual risk notes.

**When to Use This Skill**

- User asks what changed in a release or hotfix.
- User asks what regressed after a commit range.
- User needs a quick deploy-ready changelog.

**Inputs**

- Required: incident/release context, impacted components, and requested audience detail.
- Evidence: command outputs, recent changes, and validation status.
- Decision context: containment urgency and acceptable risk level.

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

- If commit range is unclear, summarize last stable window and mark assumption.
- If verification evidence is missing, report change-only summary and set confidence low.

**Output contract**

- Always include:
- **What changed**
- **Why it changed**
- **Impact scope**
- **How it was validated**
- **Residual risk**
- Keep to short bullets and avoid non-actionable text.

**Safety guardrails**

- Never claim deployment success without explicit evidence.
- Never infer service impact from file names alone without a confidence note.
- Never use markdown headers or tables.

**Cross References**

- `../AGENTS.md`
- `workspace/skills/summarize/SKILL.md`
- `workspace/AGENT.md`
