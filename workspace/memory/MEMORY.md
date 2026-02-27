**Stable Facts**

- Primary environment for current work: `pico` (PicoClaw service) plus `paper` services.
- Read-only API diagnostics are preferred.
- Public fallback URLs are required when internal Docker hostnames are unreachable.
- Core orchestration docs:
  - `/Users/gherardolattanzi/Desktop/maxextract/picoclaw-deploy/workspace/AGENT.md`
  - `/Users/gherardolattanzi/Desktop/maxextract/picoclaw-deploy/workspace/USER.md`
  - `/Users/gherardolattanzi/Desktop/maxextract/picoclaw-deploy/config/skills/maxextract-infra/SKILL.md`
- Default bot ops parameters:
  - `--mode all`
  - `--days auto`
  - `COOLIFY_CONTEXT=mycoolify`
  - `MAXEXTRACT_USE_SSH=1`
- Digest/report source policy:
  - prefer DB
  - fallback to API with `pnl_source=api_fallback`

**Known Constraints**

- Some public endpoints may return 404 depending on Coolify routing.
- Coolify deployment logs for failed builds can sometimes be unavailable.
- ROI 7d/30d may be unavailable on day-0 history; dynamic window is required.
