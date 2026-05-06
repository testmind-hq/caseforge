# ── POST /api/tokens - [unicode_fuzzing] scope zero_width ──
# case_id=TC-6b8f84d1
# case_name=POST /api/tokens - [unicode_fuzzing] scope zero_width
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": "Amelia Cummings",
  "scope": "​hello"
}
```

HTTP 400

