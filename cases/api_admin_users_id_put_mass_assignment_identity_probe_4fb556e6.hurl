# ── PUT /api/admin/users/{id} - [mass_assignment] identity probe ──
# case_id=TC-4fb556e6
# case_name=PUT /api/admin/users/{id} - [mass_assignment] identity probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

PUT {{base_url}}/api/admin/users/{id}
Content-Type: application/json
```json
{
  "createdBy": "__probe__",
  "isActive": true,
  "ownerId": "__probe__",
  "role": "super_admin",
  "userId": "__probe__",
  "user_id": "__probe__"
}
```

HTTP 400

