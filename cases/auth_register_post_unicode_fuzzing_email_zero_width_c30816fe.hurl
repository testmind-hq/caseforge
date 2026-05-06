# ── POST /auth/register - [unicode_fuzzing] email zero_width ──
# case_id=TC-c30816fe
# case_name=POST /auth/register - [unicode_fuzzing] email zero_width
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": "​hello",
  "password": "every"
}
```

HTTP 400

