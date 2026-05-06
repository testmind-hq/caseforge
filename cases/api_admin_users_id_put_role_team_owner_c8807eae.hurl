# ── PUT /api/admin/users/{id} - role = team_owner ──
# case_id=TC-c8807eae
# case_name=PUT /api/admin/users/{id} - role = team_owner
# step_id=step-main
# step_type=test
# technique=decision_table
# priority=P1

PUT {{base_url}}/api/admin/users/{id}
Content-Type: application/json
```json
{
  "isActive": true,
  "role": "team_owner"
}
```

HTTP 200

[Asserts]
duration < 2000

