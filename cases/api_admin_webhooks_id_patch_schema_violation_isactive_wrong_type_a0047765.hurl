# ── PATCH /api/admin/webhooks/:id - [schema_violation] isActive_wrong_type ──
# case_id=TC-a0047765
# case_name=PATCH /api/admin/webhooks/:id - [schema_violation] isActive_wrong_type
# step_id=step-main
# step_type=test
# technique=schema_violation
# priority=P2

PATCH {{base_url}}/api/admin/webhooks/:id
Content-Type: application/json
```json
{
  "events": [
    "whatever"
  ],
  "isActive": "not_a_boolean",
  "name": "Alexander Gordon",
  "url": "https://www.grouptechnologies.net/deliverables/web-enabled/generate/e-enable"
}
```

HTTP 422

