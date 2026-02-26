---
name: trading-dashboard
description: Query MaxExtract runtime bots for health, state, and performance with safe fallbacks.
---

**REMINDER: Never use # headers or pipe tables in your output. Use **bold** lines and bullet lists only.**

**Scope**

- Default: runtime bots only (`paper`, `live`).
- Exclude recorders and cross-arb unless explicitly requested.

**Primary Use Cases**

- current health across bots
- performance snapshot per bot
- strategy/runtime state checks
- unified operator digest

**Primary Workflow**

1. Inventory:
`/Users/gherardolattanzi/Desktop/maxextract/scripts/me_bots_inventory.sh --context mycoolify --mode all --json`
2. Realtime API state:
`/Users/gherardolattanzi/Desktop/maxextract/scripts/me_bots_api_state.sh --context mycoolify --mode all --json`
3. Historical ranking:
`/Users/gherardolattanzi/Desktop/maxextract/scripts/me_bots_db_roi.sh --mode all --days auto --json`
4. Unified digest:
`/Users/gherardolattanzi/Desktop/maxextract/scripts/me_bots_digest.sh --context mycoolify --mode all --days auto`

**Fast Paths**

- health-only:
  - inventory + API state
- ranking-only:
  - DB ROI only
- full operator view:
  - digest only

**Runtime Endpoints (Direct Checks)**

- BTC 5m:
`http://ess8wcoo0cc8gwc8s8osc84g:3000`
- BTC 15m:
`http://hkcowc8080w80kgoss8k40ss:3000`
- ETH 15m:
`http://g0o4ccw00c4gskog44o8g08w:3000`

Useful routes:

```sh
curl -s http://ess8wcoo0cc8gwc8s8osc84g:3000/api/health
curl -s http://ess8wcoo0cc8gwc8s8osc84g:3000/api/state
curl -s http://ess8wcoo0cc8gwc8s8osc84g:3000/api/strategy/state
curl -s "http://ess8wcoo0cc8gwc8s8osc84g:3000/api/polymarket/activity?limit=20"
```

**Connectivity Fallback**

If internal UUID hostnames fail (`curl` exit code `6`), use:

- `MAXEXTRACT_RUNTIME_BTC5M_URL`
- `MAXEXTRACT_RUNTIME_BTC15M_URL`
- `MAXEXTRACT_RUNTIME_ETH15M_URL`

Example:

```sh
curl -s "${MAXEXTRACT_RUNTIME_BTC5M_URL}/api/state"
```

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
