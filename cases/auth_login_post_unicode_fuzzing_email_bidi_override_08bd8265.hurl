# ── POST /auth/login - [unicode_fuzzing] email bidi_override ──
# case_id=TC-08bd8265
# case_name=POST /auth/login - [unicode_fuzzing] email bidi_override
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "email": "‮hello",
  "password": "themselves"
}
```

HTTP 400

