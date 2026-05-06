# ── PATCH /api/admin/webhooks/:id - mutation: events object instead of array ──
# case_id=TC-309789e7
# case_name=PATCH /api/admin/webhooks/:id - mutation: events object instead of array
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

PATCH {{base_url}}/api/admin/webhooks/:id
Content-Type: application/json
```json
{
  "events": {},
  "isActive": false,
  "name": "Kristin Burton",
  "url": "http://www.financialpartnerships.org/metrics/cross-platform/platforms/viral"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

