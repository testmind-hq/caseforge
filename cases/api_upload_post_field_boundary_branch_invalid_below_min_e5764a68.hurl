# ── POST /api/upload - [field_boundary] branch invalid_below_min ──
# case_id=TC-e5764a68
# case_name=POST /api/upload - [field_boundary] branch invalid_below_min
# step_id=step-main
# step_type=test
# technique=field_boundary
# priority=P2

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "",
  "commitSha": "about",
  "service": "scold",
  "specContent": "muster"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

