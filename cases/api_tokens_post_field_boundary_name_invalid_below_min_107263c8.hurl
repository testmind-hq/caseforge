# ── POST /api/tokens - [field_boundary] name invalid_below_min ──
# case_id=TC-107263c8
# case_name=POST /api/tokens - [field_boundary] name invalid_below_min
# step_id=step-main
# step_type=test
# technique=field_boundary
# priority=P2

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": "",
  "scope": "read"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

