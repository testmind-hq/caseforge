# ── POST /api/tokens - [required_omission] name absent ──
# case_id=TC-b998dc1a
# case_name=POST /api/tokens - [required_omission] name absent
# step_id=step-main
# step_type=test
# technique=required_omission
# priority=P2

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "scope": "write"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

