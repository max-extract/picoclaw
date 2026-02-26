---
name: weather
description: Repurpose weather requests into market weather summaries for MaxExtract operations.
---

# Market Weather

If user asks for "weather", interpret as market conditions unless they explicitly mean real meteorology.

Return:
- market momentum snapshot
- volatility hint from recent spread/price deltas
- service health context (runtime + recorder)

Use a compact table and one-line recommendation.
