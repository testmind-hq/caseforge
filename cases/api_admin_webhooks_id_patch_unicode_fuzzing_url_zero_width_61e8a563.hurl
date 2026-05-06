# ── PATCH /api/admin/webhooks/:id - [unicode_fuzzing] url zero_width ──
# case_id=TC-61e8a563
# case_name=PATCH /api/admin/webhooks/:id - [unicode_fuzzing] url zero_width
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
  "url": "​hello"
}
```

HTTP 400

