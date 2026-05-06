# ══════════════════════════════════════════════════
# auth chain: GET /api/diff
# case_id=TC-6af54553
# case_name=auth chain: GET /api/diff
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

# ── GET /api/diff with auth token [test] ────
# step_id=step-test
# step_type=test
# title=GET /api/diff with auth token
# depends_on=step-auth

GET {{base_url}}/api/diff
Authorization: Bearer {{authToken}}

HTTP 200

[Asserts]
duration < 2000

