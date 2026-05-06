# ══════════════════════════════════════════════════
# auth chain: POST /api/admin/webhooks
# case_id=TC-f4c0b7fc
# case_name=auth chain: POST /api/admin/webhooks
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

# ── POST /api/admin/webhooks with auth token [test] ──
# step_id=step-test
# step_type=test
# title=POST /api/admin/webhooks with auth token
# depends_on=step-auth

POST {{base_url}}/api/admin/webhooks
Authorization: Bearer {{authToken}}
Content-Type: application/json
```json
{
  "events": [
    "where"
  ],
  "name": "Lilla Henderson",
  "providerType": "shirt",
  "teamId": "1e74395d-96d5-4632-bff5-1db94dfc9c0c",
  "url": "http://www.brandengage.info/out-of-the-box/end-to-end/engineer/visualize"
}
```

HTTP 200

[Asserts]
duration < 2000

