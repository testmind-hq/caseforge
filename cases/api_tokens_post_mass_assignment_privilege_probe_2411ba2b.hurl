# ── POST /api/tokens - [mass_assignment] privilege probe ──
# case_id=TC-2411ba2b
# case_name=POST /api/tokens - [mass_assignment] privilege probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "admin": true,
  "isAdmin": true,
  "is_admin": true,
  "name": "Jalen Phillips",
  "role": "__probe__",
  "scope": "write"
}
```

HTTP 400

