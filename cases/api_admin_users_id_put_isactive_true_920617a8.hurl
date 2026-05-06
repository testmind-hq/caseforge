# ── PUT /api/admin/users/{id} - isActive = true ──
# case_id=TC-920617a8
# case_name=PUT /api/admin/users/{id} - isActive = true
# step_id=step-main
# step_type=test
# technique=decision_table
# priority=P1

PUT {{base_url}}/api/admin/users/{id}
Content-Type: application/json
```json
{
  "isActive": true,
  "role": "super_admin"
}
```

HTTP 200

[Asserts]
duration < 2000

