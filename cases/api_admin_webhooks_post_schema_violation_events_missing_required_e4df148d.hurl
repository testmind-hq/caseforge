# ── POST /api/admin/webhooks - [schema_violation] events_missing_required ──
# case_id=TC-e4df148d
# case_name=POST /api/admin/webhooks - [schema_violation] events_missing_required
# step_id=step-main
# step_type=test
# technique=schema_violation
# priority=P2

POST {{base_url}}/api/admin/webhooks
Content-Type: application/json
```json
{
  "name": "Raphael Davies",
  "providerType": "me",
  "teamId": "8afc12a7-a242-4e1f-b05b-4ade3fb01c0f",
  "url": "https://www.legacyincubate.io/seize"
}
```

HTTP 422

