# ── POST /auth/register - null injection: email ──
# case_id=TC-031620b5
# case_name=POST /auth/register - null injection: email
# step_id=step-main
# step_type=test
# technique=constraint_mutation
# priority=P2

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": null,
  "password": "mouth"
}
```

HTTP 422

