# ── POST /api/upload - [field_boundary] service invalid_below_min ──
# case_id=TC-a957f4b8
# case_name=POST /api/upload - [field_boundary] service invalid_below_min
# step_id=step-main
# step_type=test
# technique=field_boundary
# priority=P2

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "next",
  "commitSha": "none",
  "service": "",
  "specContent": "through"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

