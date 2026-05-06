# ── PATCH /api/admin/webhooks/:id - [unicode_fuzzing] name control_char ──
# case_id=TC-9fed73af
# case_name=PATCH /api/admin/webhooks/:id - [unicode_fuzzing] name control_char
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
  "name": "hello\u0000world",
  "url": "https://www.productdrive.io/grow/world-class"
}
```

HTTP 400

