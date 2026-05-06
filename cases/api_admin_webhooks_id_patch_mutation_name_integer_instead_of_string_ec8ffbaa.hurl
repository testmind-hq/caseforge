# ── PATCH /api/admin/webhooks/:id - mutation: name integer instead of string ──
# case_id=TC-ec8ffbaa
# case_name=PATCH /api/admin/webhooks/:id - mutation: name integer instead of string
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
  "name": 12345,
  "url": "http://www.financialpartnerships.org/metrics/cross-platform/platforms/viral"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

