# ══════════════════════════════════════════════════
# sequence chain: /auth/register → DELETE /api/admin/users/{id}
# case_id=TC-b3bffa74
# case_name=sequence chain: /auth/register → DELETE /api/admin/users/{id}
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
  "email": "fannystevenson@daugherty.com",
  "password": "way"
}
```

HTTP *

[Captures]
id: jsonpath "$.userId"

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

