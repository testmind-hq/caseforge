# ── PATCH /api/admin/webhooks/:id - mutation: isActive null value ──
# case_id=TC-c42eb537
# case_name=PATCH /api/admin/webhooks/:id - mutation: isActive null value
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
  "isActive": null,
  "name": "Kristin Burton",
  "url": "http://www.financialpartnerships.org/metrics/cross-platform/platforms/viral"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

