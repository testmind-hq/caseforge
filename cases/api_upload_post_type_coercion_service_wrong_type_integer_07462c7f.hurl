# ── POST /api/upload - [type_coercion] service wrong_type_integer ──
# case_id=TC-07462c7f
# case_name=POST /api/upload - [type_coercion] service wrong_type_integer
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
  "service": 123,
  "specContent": "you"
}
```

HTTP 422

