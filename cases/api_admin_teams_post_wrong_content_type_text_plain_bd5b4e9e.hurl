# ── POST /api/admin/teams - wrong content-type (text/plain) ──
# case_id=TC-bd5b4e9e
# case_name=POST /api/admin/teams - wrong content-type (text/plain)
# step_id=step-main
# step_type=test
# technique=constraint_mutation
# priority=P2

POST {{base_url}}/api/admin/teams
Content-Type: text/plain
```json
{
  "description": "Explicitly name the person before you wrap it.",
  "displayName": "should",
  "name": "Chloe Oliver"
}
```

HTTP 415

