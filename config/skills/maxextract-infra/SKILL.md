---
name: maxextract-infra
description: Manage MaxExtract runtime bots via Coolify and bot backend APIs.
---

**REMINDER: Never use # headers or pipe tables in output. Use bold lines and bullet lists.**

**Scope**

- Runtime bots only (`paper`, `live`).
- Exclude non-runtime services unless explicitly requested.

**Intent Router**

- `inventory` or `deploy status`:
  - run SSH-first inventory command
- `realtime state`:
  - run SSH-first API state command
- `roi/pnl ranking`:
  - run DB ROI script first
- `executive digest`:
  - run SSH-first digest command
- `periodic update`:
  - run SSH-first periodic report command
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

- SSH defaults:
`MAXEXTRACT_USE_SSH=1`
`MAXEXTRACT_SSH_TARGET=<user@host>`
- Telegram rendering default:
`MAXEXTRACT_OUTPUT_FORMAT=telegram`
- Inventory:
`MAXEXTRACT_USE_SSH=1 /Users/gherardolattanzi/Desktop/maxextract/scripts/me_bots_inventory.sh --context mycoolify --mode all --json`
- API state:
`MAXEXTRACT_USE_SSH=1 /Users/gherardolattanzi/Desktop/maxextract/scripts/me_bots_api_state.sh --context mycoolify --mode all --json`
- DB ranking:
`/Users/gherardolattanzi/Desktop/maxextract/scripts/me_bots_db_roi.sh --mode all --days auto --json`
- Digest:
`MAXEXTRACT_USE_SSH=1 /Users/gherardolattanzi/Desktop/maxextract/scripts/me_bots_digest.sh --context mycoolify --mode all --days auto`
- Periodic:
`MAXEXTRACT_USE_SSH=1 /Users/gherardolattanzi/Desktop/maxextract/scripts/me_bots_periodic_report.sh --context mycoolify --mode all --days auto --interval-hours 3`

**Bot-Specific Commands**

- Resolve one bot:
`MAXEXTRACT_USE_SSH=1 /Users/gherardolattanzi/Desktop/maxextract/scripts/me_bot_resolve.sh --mode paper --strategy ema-until-expiry --market btc-5m --json`
- Bot health/state:
`MAXEXTRACT_USE_SSH=1 /Users/gherardolattanzi/Desktop/maxextract/scripts/me_bot_health.sh --mode paper --strategy ema-until-expiry --market btc-5m --json`
- Bot ROI/PnL (DB-first):
`MAXEXTRACT_USE_SSH=1 /Users/gherardolattanzi/Desktop/maxextract/scripts/me_bot_roi.sh --mode paper --strategy ema-until-expiry --market btc-5m --days auto --json`
- Unified bot report:
`MAXEXTRACT_USE_SSH=1 /Users/gherardolattanzi/Desktop/maxextract/scripts/me_bot_report.sh --mode paper --strategy ema-until-expiry --market btc-5m --days auto --json`

**Rendering Modes**

- `--json`: machine parsing.
- `--telegram` (or default): bold sections + bullets, Telegram-safe.
- `--table`: legacy terminal table/flat output.

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
