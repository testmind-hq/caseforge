# ── POST /api/tokens - invalid scope: value not in enum ──
# case_id=TC-a9cdb025
# case_name=POST /api/tokens - invalid scope: value not in enum
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P2

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": "Damion Rivera",
  "scope": "__invalid_enum__"
}
```

HTTP 422

