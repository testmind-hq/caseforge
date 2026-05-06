# ══════════════════════════════════════════════════
# sequence chain: /api/tokens → DELETE /api/admin/users/{id}
# case_id=TC-60268ad8
# case_name=sequence chain: /api/tokens → DELETE /api/admin/users/{id}
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
  "name": "Rafael Hopkins",
  "scope": "write"
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

