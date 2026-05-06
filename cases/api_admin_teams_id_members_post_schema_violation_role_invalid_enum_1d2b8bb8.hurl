# ── POST /api/admin/teams/{id}/members - [schema_violation] role_invalid_enum ──
# case_id=TC-1d2b8bb8
# case_name=POST /api/admin/teams/{id}/members - [schema_violation] role_invalid_enum
# step_id=step-main
# step_type=test
# technique=schema_violation
# priority=P2

POST {{base_url}}/api/admin/teams/{id}/members
Content-Type: application/json
```json
{
  "role": "__invalid__",
  "userId": "b28b1b32-e5b1-4269-b005-d53ff9fd5a8d"
}
```

HTTP 422

