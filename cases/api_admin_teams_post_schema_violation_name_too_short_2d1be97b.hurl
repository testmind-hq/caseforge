# ── POST /api/admin/teams - [schema_violation] name_too_short ──
# case_id=TC-2d1be97b
# case_name=POST /api/admin/teams - [schema_violation] name_too_short
# step_id=step-main
# step_type=test
# technique=schema_violation
# priority=P2

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": "Alert on person thresholds then.",
  "displayName": "most",
  "name": ""
}
```

HTTP 422

