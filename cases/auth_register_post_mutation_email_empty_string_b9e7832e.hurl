# ── POST /auth/register - mutation: email empty string ──
# case_id=TC-b9e7832e
# case_name=POST /auth/register - mutation: email empty string
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": "",
  "password": "where"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

