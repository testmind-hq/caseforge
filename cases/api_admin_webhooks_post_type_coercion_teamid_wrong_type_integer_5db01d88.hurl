# ── POST /api/admin/webhooks - [type_coercion] teamId wrong_type_integer ──
# case_id=TC-5db01d88
# case_name=POST /api/admin/webhooks - [type_coercion] teamId wrong_type_integer
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

POST {{base_url}}/api/admin/webhooks
Content-Type: application/json
```json
{
  "events": [
    "improvised"
  ],
  "name": "Horace Evans",
  "providerType": "impress",
  "teamId": 123,
  "url": "https://www.productplatforms.com/impactful/turn-key/infrastructures/integrate"
}
```

HTTP 422

