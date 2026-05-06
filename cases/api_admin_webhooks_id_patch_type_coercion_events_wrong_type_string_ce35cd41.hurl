# ── PATCH /api/admin/webhooks/:id - [type_coercion] events wrong_type_string ──
# case_id=TC-ce35cd41
# case_name=PATCH /api/admin/webhooks/:id - [type_coercion] events wrong_type_string
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

PATCH {{base_url}}/api/admin/webhooks/:id
Content-Type: application/json
```json
{
  "events": "not_an_array",
  "isActive": false,
  "name": "Emile Jones",
  "url": "https://www.financeoptimize.com/transform/cross-media/technologies"
}
```

HTTP 422

