# ── POST /api/upload - [type_coercion] commitSha wrong_type_integer ──
# case_id=TC-b806224f
# case_name=POST /api/upload - [type_coercion] commitSha wrong_type_integer
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "point",
  "commitSha": 123,
  "service": "the",
  "specContent": "you"
}
```

HTTP 422

