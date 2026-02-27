---
name: skill-creator
description: Create or refactor PicoClaw skills into consistent long-form, executable runbook format.
---


**Skill Creator**

This skill provides operational guidance for `skill-creator`.

**Overview**

- Purpose: Define and enforce the canonical long-form skill authoring standard.
- Focus: structure consistency, executable snippets, and documentation quality gates.

**When to Use This Skill**

- User asks to add a new skill.
- Existing skill is short or inconsistent.
- Skill formatting causes rendering or usability issues.

**Inputs**

- Required: target files/skills and formatting or behavior expectations.
- Constraints: loader compatibility, markdown conventions, and section order policy.
- Verification: lint/structure checks and link/path correctness.

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

- If exact commands are unknown, provide safe placeholders and mark them clearly.
- If overlap with existing skill exists, merge or reference instead of duplicating.

**Output contract**

- Every created/updated skill must include:
- valid frontmatter
- executable command examples
- explicit fallback path
- explicit output contract
- explicit safety guardrails

**Safety guardrails**

- Never use markdown headers with `#` in skill content.
- Never use pipe-table formatting.
- Keep names slug-safe and aligned with folder name.

**Cross References**

```

```bash
rg -n '^#' workspace/skills/*/SKILL.md
rg -n 'pipe-table' workspace/skills/*/SKILL.md
```
