# ── POST /auth/register - [unicode_fuzzing] password zero_width ──
# case_id=TC-e4e8966c
# case_name=POST /auth/register - [unicode_fuzzing] password zero_width
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": "charityross@barber.biz",
  "password": "​hello"
}
```

HTTP 400

