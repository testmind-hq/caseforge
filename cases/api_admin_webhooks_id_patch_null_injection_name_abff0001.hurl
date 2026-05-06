# ── PATCH /api/admin/webhooks/:id - null injection: name ──
# case_id=TC-abff0001
# case_name=PATCH /api/admin/webhooks/:id - null injection: name
# step_id=step-main
# step_type=test
# technique=constraint_mutation
# priority=P2

PATCH {{base_url}}/api/admin/webhooks/:id
Content-Type: application/json
```json
{
  "events": [
    "aloof"
  ],
  "isActive": true,
  "name": null,
  "url": "http://www.dynamicmarkets.net/vertical"
}
```

HTTP 422

