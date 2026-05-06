# ── POST /auth/register - [unicode_fuzzing] password zalgo ──
# case_id=TC-ab0475dc
# case_name=POST /auth/register - [unicode_fuzzing] password zalgo
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": "charityross@barber.biz",
  "password": "z̀́̂̃̄̅̆̇a"
}
```

HTTP 400

