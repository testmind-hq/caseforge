# ── POST /api/admin/teams - name at min_minus_one_invalid boundary ──
# case_id=TC-2ccbadc2
# case_name=POST /api/admin/teams - name at min_minus_one_invalid boundary
# step_id=step-main
# step_type=test
# technique=boundary_value
# priority=P2

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": "Set a realistic target for year.",
  "displayName": "moreover",
  "name": "s"
}
```

HTTP 422

