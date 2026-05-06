# ── POST /auth/login - [mass_assignment] privilege probe ──
# case_id=TC-f4f54666
# case_name=POST /auth/login - [mass_assignment] privilege probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "admin": true,
  "email": "kriswong@koch.io",
  "isAdmin": true,
  "is_admin": true,
  "password": "us",
  "role": "__probe__"
}
```

HTTP 400

