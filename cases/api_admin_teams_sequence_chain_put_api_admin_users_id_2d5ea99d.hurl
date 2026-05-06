# ══════════════════════════════════════════════════
# sequence chain: /api/admin/teams → PUT /api/admin/users/{id}
# case_id=TC-2d5ea99d
# case_name=sequence chain: /api/admin/teams → PUT /api/admin/users/{id}
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
  "description": "Stage number behind feature flags.",
  "displayName": "sew",
  "name": "Stanley Purdy"
}
```

HTTP *

[Captures]
id: jsonpath "$.id"

[Asserts]
status < 300

# ── use via PUT /api/admin/users/{id} [test] ──
# step_id=step-test
# step_type=test
# title=use via PUT /api/admin/users/{id}
# depends_on=step-setup

PUT {{base_url}}/api/admin/users/{{id}}
Content-Type: application/json
```json
{
  "isActive": false,
  "role": "super_admin"
}
```

HTTP *

[Asserts]
status < 300

