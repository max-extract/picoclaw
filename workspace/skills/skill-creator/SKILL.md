---
name: skill-creator
description: Create or update PicoClaw skills for MaxExtract operational workflows.
---

**Skill Creator (Ops-Oriented)**

When adding a new skill:

1. Include frontmatter with `name` and `description`.
2. Define exact command/endpoint examples.
3. Define output format contract.
4. Add explicit fallback behavior.
5. Never use # headers in skill content — use **bold** lines instead.
6. Add cross references to:
- `/Users/gherardolattanzi/Desktop/maxextract/picoclaw-deploy/workspace/AGENT.md`
- `/Users/gherardolattanzi/Desktop/maxextract/picoclaw-deploy/workspace/USER.md`
- `/Users/gherardolattanzi/Desktop/maxextract/picoclaw-deploy/config/skills/maxextract-infra/SKILL.md`

**Mandatory Sections For New Skills**

- Scope
- Trigger conditions
- Primary workflow
- Fallback behavior
- Output contract
- Safety guardrails

**Quality Criteria**

- command examples must be executable as written
- every fallback path must be explicit
- avoid overlaps with existing skills unless intentionally extending
