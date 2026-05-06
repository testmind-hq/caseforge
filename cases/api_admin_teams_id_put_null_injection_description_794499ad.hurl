# ── PUT /api/admin/teams/{id} - null injection: description ──
# case_id=TC-794499ad
# case_name=PUT /api/admin/teams/{id} - null injection: description
# step_id=step-main
# step_type=test
# technique=constraint_mutation
# priority=P2

PUT {{base_url}}/api/admin/teams/{id}
Content-Type: application/json
```json
{
  "description": null,
  "displayName": "nervous"
}
```

HTTP 422

