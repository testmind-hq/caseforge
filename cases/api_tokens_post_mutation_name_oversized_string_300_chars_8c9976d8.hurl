# ── POST /api/tokens - mutation: name oversized string (300 chars) ──
# case_id=TC-8c9976d8
# case_name=POST /api/tokens - mutation: name oversized string (300 chars)
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
  "scope": "write"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

