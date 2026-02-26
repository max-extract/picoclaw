# MaxExtract Database Queries

You can query the MaxExtract PostgreSQL + TimescaleDB database directly using `psql`.
The connection string is available as the `DATABASE_URL` environment variable.

## Connection

```sh
psql "$DATABASE_URL" -c "YOUR SQL HERE"
```

For multi-line queries, use a heredoc:
```sh
psql "$DATABASE_URL" << 'SQL'
SELECT ... FROM trades ...;
SQL
```

Always use `-t` (tuples only) or `--csv` for machine-readable output when processing results.

## Schema

### trades
Trade lifecycle from entry to exit.

| Column | Type | Description |
|---|---|---|
| trade_id | SERIAL PK | Auto-incrementing ID |
| status | TEXT | `open`, `closed`, `expired` |
| market_slug | TEXT | Polymarket market identifier |
| side | TEXT | `UP` or `DOWN` |
| entry_time | TIMESTAMPTZ | When the trade was entered |
| entry_price | FLOAT | Entry price paid |
| exit_time | TIMESTAMPTZ | When the trade was closed |
| exit_price | FLOAT | Exit price received |
| size | FLOAT | Position size in USD |
| tokens | FLOAT | Number of tokens held |
| pnl_usd | FLOAT | Realized profit/loss in USD |
| pnl_pct | FLOAT | Realized profit/loss percentage |
| fees | FLOAT | Trading fees paid |
| hold_minutes | FLOAT | Duration of the trade |
| regime | TEXT | Market regime at entry |
| strike_price | FLOAT | Option strike price |
| btc_price_at_entry | FLOAT | BTC spot price at entry |
| exit_reason | TEXT | Why the trade was closed |

### orderbook_snapshots (TimescaleDB hypertable)
High-frequency orderbook data captured every 5 or 15 minutes.

| Column | Type | Description |
|---|---|---|
| ts | TIMESTAMPTZ | Snapshot timestamp |
| asset | TEXT | `btc`, `eth`, `sol` |
| timeframe | TEXT | `5m`, `15m`, `1h` |
| market | TEXT | Polymarket market slug |
| price | FLOAT | Underlying asset price (Chainlink) |
| strike_price | FLOAT | Option strike price |
| up_bids / up_asks | JSONB | UP token orderbook |
| down_bids / down_asks | JSONB | DOWN token orderbook |

### analytics_state
Key-value JSONB store for strategy state persistence.

| Column | Type | Description |
|---|---|---|
| state_type | TEXT UNIQUE | Key identifier |
| data | JSONB | Arbitrary state blob |
| updated_at | TIMESTAMPTZ | Last update time |

### runtime_metrics (TimescaleDB hypertable)
Operational health metrics per service.

| Column | Type | Description |
|---|---|---|
| ts | TIMESTAMPTZ | Metric timestamp |
| service | TEXT | Service name |
| cpu_percent | FLOAT | CPU usage |
| rss_mb | FLOAT | Memory usage (RSS) |
| snapshots_written | BIGINT | Snapshots written count |
| snapshots_skipped | BIGINT | Snapshots skipped count |

## Common Queries

### Open positions
```sql
SELECT trade_id, market_slug, side, entry_time, entry_price, size, tokens
FROM trades WHERE status = 'open' ORDER BY entry_time DESC;
```

### Recent closed trades
```sql
SELECT trade_id, market_slug, side, entry_time, exit_time,
       entry_price, exit_price, pnl_usd, pnl_pct, exit_reason
FROM trades WHERE status = 'closed'
ORDER BY exit_time DESC LIMIT 20;
```

### Daily PnL summary
```sql
SELECT date_trunc('day', exit_time) AS day,
       count(*) AS trades,
       sum(pnl_usd) AS total_pnl,
       avg(pnl_pct) AS avg_pnl_pct,
       sum(CASE WHEN pnl_usd > 0 THEN 1 ELSE 0 END)::float / count(*) AS win_rate
FROM trades WHERE status = 'closed' AND exit_time IS NOT NULL
GROUP BY 1 ORDER BY 1 DESC LIMIT 14;
```

### Total performance
```sql
SELECT count(*) AS total_trades,
       sum(pnl_usd) AS total_pnl,
       avg(pnl_usd) AS avg_pnl,
       sum(CASE WHEN pnl_usd > 0 THEN 1 ELSE 0 END)::float / NULLIF(count(*), 0) AS win_rate,
       max(pnl_usd) AS best_trade,
       min(pnl_usd) AS worst_trade
FROM trades WHERE status = 'closed';
```

### Orderbook snapshot count (today)
```sql
SELECT asset, timeframe, count(*) AS snapshots
FROM orderbook_snapshots
WHERE ts >= now() - interval '24 hours'
GROUP BY asset, timeframe ORDER BY asset, timeframe;
```

### Latest snapshot per asset
```sql
SELECT DISTINCT ON (asset, timeframe)
       asset, timeframe, ts, price, strike_price, market
FROM orderbook_snapshots
ORDER BY asset, timeframe, ts DESC;
```

### Service health metrics (last hour)
```sql
SELECT service, avg(cpu_percent) AS avg_cpu, avg(rss_mb) AS avg_mem,
       sum(snapshots_written) AS writes, sum(snapshots_skipped) AS skips
FROM runtime_metrics
WHERE ts >= now() - interval '1 hour'
GROUP BY service;
```

### Trades by market
```sql
SELECT market_slug, count(*) AS trades, sum(pnl_usd) AS pnl,
       avg(hold_minutes) AS avg_hold_min
FROM trades WHERE status = 'closed'
GROUP BY market_slug ORDER BY pnl DESC;
```

## Tips

- Use `LIMIT` on large tables (orderbook_snapshots can have millions of rows).
- For time-range queries on hypertables, always filter on `ts` for performance.
- Use `--csv` flag with psql for structured output: `psql "$DATABASE_URL" --csv -c "..."`.
- The database is the paper trading database -- all data is from simulated trades.
