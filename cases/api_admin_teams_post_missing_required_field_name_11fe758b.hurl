# ── POST /api/admin/teams - missing required field "name" ──
# case_id=TC-11fe758b
# case_name=POST /api/admin/teams - missing required field "name"
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P1

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": "Celebrate wins tied to the man.",
  "displayName": "lastly"
}
```

HTTP 422

