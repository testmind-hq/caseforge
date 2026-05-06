# ── PATCH /api/admin/webhooks/:id - [unicode_fuzzing] name zero_width ──
# case_id=TC-6bdb26ba
# case_name=PATCH /api/admin/webhooks/:id - [unicode_fuzzing] name zero_width
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

PATCH {{base_url}}/api/admin/webhooks/:id
Content-Type: application/json
```json
{
  "events": [
    "that"
  ],
  "isActive": true,
  "name": "​hello",
  "url": "https://www.productdrive.io/grow/world-class"
}
```

HTTP 400

