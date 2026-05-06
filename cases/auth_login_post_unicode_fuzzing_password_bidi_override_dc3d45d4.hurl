# ── POST /auth/login - [unicode_fuzzing] password bidi_override ──
# case_id=TC-dc3d45d4
# case_name=POST /auth/login - [unicode_fuzzing] password bidi_override
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "email": "lilyperez@allen.io",
  "password": "‮hello"
}
```

HTTP 400

