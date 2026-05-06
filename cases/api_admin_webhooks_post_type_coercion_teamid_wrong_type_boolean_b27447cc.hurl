# ── POST /api/admin/webhooks - [type_coercion] teamId wrong_type_boolean ──
# case_id=TC-b27447cc
# case_name=POST /api/admin/webhooks - [type_coercion] teamId wrong_type_boolean
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
  "teamId": true,
  "url": "https://www.productplatforms.com/impactful/turn-key/infrastructures/integrate"
}
```

HTTP 422

