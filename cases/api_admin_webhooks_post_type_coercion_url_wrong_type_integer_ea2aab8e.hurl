# ── POST /api/admin/webhooks - [type_coercion] url wrong_type_integer ──
# case_id=TC-ea2aab8e
# case_name=POST /api/admin/webhooks - [type_coercion] url wrong_type_integer
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
  "teamId": "f82e8b17-baa1-4aef-a496-6b602b3c5f43",
  "url": 123
}
```

HTTP 422

