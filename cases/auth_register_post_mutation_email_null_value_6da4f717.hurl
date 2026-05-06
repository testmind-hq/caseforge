# ── POST /auth/register - mutation: email null value ──
# case_id=TC-6da4f717
# case_name=POST /auth/register - mutation: email null value
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": null,
  "password": "where"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

