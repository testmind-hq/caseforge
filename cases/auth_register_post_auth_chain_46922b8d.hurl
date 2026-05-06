# ══════════════════════════════════════════════════
# auth chain: POST /auth/register
# case_id=TC-46922b8d
# case_name=auth chain: POST /auth/register
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

# ── POST /auth/register with auth token [test] ──
# step_id=step-test
# step_type=test
# title=POST /auth/register with auth token
# depends_on=step-auth

POST {{base_url}}/auth/register
Authorization: Bearer {{authToken}}
Content-Type: application/json
```json
{
  "email": "edbarber@reyes.name",
  "password": "nest"
}
```

HTTP 200

[Asserts]
duration < 2000

