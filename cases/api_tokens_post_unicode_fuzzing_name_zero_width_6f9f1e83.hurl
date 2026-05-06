# ── POST /api/tokens - [unicode_fuzzing] name zero_width ──
# case_id=TC-6f9f1e83
# case_name=POST /api/tokens - [unicode_fuzzing] name zero_width
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": "​hello",
  "scope": "read"
}
```

HTTP 400

