---
name: maxextract-infra
description: Manage MaxExtract services through Coolify API and internal network endpoints.
---

# MaxExtract Infrastructure Management

You have access to the MaxExtract trading infrastructure via the Coolify deployment platform.
Use the `exec` tool with `curl` to interact with the Coolify API and internal services.

## Authentication

All Coolify API calls require a Bearer token:

```
curl -s -H "Authorization: Bearer $COOLIFY_API_TOKEN" "$COOLIFY_API_URL/api/v1/..."
```

The token and URL are available as environment variables `COOLIFY_API_TOKEN` and `COOLIFY_API_URL`.

## Service Registry

- DB - Paper: UUID `zkk4wok8k08s440ggk4sso08`, type PostgreSQL + TimescaleDB, host `zkk4wok8k08s440ggk4sso08:5432`
- EMA until expiry - BTC 5m: UUID `ess8wcoo0cc8gwc8s8osc84g`, type runtime, host `ess8wcoo0cc8gwc8s8osc84g:3000`
- EMA until expiry - BTC 15m: UUID `hkcowc8080w80kgoss8k40ss`, type runtime, host `hkcowc8080w80kgoss8k40ss:3000`
- EMA until expiry - ETH 15m: UUID `g0o4ccw00c4gskog44o8g08w`, type runtime, host `g0o4ccw00c4gskog44o8g08w:3000`
- Recorder 5min: UUID `vwg4o4cw4wg8ckwk88ks0408`, type recorder, host `vwg4o4cw4wg8ckwk88ks0408:3000`
- Recorder 15min: UUID `p8g00kog08ksoo8sksok4ssw`, type recorder, host `p8g00kog08ksoo8sksok4ssw:3000`
- Cross Arb Monitor: UUID `c4c08gokgcggs08soo4088os`, type strategy monitor, host `c4c08gokgcggs08soo4088os:3000`

Important: treat this table as baseline only. For current/accurate services, always query Coolify API live.

## Dynamic Discovery (Always Prefer This)

When asked "what services do we have?" or "what can you do in MaxExtract?", first run:

```sh
curl -s -H "Authorization: Bearer $COOLIFY_API_TOKEN" "$COOLIFY_API_URL/api/v1/applications" | jq '.[] | {name, uuid, status}'
```

Build the answer from this live result (do not rely only on static skill text).

## Coolify API Operations

### Check All Services Status

```sh
curl -s -H "Authorization: Bearer $COOLIFY_API_TOKEN" "$COOLIFY_API_URL/api/v1/applications" | jq '.[] | {name, uuid, status}'
```

### Get Single Service Details

```sh
curl -s -H "Authorization: Bearer $COOLIFY_API_TOKEN" "$COOLIFY_API_URL/api/v1/applications/{uuid}"
```

### Read Service Logs

```sh
curl -s -H "Authorization: Bearer $COOLIFY_API_TOKEN" "$COOLIFY_API_URL/api/v1/applications/{uuid}/logs?since=60"
```

### Restart a Service

```sh
curl -s -X POST -H "Authorization: Bearer $COOLIFY_API_TOKEN" "$COOLIFY_API_URL/api/v1/applications/{uuid}/restart"
```

### Deploy a Service

```sh
curl -s -X POST -H "Authorization: Bearer $COOLIFY_API_TOKEN" "$COOLIFY_API_URL/api/v1/deploy?uuid={uuid}"
```

### List Deployments

```sh
curl -s -H "Authorization: Bearer $COOLIFY_API_TOKEN" "$COOLIFY_API_URL/api/v1/applications/{uuid}/deployments"
```

### List Environment Variables

```sh
curl -s -H "Authorization: Bearer $COOLIFY_API_TOKEN" "$COOLIFY_API_URL/api/v1/applications/{uuid}/envs"
```

## Quick Health Check (All Services)

To check health of runtime and recorder services directly (no auth needed, internal network):

```sh
curl -s http://ess8wcoo0cc8gwc8s8osc84g:3000/api/health  # BTC 5m
curl -s http://hkcowc8080w80kgoss8k40ss:3000/api/health   # BTC 15m
curl -s http://g0o4ccw00c4gskog44o8g08w:3000/api/health   # ETH 15m
curl -s http://vwg4o4cw4wg8ckwk88ks0408:3000/api/health   # Recorder 5min
curl -s http://p8g00kog08ksoo8sksok4ssw:3000/api/health   # Recorder 15min
```

## Safety Rules

- The paper environment uses `DRY_RUN=true` -- all trades are simulated.
- Never change `DRY_RUN` from `true` to `false` without explicit user approval.
- Always check for open trades before restarting runtime services.
- The live environment is currently empty -- do not create services there without approval.

## Rendering Rules

- Telegram-safe rendering: do not use markdown pipe tables.
- Never output markdown table separators like `|---|`.
- Use bullets or a code block with monospaced rows:
  - `SERVICE | UUID | STATUS | HEALTH`
- Add one-line summary before rows: `Summary: Healthy X/Y, Degraded Z, Unreachable W`.
- If one call fails, still return partial rows with `n/a`.
- Add final `Next action:` when any row is not healthy.
