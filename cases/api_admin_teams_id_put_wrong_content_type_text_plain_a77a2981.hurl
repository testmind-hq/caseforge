# ── PUT /api/admin/teams/{id} - wrong content-type (text/plain) ──
# case_id=TC-a77a2981
# case_name=PUT /api/admin/teams/{id} - wrong content-type (text/plain)
# step_id=step-main
# step_type=test
# technique=constraint_mutation
# priority=P2

PUT {{base_url}}/api/admin/teams/{id}
Content-Type: text/plain
```json
{
  "description": "Publish a changelog entry for the work.",
  "displayName": "nervous"
}
```

HTTP 415

