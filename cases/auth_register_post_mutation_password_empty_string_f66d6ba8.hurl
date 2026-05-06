# ── POST /auth/register - mutation: password empty string ──
# case_id=TC-f66d6ba8
# case_name=POST /auth/register - mutation: password empty string
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": "marjoriecole@donnelly.org",
  "password": ""
}
```

HTTP *

[Asserts]
status >= 400
status < 500

