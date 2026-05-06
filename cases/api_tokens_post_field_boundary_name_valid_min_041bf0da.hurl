# ── POST /api/tokens - [field_boundary] name valid_min ──
# case_id=TC-041bf0da
# case_name=POST /api/tokens - [field_boundary] name valid_min
# step_id=step-main
# step_type=test
# technique=field_boundary
# priority=P1

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": "a",
  "scope": "read"
}
```

HTTP *

[Asserts]
status >= 200
status < 300

