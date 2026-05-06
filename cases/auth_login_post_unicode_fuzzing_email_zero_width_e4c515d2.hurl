# ── POST /auth/login - [unicode_fuzzing] email zero_width ──
# case_id=TC-e4c515d2
# case_name=POST /auth/login - [unicode_fuzzing] email zero_width
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "email": "​hello",
  "password": "themselves"
}
```

HTTP 400

