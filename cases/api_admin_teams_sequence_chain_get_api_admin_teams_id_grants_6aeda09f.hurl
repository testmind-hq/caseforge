# ══════════════════════════════════════════════════
# sequence chain: /api/admin/teams → GET /api/admin/teams/{id}/grants
# case_id=TC-6aeda09f
# case_name=sequence chain: /api/admin/teams → GET /api/admin/teams/{id}/grants
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
  "description": "They ski patiently to stabilize the year.",
  "displayName": "fiercely",
  "name": "Cassandra Robbins"
}
```

HTTP *

[Captures]
id: jsonpath "$.id"

[Asserts]
status < 300

# ── use via GET /api/admin/teams/{id}/grants [test] ──
# step_id=step-test
# step_type=test
# title=use via GET /api/admin/teams/{id}/grants
# depends_on=step-setup

GET {{base_url}}/api/admin/teams/{{id}}/grants

HTTP *

[Asserts]
status < 300

