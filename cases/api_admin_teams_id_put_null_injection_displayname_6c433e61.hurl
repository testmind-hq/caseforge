# ── PUT /api/admin/teams/{id} - null injection: displayName ──
# case_id=TC-6c433e61
# case_name=PUT /api/admin/teams/{id} - null injection: displayName
# step_id=step-main
# step_type=test
# technique=constraint_mutation
# priority=P2

PUT {{base_url}}/api/admin/teams/{id}
Content-Type: application/json
```json
{
  "description": "Publish a changelog entry for the work.",
  "displayName": null
}
```

HTTP 422

