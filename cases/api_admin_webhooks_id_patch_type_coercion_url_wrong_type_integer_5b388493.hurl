# ── PATCH /api/admin/webhooks/:id - [type_coercion] url wrong_type_integer ──
# case_id=TC-5b388493
# case_name=PATCH /api/admin/webhooks/:id - [type_coercion] url wrong_type_integer
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

PATCH {{base_url}}/api/admin/webhooks/:id
Content-Type: application/json
```json
{
  "events": [
    "some"
  ],
  "isActive": false,
  "name": "Emile Jones",
  "url": 123
}
```

HTTP 422

