---
name: maxextract-infra
description: Manage MaxExtract services through Coolify API and internal network endpoints.
---

**REMINDER: Never use # headers or pipe tables in your output. Use **bold** lines and bullet lists only.**

**MaxExtract Infrastructure Management**

You have access to the MaxExtract trading infrastructure via the Coolify deployment platform.
Use the `exec` tool with `curl` to interact with the Coolify API and internal services.

**Authentication**

All Coolify API calls require a Bearer token:

```
curl -s -H "Authorization: Bearer $COOLIFY_API_TOKEN" "$COOLIFY_API_URL/api/v1/..."
```

The token and URL are available as environment variables `COOLIFY_API_TOKEN` and `COOLIFY_API_URL`.

**Service Registry (baseline — always verify live)**

- **DB - Paper:** UUID `zkk4wok8k08s440ggk4sso08`, PostgreSQL + TimescaleDB, host `zkk4wok8k08s440ggk4sso08:5432`
- **EMA until expiry - BTC 5m:** UUID `ess8wcoo0cc8gwc8s8osc84g`, runtime, host `ess8wcoo0cc8gwc8s8osc84g:3000`
- **EMA until expiry - BTC 15m:** UUID `hkcowc8080w80kgoss8k40ss`, runtime, host `hkcowc8080w80kgoss8k40ss:3000`
- **EMA until expiry - ETH 15m:** UUID `g0o4ccw00c4gskog44o8g08w`, runtime, host `g0o4ccw00c4gskog44o8g08w:3000`
- **Recorder 5min:** UUID `vwg4o4cw4wg8ckwk88ks0408`, recorder, host `vwg4o4cw4wg8ckwk88ks0408:3000`
- **Recorder 15min:** UUID `p8g00kog08ksoo8sksok4ssw`, recorder, host `p8g00kog08ksoo8sksok4ssw:3000`
- **Cross Arb Monitor:** UUID `c4c08gokgcggs08soo4088os`, strategy monitor, host `c4c08gokgcggs08soo4088os:3000`

Always verify via live API — this list may be stale.

**Dynamic Discovery (Always Prefer This)**

When asked "what services do we have?" or "status", run:

```sh
curl -s -H "Authorization: Bearer $COOLIFY_API_TOKEN" "$COOLIFY_API_URL/api/v1/applications" | jq '.[] | {name, uuid, status}'
```

Build the answer from this live result.

**Coolify API Operations**

**Check all services:**
```sh
curl -s -H "Authorization: Bearer $COOLIFY_API_TOKEN" "$COOLIFY_API_URL/api/v1/applications" | jq '.[] | {name, uuid, status}'
```

**Get single service:**
```sh
curl -s -H "Authorization: Bearer $COOLIFY_API_TOKEN" "$COOLIFY_API_URL/api/v1/applications/{uuid}"
```

**Read logs:**
```sh
curl -s -H "Authorization: Bearer $COOLIFY_API_TOKEN" "$COOLIFY_API_URL/api/v1/applications/{uuid}/logs?since=60"
```

**Restart (ask confirmation first):**
```sh
curl -s -X POST -H "Authorization: Bearer $COOLIFY_API_TOKEN" "$COOLIFY_API_URL/api/v1/applications/{uuid}/restart"
```

**Deploy:**
```sh
curl -s -X POST -H "Authorization: Bearer $COOLIFY_API_TOKEN" "$COOLIFY_API_URL/api/v1/deploy?uuid={uuid}"
```

**List deployments:**
```sh
curl -s -H "Authorization: Bearer $COOLIFY_API_TOKEN" "$COOLIFY_API_URL/api/v1/applications/{uuid}/deployments"
```

**List env vars:**
```sh
curl -s -H "Authorization: Bearer $COOLIFY_API_TOKEN" "$COOLIFY_API_URL/api/v1/applications/{uuid}/envs"
```

**Quick Health Check (no auth, internal network)**

```sh
curl -s http://ess8wcoo0cc8gwc8s8osc84g:3000/api/health
curl -s http://hkcowc8080w80kgoss8k40ss:3000/api/health
curl -s http://g0o4ccw00c4gskog44o8g08w:3000/api/health
curl -s http://vwg4o4cw4wg8ckwk88ks0408:3000/api/health
curl -s http://p8g00kog08ksoo8sksok4ssw:3000/api/health
```

**Safety Rules**

- Paper environment uses `DRY_RUN=true` — all trades are simulated.
- Never change `DRY_RUN` from `true` to `false` without explicit user approval.
- Always check for open trades before restarting runtime services.
- Live environment is currently empty — do not create services there without approval.

**Output Formatting**

- Never use # headers — use **bold** lines.
- Never use pipe tables or `|---|` separators.
- Use bullet lists or monospaced code blocks.
- One status emoji per line max: 🟢 healthy, 🟡 degraded, 🔴 unreachable, ⚠️ action.
- **Bold** only key labels: **Summary**, **Services**, **Next action**.
- Do not repeat the same service in both intro and data rows.
- If a call fails, show `n/a` and continue.
- End with **Next action:** when any row is not healthy.
