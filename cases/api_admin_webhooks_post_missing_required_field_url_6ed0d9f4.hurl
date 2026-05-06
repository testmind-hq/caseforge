# ── POST /api/admin/webhooks - missing required field "url" ──
# case_id=TC-6ed0d9f4
# case_name=POST /api/admin/webhooks - missing required field "url"
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P1

POST {{base_url}}/api/admin/webhooks
Content-Type: application/json
```json
{
  "events": [
    "the"
  ],
  "name": "Carey Jimenez",
  "providerType": "hourly",
  "teamId": "68326c3d-2def-4030-9c4f-dfcb153eda58"
}
```

HTTP 422

