# ══════════════════════════════════════════════════
# sequence chain: /api/admin/teams → GET /api/admin/teams/{id}/services
# case_id=TC-3642a068
# case_name=sequence chain: /api/admin/teams → GET /api/admin/teams/{id}/services
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
  "description": "Track problem over time weekly.",
  "displayName": "of",
  "name": "Owen Perez"
}
```

HTTP *

[Captures]
id: jsonpath "$.id"

[Asserts]
status < 300

# ── use via GET /api/admin/teams/{id}/services [test] ──
# step_id=step-test
# step_type=test
# title=use via GET /api/admin/teams/{id}/services
# depends_on=step-setup

GET {{base_url}}/api/admin/teams/{{id}}/services

HTTP *

[Asserts]
status < 300

