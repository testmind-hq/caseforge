# ══════════════════════════════════════════════════
# sequence chain: /api/admin/teams → POST /api/admin/teams/{id}/members
# case_id=TC-210690e6
# case_name=sequence chain: /api/admin/teams → POST /api/admin/teams/{id}/members
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
  "description": "Weekends reserve time for Animation and fact.",
  "displayName": "today",
  "name": "Jeffrey Lyons"
}
```

HTTP *

[Captures]
id: jsonpath "$.id"

[Asserts]
status < 300

# ── use via POST /api/admin/teams/{id}/members [test] ──
# step_id=step-test
# step_type=test
# title=use via POST /api/admin/teams/{id}/members
# depends_on=step-setup

POST {{base_url}}/api/admin/teams/{{id}}/members
Content-Type: application/json
```json
{
  "role": "owner",
  "userId": "45f53f9f-487d-4010-8fff-c2d438433278"
}
```

HTTP *

[Asserts]
status < 300

