# ── POST /api/tokens - mutation: scope empty string ──
# case_id=TC-c8cd2aed
# case_name=POST /api/tokens - mutation: scope empty string
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": "Clifford Ruiz",
  "scope": ""
}
```

HTTP *

[Asserts]
status >= 400
status < 500

