# ── PUT /api/admin/users/{id} - isActive = false ──
# case_id=TC-307b2101
# case_name=PUT /api/admin/users/{id} - isActive = false
# step_id=step-main
# step_type=test
# technique=decision_table
# priority=P1

PUT {{base_url}}/api/admin/users/{id}
Content-Type: application/json
```json
{
  "isActive": false,
  "role": "team_member"
}
```

HTTP 200

[Asserts]
duration < 2000

