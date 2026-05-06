# ── POST /api/tokens - mutation: name null value ──
# case_id=TC-816809db
# case_name=POST /api/tokens - mutation: name null value
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": null,
  "scope": "write"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

