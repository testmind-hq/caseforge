# ── PUT /api/admin/users/{id} - role = guest ──
# case_id=TC-d671319d
# case_name=PUT /api/admin/users/{id} - role = guest
# step_id=step-main
# step_type=test
# technique=decision_table
# priority=P1

PUT {{base_url}}/api/admin/users/{id}
Content-Type: application/json
```json
{
  "isActive": false,
  "role": "guest"
}
```

HTTP 200

[Asserts]
duration < 2000

