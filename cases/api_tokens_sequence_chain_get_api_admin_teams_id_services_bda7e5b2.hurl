# ══════════════════════════════════════════════════
# sequence chain: /api/tokens → GET /api/admin/teams/{id}/services
# case_id=TC-bda7e5b2
# case_name=sequence chain: /api/tokens → GET /api/admin/teams/{id}/services
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
  "name": "Bernadine Murray",
  "scope": "write"
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

