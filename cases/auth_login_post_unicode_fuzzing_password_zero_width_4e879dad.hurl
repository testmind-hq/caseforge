# ── POST /auth/login - [unicode_fuzzing] password zero_width ──
# case_id=TC-4e879dad
# case_name=POST /auth/login - [unicode_fuzzing] password zero_width
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "email": "lilyperez@allen.io",
  "password": "​hello"
}
```

HTTP 400

