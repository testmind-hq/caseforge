# ══════════════════════════════════════════════════
# sequence chain: /api/tokens → PUT /api/admin/users/{id}
# case_id=TC-3028e37b
# case_name=sequence chain: /api/tokens → PUT /api/admin/users/{id}
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
  "name": "Dante Kennedy",
  "scope": "write"
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
  "isActive": true,
  "role": "super_admin"
}
```

HTTP *

[Asserts]
status < 300

