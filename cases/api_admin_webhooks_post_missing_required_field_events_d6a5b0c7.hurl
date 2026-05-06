# ── POST /api/admin/webhooks - missing required field "events" ──
# case_id=TC-d6a5b0c7
# case_name=POST /api/admin/webhooks - missing required field "events"
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P1

POST {{base_url}}/api/admin/webhooks
Content-Type: application/json
```json
{
  "name": "Rebecca Mann",
  "providerType": "painter",
  "teamId": "1485872f-38ec-4ac0-88b9-3d10f551b3a4",
  "url": "https://www.chiefsyndicate.biz/utilize/deliverables/innovate/transition"
}
```

HTTP 422

