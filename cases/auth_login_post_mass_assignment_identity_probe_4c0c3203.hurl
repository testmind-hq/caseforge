# ── POST /auth/login - [mass_assignment] identity probe ──
# case_id=TC-4c0c3203
# case_name=POST /auth/login - [mass_assignment] identity probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "createdBy": "__probe__",
  "email": "kriswong@koch.io",
  "ownerId": "__probe__",
  "password": "us",
  "userId": "__probe__",
  "user_id": "__probe__"
}
```

HTTP 400

