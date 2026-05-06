# ── POST /api/upload - [type_coercion] specContent wrong_type_integer ──
# case_id=TC-bbde20a6
# case_name=POST /api/upload - [type_coercion] specContent wrong_type_integer
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "point",
  "commitSha": "throw",
  "service": "the",
  "specContent": 123
}
```

HTTP 422

