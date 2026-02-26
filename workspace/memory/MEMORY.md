**Stable Facts**

- Primary environment for current work: `pico` (PicoClaw service) plus `paper` services.
- Read-only API diagnostics are preferred.
- Public fallback URLs are required when internal Docker hostnames are unreachable.

**Known Constraints**

- Cross Arb public endpoint may return 404 depending on Coolify routing.
- Coolify deployment logs for failed builds can sometimes be unavailable.
