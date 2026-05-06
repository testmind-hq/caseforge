# ── POST /auth/register - [unicode_fuzzing] password bidi_override ──
# case_id=TC-28ca4955
# case_name=POST /auth/register - [unicode_fuzzing] password bidi_override
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": "charityross@barber.biz",
  "password": "‮hello"
}
```

HTTP 400

