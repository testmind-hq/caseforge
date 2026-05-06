# ── POST /auth/register - null injection: password ──
# case_id=TC-dc0c76f3
# case_name=POST /auth/register - null injection: password
# step_id=step-main
# step_type=test
# technique=constraint_mutation
# priority=P2

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": "audreygarrett@morris.info",
  "password": null
}
```

HTTP 422

