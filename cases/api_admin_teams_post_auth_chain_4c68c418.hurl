# ══════════════════════════════════════════════════
# auth chain: POST /api/admin/teams
# case_id=TC-4c68c418
# case_name=auth chain: POST /api/admin/teams
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

# ── POST /api/admin/teams with auth token [test] ──
# step_id=step-test
# step_type=test
# title=POST /api/admin/teams with auth token
# depends_on=step-auth

POST {{base_url}}/api/admin/teams
Authorization: Bearer {{authToken}}
Content-Type: application/json
```json
{
  "description": "The government should confusing.",
  "displayName": "yours",
  "name": "Lee Burton"
}
```

HTTP 200

[Asserts]
duration < 2000

