# ── POST /api/upload - [mass_assignment] status probe ──
# case_id=TC-0310fa1a
# case_name=POST /api/upload - [mass_assignment] status probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "approved": true,
  "banned": false,
  "branch": "oops",
  "commitSha": "mustering",
  "disabled": false,
  "service": "I",
  "specContent": "cut",
  "verified": true
}
```

HTTP 400

