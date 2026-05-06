# ── POST /auth/register - [unicode_fuzzing] email control_char ──
# case_id=TC-619e4131
# case_name=POST /auth/register - [unicode_fuzzing] email control_char
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": "hello\u0000world",
  "password": "every"
}
```

HTTP 400

