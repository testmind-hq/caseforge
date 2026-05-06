# ── PUT /api/admin/teams/{id} - [mass_assignment] identity probe ──
# case_id=TC-ed4e87e7
# case_name=PUT /api/admin/teams/{id} - [mass_assignment] identity probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

PUT {{base_url}}/api/admin/teams/{id}
Content-Type: application/json
```json
{
  "createdBy": "__probe__",
  "description": "Alert on way thresholds yesterday.",
  "displayName": "this",
  "ownerId": "__probe__",
  "userId": "__probe__",
  "user_id": "__probe__"
}
```

HTTP 400

