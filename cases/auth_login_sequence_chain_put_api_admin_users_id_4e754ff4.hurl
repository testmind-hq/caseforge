# ══════════════════════════════════════════════════
# sequence chain: /auth/login → PUT /api/admin/users/{id}
# case_id=TC-4e754ff4
# case_name=sequence chain: /auth/login → PUT /api/admin/users/{id}
# case_kind=chain
# priority=P1
# ══════════════════════════════════════════════════

# ── create via POST /auth/login [setup] ─────
# step_id=step-setup
# step_type=setup
# title=create via POST /auth/login

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "email": "amparoknight@evans.biz",
  "password": "always"
}
```

HTTP *

[Captures]
id: jsonpath "$.userId"

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
  "role": "team_owner"
}
```

HTTP *

[Asserts]
status < 300

