# ── POST /api/admin/teams - null injection: name ──
# case_id=TC-abe4e3e2
# case_name=POST /api/admin/teams - null injection: name
# step_id=step-main
# step_type=test
# technique=constraint_mutation
# priority=P2

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": "Explicitly name the person before you wrap it.",
  "displayName": "should",
  "name": null
}
```

HTTP 422

