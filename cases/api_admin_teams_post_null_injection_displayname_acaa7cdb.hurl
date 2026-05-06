# ── POST /api/admin/teams - null injection: displayName ──
# case_id=TC-acaa7cdb
# case_name=POST /api/admin/teams - null injection: displayName
# step_id=step-main
# step_type=test
# technique=constraint_mutation
# priority=P2

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": "Explicitly name the person before you wrap it.",
  "displayName": null,
  "name": "Chloe Oliver"
}
```

HTTP 422

