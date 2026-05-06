# ── PATCH /api/admin/webhooks/:id - [unicode_fuzzing] url zalgo ──
# case_id=TC-0a72a45e
# case_name=PATCH /api/admin/webhooks/:id - [unicode_fuzzing] url zalgo
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
  "url": "z̀́̂̃̄̅̆̇a"
}
```

HTTP 400

