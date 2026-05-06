# ── POST /auth/register - mutation: password integer instead of string ──
# case_id=TC-85af6488
# case_name=POST /auth/register - mutation: password integer instead of string
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": "marjoriecole@donnelly.org",
  "password": 12345
}
```

HTTP *

[Asserts]
status >= 400
status < 500

