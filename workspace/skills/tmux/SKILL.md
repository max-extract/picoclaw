---
name: tmux
description: Lightweight terminal-session checks for operational debugging.
---

**tmux Checks**

If requested, use tmux only for:
- session listing
- checking if a command is still running
- collecting minimal live output for a stuck script

Do not depend on tmux for normal MaxExtract operations.

Preferred normal path:
- run canonical scripts in `/Users/gherardolattanzi/Desktop/maxextract/scripts/`
- return compact parsed output

**When Not To Use tmux**

- routine bot status queries
- ROI ranking queries
- digest/periodic report generation

**When To Use tmux**

- long-running command appears hung
- command must keep running while diagnostics continue
- quick session introspection requested by user
