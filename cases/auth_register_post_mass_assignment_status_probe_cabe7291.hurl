# ── POST /auth/register - [mass_assignment] status probe ──
# case_id=TC-cabe7291
# case_name=POST /auth/register - [mass_assignment] status probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "approved": true,
  "banned": false,
  "disabled": false,
  "email": "waynedaniels@farrell.io",
  "password": "instead",
  "verified": true
}
```

HTTP 400

