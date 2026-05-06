# ── POST /api/tokens - [unicode_fuzzing] name control_char ──
# case_id=TC-fc869137
# case_name=POST /api/tokens - [unicode_fuzzing] name control_char
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": "hello\u0000world",
  "scope": "read"
}
```

HTTP 400

