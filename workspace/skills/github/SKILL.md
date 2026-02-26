---
name: github
description: Support release and incident communication workflows for the MaxExtract fork.
---

**GitHub Workflow (Ops Focus)**

Use this skill when asked to summarize what changed, what was deployed, or what regressed.

Prefer concise changelog style:
- `commit hash`
- `purpose`
- `impact`
- `verification`

**Output style**
- 3-6 bullet points max.
- Mention deployment UUID if available.
- Mention affected MaxExtract area (`runtime`, `db`, `strategies`, `scripts`, `picoclaw-deploy`).
- Never use # headers or pipe tables in output.

**Release Summary Checklist**

- what changed
- why it changed
- how it was validated
- deployment target (`paper`, `live`, `pico`)
- residual risk or follow-up

**Incident Diff Checklist**

- first bad commit (if known)
- services impacted
- rollback candidate
- hotfix candidate

**Cross References**

- Repo policy:
`/Users/gherardolattanzi/Desktop/maxextract/AGENTS.md`
- Incident summary format:
`/Users/gherardolattanzi/Desktop/maxextract/picoclaw-deploy/workspace/skills/summarize/SKILL.md`
