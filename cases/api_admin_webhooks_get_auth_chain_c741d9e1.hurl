# ══════════════════════════════════════════════════
# auth chain: GET /api/admin/webhooks
# case_id=TC-c741d9e1
# case_name=auth chain: GET /api/admin/webhooks
# case_kind=chain
# priority=P1
# ══════════════════════════════════════════════════

# ── authenticate via POST /api/tokens [setup] ──
# step_id=step-auth
# step_type=setup
# title=authenticate via POST /api/tokens

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": "Jakob Jensen",
  "scope": "write"
}
```

HTTP *

[Captures]
authToken: jsonpath "$.token"

[Asserts]
status < 300

# ── GET /api/admin/webhooks with auth token [test] ──
# step_id=step-test
# step_type=test
# title=GET /api/admin/webhooks with auth token
# depends_on=step-auth

GET {{base_url}}/api/admin/webhooks
Authorization: Bearer {{authToken}}

HTTP 200

[Asserts]
duration < 2000

