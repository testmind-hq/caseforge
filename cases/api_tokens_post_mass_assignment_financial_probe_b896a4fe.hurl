# ── POST /api/tokens - [mass_assignment] financial probe ──
# case_id=TC-b896a4fe
# case_name=POST /api/tokens - [mass_assignment] financial probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "balance": 1,
  "credits": 1,
  "discount": 0,
  "name": "Jalen Phillips",
  "price": 1,
  "scope": "write"
}
```

HTTP 400

