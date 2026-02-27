---
name: maxextract-strategy-map
description: Build a clear map of active MaxExtract strategies, markets, and runtime deployments.
---


**MaxExtract Strategy Map**

This skill provides operational guidance for `maxextract-strategy-map`.

**Overview**

- Purpose: Build a strategy-to-market and runner-to-deployment map for quick orientation.
- Focus: config discovery, runner linkage, and deployment coverage gaps.

**When to Use This Skill**

- Use this skill when strategy-level operational questions require focused checks.
- Use this skill when user asks for status, drift, config, or rollout confidence.

- User asks what strategies are active.
- User needs a strategy-to-market view before making changes.
- User wants a quick map of paper versus live coverage.

**Inputs**

- Required: target scope (`paper`, `live`, or `all`) and query intent (health, ranking, digest, rollout).
- Data sources: runtime APIs, script outputs, and DB ranking path when available.
- Runtime context: current bot inventory and freshness of the latest sample.

**Primary workflow**

1. List strategy folders and strategy YAML files.
2. Map runners to strategy slugs.
3. Pull runtime inventory from scripts.
4. Produce a compact strategy map with gaps.

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

- If remote inventory is unavailable, report strategy map from repository files only.
- If runner discovery is ambiguous, report unresolved slug names explicitly.

**Output contract**

- Always return:
- **Strategies discovered**
- **Markets per strategy**
- **Runner coverage**
- **Deployed bots**
- **Gaps**
- Never fabricate bot instances.

**Safety guardrails**

- Do not modify deployment state.
- Do not infer live deployment from paper-only evidence.

**Cross References**

- `workspace/AGENT.md`
- `workspace/USER.md`
- `config/skills/maxextract-infra/SKILL.md`
