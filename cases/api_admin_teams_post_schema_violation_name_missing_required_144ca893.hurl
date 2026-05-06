# ── POST /api/admin/teams - [schema_violation] name_missing_required ──
# case_id=TC-144ca893
# case_name=POST /api/admin/teams - [schema_violation] name_missing_required
# step_id=step-main
# step_type=test
# technique=schema_violation
# priority=P2

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": "Alert on person thresholds then.",
  "displayName": "most"
}
```

HTTP 422

