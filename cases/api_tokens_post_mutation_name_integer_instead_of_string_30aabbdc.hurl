# ── POST /api/tokens - mutation: name integer instead of string ──
# case_id=TC-30aabbdc
# case_name=POST /api/tokens - mutation: name integer instead of string
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": 12345,
  "scope": "write"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

