# ── POST /api/admin/webhooks - [type_coercion] name wrong_type_integer ──
# case_id=TC-39c60504
# case_name=POST /api/admin/webhooks - [type_coercion] name wrong_type_integer
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
  "name": 123,
  "providerType": "impress",
  "teamId": "f82e8b17-baa1-4aef-a496-6b602b3c5f43",
  "url": "https://www.productplatforms.com/impactful/turn-key/infrastructures/integrate"
}
```

HTTP 422

