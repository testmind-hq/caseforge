# ── [OWASP-API3] PATCH /api/admin/webhooks/:id — BOPLA property-level access ──
# case_id=TC-d7a97bb7
# case_name=[OWASP-API3] PATCH /api/admin/webhooks/:id — BOPLA property-level access
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

PATCH {{base_url}}/api/admin/webhooks/:id
Content-Type: application/json
```json
{
  "events": [
    "leap"
  ],
  "isActive": true,
  "is_admin": true,
  "name": "Lacy Mccarthy",
  "role": "admin",
  "url": "http://www.mainrobust.net/user-centric/empower"
}
```

HTTP 200

[Asserts]
jsonpath "$.is_admin" != true
jsonpath "$.role" != "admin"

