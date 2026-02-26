---
name: maxextract-infra
description: Manage MaxExtract runtime bots via Coolify and bot backend APIs.
---

**REMINDER: Never use # headers or pipe tables in output. Use bold lines and bullet lists.**

**Scope**

- Runtime bots only (`paper`, `live`): EMA Until Expiry, Conviction, Latency Lite.
- Exclude recorders, DB services, and cross-arb unless explicitly requested.

**Intent Router**

- `inventory` or `deploy status`:
  - run inventory script
- `realtime state`:
  - run API state script
- `roi/pnl ranking`:
  - run DB ROI script first
- `executive digest`:
  - run digest script
- `periodic update`:
  - run periodic report script
- `mutating action` (`restart`, `deploy`, `switch`):
  - run diagnostics first
  - ask confirmation before action

**Operational Contract**

1. Route from intent to minimal command set.
2. Prefer DB for realized ranking; fallback to API if DB unavailable.
3. Keep diagnostics read-only by default.
4. Ask explicit confirmation before deploy/restart/stop.
5. Return Telegram-safe, compact output.

**Orchestration References**

- Agent policy:
`/Users/gherardolattanzi/Desktop/maxextract/picoclaw-deploy/workspace/AGENT.md`
- User intent profile:
`/Users/gherardolattanzi/Desktop/maxextract/picoclaw-deploy/workspace/USER.md`
- Persistent context:
`/Users/gherardolattanzi/Desktop/maxextract/picoclaw-deploy/workspace/memory/MEMORY.md`
- MaxExtract infrastructure baseline:
`/Users/gherardolattanzi/Desktop/maxextract/AGENTS.md`
- Strategy catalog:
`/Users/gherardolattanzi/Desktop/maxextract/strategies/`

**Mandatory Commands**

- Inventory:
`/Users/gherardolattanzi/Desktop/maxextract/scripts/me_bots_inventory.sh --context mycoolify --mode all --json`
- API state:
`/Users/gherardolattanzi/Desktop/maxextract/scripts/me_bots_api_state.sh --context mycoolify --mode all --json`
- DB ranking:
`/Users/gherardolattanzi/Desktop/maxextract/scripts/me_bots_db_roi.sh --mode all --days auto --json`
- Digest:
`/Users/gherardolattanzi/Desktop/maxextract/scripts/me_bots_digest.sh --context mycoolify --mode all --days auto`
- Periodic:
`/Users/gherardolattanzi/Desktop/maxextract/scripts/me_bots_periodic_report.sh --context mycoolify --mode all --days auto --interval-hours 3`

**Mutation Gate**

Before any write action, always provide:

- affected bots
- current health snapshot
- source confidence (`db` or `api_fallback`)
- rollback direction
- explicit confirmation prompt

**Response Contract**

- Always state source: `db` or `api_fallback`.
- Always state window: `auto` or explicit days.
- If data is partial, mark missing fields as `n/a`, never fabricate.
- For mutating requests, include:
- affected bots
- risk/blast radius
- rollback direction
- explicit confirmation prompt

**Output Skeleton**

- **Summary**
- **Source**
- **Window**
- **Bots**
- **Next action** (only if actionable)

**Failure Handling**

- If DB unavailable:
  - continue with `api_fallback`
  - show explicit DB error reason
- If one bot API fails:
  - keep partial output
  - mark bot fields as `n/a`
- If Coolify API fails:
  - report failure
  - do not execute mutating actions
