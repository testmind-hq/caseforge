# ── POST /api/admin/teams/{id}/members - invalid role: value not in enum ──
# case_id=TC-54b6ea73
# case_name=POST /api/admin/teams/{id}/members - invalid role: value not in enum
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P2

POST {{base_url}}/api/admin/teams/{id}/members
Content-Type: application/json
```json
{
  "role": "__invalid_enum__",
  "userId": "45cf0fb5-a53d-4f38-94af-85fabe94e394"
}
```

HTTP 422

