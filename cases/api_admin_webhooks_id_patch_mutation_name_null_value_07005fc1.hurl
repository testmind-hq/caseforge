# ── PATCH /api/admin/webhooks/:id - mutation: name null value ──
# case_id=TC-07005fc1
# case_name=PATCH /api/admin/webhooks/:id - mutation: name null value
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

PATCH {{base_url}}/api/admin/webhooks/:id
Content-Type: application/json
```json
{
  "events": [
    "might"
  ],
  "isActive": false,
  "name": null,
  "url": "http://www.financialpartnerships.org/metrics/cross-platform/platforms/viral"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

