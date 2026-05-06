# ══════════════════════════════════════════════════
# sequence chain: /auth/register → DELETE /api/admin/grants/{id}
# case_id=TC-465a3cf5
# case_name=sequence chain: /auth/register → DELETE /api/admin/grants/{id}
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
  "email": "hollybarker@garza.com",
  "password": "who"
}
```

HTTP *

[Captures]
id: jsonpath "$.userId"

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

