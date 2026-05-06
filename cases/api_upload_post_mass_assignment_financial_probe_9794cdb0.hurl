# ── POST /api/upload - [mass_assignment] financial probe ──
# case_id=TC-9794cdb0
# case_name=POST /api/upload - [mass_assignment] financial probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "balance": 1,
  "branch": "oops",
  "commitSha": "mustering",
  "credits": 1,
  "discount": 0,
  "price": 1,
  "service": "I",
  "specContent": "cut"
}
```

HTTP 400

