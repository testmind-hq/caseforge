# ── PATCH /api/admin/webhooks/:id - [unicode_fuzzing] url control_char ──
# case_id=TC-ed68863e
# case_name=PATCH /api/admin/webhooks/:id - [unicode_fuzzing] url control_char
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
  "name": "Nicole Heller",
  "url": "hello\u0000world"
}
```

HTTP 400

