**User Preferences**

- Preferred language: Italian (fallback concise English).
- Wants direct operational answers, minimal back-and-forth.
- Primary focus: runtime bots in `paper` and `live`.
- For profitability questions, prioritize DB ROI ranking over ad-hoc estimates.

**Typical Requests**

- "fammi ranking ROI bot paper/live"
- "quali bot live sono unhealthy"
- "dammi top e bottom pnl (finestra dinamica)"
- "fammi digest unico stato + performance"
- "mandami report periodico ogni 3 ore"
- "fammi vedere la ruota trading e switch strategia"

**Mutation Preference**

- For switch/restart/deploy requests: require a short impact summary and explicit confirmation before action.

**Actionable Reply Expectations**

- Include bot scope explicitly: `paper`, `live`, or `all`.
- Include source explicitly: `db` or `api_fallback`.
- Include window explicitly: `auto` or `<n>d`.
- For strategy-switch ideas, always mention affected strategy slugs from:
`${MAXEXTRACT_ROOT}/strategies/`
