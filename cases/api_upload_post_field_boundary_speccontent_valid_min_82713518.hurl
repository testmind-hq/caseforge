# ── POST /api/upload - [field_boundary] specContent valid_min ──
# case_id=TC-82713518
# case_name=POST /api/upload - [field_boundary] specContent valid_min
# step_id=step-main
# step_type=test
# technique=field_boundary
# priority=P1

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "drink",
  "commitSha": "his",
  "service": "few",
  "specContent": "a"
}
```

HTTP *

[Asserts]
status >= 200
status < 300

