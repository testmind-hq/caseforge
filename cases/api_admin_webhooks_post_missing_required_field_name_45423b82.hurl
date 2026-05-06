# ── POST /api/admin/webhooks - missing required field "name" ──
# case_id=TC-45423b82
# case_name=POST /api/admin/webhooks - missing required field "name"
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P1

POST {{base_url}}/api/admin/webhooks
Content-Type: application/json
```json
{
  "events": [
    "Plutonian"
  ],
  "providerType": "choir",
  "teamId": "5289bf89-a443-44f7-a319-2a66891988ac",
  "url": "https://www.humandeploy.io/magnetic/roi/maximize/embrace"
}
```

HTTP 422

