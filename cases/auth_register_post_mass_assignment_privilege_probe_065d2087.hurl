# ── POST /auth/register - [mass_assignment] privilege probe ──
# case_id=TC-065d2087
# case_name=POST /auth/register - [mass_assignment] privilege probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "admin": true,
  "email": "waynedaniels@farrell.io",
  "isAdmin": true,
  "is_admin": true,
  "password": "instead",
  "role": "__probe__"
}
```

HTTP 400

