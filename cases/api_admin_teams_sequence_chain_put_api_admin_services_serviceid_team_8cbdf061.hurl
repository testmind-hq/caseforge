# ══════════════════════════════════════════════════
# sequence chain: /api/admin/teams → PUT /api/admin/services/{serviceId}/team
# case_id=TC-8cbdf061
# case_name=sequence chain: /api/admin/teams → PUT /api/admin/services/{serviceId}/team
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
  "description": "Optimize company for light clarity.",
  "displayName": "many",
  "name": "Christina Patterson"
}
```

HTTP *

[Captures]
serviceId: jsonpath "$.id"

[Asserts]
status < 300

# ── use via PUT /api/admin/services/{serviceId}/team [test] ──
# step_id=step-test
# step_type=test
# title=use via PUT /api/admin/services/{serviceId}/team
# depends_on=step-setup

PUT {{base_url}}/api/admin/services/{{serviceId}}/team
Content-Type: application/json
```json
{
  "teamId": "40d2db88-109b-49a0-8983-e2740333822a"
}
```

HTTP *

[Asserts]
status < 300

