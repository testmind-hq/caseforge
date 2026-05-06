# ── POST /api/admin/teams - null injection: description ──
# case_id=TC-5294fe7b
# case_name=POST /api/admin/teams - null injection: description
# step_id=step-main
# step_type=test
# technique=constraint_mutation
# priority=P2

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": null,
  "displayName": "should",
  "name": "Chloe Oliver"
}
```

HTTP 422

