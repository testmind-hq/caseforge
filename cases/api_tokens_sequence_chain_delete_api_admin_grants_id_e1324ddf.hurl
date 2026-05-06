# ══════════════════════════════════════════════════
# sequence chain: /api/tokens → DELETE /api/admin/grants/{id}
# case_id=TC-e1324ddf
# case_name=sequence chain: /api/tokens → DELETE /api/admin/grants/{id}
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
  "name": "Bernardo Auer",
  "scope": "write"
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

