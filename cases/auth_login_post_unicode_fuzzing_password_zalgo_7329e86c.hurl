# ── POST /auth/login - [unicode_fuzzing] password zalgo ──
# case_id=TC-7329e86c
# case_name=POST /auth/login - [unicode_fuzzing] password zalgo
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "email": "lilyperez@allen.io",
  "password": "z̀́̂̃̄̅̆̇a"
}
```

HTTP 400

