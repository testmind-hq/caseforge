# ── PUT /api/admin/users/{id} - wrong content-type (text/plain) ──
# case_id=TC-69ba511c
# case_name=PUT /api/admin/users/{id} - wrong content-type (text/plain)
# step_id=step-main
# step_type=test
# technique=constraint_mutation
# priority=P2

PUT {{base_url}}/api/admin/users/{id}
Content-Type: text/plain
```json
{
  "isActive": false,
  "role": "super_admin"
}
```

HTTP 415

