# ── PUT /api/admin/users/{id} - role = super_admin ──
# case_id=TC-72c28c85
# case_name=PUT /api/admin/users/{id} - role = super_admin
# step_id=step-main
# step_type=test
# technique=decision_table
# priority=P1

PUT {{base_url}}/api/admin/users/{id}
Content-Type: application/json
```json
{
  "isActive": false,
  "role": "super_admin"
}
```

HTTP 200

[Asserts]
duration < 2000

