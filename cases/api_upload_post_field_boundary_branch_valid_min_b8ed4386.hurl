# ── POST /api/upload - [field_boundary] branch valid_min ──
# case_id=TC-b8ed4386
# case_name=POST /api/upload - [field_boundary] branch valid_min
# step_id=step-main
# step_type=test
# technique=field_boundary
# priority=P1

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "a",
  "commitSha": "girl",
  "service": "those",
  "specContent": "many"
}
```

HTTP *

[Asserts]
status >= 200
status < 300

