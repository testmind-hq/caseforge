# ══════════════════════════════════════════════════
# sequence chain: /api/tokens → POST /api/admin/teams/{id}/members
# case_id=TC-714b8b84
# case_name=sequence chain: /api/tokens → POST /api/admin/teams/{id}/members
# case_kind=chain
# priority=P1
# ══════════════════════════════════════════════════

# ── create via POST /api/tokens [setup] ─────
# step_id=step-setup
# step_type=setup
# title=create via POST /api/tokens

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": "Athena Fernandez",
  "scope": "read"
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
  "role": "member",
  "userId": "02ef8546-0050-41de-be11-ab585b23ac54"
}
```

HTTP *

[Asserts]
status < 300

