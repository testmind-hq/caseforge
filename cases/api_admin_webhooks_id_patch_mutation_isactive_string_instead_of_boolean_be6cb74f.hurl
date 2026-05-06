# ── PATCH /api/admin/webhooks/:id - mutation: isActive string instead of boolean ──
# case_id=TC-be6cb74f
# case_name=PATCH /api/admin/webhooks/:id - mutation: isActive string instead of boolean
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
  "isActive": "yes",
  "name": "Kristin Burton",
  "url": "http://www.financialpartnerships.org/metrics/cross-platform/platforms/viral"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

