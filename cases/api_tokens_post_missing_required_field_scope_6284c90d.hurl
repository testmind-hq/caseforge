# ── POST /api/tokens - missing required field "scope" ──
# case_id=TC-6284c90d
# case_name=POST /api/tokens - missing required field "scope"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": "Damion Rivera"
}
```

HTTP 422

