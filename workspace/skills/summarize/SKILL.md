---
name: summarize
description: Summarize incidents and deployment outcomes for MaxExtract operations.
---

**Incident Summary Template**

- **What happened:** ...
- **Current status:** ...
- **Root cause:** known / unknown
- **Fix applied:** ...
- **Next action:** ...
- **Owner:** who should act next
- **ETA:** expected follow-up time

Keep it short and actionable. Never use # headers or pipe tables.

**MaxExtract Context Hints**

- Mention impacted bot scope: `paper`, `live`, or `all`.
- Mention source confidence: `db` or `api_fallback`.
- Mention if issue is backend API, DB config, or Coolify routing.

**Operator Quality Bar**

- include one clear decision at the end:
  - monitor
  - fix now
  - rollback
- avoid generic text; include concrete command or owner for next step
