# ── PUT /api/admin/teams/{id}/members/{userId} - invalid role: value not in enum ──
# case_id=TC-1385a015
# case_name=PUT /api/admin/teams/{id}/members/{userId} - invalid role: value not in enum
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P2

PUT {{base_url}}/api/admin/teams/{id}/members/{userId}
Content-Type: application/json
```json
{
  "role": "__invalid_enum__"
}
```

HTTP 422

