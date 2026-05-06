# ── POST /api/upload - [type_coercion] branch wrong_type_integer ──
# case_id=TC-6a08feec
# case_name=POST /api/upload - [type_coercion] branch wrong_type_integer
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": 123,
  "commitSha": "throw",
  "service": "the",
  "specContent": "you"
}
```

HTTP 422

