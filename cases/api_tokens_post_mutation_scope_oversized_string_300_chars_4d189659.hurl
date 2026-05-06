# ── POST /api/tokens - mutation: scope oversized string (300 chars) ──
# case_id=TC-4d189659
# case_name=POST /api/tokens - mutation: scope oversized string (300 chars)
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": "Clifford Ruiz",
  "scope": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

