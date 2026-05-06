# ── POST /auth/login - [unicode_fuzzing] email control_char ──
# case_id=TC-ce646cde
# case_name=POST /auth/login - [unicode_fuzzing] email control_char
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "email": "hello\u0000world",
  "password": "themselves"
}
```

HTTP 400

