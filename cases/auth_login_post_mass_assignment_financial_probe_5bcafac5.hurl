# ── POST /auth/login - [mass_assignment] financial probe ──
# case_id=TC-5bcafac5
# case_name=POST /auth/login - [mass_assignment] financial probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "balance": 1,
  "credits": 1,
  "discount": 0,
  "email": "kriswong@koch.io",
  "password": "us",
  "price": 1
}
```

HTTP 400

