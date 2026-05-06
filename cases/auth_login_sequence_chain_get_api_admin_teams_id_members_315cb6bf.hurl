# ══════════════════════════════════════════════════
# sequence chain: /auth/login → GET /api/admin/teams/{id}/members
# case_id=TC-315cb6bf
# case_name=sequence chain: /auth/login → GET /api/admin/teams/{id}/members
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
  "email": "manuelcasper@owen.net",
  "password": "herself"
}
```

HTTP *

[Captures]
id: jsonpath "$.userId"

[Asserts]
status < 300

# ── use via GET /api/admin/teams/{id}/members [test] ──
# step_id=step-test
# step_type=test
# title=use via GET /api/admin/teams/{id}/members
# depends_on=step-setup

GET {{base_url}}/api/admin/teams/{{id}}/members

HTTP *

[Asserts]
status < 300

