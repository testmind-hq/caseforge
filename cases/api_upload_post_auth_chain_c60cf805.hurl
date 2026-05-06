# ══════════════════════════════════════════════════
# auth chain: POST /api/upload
# case_id=TC-c60cf805
# case_name=auth chain: POST /api/upload
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

# ── POST /api/upload with auth token [test] ──
# step_id=step-test
# step_type=test
# title=POST /api/upload with auth token
# depends_on=step-auth

POST {{base_url}}/api/upload
Authorization: Bearer {{authToken}}
Content-Type: application/json
```json
{
  "branch": "they",
  "commitSha": "sometimes",
  "service": "Darwinian",
  "specContent": "i.e."
}
```

HTTP 200

[Asserts]
duration < 2000

