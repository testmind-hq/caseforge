# ── POST /auth/register - [unicode_fuzzing] password control_char ──
# case_id=TC-cd54b4b0
# case_name=POST /auth/register - [unicode_fuzzing] password control_char
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": "charityross@barber.biz",
  "password": "hello\u0000world"
}
```

HTTP 400

