# ══════════════════════════════════════════════════
# sequence chain: /api/admin/teams → GET /api/admin/teams/{id}/members
# case_id=TC-0cb6ef87
# case_name=sequence chain: /api/admin/teams → GET /api/admin/teams/{id}/members
# case_kind=chain
# priority=P1
# ══════════════════════════════════════════════════

# ── create via POST /api/admin/teams [setup] ──
# step_id=step-setup
# step_type=setup
# title=create via POST /api/admin/teams

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": "Mind the hand, then celebrate!",
  "displayName": "ride",
  "name": "Dolores Grady"
}
```

HTTP *

[Captures]
id: jsonpath "$.id"

[Asserts]
status < 300

# ── use via GET /api/admin/teams/{id}/members [test] ──
# step_id=step-test
# step_type=test
# title=use via GET /api/admin/teams/{id}/members
# depends_on=step-setup

GET {{base_url}}/api/admin/teams/{{id}}/members

HTTP *

[Asserts]
status < 300

