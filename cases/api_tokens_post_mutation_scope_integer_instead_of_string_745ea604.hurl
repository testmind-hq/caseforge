# ── POST /api/tokens - mutation: scope integer instead of string ──
# case_id=TC-745ea604
# case_name=POST /api/tokens - mutation: scope integer instead of string
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": "Clifford Ruiz",
  "scope": 12345
}
```

HTTP *

[Asserts]
status >= 400
status < 500

