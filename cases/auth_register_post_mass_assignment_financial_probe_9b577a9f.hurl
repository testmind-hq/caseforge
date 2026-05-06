# ── POST /auth/register - [mass_assignment] financial probe ──
# case_id=TC-9b577a9f
# case_name=POST /auth/register - [mass_assignment] financial probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "balance": 1,
  "credits": 1,
  "discount": 0,
  "email": "waynedaniels@farrell.io",
  "password": "instead",
  "price": 1
}
```

HTTP 400

