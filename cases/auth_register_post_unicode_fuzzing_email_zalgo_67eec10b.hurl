# ── POST /auth/register - [unicode_fuzzing] email zalgo ──
# case_id=TC-67eec10b
# case_name=POST /auth/register - [unicode_fuzzing] email zalgo
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": "z̀́̂̃̄̅̆̇a",
  "password": "every"
}
```

HTTP 400

