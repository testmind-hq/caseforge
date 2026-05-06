# ── POST /api/tokens - [mass_assignment] status probe ──
# case_id=TC-248852e9
# case_name=POST /api/tokens - [mass_assignment] status probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "approved": true,
  "banned": false,
  "disabled": false,
  "name": "Jalen Phillips",
  "scope": "write",
  "verified": true
}
```

HTTP 400

