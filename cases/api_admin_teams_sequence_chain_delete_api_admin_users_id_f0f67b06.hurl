# ══════════════════════════════════════════════════
# sequence chain: /api/admin/teams → DELETE /api/admin/users/{id}
# case_id=TC-f0f67b06
# case_name=sequence chain: /api/admin/teams → DELETE /api/admin/users/{id}
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
  "description": "Invite review for the group in Birmingham.",
  "displayName": "eventually",
  "name": "Robyn Williams"
}
```

HTTP *

[Captures]
id: jsonpath "$.id"

[Asserts]
status < 300

# ── use via DELETE /api/admin/users/{id} [test] ──
# step_id=step-test
# step_type=test
# title=use via DELETE /api/admin/users/{id}
# depends_on=step-setup

DELETE {{base_url}}/api/admin/users/{{id}}

HTTP *

[Asserts]
status < 300

