# ── PATCH /api/admin/webhooks/:id - mutation: name empty string ──
# case_id=TC-48b3b8ee
# case_name=PATCH /api/admin/webhooks/:id - mutation: name empty string
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
  "name": "",
  "url": "http://www.financialpartnerships.org/metrics/cross-platform/platforms/viral"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

