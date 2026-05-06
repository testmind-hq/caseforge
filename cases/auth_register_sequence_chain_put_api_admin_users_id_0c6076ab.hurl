# ══════════════════════════════════════════════════
# sequence chain: /auth/register → PUT /api/admin/users/{id}
# case_id=TC-0c6076ab
# case_name=sequence chain: /auth/register → PUT /api/admin/users/{id}
# case_kind=chain
# priority=P1
# ══════════════════════════════════════════════════

# ── create via POST /auth/register [setup] ──
# step_id=step-setup
# step_type=setup
# title=create via POST /auth/register

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": "sharonwright@dietrich.org",
  "password": "it"
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
  "isActive": false,
  "role": "team_owner"
}
```

HTTP *

[Asserts]
status < 300

