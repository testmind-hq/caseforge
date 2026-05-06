# ── PATCH /api/admin/webhooks/:id - [unicode_fuzzing] name zalgo ──
# case_id=TC-a31d1299
# case_name=PATCH /api/admin/webhooks/:id - [unicode_fuzzing] name zalgo
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
  "name": "z̀́̂̃̄̅̆̇a",
  "url": "https://www.productdrive.io/grow/world-class"
}
```

HTTP 400

