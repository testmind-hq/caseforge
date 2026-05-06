# ── POST /api/upload - [field_boundary] service valid_min ──
# case_id=TC-db5c5368
# case_name=POST /api/upload - [field_boundary] service valid_min
# step_id=step-main
# step_type=test
# technique=field_boundary
# priority=P1

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "it",
  "commitSha": "why",
  "service": "a",
  "specContent": "all"
}
```

HTTP *

[Asserts]
status >= 200
status < 300

