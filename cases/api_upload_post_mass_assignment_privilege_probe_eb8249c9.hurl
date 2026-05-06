# ── POST /api/upload - [mass_assignment] privilege probe ──
# case_id=TC-eb8249c9
# case_name=POST /api/upload - [mass_assignment] privilege probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "admin": true,
  "branch": "oops",
  "commitSha": "mustering",
  "isAdmin": true,
  "is_admin": true,
  "role": "__probe__",
  "service": "I",
  "specContent": "cut"
}
```

HTTP 400

