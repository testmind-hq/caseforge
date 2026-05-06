# ── POST /auth/register - mutation: password null value ──
# case_id=TC-8df134ff
# case_name=POST /auth/register - mutation: password null value
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": "marjoriecole@donnelly.org",
  "password": null
}
```

HTTP *

[Asserts]
status >= 400
status < 500

