# ── POST /auth/login - [mass_assignment] status probe ──
# case_id=TC-f197447f
# case_name=POST /auth/login - [mass_assignment] status probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "approved": true,
  "banned": false,
  "disabled": false,
  "email": "kriswong@koch.io",
  "password": "us",
  "verified": true
}
```

HTTP 400

