# Periodic Ops Checks

- Every cycle, run a lightweight health check on runtime and recorder APIs.
- If any service is degraded/unreachable, produce a compact status table and one next action.
- If all services are healthy, respond with `HEARTBEAT_OK` only.
