**Periodic Bot Reports**

- Every heartbeat cycle run:
`/Users/gherardolattanzi/Desktop/maxextract/scripts/me_bots_periodic_report.sh --context mycoolify --mode all --days auto --interval-hours 3`
- Expected behavior:
- sends a compact report every 3 hours
- returns `HEARTBEAT_OK` between intervals
- failure handling:
- if DB unavailable, keep report with `pnl_source=api_fallback`
- if API partial, keep available bots and mark `n/a` per failed bot

**Operational Checks**

- Before sending report:
- read `/Users/gherardolattanzi/Desktop/maxextract/picoclaw-deploy/workspace/USER.md`
- read `/Users/gherardolattanzi/Desktop/maxextract/picoclaw-deploy/workspace/memory/MEMORY.md`
- Output order:
- **Summary**
- **Source**
- **Window**
- **Bots**
- **Next action** (only if actionable)
