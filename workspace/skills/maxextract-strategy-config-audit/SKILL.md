---
name: maxextract-strategy-config-audit
description: Audit MaxExtract strategy YAML configuration for consistency and risky drift.
---


**MaxExtract Strategy Config Audit**

This skill provides operational guidance for `maxextract-strategy-config-audit`.

**Overview**

- Purpose: Detect config drift and risky parameter deltas across strategy YAML files.
- Focus: base/mode/market consistency, missing files, and override sanity.

**When to Use This Skill**

- Use this skill when strategy-level operational questions require focused checks.
- Use this skill when user asks for status, drift, config, or rollout confidence.

- User suspects config drift or inconsistent behavior by market.
- Strategy change is planned and baseline consistency is needed.
- Runtime output suggests a wrong parameter set.

**Inputs**

- Required: target scope (`paper`, `live`, or `all`) and query intent (health, ranking, digest, rollout).
- Data sources: runtime APIs, script outputs, and DB ranking path when available.
- Runtime context: current bot inventory and freshness of the latest sample.

**Primary workflow**

1. Verify required files exist for each strategy.
2. Compare `strategy.base.yaml`, `paper.yaml`, `live.yaml`.
3. Compare market-level overrides against base and mode files.
4. Report suspicious deltas and missing files.

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

- If profile tests cannot run, perform static file diff and mark validation as partial.
- If a strategy has non-standard layout, report it without forcing normalization.

**Output contract**

- Always include:
- **Missing files**
- **High-risk parameter drift**
- **Mode mismatch notes**
- **Validation source**
- Use `n/a` where a file is absent.

**Safety guardrails**

- Do not auto-correct YAML values.
- Treat `live.yaml` differences as intentional until proven otherwise.

**Cross References**

- `workspace/AGENT.md`
- `workspace/USER.md`
- `config/skills/maxextract-infra/SKILL.md`
