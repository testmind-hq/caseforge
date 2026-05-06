# ══════════════════════════════════════════════════
# sequence chain: /api/admin/teams → DELETE /api/admin/grants/{id}
# case_id=TC-70b060a1
# case_name=sequence chain: /api/admin/teams → DELETE /api/admin/grants/{id}
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
  "description": "Ours child be ready for irritation.",
  "displayName": "daily",
  "name": "Cordell Marshall"
}
```

HTTP *

[Captures]
id: jsonpath "$.id"

[Asserts]
status < 300

# ── use via DELETE /api/admin/grants/{id} [test] ──
# step_id=step-test
# step_type=test
# title=use via DELETE /api/admin/grants/{id}
# depends_on=step-setup

DELETE {{base_url}}/api/admin/grants/{{id}}

HTTP *

[Asserts]
status < 300

