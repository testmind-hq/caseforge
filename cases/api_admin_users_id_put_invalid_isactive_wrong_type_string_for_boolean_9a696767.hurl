# ── PUT /api/admin/users/{id} - invalid isActive: wrong type (string for boolean) ──
# case_id=TC-9a696767
# case_name=PUT /api/admin/users/{id} - invalid isActive: wrong type (string for boolean)
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P2

PUT {{base_url}}/api/admin/users/{id}
Content-Type: application/json
```json
{
  "isActive": "not_a_boolean",
  "role": "super_admin"
}
```

HTTP 422

