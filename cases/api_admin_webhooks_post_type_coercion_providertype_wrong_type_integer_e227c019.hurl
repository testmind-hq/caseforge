# ── POST /api/admin/webhooks - [type_coercion] providerType wrong_type_integer ──
# case_id=TC-e227c019
# case_name=POST /api/admin/webhooks - [type_coercion] providerType wrong_type_integer
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
  "providerType": 123,
  "teamId": "f82e8b17-baa1-4aef-a496-6b602b3c5f43",
  "url": "https://www.productplatforms.com/impactful/turn-key/infrastructures/integrate"
}
```

HTTP 422

