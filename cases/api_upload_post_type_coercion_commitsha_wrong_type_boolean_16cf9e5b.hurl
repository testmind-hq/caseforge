# ── POST /api/upload - [type_coercion] commitSha wrong_type_boolean ──
# case_id=TC-16cf9e5b
# case_name=POST /api/upload - [type_coercion] commitSha wrong_type_boolean
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "point",
  "commitSha": true,
  "service": "the",
  "specContent": "you"
}
```

HTTP 422

