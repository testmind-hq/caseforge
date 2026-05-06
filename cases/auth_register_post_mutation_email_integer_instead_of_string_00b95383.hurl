# ── POST /auth/register - mutation: email integer instead of string ──
# case_id=TC-00b95383
# case_name=POST /auth/register - mutation: email integer instead of string
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": 12345,
  "password": "where"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

