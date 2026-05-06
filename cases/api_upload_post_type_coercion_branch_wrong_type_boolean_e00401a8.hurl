# ── POST /api/upload - [type_coercion] branch wrong_type_boolean ──
# case_id=TC-e00401a8
# case_name=POST /api/upload - [type_coercion] branch wrong_type_boolean
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": true,
  "commitSha": "throw",
  "service": "the",
  "specContent": "you"
}
```

HTTP 422

