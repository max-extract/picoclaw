---
name: trading-dashboard
description: Query runtime and recorder APIs for market state, trades, PnL, and service health.
---

# MaxExtract Trading Dashboard

You can query the MaxExtract trading bot runtime services and orderbook recorders
over the internal Docker network. No authentication is needed for these endpoints.

## Runtime Services (Trading Bots)

Three runtime instances run the "EMA Until Expiry" strategy on different markets:

| Instance | Internal Base URL |
|---|---|
| BTC 5-minute | `http://ess8wcoo0cc8gwc8s8osc84g:3000` |
| BTC 15-minute | `http://hkcowc8080w80kgoss8k40ss:3000` |
| ETH 15-minute | `http://g0o4ccw00c4gskog44o8g08w:3000` |

## Cross Arb Monitor (Strategy 5)

Cross Arb has a separate monitor service:

| Service | Internal Base URL |
|---|---|
| Cross Arb Monitor | `http://c4c08gokgcggs08soo4088os:3000` |

### Cross Arb Endpoints

```sh
curl -s http://c4c08gokgcggs08soo4088os:3000/api/health
curl -s http://c4c08gokgcggs08soo4088os:3000/api/state/all
curl -s "http://c4c08gokgcggs08soo4088os:3000/api/ui/state?symbol=btc"
curl -s "http://c4c08gokgcggs08soo4088os:3000/api/ui/state?symbol=eth"
```

### ROI Query (Cross Arb)

Use `api/ui/state` and read:
- `baseState.performance.roi`
- `baseState.performance.totalPnl`
- `baseState.performance.trades`
- `baseState.performance.winRate`

Example:

```sh
curl -s "http://c4c08gokgcggs08soo4088os:3000/api/ui/state?symbol=btc" | jq '{roi: .baseState.performance.roi, totalPnl: .baseState.performance.totalPnl, trades: .baseState.performance.trades, winRate: .baseState.performance.winRate}'
```

### Key Endpoints

**Full dashboard state** (positions, trades, performance, market data):
```sh
curl -s http://ess8wcoo0cc8gwc8s8osc84g:3000/api/state
```
Returns JSON with: `market`, `positions`, `recentTrades`, `performance`, `analytics`, `uptime`, etc.

**Health check:**
```sh
curl -s http://ess8wcoo0cc8gwc8s8osc84g:3000/api/health
```

**Strategy state** (current parameters, thresholds, regime):
```sh
curl -s http://ess8wcoo0cc8gwc8s8osc84g:3000/api/strategy/state
```

**Market snapshot** (current orderbook, prices, spreads):
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

### Trading Controls (use with caution)

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

## Orderbook Recorders

Two recorders capture orderbook snapshots at different intervals:

| Recorder | Internal Base URL |
|---|---|
| 5-minute | `http://vwg4o4cw4wg8ckwk88ks0408:3000` |
| 15-minute | `http://p8g00kog08ksoo8sksok4ssw:3000` |

### Key Endpoints

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

## Tips

- To get a quick overview, query `/api/state` on each runtime and `/api/stats` on each recorder.
- The `performance` field in `/api/state` contains PnL, win rate, and trade count.
- The `recentTrades` array shows the last trades with entry/exit prices and PnL.
- All runtimes are in paper mode (`DRY_RUN=true`), so trades are simulated.
- Replace the hostname UUID in the URLs above to query different instances.
- For Cross Arb ROI, do not use DB assumptions; always query `/api/ui/state`.

## Response Formatting Rules

When presenting service data:

- Use compact markdown tables with exactly these columns when possible:
  - `Service | Status | Key Metrics | Notes`
- Keep values short (one line per row), round numeric metrics to 2 decimals.
- If a metric is unavailable, use `n/a` instead of long explanations.
- After the table, add a short `Next action:` line if any service is unhealthy.
