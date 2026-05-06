# ══════════════════════════════════════════════════
# sequence chain: /api/tokens → GET /api/admin/teams/{id}/grants
# case_id=TC-f107e18d
# case_name=sequence chain: /api/tokens → GET /api/admin/teams/{id}/grants
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
  "name": "Janie Stone",
  "scope": "write"
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

