---
name: database
description: Access MaxExtract bot data using APIs first, DB only when configured for ranking queries.
---

**REMINDER: Never use # headers or pipe tables in your output. Use **bold** lines and bullet lists only.**

**Default Mode**

Use runtime bot APIs first for state and health.
Use DB only for historical ROI/PnL ranking when available via config/env.
Default scope is runtime bots, not recorders.

**When To Use This Skill**

- user asks ROI/PnL ranking
- user asks trade summary from historical data
- DB connectivity is uncertain and fallback is needed
- you need source-confidence labeling (`db` vs `api_fallback`)

**Runtime API Endpoints**

Internal hosts:
- **BTC 5m:** `http://ess8wcoo0cc8gwc8s8osc84g:3000`
- **BTC 15m:** `http://hkcowc8080w80kgoss8k40ss:3000`
- **ETH 15m:** `http://g0o4ccw00c4gskog44o8g08w:3000`

Key routes:

```sh
curl -s http://ess8wcoo0cc8gwc8s8osc84g:3000/api/state
curl -s http://ess8wcoo0cc8gwc8s8osc84g:3000/api/strategy/state
curl -s http://ess8wcoo0cc8gwc8s8osc84g:3000/api/polymarket/activity?limit=20
curl -s http://ess8wcoo0cc8gwc8s8osc84g:3000/api/health
```

**Ranking Path (Preferred)**

Use:
`/Users/gherardolattanzi/Desktop/maxextract/scripts/me_bots_db_roi.sh --mode all --days auto --json`

If unavailable, fallback:
`/Users/gherardolattanzi/Desktop/maxextract/scripts/me_bots_api_state.sh --context mycoolify --mode all --json`

**Decision Rules**

- Ranking request:
  - try DB first
  - if DB fails, fallback to API and label `api_fallback`
- Pure realtime status request:
  - API only, skip DB
- If history window missing:
  - use `--days auto`

**Guidance**

- For PnL/trade summaries, prefer `/api/state` and `/api/polymarket/activity`.
- For ranking questions, return source as `db` or `api_fallback`.
- If DB is not configured, report explicit reason and continue with fallback.
- Cross-reference infra policy:
`/Users/gherardolattanzi/Desktop/maxextract/picoclaw-deploy/config/skills/maxextract-infra/SKILL.md`

**Output Minimum**

- **Summary**
- **Source**
- **Window**
- **Top bots**
- **Bottom bots** (if requested)

**Quality Guardrails**

- never mix DB and API ranking without declaring authority
- never infer ROI from missing trade history
- use `n/a` for unavailable metrics
