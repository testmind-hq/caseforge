# ── POST /api/tokens - [unicode_fuzzing] scope control_char ──
# case_id=TC-0d728fca
# case_name=POST /api/tokens - [unicode_fuzzing] scope control_char
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": "Amelia Cummings",
  "scope": "hello\u0000world"
}
```

HTTP 400

