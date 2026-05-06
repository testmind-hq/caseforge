# ── PATCH /api/admin/webhooks/:id - mutation: events null value ──
# case_id=TC-2d09c873
# case_name=PATCH /api/admin/webhooks/:id - mutation: events null value
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

PATCH {{base_url}}/api/admin/webhooks/:id
Content-Type: application/json
```json
{
  "events": null,
  "isActive": false,
  "name": "Kristin Burton",
  "url": "http://www.financialpartnerships.org/metrics/cross-platform/platforms/viral"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

