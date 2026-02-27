---
name: trading-dashboard
description: Query MaxExtract runtime bots for health, state, and performance with safe fallbacks.
---

**REMINDER: Never use # headers or pipe tables in your output. Use **bold** lines and bullet lists only.**

**Scope**

- Default: runtime bots only (`paper`, `live`).
- Exclude non-runtime services unless explicitly requested.

**Primary Use Cases**

- current health across bots
- performance snapshot per bot
- strategy/runtime state checks
- unified operator digest

**Primary Workflow**

1. Inventory:
`MAXEXTRACT_USE_SSH=1 /Users/gherardolattanzi/Desktop/maxextract/scripts/me_bots_inventory.sh --context mycoolify --mode all --json`
2. Realtime API state:
`MAXEXTRACT_USE_SSH=1 /Users/gherardolattanzi/Desktop/maxextract/scripts/me_bots_api_state.sh --context mycoolify --mode all --json`
3. Historical ranking:
`/Users/gherardolattanzi/Desktop/maxextract/scripts/me_bots_db_roi.sh --mode all --days auto --json`
4. Unified digest:
`MAXEXTRACT_USE_SSH=1 /Users/gherardolattanzi/Desktop/maxextract/scripts/me_bots_digest.sh --context mycoolify --mode all --days auto`

**Fast Paths**

- health-only:
  - inventory + API state
- ranking-only:
  - DB ROI only
- full operator view:
  - digest only

**Bot Fast Paths**

- one bot health/state:
`MAXEXTRACT_USE_SSH=1 /Users/gherardolattanzi/Desktop/maxextract/scripts/me_bot_health.sh --mode paper --strategy ema-until-expiry --market btc-5m --json`
- one bot ROI/PnL:
`MAXEXTRACT_USE_SSH=1 /Users/gherardolattanzi/Desktop/maxextract/scripts/me_bot_roi.sh --mode paper --strategy ema-until-expiry --market btc-5m --days auto --json`
- one bot unified report:
`MAXEXTRACT_USE_SSH=1 /Users/gherardolattanzi/Desktop/maxextract/scripts/me_bot_report.sh --mode paper --strategy ema-until-expiry --market btc-5m --days auto --json`

**Runtime Endpoints (Direct Checks)**

Do not hardcode individual bot names/timeframes in docs or replies.
Use inventory output to discover active bots and their endpoints, then query:

- `/api/health`
- `/api/state`
- `/api/strategy/state`
- `/api/polymarket/activity?limit=20`

**Connectivity Fallback**

If internal endpoints fail, use configured public runtime URLs from env/config.
Do not expose hardcoded per-bot URL names in user-facing output.

**State Fields To Extract**

- bot identity: `mode`, `market`, `strategy`
- health: `ok`, `uptime`, `errors`
- performance: `pnl`, `roi`, `trades`, `win_rate`
- freshness: timestamp of source query

**Output Contract**

- Always include:
- **Summary**
- **Source** (`db` or `api_fallback`)
- **Window** (`auto` or explicit days)
- **Bots**
- **Next action** (only if actionable)
- Use `n/a` for missing metrics.
- Never fabricate values.

**Safety Rules**

- do not execute test-entry/test-exit/reset endpoints unless explicitly requested
- if user requests mutation, route through infra mutation gate first
- if data sources disagree, prefer DB for historical ranking and explain conflict briefly

**Cross References**

- Infra policy:
`/Users/gherardolattanzi/Desktop/maxextract/picoclaw-deploy/config/skills/maxextract-infra/SKILL.md`
- User intent profile:
`/Users/gherardolattanzi/Desktop/maxextract/picoclaw-deploy/workspace/USER.md`
- MaxExtract architecture:
`/Users/gherardolattanzi/Desktop/maxextract/AGENTS.md`
