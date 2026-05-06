# ── POST /auth/register - [unicode_fuzzing] email bidi_override ──
# case_id=TC-cd50c303
# case_name=POST /auth/register - [unicode_fuzzing] email bidi_override
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": "‮hello",
  "password": "every"
}
```

HTTP 400

