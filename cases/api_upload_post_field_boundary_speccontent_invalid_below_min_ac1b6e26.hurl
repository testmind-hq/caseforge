# ── POST /api/upload - [field_boundary] specContent invalid_below_min ──
# case_id=TC-ac1b6e26
# case_name=POST /api/upload - [field_boundary] specContent invalid_below_min
# step_id=step-main
# step_type=test
# technique=field_boundary
# priority=P2

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "whom",
  "commitSha": "to",
  "service": "constantly",
  "specContent": ""
}
```

HTTP *

[Asserts]
status >= 400
status < 500

