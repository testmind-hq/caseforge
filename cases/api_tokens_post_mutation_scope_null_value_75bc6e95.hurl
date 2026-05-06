# ── POST /api/tokens - mutation: scope null value ──
# case_id=TC-75bc6e95
# case_name=POST /api/tokens - mutation: scope null value
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": "Clifford Ruiz",
  "scope": null
}
```

HTTP *

[Asserts]
status >= 400
status < 500

