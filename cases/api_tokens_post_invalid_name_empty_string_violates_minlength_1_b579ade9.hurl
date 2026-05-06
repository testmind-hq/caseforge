# ── POST /api/tokens - invalid name: empty string violates minLength 1 ──
# case_id=TC-b579ade9
# case_name=POST /api/tokens - invalid name: empty string violates minLength 1
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P2

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": "",
  "scope": "read"
}
```

HTTP 422

