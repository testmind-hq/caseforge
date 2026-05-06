# ── POST /auth/login - [unicode_fuzzing] email zalgo ──
# case_id=TC-1091cce6
# case_name=POST /auth/login - [unicode_fuzzing] email zalgo
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "email": "z̀́̂̃̄̅̆̇a",
  "password": "themselves"
}
```

HTTP 400

