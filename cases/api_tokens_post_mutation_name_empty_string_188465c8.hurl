# ── POST /api/tokens - mutation: name empty string ──
# case_id=TC-188465c8
# case_name=POST /api/tokens - mutation: name empty string
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": "",
  "scope": "write"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

