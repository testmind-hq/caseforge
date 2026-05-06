# ── POST /api/admin/teams - name at min_valid boundary ──
# case_id=TC-084178e7
# case_name=POST /api/admin/teams - name at min_valid boundary
# step_id=step-main
# step_type=test
# technique=boundary_value
# priority=P1

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": "Set a realistic target for year.",
  "displayName": "moreover",
  "name": "X"
}
```

HTTP 200

[Asserts]
duration < 2000

