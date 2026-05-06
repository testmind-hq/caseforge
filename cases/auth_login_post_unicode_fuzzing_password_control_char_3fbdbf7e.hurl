# ── POST /auth/login - [unicode_fuzzing] password control_char ──
# case_id=TC-3fbdbf7e
# case_name=POST /auth/login - [unicode_fuzzing] password control_char
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "email": "lilyperez@allen.io",
  "password": "hello\u0000world"
}
```

HTTP 400

