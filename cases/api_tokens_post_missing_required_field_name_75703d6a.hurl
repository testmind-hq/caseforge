# ── POST /api/tokens - missing required field "name" ──
# case_id=TC-75703d6a
# case_name=POST /api/tokens - missing required field "name"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "scope": "read"
}
```

HTTP 422

