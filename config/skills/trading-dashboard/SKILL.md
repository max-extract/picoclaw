---
name: trading-dashboard
description: Query runtime and recorder APIs for market state, trades, PnL, and service health.
---

**REMINDER: Never use # headers or pipe tables in your output. Use **bold** lines and bullet lists only.**

**MaxExtract Trading Dashboard**

Query the MaxExtract trading bot runtime services and orderbook recorders
over the internal Docker network. No authentication is needed for these endpoints.

**Connectivity Fallback (Important)**

If internal UUID hostnames fail with `curl` exit code 6, use the public base URLs
provided via environment variables:

- `MAXEXTRACT_RUNTIME_BTC5M_URL`
- `MAXEXTRACT_RUNTIME_BTC15M_URL`
- `MAXEXTRACT_RUNTIME_ETH15M_URL`
- `MAXEXTRACT_RECORDER_5M_URL`
- `MAXEXTRACT_RECORDER_15M_URL`
- `MAXEXTRACT_CROSS_ARB_URL`

Example fallback:

```sh
curl -s "${MAXEXTRACT_RUNTIME_BTC5M_URL}/api/state"
curl -s "${MAXEXTRACT_RECORDER_5M_URL}/api/stats"
```

**Runtime Services (Trading Bots)**

Three runtime instances run the "EMA Until Expiry" strategy:

- **BTC 5m:** `http://ess8wcoo0cc8gwc8s8osc84g:3000`
- **BTC 15m:** `http://hkcowc8080w80kgoss8k40ss:3000`
- **ETH 15m:** `http://g0o4ccw00c4gskog44o8g08w:3000`

**Cross Arb Monitor (Strategy 5)**

- **Cross Arb Monitor:** `http://c4c08gokgcggs08soo4088os:3000`

Cross Arb endpoints:

```sh
curl -s http://c4c08gokgcggs08soo4088os:3000/api/health
curl -s http://c4c08gokgcggs08soo4088os:3000/api/state/all
curl -s "http://c4c08gokgcggs08soo4088os:3000/api/ui/state?symbol=btc"
curl -s "http://c4c08gokgcggs08soo4088os:3000/api/ui/state?symbol=eth"
```

**ROI Query (Cross Arb)**

Use `api/ui/state` and read:
- `baseState.performance.roi`
- `baseState.performance.totalPnl`
- `baseState.performance.trades`
- `baseState.performance.winRate`

```sh
curl -s "http://c4c08gokgcggs08soo4088os:3000/api/ui/state?symbol=btc" | jq '{roi: .baseState.performance.roi, totalPnl: .baseState.performance.totalPnl, trades: .baseState.performance.trades, winRate: .baseState.performance.winRate}'
```

Public fallback:

```sh
curl -s "${MAXEXTRACT_CROSS_ARB_URL}/api/ui/state?symbol=btc" | jq '{roi: .baseState.performance.roi, totalPnl: .baseState.performance.totalPnl}'
```

If the endpoint returns `404 page not found`, report that Cross Arb API is not exposed publicly and suggest checking Coolify routing.

**Key Runtime Endpoints**

**Full dashboard state** (positions, trades, performance, market data):
```sh
curl -s http://ess8wcoo0cc8gwc8s8osc84g:3000/api/state
```

**Health check:**
```sh
curl -s http://ess8wcoo0cc8gwc8s8osc84g:3000/api/health
```

**Strategy state** (parameters, thresholds, regime):
```sh
curl -s http://ess8wcoo0cc8gwc8s8osc84g:3000/api/strategy/state
```

**Market snapshot** (orderbook, prices, spreads):
```sh
curl -s "http://ess8wcoo0cc8gwc8s8osc84g:3000/api/mind-conviction/snapshot?asset=btc&timeframe=5m&bankroll=200"
```

**Recorder status** (from runtime's perspective):
```sh
curl -s http://ess8wcoo0cc8gwc8s8osc84g:3000/api/mind-conviction/recorder-status
```

**Polymarket activity feed:**
```sh
curl -s "http://ess8wcoo0cc8gwc8s8osc84g:3000/api/polymarket/activity?limit=20"
```

**Trading Controls (use with caution)**

**Manual entry signal** (paper mode only):
```sh
curl -s -X POST -H "Content-Type: application/json" -d '{"side":"UP","size":10}' http://ess8wcoo0cc8gwc8s8osc84g:3000/api/test-entry-signal
```

**Manual exit signal** (paper mode only):
```sh
curl -s -X POST http://ess8wcoo0cc8gwc8s8osc84g:3000/api/test-exit-signal
```

**Reset all state** (clears paper trader, use with caution):
```sh
curl -s -X POST http://ess8wcoo0cc8gwc8s8osc84g:3000/api/reset-all
```

**Orderbook Recorders**

Two recorders capture orderbook snapshots:

- **Recorder 5m:** `http://vwg4o4cw4wg8ckwk88ks0408:3000`
- **Recorder 15m:** `http://p8g00kog08ksoo8sksok4ssw:3000`

**Service stats** (assets tracked, snapshot counts, uptime):
```sh
curl -s http://vwg4o4cw4wg8ckwk88ks0408:3000/api/stats
```

**Health check:**
```sh
curl -s http://vwg4o4cw4wg8ckwk88ks0408:3000/api/health
```

**List data files:**
```sh
curl -s http://vwg4o4cw4wg8ckwk88ks0408:3000/api/files
```

**Cost metrics** (5min recorder only):
```sh
curl -s http://vwg4o4cw4wg8ckwk88ks0408:3000/api/cost-metrics
```

**Data footprint** (5min recorder only):
```sh
curl -s http://vwg4o4cw4wg8ckwk88ks0408:3000/api/data-footprint
```

**Tips**

- Query `/api/state` on each runtime and `/api/stats` on each recorder for a quick overview.
- The `performance` field in `/api/state` has PnL, win rate, trade count.
- The `recentTrades` array shows last trades with entry/exit prices and PnL.
- All runtimes are in paper mode (`DRY_RUN=true`), trades are simulated.
- Replace the hostname UUID in URLs to query different instances.
- For Cross Arb ROI, always query `/api/ui/state`, do not assume DB schema.

**Output Formatting**

- Never use # headers — use **bold** lines.
- Never use pipe tables or `|---|` separators.
- Use bullet lists or monospaced code blocks.
- Row format: `SERVICE | STATUS | KEY_METRICS | NOTES`
- One status emoji per line max: 🟢 healthy, 🟡 degraded, 🔴 unreachable, ⚠️ action.
- **Bold** only key labels: **Summary**, **Key metrics**, **Next action**.
- Do not print duplicate rows for the same service.
- Do not mix bullet and code-block forms for the same data.
- Round numeric metrics to 2 decimal places.
- Use `n/a` for unavailable metrics.
