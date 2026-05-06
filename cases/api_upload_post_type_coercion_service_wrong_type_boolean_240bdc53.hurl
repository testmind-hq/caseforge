# ── POST /api/upload - [type_coercion] service wrong_type_boolean ──
# case_id=TC-240bdc53
# case_name=POST /api/upload - [type_coercion] service wrong_type_boolean
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
  "service": true,
  "specContent": "you"
}
```

HTTP 422

