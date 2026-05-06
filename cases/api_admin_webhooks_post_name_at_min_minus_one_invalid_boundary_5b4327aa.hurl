# ── POST /api/admin/webhooks - name at min_minus_one_invalid boundary ──
# case_id=TC-5b4327aa
# case_name=POST /api/admin/webhooks - name at min_minus_one_invalid boundary
# step_id=step-main
# step_type=test
# technique=boundary_value
# priority=P2

POST {{base_url}}/api/admin/webhooks
Content-Type: application/json
```json
{
  "events": [
    "these"
  ],
  "name": "b",
  "providerType": "infrequently",
  "teamId": "4a6f39f6-5059-431c-b5eb-9711769c6023",
  "url": "http://www.juniorexpedite.com/partnerships"
}
```

HTTP 422

