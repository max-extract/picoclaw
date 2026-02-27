**Periodic Bot Reports**

- Every heartbeat cycle run:
`MAXEXTRACT_USE_SSH=1 ${MAXEXTRACT_ROOT}/scripts/me_bots_periodic_report.sh --context mycoolify --mode all --days auto --interval-hours 3`
- Expected behavior:
- sends a compact report every 3 hours
- returns `HEARTBEAT_OK` between intervals
- failure handling:
- if DB unavailable, keep report with `pnl_source=api_fallback`
- if API partial, keep available bots and mark `n/a` per failed bot

**Operational Checks**

- Before sending report:
- read `workspace/USER.md`
- read `workspace/memory/MEMORY.md`
- Output order:
- **Summary**
- **Source**
- **Window**
- **Bots**
- **Next action** (only if actionable)
