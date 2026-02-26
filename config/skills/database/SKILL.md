---
name: database
description: Access MaxExtract data via runtime and recorder APIs without direct database credentials.
---

# MaxExtract Data Access (API-Only)

Use runtime and recorder HTTP APIs for trading data and metrics.
Do not use direct PostgreSQL access from PicoClaw.

## Runtime API Endpoints

Use these internal hosts:
- BTC 5m: `http://ess8wcoo0cc8gwc8s8osc84g:3000`
- BTC 15m: `http://hkcowc8080w80kgoss8k40ss:3000`
- ETH 15m: `http://g0o4ccw00c4gskog44o8g08w:3000`

Key routes:

```sh
curl -s http://ess8wcoo0cc8gwc8s8osc84g:3000/api/state
curl -s http://ess8wcoo0cc8gwc8s8osc84g:3000/api/strategy/state
curl -s http://ess8wcoo0cc8gwc8s8osc84g:3000/api/polymarket/activity?limit=20
curl -s http://ess8wcoo0cc8gwc8s8osc84g:3000/api/health
```

## Recorder API Endpoints

Use these internal hosts:
- Recorder 5m: `http://vwg4o4cw4wg8ckwk88ks0408:3000`
- Recorder 15m: `http://p8g00kog08ksoo8sksok4ssw:3000`

Key routes:

```sh
curl -s http://vwg4o4cw4wg8ckwk88ks0408:3000/api/stats
curl -s http://vwg4o4cw4wg8ckwk88ks0408:3000/api/cost-metrics
curl -s http://vwg4o4cw4wg8ckwk88ks0408:3000/api/data-footprint
curl -s http://vwg4o4cw4wg8ckwk88ks0408:3000/api/health
```

## Guidance

- For PnL/trade summaries, prefer `/api/state` and `/api/polymarket/activity`.
- For ingestion/recorder health, use `/api/stats` and `/api/health`.
- If API data is insufficient, ask the user before introducing DB credentials.
