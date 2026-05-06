# ── POST /auth/register - [mass_assignment] identity probe ──
# case_id=TC-be5d4ca2
# case_name=POST /auth/register - [mass_assignment] identity probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "createdBy": "__probe__",
  "email": "waynedaniels@farrell.io",
  "ownerId": "__probe__",
  "password": "instead",
  "userId": "__probe__",
  "user_id": "__probe__"
}
```

HTTP 400

