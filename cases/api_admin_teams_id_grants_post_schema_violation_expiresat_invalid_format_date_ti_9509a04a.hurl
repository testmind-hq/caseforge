# ── POST /api/admin/teams/{id}/grants - [schema_violation] expiresAt_invalid_format_date-time ──
# case_id=TC-9509a04a
# case_name=POST /api/admin/teams/{id}/grants - [schema_violation] expiresAt_invalid_format_date-time
# step_id=step-main
# step_type=test
# technique=schema_violation
# priority=P2

POST {{base_url}}/api/admin/teams/{id}/grants
Content-Type: application/json
```json
{
  "branches": [
    "am"
  ],
  "expiresAt": "not-a-date",
  "granteeTeamId": "7a8e7c06-efab-4a89-8471-23bbf2a20eea",
  "granteeUserId": "55b411ae-4ae9-4cf6-802a-a4a242203443",
  "serviceId": "435a1f1c-09a1-4465-b8ad-2053fa825257"
}
```

HTTP 422

